package network_delegation

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	netwkDeleg "github.com/Oneledger/protocol/data/network_delegation"
)

type Withdraw struct {
	Delegator action.Address `json:"delegator"`
	Amount    action.Amount  `json:"amount"`
}

func (w Withdraw) Signers() []action.Address {
	return []action.Address{w.Delegator}
}

func (w Withdraw) Type() action.Type {
	return action.REWARDS_WITHDRAW_NETWORK_DELEGATE
}

func (w Withdraw) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(w.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.Delegator"),
		Value: w.Delegator.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.Amount"),
		Value: []byte(w.Amount.String()),
	}
	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (w Withdraw) Marshal() ([]byte, error) {
	return json.Marshal(w)
}

func (w *Withdraw) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, w)
}

type delegWithdrawRewardsTx struct {
}

func (delegWithdrawRewardsTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	withdraw := Withdraw{}
	err := withdraw.Unmarshal(signedTx.Data)
	if err != nil {
		return false, err
	}

	err = action.ValidateBasic(signedTx.RawBytes(), withdraw.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	currency, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return false, errors.Wrap(action.ErrInvalidCurrency, withdraw.Amount.Currency)
	}
	if currency.Name != withdraw.Amount.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, withdraw.Amount.String())
	}

	err = withdraw.Delegator.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}
	return true, nil
}

func (delegWithdrawRewardsTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runDeleWithdraw(ctx, tx)
}

func (delegWithdrawRewardsTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runDeleWithdraw(ctx, tx)
}

func (delegWithdrawRewardsTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

var _ action.Msg = &Withdraw{}
var _ action.Tx = &delegWithdrawRewardsTx{}

func runDeleWithdraw(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	withdraw := Withdraw{}
	err := withdraw.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, withdraw.Tags(), err)
	}

	height := ctx.Header.GetHeight()
	options, err := ctx.GovernanceStore.GetNetworkDelegOptions()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, netwkDeleg.ErrGettingDelgOption, withdraw.Tags(), err)
	}

	// initiate a withdrawal which matures at block [height+RewardsMaturityTime]
	coinAmt := withdraw.Amount.ToCoin(ctx.Currencies)
	err = ctx.NetwkDelegators.Rewards.Withdraw(withdraw.Delegator, coinAmt.Amount, height+options.RewardsMaturityTime)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, netwkDeleg.ErrInitiateWithdrawal, withdraw.Tags(), err)
	}

	ctx.Logger.Debugf("Successfully initiated withdrawal, delegator= %s, amount= %s, mature_height= %d",
		withdraw.Delegator.String(), coinAmt, height+options.RewardsMaturityTime)
	return helpers.LogAndReturnTrue(ctx.Logger, withdraw.Tags(), "Success")
}
