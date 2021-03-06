package network_delegation

// below is removed since finalize withdraw rewards logic is moved to block beginner, OLP-1266
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/Oneledger/protocol/action"
//	"github.com/Oneledger/protocol/action/helpers"
//	"github.com/Oneledger/protocol/data/keys"
//	net_delg "github.com/Oneledger/protocol/data/network_delegation"
//	"github.com/pkg/errors"
//	"github.com/tendermint/tendermint/libs/kv"
//)
//
//type DeleWithdrawRewards struct {
//	Delegator keys.Address  `json:"delegator"`
//	Amount    action.Amount `json:"amount"`
//}
//
//var _ action.Msg = &DeleWithdrawRewards{}
//
//func (wr DeleWithdrawRewards) Signers() []action.Address {
//	return []action.Address{wr.Delegator.Bytes()}
//}
//
//func (wr DeleWithdrawRewards) Type() action.Type {
//	return action.REWARDS_FINALIZE_NETWORK_DELEGATE
//}
//
//func (wr DeleWithdrawRewards) Tags() kv.Pairs {
//	tags := make([]kv.Pair, 0)
//
//	tag := kv.Pair{
//		Key:   []byte("tx.type"),
//		Value: []byte(wr.Type().String()),
//	}
//	tag2 := kv.Pair{
//		Key:   []byte("tx.delegator"),
//		Value: wr.Delegator.Bytes(),
//	}
//	tag3 := kv.Pair{
//		Key:   []byte("tx.amount"),
//		Value: []byte(wr.Amount.String()),
//	}
//
//	tags = append(tags, tag, tag2, tag3)
//	return tags
//}
//
//func (wr *DeleWithdrawRewards) Marshal() ([]byte, error) {
//	return json.Marshal(wr)
//}
//
//func (wr *DeleWithdrawRewards) Unmarshal(data []byte) error {
//	return json.Unmarshal(data, wr)
//}
//
//type finalizeWithdrawRewardsTx struct{}
//
//var _ action.Tx = &finalizeWithdrawRewardsTx{}
//
//func (wt finalizeWithdrawRewardsTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
//	ctx.Logger.Debug("Validate DeleFinalizeRewardsTx transaction for CheckTx", tx)
//	w := &DeleWithdrawRewards{}
//	err := w.Unmarshal(tx.Data)
//	if err != nil {
//		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
//	}
//
//	// validate basic signature
//	err = action.ValidateBasic(tx.RawBytes(), w.Signers(), tx.Signatures)
//	if err != nil {
//		return false, err
//	}
//
//	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
//	if err != nil {
//		return false, err
//	}
//
//	// validate params
//	if err = w.Delegator.Err(); err != nil {
//		return false, action.ErrInvalidAddress
//	}
//
//	return true, nil
//}
//
//func (wt finalizeWithdrawRewardsTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	ctx.Logger.Debug("ProcessCheck DeleFinalizeRewardsTx transaction for CheckTx", tx)
//	return runFinalizeWithdraw(ctx, tx)
//}
//
//func (wt finalizeWithdrawRewardsTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
//	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
//}
//
//func (wt finalizeWithdrawRewardsTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	ctx.Logger.Debug("ProcessDeliver DeleFinalizeRewardsTx transaction for DeliverTx", tx)
//	return runFinalizeWithdraw(ctx, tx)
//}
//
//func runFinalizeWithdraw(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
//	w := &DeleWithdrawRewards{}
//	err := w.Unmarshal(tx.Data)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, w.Tags(), err)
//	}
//
//	// cut the amount from matured rewards if there is enough to withdraw
//	ds := ctx.NetwkDelegators.Rewards
//	rewards, err := ds.GetMaturedRewards(w.Delegator)
//	fmt.Println("rewards", rewards)
//	fmt.Println("err", err)
//	err = ds.Finalize(w.Delegator, &w.Amount.Value)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrFinalizingDelgRewards, w.Tags(), err)
//	}
//
//	rewards, err = ds.GetMaturedRewards(w.Delegator)
//	fmt.Println("rewards after finalize", rewards)
//	fmt.Println("err", err)
//
//	// add the amount to delegator address
//	withdrawCoin := w.Amount.ToCoin(ctx.Currencies)
//	err = ctx.Balances.AddToAddress(w.Delegator.Bytes(), withdrawCoin)
//	if err != nil {
//		return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrAddingWithdrawAmountToBalance, w.Tags(), err)
//	}
//	return true, action.Response{Events: action.GetEvent(w.Tags(), "delegation_rewards_withdraw_success")}
//}
