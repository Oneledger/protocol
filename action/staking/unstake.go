package staking

import (
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &Unstake{}

type Unstake struct {
	ValidatorAddress keys.Address
	StakeAddress     keys.Address
	Stake            action.Amount
}

func (ust Unstake) Marshal() ([]byte, error) {
	return json.Marshal(ust)
}

func (ust *Unstake) Unmarshal(data []byte) error {
	return json.Unmarshal(data, ust)
}

func (ust Unstake) Signers() []action.Address {
	return []action.Address{ust.StakeAddress.Bytes(), ust.ValidatorAddress.Bytes()}
}

func (ust Unstake) Type() action.Type {
	return action.UNSTAKE
}

func (ust Unstake) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(ust.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: ust.ValidatorAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.delegator"),
		Value: ust.StakeAddress.Bytes(),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.amount"),
		Value: ust.Stake.Value.BigInt().Bytes(),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

var _ action.Tx = unstakeTx{}

type unstakeTx struct{}

func (us unstakeTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	ust := &Unstake{}
	err := ust.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	err = action.ValidateBasic(tx.RawBytes(), ust.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if err := ust.StakeAddress.Err(); err != nil {
		return false, err
	}

	if ust.ValidatorAddress == nil {
		return false, action.ErrMissingData
	}

	coin := ust.Stake.ToCoinWithBase(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromInt(0)) {
		return false, action.ErrInvalidAmount
	}

	if coin.Currency.Name != "OLT" {
		return false, action.ErrInvalidCurrency
	}

	return true, nil
}

func (us unstakeTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing Apply unstake Transaction for CheckTx", tx)
	return runCheckUnstake(ctx, tx)
}

func (us unstakeTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("Processing Apply unstake Transaction for DeliverTx", tx)
	return runCheckUnstake(ctx, tx)
}

func (us unstakeTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	// TODO: Add fee processing
	return action.BasicFeeHandling(ctx, signedTx, start, size, 2)
}

func runCheckUnstake(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ust := &Unstake{}
	err := ust.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	unstake := identity.Unstake{
		Address: ust.ValidatorAddress,
		Amount:  ust.Stake.Value,
	}

	height := ctx.Header.GetHeight()

	options, err := ctx.Govern.GetStakingOptions()
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, ust.StakeAddress.String()).Error()}
	}

	err = ctx.Delegators.Unstake(ust.ValidatorAddress, ust.StakeAddress, ust.Stake.Value, height+options.MaturityTime)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, ust.StakeAddress.String()).Error()}
	}

	err = ctx.Validators.HandleUnstake(unstake)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	return true, action.Response{Events: action.GetEvent(ust.Tags(), "unstake")}
}
