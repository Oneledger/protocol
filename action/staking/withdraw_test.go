package staking

import (
	"testing"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

func assemblyWithdrawData(stake int64, feeAmt int64) action.SignedTx {

	amt := balance.NewAmountFromInt(stake)
	amount := &action.Amount{
		Currency: "OLT",
		Value:    *amt,
	}

	av := &Withdraw{
		StakeAddress:     from.Bytes(),
		Stake:            *amount,
		ValidatorAddress: from.Bytes(),
	}
	fee := action.Fee{
		Price: action.Amount{
			Currency: "OLT",
			Value:    *balance.NewAmount(feeAmt),
		},
		Gas: 10,
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

func initCheckWithdraw(t *testing.T, ctx *action.Context, balToValidate int64) {
	ctx.Delegators.SetDelegatorBoundedAmount(from.Bytes(), *balance.NewAmountFromInt(1))

	amt, _ := ctx.Delegators.GetValidatorDelegationAmount(from.Bytes(), from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetValidatorDelegationAmount: %v", amt)

	amt, _ = ctx.Delegators.GetValidatorAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetValidatorAmount: %v", amt)

	amt, _ = ctx.Delegators.GetDelegatorEffectiveAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(0)), "Got GetDelegatorEffectiveAmount: %v", amt)

	amt, _ = ctx.Delegators.GetDelegatorBoundedAmount(from.Bytes())
	assert.True(t, amt.Equals(*balance.NewAmountFromInt(1)), "Got GetDelegatorBoundedAmount: %v", amt)

	options, _ := ctx.Govern.GetStakingOptions()
	assert.True(t, options.MinSelfDelegationAmount.Equals(*balance.NewAmountFromInt(1)), "Got MinSelfDelegationAmount: %v", amt)

	val := getBalanceFromAddress(ctx, from)
	assert.True(t, val.Equals(*balance.NewAmountFromInt(balToValidate)), "Got balance on address %s  - %s, required - %d", from.String(), val.String(), balToValidate)

	validatorSet, _ := ctx.Validators.GetValidatorSet()
	assert.True(t, len(validatorSet) == 0)
}

func TestWithdrawTx_ProcessDeliver_OK(t *testing.T) {
	ast := &withdrawTx{}
	ctx := &action.Context{}

	t.Run("withdraw with valid and existing amount, should return ok", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyApplyUnstakeData(1, 10000000000)

		// check init data
		initCheckWithdraw(t, ctx, initAmt)

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
		assert.True(t, options.MinSelfDelegationAmount.Equals(*balance.NewAmountFromInt(1)), "Got MinSelfDelegationAmount: %v", amt)

		val := getBalanceFromAddress(ctx, from)
		requiredBal := initAmt + 1
		assert.True(t, val.Equals(*balance.NewAmountFromInt(requiredBal)), "Got balance on address %s  - %s, required - %d", from.String(), val.String(), requiredBal)

		validatorSet, _ := ctx.Validators.GetValidatorSet()
		assert.True(t, len(validatorSet) == 0)
	})

}

func TestWithdrawTx_ProcessDeliver_Error(t *testing.T) {
	ast := &withdrawTx{}
	ctx := &action.Context{}

	t.Run("withdraw with an amount greater than on balance, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyApplyUnstakeData(100, 10)

		// check init data
		initCheckWithdraw(t, ctx, initAmt)

		// simulate validate
		ok, _ := ast.Validate(ctx, tx)
		assert.False(t, ok)
	})

	t.Run("withdraw with an zero amount, should return error", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)

		var initAmt int64 = 10

		ctx = assemblyCtxData("OLT", 0, true, true, initAmt)
		tx := assemblyApplyUnstakeData(0, 10)

		// check init data
		initCheckWithdraw(t, ctx, initAmt)

		// simulate validate
		ok, _ := ast.Validate(ctx, tx)
		assert.False(t, ok)
	})
}
