package staking

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
)

var _ action.Msg = &Stake{}

type Stake struct {
	ValidatorAddress     keys.Address
	StakeAddress         keys.Address
	ValidatorPubKey      keys.PublicKey
	ValidatorECDSAPubKey keys.PublicKey
	NodeName             string
	Stake                action.Amount
}

func (st Stake) Marshal() ([]byte, error) {
	return json.Marshal(st)
}

func (st *Stake) Unmarshal(data []byte) error {
	return json.Unmarshal(data, st)
}

func (st Stake) Signers() []action.Address {
	return []action.Address{st.StakeAddress.Bytes(), st.ValidatorAddress.Bytes()}
}

func (st Stake) Type() action.Type {
	return action.STAKE
}

func (st Stake) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(st.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: st.ValidatorAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.delegator"),
		Value: st.StakeAddress.Bytes(),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.amount"),
		Value: st.Stake.Value.BigInt().Bytes(),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

var _ action.Tx = stakeTx{}

type stakeTx struct{}

func (s stakeTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	st := &Stake{}
	err := st.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	err = action.ValidateBasic(tx.RawBytes(), st.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if err := st.StakeAddress.Err(); err != nil {
		return false, err
	}

	if err := st.ValidatorAddress.Err(); err != nil {
		return false, err
	}

	_, err = st.ValidatorPubKey.GetHandler()
	if err != nil {
		return false, action.ErrInvalidPubkey
	}

	coin := st.Stake.ToCoinWithBase(ctx.Currencies)
	if !coin.IsValid() {
		return false, errors.Wrap(action.ErrInvalidAmount, coin.String())
	}

	if coin.Currency.Name != "OLT" {
		return false, action.ErrInvalidCurrency
	}

	_, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return false, action.ErrInvalidCurrency
	}

	err = ctx.Balances.CheckBalanceFromAddress(st.StakeAddress, coin)
	if err != nil {
		return false, action.ErrNotEnoughFund
	}

	return true, nil
}

func (s stakeTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing Apply stake Transaction for ProcessCheck", tx)
	ok, result = runCheckStake(ctx, tx)
	ctx.Logger.Detail("Result Apply stake Transaction for ProcessCheck", ok, result)
	return
}

func (s stakeTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing Apply stake Transaction for ProcessDeliver", tx)
	ok, result = runCheckStake(ctx, tx)
	ctx.Logger.Detail("Result Apply stake Transaction for ProcessDeliver", ok, result)
	return
}

func (s stakeTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.Logger.Detail("Processing Apply stake Transaction for ProcessFee", signedTx)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 2)
}

// allow validator's stake address update only when it's clean
func isStakeAddressClean(ctx *action.Context, v *identity.Validator) (bool, error) {
	options, err := ctx.GovernanceStore.GetStakingOptions()
	if err != nil {
		return false, err
	}
	zero := balance.NewAmountFromInt(0)

	// check locked amount
	lockedAmt, err := ctx.Delegators.GetValidatorDelegationAmount(v.Address, v.StakeAddress)
	if err != nil {
		return false, err
	}
	if !lockedAmt.Equals(*zero) {
		return false, nil
	}

	// check pending amount
	pendingAmounts := ctx.Delegators.GetMaturedPendingAmount(v.StakeAddress, ctx.Header.Height, options.MaturityTime+1)
	if len(pendingAmounts) != 0 {
		return false, nil
	}

	// check bounded amount
	boundCoin, err := ctx.Delegators.GetDelegatorBoundedAmount(v.StakeAddress)
	if err != nil {
		return false, err
	}
	if !boundCoin.Equals(*zero) {
		return false, nil
	}

	return true, nil
}

func runCheckStake(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	st := &Stake{}
	err := st.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if ctx.EvidenceStore.IsFrozenValidator(st.ValidatorAddress) {
		return false, action.Response{Log: evidence.ErrFrozenValidator.Error()}
	}

	stake := identity.Stake{
		ValidatorAddress: st.ValidatorAddress,
		StakeAddress:     st.StakeAddress,
		Pubkey:           st.ValidatorPubKey,
		ECDSAPubKey:      st.ValidatorECDSAPubKey,
		Name:             st.NodeName,
		Amount:           st.Stake.Value,
	}

	// trying to update existing validator's stake address
	updateStakeAddress := false
	if ctx.Validators.Exists(st.ValidatorAddress) {
		validator, err := ctx.Validators.Get(st.ValidatorAddress)
		if err != nil {
			return false, action.Response{Log: errors.Wrap(err, st.StakeAddress.String()).Error()}
		}

		clean, err := isStakeAddressClean(ctx, validator)
		if err != nil {
			return false, action.Response{Log: errors.Wrap(err, st.StakeAddress.String()).Error()}
		}

		// update not allowed if existing stake address not cleaned up
		if !validator.StakeAddress.Equal(st.StakeAddress) {
			if !clean {
				return false, action.Response{Log: action.ErrStakeAddressInUse.Marshal()}
			}
			updateStakeAddress = true
		}
	}

	coin := st.Stake.ToCoinWithBase(ctx.Currencies)

	err = ctx.Balances.MinusFromAddress(st.StakeAddress, coin)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, st.StakeAddress.String()).Error()}
	}

	err = ctx.Delegators.Stake(st.ValidatorAddress, st.StakeAddress, st.Stake.Value)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, st.StakeAddress.String()).Error()}
	}

	err = ctx.Validators.HandleStake(stake, updateStakeAddress, ctx.Header.Height)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	return true, action.Response{Events: action.GetEvent(st.Tags(), "apply_stake")}
}
