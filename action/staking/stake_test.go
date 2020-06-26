package staking

import (
	"errors"
	"math/big"
	"os"
	"testing"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/bytes"
	db "github.com/tendermint/tm-db"
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
	from, fromPubkey, fromPrikey = generateKeyPair()
)

func assemblyCtxData(currencyName string, currencyDecimal int, setStore bool, setLogger bool, setCoin int64) *action.Context {

	ctx := &action.Context{}
	db := db.NewDB("test", db.MemDBBackend, "")
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
	store.State.Commit()

	ctx.FeeOpt = &fees.FeeOption{
		FeeCurrency:   currency,
		MinFeeDecimal: 9,
	}
	ctx.FeePool = fees.NewStore("tf", cs)
	ctx.FeePool.SetupOpt(ctx.FeeOpt)
	ctx.Govern = governance.NewStore("tg", cs)
	ctx.Delegators = delegation.NewDelegationStore("tst", cs)
	ctx.Validators = identity.NewValidatorStore("tv", cs)

	validator := identity.NewValidator(
		from.Bytes(),
		from.Bytes(),
		keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
		keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
		*balance.NewAmountFromInt(0),
		"test_node",
	)
	_ = ctx.Validators.Set(*validator)

	delop := delegation.Options{
		MinSelfDelegationAmount: *balance.NewAmountFromInt(1),
	}
	ctx.Govern.SetStakingOptions(delop)
	return ctx
}

func assemblyStakeData(stake int64, feeAmt int64) action.SignedTx {
	amt := balance.NewAmountFromInt(stake)
	amount := &action.Amount{
		Currency: "OLT",
		Value:    *amt,
	}

	av := &Stake{
		StakeAddress:     from.Bytes(),
		Stake:            *amount,
		NodeName:         "test_node",
		ValidatorAddress: from.Bytes(),
		ValidatorPubKey:  keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
	}
	fee := action.Fee{
		Price: action.Amount{"OLT", *balance.NewAmount(feeAmt)},
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
			{
				Signer: keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
				Signed: signature,
			},
		},
	}
	return signed
}

func getBalanceFromAddress(ctx *action.Context, addr bytes.HexBytes) balance.Amount {
	cs := ctx.Currencies.GetCurrencies().GetCurrencySet()
	currency, _ := cs.GetCurrencyByName("OLT")
	bal, _ := ctx.Balances.GetBalance(from.Bytes(), cs)
	val := bal.GetCoin(currency).Amount
	finalAmt := big.NewInt(0).Div(val.BigInt(), currency.Base())
	val = (*balance.Amount)(finalAmt)
	return *val
}

func initCheckStake(t *testing.T, ctx *action.Context, balToValidate int64) {
	amt, _ := ctx.Delegators.GetValidatorDelegationAmount(from.Bytes(), from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetValidatorDelegationAmount: %v", amt)

	amt, _ = ctx.Delegators.GetValidatorAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetValidatorAmount: %v", amt)

	amt, _ = ctx.Delegators.GetDelegatorEffectiveAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorEffectiveAmount: %v", amt)

	amt, _ = ctx.Delegators.GetDelegatorBoundedAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorBoundedAmount: %v", amt)

	tt, _ := ctx.Govern.GetStakingOptions()
	assert.True(t, tt.MinSelfDelegationAmount.Equals(*balance.NewAmountFromInt(1)), "Got MinSelfDelegationAmount: %v", amt)

	val := getBalanceFromAddress(ctx, from)
	assert.True(t, val.Equals(*balance.NewAmountFromInt(balToValidate)), "Got balance on address %s  - %s, required - %d", from.String(), val.String(), balToValidate)

	validator, _ := ctx.Validators.Get(from.Bytes())
	assert.True(t, validator.Power == 0)
}

func TestStakeTx_ProcessDeliver_OK(t *testing.T) {
	ast := &stakeTx{}
	ctx := &action.Context{}

	t.Run("stake with valid and existing amount, should return ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyStakeData(1, 10000000000)

		// check init data
		initCheckStake(t, ctx, initAmt)

		// simulate validate
		ok, err := ast.Validate(ctx, tx)
		assert.True(t, ok, err)

		// simulate sc execution
		ok, resp := ast.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok, resp)

		// check post data

		// diff amount on balances
		var diff int64 = 1

		amt, _ := ctx.Delegators.GetValidatorDelegationAmount(from.Bytes(), from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(diff)), "Got GetValidatorDelegationAmount: %v", amt)

		amt, _ = ctx.Delegators.GetValidatorAmount(from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(diff)), "Got GetValidatorAmount: %v", amt)

		amt, _ = ctx.Delegators.GetDelegatorEffectiveAmount(from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(diff)), "Got GetDelegatorEffectiveAmount: %v", amt)

		amt, _ = ctx.Delegators.GetDelegatorBoundedAmount(from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorBoundedAmount: %v", amt)

		options, _ := ctx.Govern.GetStakingOptions()
		assert.True(t, options.MinSelfDelegationAmount.Equals(*balance.NewAmountFromInt(1)), "Got MinSelfDelegationAmount: %v", amt)

		val := getBalanceFromAddress(ctx, from)
		requiredBal := initAmt - diff
		assert.True(t, val.Equals(*balance.NewAmountFromInt(requiredBal)), "Got balance on address %s  - %s, required - %d", from.String(), val.String(), requiredBal)

		validator, _ := ctx.Validators.Get(from.Bytes())
		assert.True(t, validator.Power == diff)
	})

}

func TestStakeTx_ProcessDeliver_Error(t *testing.T) {
	ast := &stakeTx{}
	ctx := &action.Context{}

	t.Run("stake with an amount greater than on balance, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyStakeData(100, 10)

		// check init data
		initCheckStake(t, ctx, initAmt)

		// simulate validate
		ok, _ := ast.Validate(ctx, tx)
		assert.False(t, ok)
	})

	t.Run("stake with an zero amount, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyStakeData(0, 10)

		// check init data
		initCheckStake(t, ctx, initAmt)

		// simulate validate
		ok, _ := ast.Validate(ctx, tx)
		assert.False(t, ok)
	})
}
