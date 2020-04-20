package ons

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/magiconair/properties/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	db2 "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

var (
	testCases map[int]Case
	router    action.Router
	ctx       *action.Context

	createPrice = "200000000000000000000"

	pub1, priv1, _ = keys.NewKeyPairFromTendermint()
	h1, _          = pub1.GetHandler()
)

type Case struct {
	//input values
	ctx      *action.Context
	tx       *action.RawTx
	signedTx action.SignedTx
	startGas action.Gas
	endGas   action.Gas
	txType   action.Type

	//expected output
	validateResp bool
	checkResp    bool
	deliverResp  bool
	feeResp      bool
}

func setup() {
	testCases = make(map[int]Case)
	router = action.NewRouter("test")
	_ = EnableONS(router)

	//create and initialize new context
	{
		db := db2.NewDB("test", db2.MemDBBackend, "")
		cs := storage.NewState(storage.NewChainState("test", db))

		domain, err := ons.NewDomain(keys.Address("abcd"), nil, "name.ol", "", 22, "", 10)

		fmt.Println(err)

		ds := ons.NewDomainStore("d", cs)
		ds.Set(domain)

		bs := balance.NewStore("b", cs)
		feePool := fees.NewStore("f", cs)

		logger := log.NewLoggerWithPrefix(os.Stdout, "test_action_ons")

		ct, _ := chain.TypeFromName("OneLedger")

		olt := balance.Currency{
			0, "olt", ct, 18, "ones",
		}

		currencies := balance.NewCurrencySet()
		err = currencies.Register(olt)
		fmt.Println(err)

		feeOpt := &fees.FeeOption{FeeCurrency: olt, MinFeeDecimal: 9}
		feePool.SetupOpt(feeOpt)

		opt := &ons.Options{PerBlockFees: 50000, FirstLevelDomains: []string{"ol"}, BaseDomainPrice: olt.NewCoinFromUnit(1000000)}

		bs.AddToAddress(h1.Address(), olt.NewCoinFromInt(500))

		header := &abci.Header{Height: 0}
		ctx = action.NewContext(nil, header, cs, nil, bs, currencies,
			feeOpt, feePool, nil, ds, nil, nil, nil, nil,
			nil, "", "", logger, opt)
	}
}

const (
	create = iota
	update
)

func makeCreateRawTx(name ons.Name, amount balance.Amount, buyingPrice int64) action.RawTx {

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()

	Amount := action.NewAmount("olt", amount)

	tx := &DomainCreate{
		h1.Address(),
		h2.Address(),
		name,
		*Amount,
		"http://hashard.ol/lookatme",
		buyingPrice,
	}

	msg, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("Error marshalling transaction: ", err.Error())
		return action.RawTx{}
	}

	rawTx := action.RawTx{
		Type: action.DOMAIN_CREATE,
		Data: msg,
		Fee: action.Fee{
			Price: *action.NewAmount("olt", *balance.NewAmount(int64(10000000000))),
			Gas:   0,
		},
		Memo: "test1",
	}
	return rawTx
}

func makeUpdateRawTx(name ons.Name, extendPrice int64) (action.RawTx, error) {
	pubBen, _, err := keys.NewKeyPairFromTendermint()
	if err != nil {
		return action.RawTx{}, nil
	}
	pubBenH, err := pubBen.GetHandler()
	if err != nil {
		return action.RawTx{}, nil
	}
	updateTx := DomainUpdate{
		Owner:        h1.Address(),
		Beneficiary:  pubBenH.Address(),
		Name:         name,
		Active:       false,
		ExtendExpiry: extendPrice,
	}
	msg, err := json.Marshal(updateTx)
	if err != nil {
		return action.RawTx{}, nil
	}
	rawTx := action.RawTx{
		Type: action.DOMAIN_UPDATE,
		Data: msg,
		Fee: action.Fee{
			Price: *action.NewAmount("olt", *balance.NewAmount(int64(10000000000))),
			Gas:   0,
		},
		Memo: "testUpdate",
	}
	return rawTx, nil

}

