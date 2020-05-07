package staking

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
)

var _ action.Msg = &Withdraw{}

type Withdraw struct {
	ValidatorAddress keys.Address
	StakeAddress     keys.Address
	Stake            action.Amount
}

func (s Withdraw) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Withdraw) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s Withdraw) Signers() []action.Address {
	return []action.Address{s.StakeAddress.Bytes(), s.ValidatorAddress.Bytes()}
}

func (s Withdraw) Type() action.Type {
	return action.WITHDRAW
}

func (s Withdraw) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(s.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: s.ValidatorAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.delegator"),
		Value: s.StakeAddress.Bytes(),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.amount"),
		Value: s.Stake.Value.BigInt().Bytes(),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

var _ action.Tx = withdrawTx{}

type withdrawTx struct {
}

func (withdrawTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	apply := &Withdraw{}
	err := apply.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), apply.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if err := apply.StakeAddress.Err(); err != nil {
		return false, err
	}

	if apply.ValidatorAddress == nil {
		return false, action.ErrMissingData
	}

	coin := apply.Stake.ToCoinWithBase(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromInt(0)) {
		return false, action.ErrInvalidAmount
	}

	if coin.Currency.Name != "OLT" {
		return false, action.ErrInvalidAmount
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
	return action.BasicFeeHandling(ctx, signedTx, start, size, 2)
}

func runWithdraw(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	draw := &Withdraw{}
	err := draw.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	coin := draw.Stake.ToCoinWithBase(ctx.Currencies)

	err = ctx.Delegators.Withdraw(draw.ValidatorAddress, draw.StakeAddress, draw.Stake.Value)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, draw.StakeAddress.String()).Error()}
	}

	err = ctx.Balances.AddToAddress(draw.StakeAddress, coin)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "add to balance").Error()}
	}

	return true, action.Response{Events: action.GetEvent(draw.Tags(), "apply_withdraw"), Info: coin.String()}
}
