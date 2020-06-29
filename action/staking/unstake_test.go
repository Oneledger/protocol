package staking

import (
	"testing"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
)

func assemblyApplyUnstakeData(stake int64, feeAmt int64) action.SignedTx {

	amt := balance.NewAmountFromInt(stake)
	amount := &action.Amount{
		Currency: "OLT",
		Value:    *amt,
	}

	av := &Unstake{
		StakeAddress:     from.Bytes(),
		Stake:            *amount,
		ValidatorAddress: from.Bytes(),
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

func initCheckUnstake(t *testing.T, ctx *action.Context, balToValidate int64) {
	ctx.Delegators.SetValidatorDelegationAmount(from.Bytes(), from.Bytes(), *balance.NewAmountFromInt(1))
	ctx.Delegators.SetDelegatorEffectiveAmount(from.Bytes(), *balance.NewAmountFromInt(1))
	ctx.Delegators.SetValidatorAmount(from.Bytes(), *balance.NewAmountFromInt(1))
	// 1 block
	ctx.Govern.SetStakingOptions(delegation.Options{
		MaturityTime:            1,
		MinSelfDelegationAmount: *balance.NewAmountFromInt(1),
	})

	amt, _ := ctx.Delegators.GetValidatorDelegationAmount(from.Bytes(), from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(1)), "Got GetValidatorDelegationAmount: %v", amt)

	amt, _ = ctx.Delegators.GetValidatorAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(1)), "Got GetValidatorAmount: %v", amt)

	amt, _ = ctx.Delegators.GetDelegatorEffectiveAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(1)), "Got GetDelegatorEffectiveAmount: %v", amt)

	amt, _ = ctx.Delegators.GetDelegatorBoundedAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorBoundedAmount: %v", amt)

	val := getBalanceFromAddress(ctx, from)
	assert.True(t, val.Equals(*balance.NewAmountFromInt(balToValidate)), "Got balance on address %s  - %s, required - %d", from.String(), val.String(), balToValidate)

	height := ctx.Header.GetHeight()
	assert.True(t, height == 0)

	tt, _ := ctx.Govern.GetStakingOptions()
	assert.True(t, tt.MaturityTime == 1)
	assert.True(t, tt.MinSelfDelegationAmount.Equals(*balance.NewAmountFromInt(1)), "Got MinSelfDelegationAmount: %v", amt)

	mature, _ := ctx.Delegators.GetMatureAmounts(height + tt.MaturityTime)
	assert.True(t, len(mature.Data) == 0)

	validator, _ := ctx.Validators.Get(from.Bytes())
	assert.True(t, validator.Power == 0)
	assert.True(t, validator.Staking.Equals(*balance.NewAmountFromInt(0)))

	validator.Power = 1
	validator.Staking = *balance.NewAmountFromInt(1)
	ctx.Validators.Set(*validator)
}

