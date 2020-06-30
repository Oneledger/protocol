package rewards

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
)

type Withdraw struct {
	ValidatorAddress        action.Address `json:"validatorAddress"`
	ValidatorSigningAddress action.Address `json:"validatorSigningAddress"`
}

func (w Withdraw) Signers() []action.Address {
	return []action.Address{w.ValidatorSigningAddress}
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
	err = action.ValidateBasic(signedTx.RawBytes(), withdraw.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}
	if !ctx.Validators.IsValidatorAddress(withdraw.ValidatorAddress) {
		return false, action.ErrInvalidValidatorAddr
	}
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
	fmt.Println("Validator :", withdraw.ValidatorAddress)

	//1. Check the cumulative rewards DB
	_ = ctx.RewardMasterStore.RewardCm
	//2. Get the difference of amount earned vs amount withdrawn for the validator issuing this transaction

	//3. Check how much he is eligible to withdraw (which is calculated in step2)
	//4. If the amount withdrawn is less than or equal to amount eligible to be withdrawn, make the transaction success.
	//5. In case of no failure, add this amount the person withdrew, to total withdrawn amount in cumulative rewards db
	//6. Update the balance db with the withdrawn amount for that validator
	return helpers.LogAndReturnTrue(ctx.Logger, withdraw.Tags(), "Success")
}
