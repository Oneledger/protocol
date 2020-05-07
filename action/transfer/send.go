package transfer

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
)

var _ action.Msg = &Send{}

type Send struct {
	From   action.Address `json:"from"`
	To     action.Address `json:"to"`
	Amount action.Amount  `json:"amount"`
}

func (s Send) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Send) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s Send) Signers() []action.Address {
	return []action.Address{s.From.Bytes()}
}

func (s Send) Type() action.Type {
	return action.SEND
}

func (s Send) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(s.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: s.From.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.to"),
		Value: s.To.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

var _ action.Tx = sendTx{}

type sendTx struct {
}

func (sendTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	send := &Send{}
	err := send.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), send.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	if !send.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, send.Amount.String())
	}

	if send.From.Err() != nil || send.To.Err() != nil {
		return false, action.ErrInvalidAddress
	}
	return true, nil
}

func (s sendTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing Send Transaction for CheckTx", tx)
	ok, result = runTx(ctx, tx)
	return
}

func (s sendTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing Send Transaction for DeliverTx", tx)
	ok, result = runTx(ctx, tx)
	return
}

func (sendTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runTx(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	balances := ctx.Balances

	send := &Send{}
	err := send.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !send.Amount.IsValid(ctx.Currencies) {
		log := fmt.Sprint("amount is invalid", send.Amount, ctx.Currencies)
		return false, action.Response{Log: log}
	}

	coin := send.Amount.ToCoin(ctx.Currencies)

	err = balances.MinusFromAddress(send.From.Bytes(), coin)
	if err != nil {
		log := fmt.Sprint("error debiting balance in send transaction ", send.From, "err", err)
		return false, action.Response{Log: log}
	}

	err = balances.AddToAddress(send.To.Bytes(), coin)
	if err != nil {
		log := fmt.Sprint("error crediting balance in send transaction ", send.From, "err", err)
		return false, action.Response{Log: log}
	}

	return true, action.Response{Events: action.GetEvent(send.Tags(), "send_tx")}
}