func createTestCase(rawTx action.RawTx, t action.Type, testcaseType int) {
	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(rawTx)
	if err != nil {
		fmt.Println(err.Error())
	}

	a, _ := priv1.GetHandler()
	signed, _ := a.Sign(packet)
	signature := action.Signature{Signed: signed, Signer: pub1}

	signedTx := action.SignedTx{
		RawTx:      rawTx,
		Signatures: []action.Signature{signature},
	}

	testCases[testcaseType] = Case{ctx, &rawTx, signedTx, action.Gas(50000), action.Gas(100000), t, true, true, true, true}
}
func setupCreateDomain() {

	price, _ := balance.NewAmountFromString(createPrice, 10)
	rawTx := makeCreateRawTx("harshad.ol", *price, int64(500000000))
	createTestCase(rawTx, action.DOMAIN_CREATE, create)
}

func setupUpdateDomain() {
	rawTx, err := makeUpdateRawTx("harshad.ol", int64(800000000))
	if err != nil {
		return
	}
	createTestCase(rawTx, action.DOMAIN_UPDATE, update)
}

func init() {
	setup()
	setupCreateDomain()
	setupUpdateDomain()
}

func TestONSTx(t *testing.T) {

	for i := 0; i < len(testCases); i++ {
		testCase := testCases[i]
		t.Run("Testing case "+strconv.Itoa(i), func(t *testing.T) {
			//Get handler
			handler := router.Handler(testCase.txType)
			fmt.Println("-----------------------------------------------------------")
			//Call transaction Validate, Then assert response.
			response, err := handler.Validate(testCase.ctx, testCase.signedTx)
			assert.Equal(t, response, testCase.validateResp)
			if err != nil {
				fmt.Println("Validate", err.Error())
			}

			//response, resp := handler.ProcessCheck(testCase.ctx, *testCase.tx)
			//assert.Equal(t, response, testCase.checkResp)
			//fmt.Println("ProcessCheck :" , resp.Log)

			response, resp := handler.ProcessDeliver(testCase.ctx, *testCase.tx)
			assert.Equal(t, response, testCase.deliverResp)
			fmt.Println("ProcessDeliver :", resp.Log)

			ctx.State.Commit()
			//fmt.Println("height", hex.EncodeToString(ctx.State.RootHash()), ctx.State.Version())
		})
	}
}

func TestUpdate(t *testing.T) {
	testCase := testCases[0]

	handler := router.Handler(testCase.txType)
	//Call transaction Validate, Then assert response.
	response, err := handler.Validate(ctx, testCase.signedTx)
	assert.Equal(t, response, testCase.validateResp, "validate")
	if err != nil {
		fmt.Println("Validate", err.Error())
	}

	//response, resp := handler.ProcessCheck(testCase.ctx, *testCase.tx)
	//assert.Equal(t, response, testCase.checkResp)
	//fmt.Println("ProcessCheck :" , resp.Log)

	response, resp := handler.ProcessDeliver(ctx, *testCase.tx)
	assert.Equal(t, response, testCase.deliverResp, "deliver")
	fmt.Println("ProcessDeliver :", resp.Log)
	ctx.State.Commit()
	fmt.Println("height", hex.EncodeToString(ctx.State.RootHash()), ctx.State.Version())
	fmt.Println("-----------------------------------------------------------------------------------------------------")
	testCase = testCases[1]

	ctx.Header.Height = ctx.State.Version() + 1
	handler = router.Handler(testCase.txType)
	//Call transaction Validate, Then assert response.
	//t.Run("Testing case "+strconv.Itoa(1), func(t *testing.T) {
	response, err = handler.Validate(ctx, testCase.signedTx)
	assert.Equal(t, response, testCase.validateResp)
	if err != nil {
		fmt.Println("Validate", err.Error())
	}

	//response, resp = handler.ProcessCheck(ctx, *testCase.tx)
	//assert.Equal(t, response, testCase.checkResp)
	//fmt.Println("ProcessCheck :" , resp.Log)

	response, resp = handler.ProcessDeliver(ctx, *testCase.tx)
	assert.Equal(t, response, testCase.deliverResp)
	fmt.Println("ProcessDeliver :", resp.Log)
	ctx.State.Commit()
	fmt.Println("height", hex.EncodeToString(ctx.State.RootHash()), ctx.State.Version())

}
