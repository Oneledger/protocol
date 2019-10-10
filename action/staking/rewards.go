package staking

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

var _ action.Msg = &Withdraw{}

type Withdraw struct {
	// staking account from which the reward is withdraw
	From action.Address `json:"from"`
	// beneficiary account to which the reward is withdraw to
	To action.Address `json:"to"`
}

func (s Withdraw) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Withdraw) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s Withdraw) Signers() []action.Address {
	return []action.Address{s.From.Bytes()}
}

func (s Withdraw) Type() action.Type {
	return action.WITHDRAW
}

func (s Withdraw) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(s.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: s.From.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.to"),
		Value: s.To.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

var _ action.Tx = withdrawTx{}

type withdrawTx struct {
}

func (withdrawTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	Withdraw := &Withdraw{}
	err := Withdraw.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), Withdraw.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeeOpt, tx.Fee)
	if err != nil {
		return false, err
	}

	if tx.Fee.Price.Currency != ctx.FeeOpt.FeeCurrency.Name {
		return false, action.ErrInvalidFeeCurrency
	}
	if !ctx.FeeOpt.MinFee().LessThanEqualCoin(tx.Fee.Price.ToCoin(ctx.Currencies)) {
		return false, action.ErrInvalidFeePrice
	}

	return true, nil
}

func (s withdrawTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing Withdraw Transaction for CheckTx", tx)
	ok, result = runWithdraw(ctx, tx)
	return
}

func (s withdrawTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing Withdraw Transaction for DeliverTx", tx)
	ok, result = runWithdraw(ctx, tx)
	return
}

func (withdrawTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runWithdraw(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	balances := ctx.Balances
	feePool := ctx.FeePool

	draw := &Withdraw{}
	err := draw.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	allow := feePool.GetAllowedWithdraw(draw.From)
	if allow.LessThanCoin(ctx.FeeOpt.MinFee()) {
		return false, action.Response{Log: "No reward is allowed to withdraw"}
	}

	err = feePool.MinusFromAddress(draw.From, allow)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "minus from pool").Error()}
	}

	err = balances.AddToAddress(draw.To, allow)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "add to balance").Error()}
	}

	return true, action.Response{Tags: draw.Tags()}
}
