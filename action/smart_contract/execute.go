package smart_contract

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type Execute struct {
	From   action.Address `json:"from"`
	To     action.Address `json:"to"`
	Amount action.Amount  `json:"amount"`
	Data   []byte         `json:"data"`
}

func (e Execute) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Execute) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e Execute) Signers() []action.Address {
	return []action.Address{e.From.Bytes()}
}

func (e Execute) Type() action.Type {
	return action.SC_EXECUTE
}

func (e Execute) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(e.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.origin"),
		Value: e.From.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.contract"),
		Value: e.To.Bytes(),
	}
	tags = append(tags, tag, tag2, tag3)
	return tags
}

var _ action.Tx = scExecuteTx{}

type scExecuteTx struct {
}

func (s scExecuteTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	execute := &Execute{}
	err := execute.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), execute.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	if !execute.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, execute.Amount.String())
	}

	if execute.From.Err() != nil || execute.To.Err() != nil {
		return false, action.ErrInvalidAddress
	}

	if len(execute.Data) == 0 {
		return false, action.ErrMissingData
	}
	return true, nil
}

func (s scExecuteTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for CheckTx", tx)
	ok, result = runSCExecute(ctx, tx)
	return
}

func (s scExecuteTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for DeliverTx", tx)
	ok, result = runSCExecute(ctx, tx)
	return
}

func (s scExecuteTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runSCExecute(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	execute := &Execute{}
	err := execute.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !execute.Amount.IsValid(ctx.Currencies) {
		log := fmt.Sprint("amount is invalid", execute.Amount, ctx.Currencies)
		return false, action.Response{Log: log}
	}

	// coin := execute.Amount.ToCoin(ctx.Currencies)

	// TODO: Add logic

	return true, action.Response{Events: action.GetEvent(execute.Tags(), "smart_contract_execute")}
}
