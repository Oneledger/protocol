package penalization

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

// global setup
func setup() []string {
	testDBs := []string{"test_dbpath"}
	return testDBs
}

// remove test db dir after
func teardown(dbPaths []string) {
	for _, v := range dbPaths {
		err := os.RemoveAll(v)
		if err != nil {
			errors.New("Remove test db file error")
		}
	}
}

func generateKeyPair() (crypto.Address, crypto.PubKey, ed25519.PrivKeyEd25519) {

	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	return addr, pubkey, prikey
}

var (
	from, fromPubkey, fromPrikey          = generateKeyPair()
	fromMal, fromPubkeyMal, fromPrikeyMal = generateKeyPair()
)

func assemblyCtxData(currencyName string, setCoin int64) *action.Context {
	setStore := true
	setLogger := true
	ctx := &action.Context{}
	// header
	ctx.Header = &abci.Header{Height: 1}
	db := db.NewDB("test", db.MemDBBackend, "test_dbpath")
	cs := storage.NewState(storage.NewChainState("balance", db))
	// store
	var store *balance.Store
	if setStore {
		store = balance.NewStore("tb", cs)
		ctx.Balances = store
	}
	// logger
	if setLogger {
		ctx.Logger = new(log.Logger)
	}
	// register new token OTT
	currencyList := balance.NewCurrencySet()
	currency := balance.Currency{
		Id:      0,
		Name:    currencyName,
		Chain:   chain.ONELEDGER,
		Decimal: 18,
		Unit:    "nue",
	}
	err := currencyList.Register(currency)
	if err != nil {
		errors.New("register new token error")
	}
	ctx.Currencies = currencyList

	// set coin for account
	coin := currency.NewCoinFromInt(setCoin)
	err = store.AddToAddress(from.Bytes(), coin)
	if err != nil {
		errors.New("setup testing token balance error")
	}
	err = store.AddToAddress(fromMal.Bytes(), coin)
	if err != nil {
		errors.New("setup testing token balance error")
	}
	store.State.Commit()

	evidenceOption := evidence.Options{
		MinVotesRequired: 2,
		BlockVotesDiff:   4,

		PenaltyBasePercentage: 30,
		PenaltyBaseDecimals:   100,

		PenaltyBountyPercentage: 50,
		PenaltyBountyDecimals:   100,

		PenaltyBurnPercentage: 50,
		PenaltyBurnDecimals:   100,

		ValidatorReleaseTime:    5,
		ValidatorVotePercentage: 50,
		ValidatorVoteDecimals:   100,

		AllegationPercentage: 50,
		AllegationDecimals:   100,
	}

	ctx.FeeOpt = &fees.FeeOption{
		FeeCurrency:   currency,
		MinFeeDecimal: 9,
	}
	ctx.FeePool = &fees.Store{}
	ctx.FeePool.SetupOpt(ctx.FeeOpt)
	ctx.GovernanceStore = governance.NewStore("tg", cs)
	ctx.Delegators = delegation.NewDelegationStore("tst", cs)
	ctx.Validators = identity.NewValidatorStore("tv", "purged", cs)
	ctx.EvidenceStore = evidence.NewEvidenceStore("tes", cs)
	ctx.GovernanceStore.SetFeeOption(*ctx.FeeOpt)
	ctx.GovernanceStore.SetEvidenceOptions(evidenceOption)
	ctx.GovernanceStore.WithHeight(0).SetAllLUH()
	validator := identity.NewValidator(
		from.Bytes(),
		from.Bytes(),
		keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
		keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
		*balance.NewAmountFromInt(0),
		"test_node",
	)
	_ = ctx.Validators.Set(*validator)
	_ = ctx.EvidenceStore.SetValidatorStatus(validator.Address, true, 1)

	malVal := identity.NewValidator(
		fromMal.Bytes(),
		fromMal.Bytes(),
		keys.PublicKey{keys.ED25519, fromPubkeyMal.Bytes()[5:]},
		keys.PublicKey{keys.ED25519, fromPubkeyMal.Bytes()[5:]},
		*balance.NewAmountFromInt(0),
		"test_node_2",
	)
	_ = ctx.Validators.Set(*malVal)
	_ = ctx.EvidenceStore.SetValidatorStatus(malVal.Address, true, 1)
	return ctx
}

func assemblyAllegationData(requestID string, blockHeight int64, proofMsg string) action.SignedTx {
	av := &Allegation{
		RequestID:        requestID,
		ValidatorAddress: from.Bytes(),
		MaliciousAddress: fromMal.Bytes(),
		BlockHeight:      blockHeight,
		ProofMsg:         proofMsg,
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(10000000000)},
		Gas:   10,
	}
	data, _ := av.Marshal()
	tx := action.RawTx{
		Type: av.Type(),
		Data: data,
		Fee:  fee,
		Memo: "test_memo",
	}
	signature, _ := fromPrikey.Sign(tx.RawBytes())
	signed := action.SignedTx{
		RawTx: tx,
		Signatures: []action.Signature{
			{
				Signer: keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
				Signed: signature,
			},
		},
	}
	return signed
}

func TestAllegationTx_ProcessDeliver_OK(t *testing.T) {
	atx := &allegationTx{}
	ctx := &action.Context{}

	t.Run("allegation for non frozen validator and receive a request", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		ID := "test"
		ctx = assemblyCtxData("OLT", 1000)
		tx := assemblyAllegationData(ID, 2, "test")

		ar, err := ctx.EvidenceStore.GetAllegationRequest(ID)
		assert.Nil(t, ar)

		at, _ := ctx.EvidenceStore.GetAllegationTracker()
		assert.Equal(t, 0, len(at.Requests))

		ok, err := atx.Validate(ctx, tx)
		assert.NoError(t, err)
		assert.True(t, ok)

		ok, resp := atx.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok, resp)

		ar, err = ctx.EvidenceStore.GetAllegationRequest(ID)
		assert.NoError(t, err)
		assert.Equal(t, "test", ar.ProofMsg)
		assert.Equal(t, int64(2), ar.BlockHeight)

		at, _ = ctx.EvidenceStore.GetAllegationTracker()
		assert.Equal(t, 1, len(at.Requests))
		assert.True(t, at.Requests[ID])
	})
}
