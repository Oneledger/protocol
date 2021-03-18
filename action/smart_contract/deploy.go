package smart_contract

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type Deploy struct {
	From   action.Address `json:"from"`
	Amount action.Amount  `json:"amount"`
	Data   []byte         `json:"data"`
}

func (d Deploy) Marshal() ([]byte, error) {
	return json.Marshal(d)
}

func (d *Deploy) Unmarshal(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d Deploy) Signers() []action.Address {
	return []action.Address{d.From.Bytes()}
}

func (d Deploy) Type() action.Type {
	return action.SC_DEPLOY
}

func (d Deploy) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(d.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: d.From.Bytes(),
	}
	tags = append(tags, tag, tag2)
	return tags
}

var _ action.Tx = scDeployTx{}

type scDeployTx struct {
}

func (s scDeployTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	deploy := &Deploy{}
	err := deploy.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), deploy.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	//validate transaction specific field
	if !deploy.Amount.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, deploy.Amount.String())
	}

	if deploy.From.Err() != nil {
		return false, action.ErrInvalidAddress
	}

	if len(deploy.Data) == 0 {
		return false, action.ErrMissingData
	}
	return true, nil
}

func (s scDeployTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for CheckTx", tx)
	ok, result = runSCDeploy(ctx, tx)
	return
}

func (s scDeployTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing SC Deploy Transaction for DeliverTx", tx)
	ok, result = runSCDeploy(ctx, tx)
	return
}

func (s scDeployTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runSCDeploy(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	deploy := &Deploy{}
	err := deploy.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !deploy.Amount.IsValid(ctx.Currencies) {
		log := fmt.Sprint("amount is invalid", deploy.Amount, ctx.Currencies)
		return false, action.Response{Log: log}
	}

	// coin := deploy.Amount.ToCoin(ctx.Currencies)

	// TODO: Add logic

	return true, action.Response{Events: action.GetEvent(deploy.Tags(), "smart_contract_deploy")}
}