func TestUnstakeTx_ProcessDeliver_OK(t *testing.T) {
	ast := &unstakeTx{}
	ctx := &action.Context{}

	t.Run("unstake with valid and existing amount, should return ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyApplyUnstakeData(1, 10000000000)

		// check init data
		initCheckUnstake(t, ctx, initAmt)

		// simulate validate
		ok, err := ast.Validate(ctx, tx)
		assert.True(t, ok, err)

		// simulate sc execution
		ok, resp := ast.ProcessDeliver(ctx, tx.RawTx)
		assert.True(t, ok, resp)

		// check post data

		amt, _ := ctx.Delegators.GetValidatorDelegationAmount(from.Bytes(), from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetValidatorDelegationAmount: %v", amt)

		amt, _ = ctx.Delegators.GetValidatorAmount(from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetValidatorAmount: %v", amt)

		amt, _ = ctx.Delegators.GetDelegatorEffectiveAmount(from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorEffectiveAmount: %v", amt)

		amt, _ = ctx.Delegators.GetDelegatorBoundedAmount(from.Bytes())
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorBoundedAmount: %v", amt)

		options, _ := ctx.Govern.GetStakingOptions()
		assert.True(t, options.MinSelfDelegationAmount.Equals(*balance.NewAmountFromInt(1)), "Got inSelfDelegationAmount: %v", amt)

		val := getBalanceFromAddress(ctx, from)
		requiredBal := initAmt
		assert.True(t, val.Equals(*balance.NewAmountFromInt(requiredBal)), "Got balance on address %s  - %s, required - %d", from.String(), val.String(), requiredBal)

		height := ctx.Header.GetHeight()
		assert.True(t, height == 0)

		tt, _ := ctx.Govern.GetStakingOptions()
		assert.True(t, tt.MaturityTime == 1)

		mature, _ := ctx.Delegators.GetMatureAmounts(height + tt.MaturityTime)
		amt = &mature.Data[0].Amount
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(1)), "Got GetMatureAmount: %v", amt)

		validator, _ := ctx.Validators.Get(from.Bytes())
		assert.True(t, validator.Power == 0)
		assert.True(t, validator.Staking.Equals(*balance.NewAmountFromInt(0)))
	})

}

func TestUnstakeTx_ProcessDeliver_Error(t *testing.T) {
	ast := &unstakeTx{}
	ctx := &action.Context{}

	t.Run("unstake with an amount greater than on balance, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyApplyUnstakeData(100, 10)

		// check init data
		initCheckUnstake(t, ctx, initAmt)

		// simulate validate
		ok, _ := ast.Validate(ctx, tx)
		assert.False(t, ok)
	})

	t.Run("unstake with an zero amount, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyApplyUnstakeData(0, 10)

		// check init data
		initCheckUnstake(t, ctx, initAmt)

		// simulate validate
		ok, _ := ast.Validate(ctx, tx)
		assert.False(t, ok)
	})
}

func TestUpdateWithdrawRewardFunc(t *testing.T) {
	t.Run("check UpdateWithdrawReward to work properly, should return ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		db := db.NewDB("test", db.MemDBBackend, "")
		cs := storage.NewState(storage.NewChainState("balance", db))
		delegators := delegation.NewDelegationStore("tst", cs)

		height := int64(3)

		var addr keys.Address = from.Bytes()

		mature := &delegation.MatureBlock{
			Height: height,
			Data: []*delegation.MatureData{
				{Address: addr, Amount: *balance.NewAmountFromInt(1), Height: height},
			},
		}
		delegators.SetMatureAmounts(height, mature)

		mature, _ = delegators.GetMatureAmounts(height)
		assert.True(t, len(mature.Data) == 1)

		amt, _ := delegators.GetDelegatorBoundedAmount(addr)
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorBoundedAmount: %v", amt)

		delegators.UpdateWithdrawReward(height)

		mature, _ = delegators.GetMatureAmounts(height)
		assert.True(t, len(mature.Data) == 0, "Got GetMatureAmount: %v", amt)

		amt, _ = delegators.GetDelegatorBoundedAmount(addr)
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(1)), "Got GetDelegatorBoundedAmount: %v", amt)
	})

	t.Run("check UpdateWithdrawReward without tokens, should return ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		db := db.NewDB("test", db.MemDBBackend, "")
		cs := storage.NewState(storage.NewChainState("balance", db))
		delegators := delegation.NewDelegationStore("tst", cs)

		height := int64(3)

		var addr keys.Address = from.Bytes()

		mature, _ := delegators.GetMatureAmounts(height)
		assert.True(t, len(mature.Data) == 0)

		amt, _ := delegators.GetDelegatorBoundedAmount(addr)
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorBoundedAmount: %v", amt)

		delegators.UpdateWithdrawReward(height)

		mature, _ = delegators.GetMatureAmounts(height)
		assert.True(t, len(mature.Data) == 0)

		amt, _ = delegators.GetDelegatorBoundedAmount(addr)
		assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorBoundedAmount: %v", amt)
	})
}
