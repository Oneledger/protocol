package rewards

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
)

type Withdraw struct {
	ValidatorAddress action.Address `json:"validatorAddress"`
}

func (w Withdraw) Signers() []action.Address {
	return []action.Address{w.ValidatorAddress}
}

func (w Withdraw) Type() action.Type {
	return action.WITHDRAW_REWARD
}

func (w Withdraw) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(w.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.Validator"),
		Value: w.ValidatorAddress.Bytes(),
	}
	tags = append(tags, tag, tag2)
	return tags
}

func (w Withdraw) Marshal() ([]byte, error) {
	return json.Marshal(w)
}

func (w *Withdraw) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, w)
}

type withdrawTx struct {
}

func (withdrawTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	withdraw := Withdraw{}
	err := withdraw.Unmarshal(signedTx.Data)
	if err != nil {
		return false, err
	}
	fmt.Println("Validator address : ", withdraw.ValidatorAddress)
	return true, nil
}

func (withdrawTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runWithdraw(ctx, tx)
}

func (withdrawTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runWithdraw(ctx, tx)
}

func (withdrawTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

var _ action.Msg = &Withdraw{}
var _ action.Tx = &withdrawTx{}

func runWithdraw(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	withdraw := Withdraw{}
	err := withdraw.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, withdraw.Tags(), err)
	}
	return helpers.LogAndReturnTrue(ctx.Logger, withdraw.Tags(), "Success")
}
