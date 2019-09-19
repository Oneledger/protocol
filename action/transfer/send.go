package transfer

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/action"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
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

func (s Send) Tags() common.KVPairs {
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
	ok, err := action.ValidateBasic(tx.RawBytes(), send.Signers(), tx.Signatures)
	if err != nil {
		return ok, err
	}

	//validate transaction specific field

	if !send.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, send.Amount.String())
	}
	if send.From == nil || send.To == nil {
		return false, action.ErrMissingData
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
	ctx.State.ConsumeVerifySigGas(1)
	ctx.State.ConsumeStorageGas(size)
	// check the used gas for the tx
	final := ctx.Balances.State.ConsumedGas()
	used := int64(final - start)
	if used > signedTx.Fee.Gas {
		return false, action.Response{Log: action.ErrGasOverflow.Error(), GasWanted: signedTx.Fee.Gas, GasUsed: signedTx.Fee.Gas}
	}
	// only charge the first signer
	signer := signedTx.Signatures[0].Signer
	h, err := signer.GetHandler()
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	addr := h.Address()

	charge := signedTx.Fee.Price.ToCoin(ctx.Currencies).MultiplyInt64(int64(used))
	bal, _ := ctx.Balances.Get(addr)
	if _, err := bal.MinusCoin(charge); err != nil {
		return false, action.Response{Log: err.Error()}
	}
	err = ctx.Balances.Set(h.Address(), *bal)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	fmt.Println("fee charged", charge)
	err = ctx.FeePool.AddToPool(charge)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	return true, action.Response{GasWanted: signedTx.Fee.Gas, GasUsed: used}
}

func runTx(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	balances := ctx.Balances

	send := &Send{}
	err := send.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	from, err := balances.Get(send.From.Bytes())
	if err != nil {
		log := fmt.Sprint("Failed to get the balance of the owner ", send.From, "err", err)
		return false, action.Response{Log: log}
	}

	if !send.Amount.IsValid(ctx.Currencies) {
		log := fmt.Sprint("amount is invalid", send.Amount, ctx.Currencies)
		return false, action.Response{Log: log}
	}

	coin := send.Amount.ToCoin(ctx.Currencies)

	//change owner balance
	fromFinal, err := from.MinusCoin(coin)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = balances.Set(send.From.Bytes(), *fromFinal)
	if err != nil {
		log := fmt.Sprint("error updating balance in send transaction", err)
		return false, action.Response{Log: log}
	}

	//change receiver balance
	to, err := balances.Get(send.To.Bytes())
	if err != nil {
		ctx.Logger.Error("failed to get the balance of the receipient", err)
	}
	if to == nil {
		to = balance.NewBalance()
	}

	to.AddCoin(coin)
	err = balances.Set(send.To.Bytes(), *to)
	if err != nil {
		_ = balances.Set(send.From.Bytes(), *from)
		return false, action.Response{Log: "balance set failed"}
	}
	return true, action.Response{Tags: send.Tags()}
}
