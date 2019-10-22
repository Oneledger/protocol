package staking

import (
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/tendermint/tendermint/libs/common"
)

var _ action.Msg = &ApplyValidator{}

type ApplyValidator struct {
	StakeAddress         action.Address
	Stake                action.Amount
	NodeName             string
	ValidatorAddress     action.Address
	ValidatorPubKey      keys.PublicKey
	ValidatorECDSAPubKey keys.PublicKey
	Purge                bool
}

func (apply ApplyValidator) Marshal() ([]byte, error) {
	return json.Marshal(apply)
}

func (apply *ApplyValidator) Unmarshal(data []byte) error {
	return json.Unmarshal(data, apply)
}

func (apply ApplyValidator) Signers() []action.Address {
	return []action.Address{apply.StakeAddress.Bytes()}
}

func (apply ApplyValidator) Type() action.Type {
	return action.APPLYVALIDATOR
}

var _ action.Tx = applyTx{}

type applyTx struct {
}

func (a applyTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {

	apply := &ApplyValidator{}
	err := apply.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	err = action.ValidateBasic(tx.RawBytes(), apply.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeeOpt, tx.Fee)
	if err != nil {
		return false, err
	}

	if len(apply.StakeAddress) == 0 {
		return false, action.ErrMissingData
	}

	if apply.ValidatorAddress == nil {
		return false, action.ErrMissingData
	}
	_, err = apply.ValidatorPubKey.GetHandler()
	if err != nil {
		return false, action.ErrInvalidPubkey
	}

	coin := apply.Stake.ToCoin(ctx.Currencies)
	if coin.LessThanCoin(coin.Currency.NewCoinFromInt(0)) {
		return false, action.ErrInvalidAmount
	}

	if coin.Currency.Name != "VT" {
		return false, action.ErrInvalidAmount
	}

	return true, nil
}

func (a applyTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing Apply Validator Transaction for CheckTx", tx)
	return runApply(ctx, tx)
}

func checkBalances(ctx *action.Context, address action.Address, stake action.Amount) (bool, error) {

	balances := ctx.Balances

	_, ok := ctx.Currencies.GetCurrencyByName("VT")
	if !ok {
		return false, action.ErrInvalidAmount
	}

	err := balances.CheckBalanceFromAddress(address.Bytes(), stake.ToCoin(ctx.Currencies))
	if err != nil {
		return false, action.ErrNotEnoughFund
	}
	return true, nil
}

func (a applyTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing Apply Validator Transaction for DeliverTx", tx)
	return runApply(ctx, tx)
}

func (a applyTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (apply ApplyValidator) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(apply.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: apply.StakeAddress.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func runApply(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	apply := &ApplyValidator{}
	err := apply.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	_, err = checkBalances(ctx, apply.StakeAddress, apply.Stake)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	validators := ctx.Validators

	balances := ctx.Balances

	err = balances.MinusFromAddress(apply.StakeAddress.Bytes(), apply.Stake.ToCoin(ctx.Currencies))
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !apply.Purge {
		stake := identity.Stake{
			ValidatorAddress: apply.ValidatorAddress,
			StakeAddress:     apply.StakeAddress,
			Pubkey:           apply.ValidatorPubKey,
			ECDSAPubKey:      apply.ValidatorECDSAPubKey,
			Name:             apply.NodeName,
			Amount:           apply.Stake.Value,
		}
		err = validators.HandleStake(stake)
	} else {
		unstake := identity.Unstake{
			Address: apply.ValidatorAddress,
			Amount:  apply.Stake.Value,
		}
		err = validators.HandleUnstake(unstake)
	}
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	return true, action.Response{Tags: apply.Tags()}
}
