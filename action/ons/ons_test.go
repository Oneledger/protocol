package ons

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/magiconair/properties/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	db2 "github.com/tendermint/tendermint/libs/db"
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

		opt := &ons.Options{PerBlockFees: 50000, FirstLevelDomain: []string{"ol"}, BaseDomainPrice: olt.NewCoinFromUnit(1000000)}

		bs.AddToAddress(h1.Address(), olt.NewCoinFromInt(500))

		header := &abci.Header{Height: 0}
		ctx = action.NewContext(nil, header, cs, nil, bs, currencies,
			feeOpt, feePool, nil, ds, nil, nil, nil, nil,
			nil, "", "", logger, opt)
	}
}

const (
	create1 = iota
	create2
	create3
)

func makeCreateRawTx(name ons.Name, amount balance.Amount, expiryHeight int64) action.RawTx {

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()

	Amount := action.NewAmount("olt", amount)

	tx := &DomainCreate{
		h1.Address(),
		h2.Address(),
		name,
		*Amount,
		"http://hashard.ol/lookatme",
		expiryHeight,
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

func setupCreateDomain() {

	price, _ := balance.NewAmountFromString(createPrice, 10)
	rawTx := makeCreateRawTx("hashard.ol", *price, int64(500000000))

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

	testCases[create1] = Case{ctx, &rawTx, signedTx, action.Gas(50000), action.Gas(100000), action.DOMAIN_CREATE, true, true, true, true}
}

func init() {
	setup()
	setupCreateDomain()
}

func TestONSTx(t *testing.T) {

	for i, testCase := range testCases {
		t.Run("Testing case "+strconv.Itoa(i), func(t *testing.T) {
			//Get handler
			handler := router.Handler(testCase.txType)

			//Call transaction Validate, Then assert response.
			response, err := handler.Validate(testCase.ctx, testCase.signedTx)
			assert.Equal(t, response, testCase.validateResp)
			if err != nil {
				fmt.Println(err.Error())
			}

			response, resp := handler.ProcessCheck(testCase.ctx, *testCase.tx)
			assert.Equal(t, response, testCase.checkResp)
			fmt.Println(resp.Log)

			response, resp = handler.ProcessDeliver(testCase.ctx, *testCase.tx)
			assert.Equal(t, response, testCase.deliverResp)
			fmt.Println(resp.Log)

		})
	}
}
