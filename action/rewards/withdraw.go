package rewards

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/rewards"
)

type Withdraw struct {
	ValidatorAddress action.Address `json:"validatorAddress"`
	SignerAddress    action.Address `json:"signerAddress"`
	WithdrawAmount   action.Amount  `json:"withdrawAmount"`
}

func (w Withdraw) Signers() []action.Address {
	return []action.Address{w.SignerAddress}
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

	currency, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		panic("no default currency available in the network")
	}
	if currency.Name != withdraw.WithdrawAmount.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, withdraw.WithdrawAmount.String())
	}
	err = withdraw.ValidatorAddress.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}
	return true, nil
}

func (withdrawTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runWithdraw(ctx, tx)
}

func (withdrawTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runWithdraw(ctx, tx)
}

func (withdrawTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
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

	withDrawCoin := withdraw.WithdrawAmount.ToCoinWithBase(ctx.Currencies)
	err = ctx.RewardMasterStore.RewardCm.WithdrawRewards(withdraw.ValidatorAddress, withDrawCoin.Amount)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, rewards.UnableToWithdraw, withdraw.Tags(), err)
	}
	if ctx.Validators.Exists(withdraw.ValidatorAddress) {
		validator, err := ctx.Validators.Get(withdraw.ValidatorAddress)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidValidatorAddr, withdraw.Tags(), err)
		}
		if !bytes.Equal(validator.StakeAddress, withdraw.SignerAddress) {
			return helpers.LogAndReturnFalse(ctx.Logger, action.ErrStakeAddressMismatch, withdraw.Tags(), err)
		}
	}
	if !ctx.Validators.Exists(withdraw.ValidatorAddress) {
		ctx.Logger.Info("Validator Does not exist , Allowing withdraw to address :", withdraw.SignerAddress)
	}

	//6. Update the balance db with the withdrawn amount for that validator

	rewardsPool := action.Address(ctx.RewardMasterStore.GetOptions().RewardPoolAddress)

	err = ctx.Balances.MinusFromAddress(rewardsPool, withDrawCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorMinusFailed, withdraw.Tags(), err)
	}
	err = ctx.Balances.AddToAddress(withdraw.SignerAddress.Bytes(), withDrawCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, withdraw.Tags(), err)
	}
	ctx.Logger.Debugf("Successfully withdrawn %s to Address %s ", withDrawCoin, withdraw.SignerAddress.String())
	return helpers.LogAndReturnTrue(ctx.Logger, withdraw.Tags(), "Success")
}
