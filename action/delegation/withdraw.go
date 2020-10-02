package delegation

import (
"encoding/json"
"github.com/Oneledger/protocol/action"
"github.com/Oneledger/protocol/action/helpers"
"github.com/Oneledger/protocol/data/keys"
"github.com/pkg/errors"
"github.com/tendermint/tendermint/libs/kv"
)

type DeleWithdrawRewards struct {
	Delegator keys.Address `json:"delegator"`
	Amount action.Amount `json:"amount"`
}

var _ action.Msg = &DeleWithdrawRewards{}

func (wr DeleWithdrawRewards) Signers() []action.Address {
	return []action.Address{wr.Delegator.Bytes()}
}

func (wr DeleWithdrawRewards) Type() action.Type {
	return action.DELEGATION_REWARDS_WITHDRAW
}

func (wr DeleWithdrawRewards) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(wr.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.delegator"),
		Value: wr.Delegator.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.amount"),
		Value: []byte(wr.Amount.String()),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (wr *DeleWithdrawRewards) Marshal() ([]byte, error) {
	return json.Marshal(wr)
}

func (wr *DeleWithdrawRewards) Unmarshal(data []byte) error {
	return json.Unmarshal(data, wr)
}

type DeleWithdrawRewardsTx struct {}

var _ action.Tx = &DeleWithdrawRewardsTx{}

func (wt DeleWithdrawRewardsTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	ctx.Logger.Debug("Validate DeleWithdrawRewardsTx transaction for CheckTx", tx)
	w := &DeleWithdrawRewards{}
	err := w.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), w.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	// validate params
	if err = w.Delegator.Err(); err != nil {
		return false, ErrInvalidAddress
	}

	return true, nil
}

func (wt DeleWithdrawRewardsTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck DeleWithdrawRewardsTx transaction for CheckTx", tx)
	return runWithdraw(ctx, tx)
}

func (wt DeleWithdrawRewardsTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (wt DeleWithdrawRewardsTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver DeleWithdrawRewardsTx transaction for DeliverTx", tx)
	return runWithdraw(ctx, tx)
}

func runWithdraw(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	w := &DeleWithdrawRewards{}
	err := w.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, w.Tags(), err)
	}

	// check if there is enough reward to withdraw
	ds := ctx.DelegationRewardStore
	maturedAmount, err := ds.WithPrefixType(maturedPrefix).Get(w.Delegator)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrGettingMaturedDelegationRewards, w.Tags(), err)
	}

	maturedAmountMock := action.Amount{}
	maturedValue := maturedAmountMock.Value
	withdrawValue := w.Amount.Value

	// cut the amount from matured reward store
	remainValue, err := maturedValue.Minus(withdrawValue)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrDeductingMaturedAmount, w.Tags(), err)
	}
	err = ds.WithPrefixType(maturedPrefix).Set(w.Delegator, remainValue)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSettingMaturedDelegationAmount, w.Tags(), err)
	}

	// add the amount to delegator address
	withdrawCoin := w.Amount.ToCoin(ctx.Currencies)
	err = ctx.Balances.AddToAddress(w.Delegator.Bytes(), withdrawCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrAddingWithdrawAmountToBalance, w.Tags(), err)
	}
	return true, action.Response{Events: action.GetEvent(w.Tags(), "delegation_rewards_withdraw_success")}
}
