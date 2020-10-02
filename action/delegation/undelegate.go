package delegation

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type Undelegate struct {
	Delegator keys.Address `json:"delegator"`
	Amount action.Amount `json:"amount"`
}

var _ action.Msg = &Undelegate{}

func (ud Undelegate) Signers() []action.Address {
	return []action.Address{ud.Delegator.Bytes()}
}

func (ud Undelegate) Type() action.Type {
	return action.DELEGATION_UNDELEGATE
}

func (ud Undelegate) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(ud.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.delegator"),
		Value: ud.Delegator.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.amount"),
		Value: []byte(ud.Amount.String()),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (ud *Undelegate) Marshal() ([]byte, error) {
	return json.Marshal(ud)
}

func (ud *Undelegate) Unmarshal(data []byte) error {
	return json.Unmarshal(data, ud)
}

type UndelegateTx struct {}

var _ action.Tx = &UndelegateTx{}

func (u UndelegateTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	ctx.Logger.Debug("Validate UndelegateTx transaction for CheckTx", tx)
	ud := &Undelegate{}
	err := ud.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), ud.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	// validate params
	if err = ud.Delegator.Err(); err != nil {
		return false, ErrInvalidAddress
	}

	return true, nil
}

func (u UndelegateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck UndelegateTx transaction for CheckTx", tx)
	return runUndelegate(ctx, tx)
}

func (u UndelegateTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (u UndelegateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver UndelegateTx transaction for DeliverTx", tx)
	return runUndelegate(ctx, tx)
}

func runUndelegate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ud := &Undelegate{}
	err := ud.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, ud.Tags(), err)
	}

	// check if there is enough amount to undelegate
	ds := ctx.DelegationStore
	delegationAmount, err := ds.WithPrefixType(activePrefix).Get(ud.Delegator)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrGettingActiveDelegationAmount, ud.Tags(), err)
	}

	delegationAmountMock := action.Amount{}
	delegationValue := delegationAmountMock.Value
	undelegateValue := ud.Amount.Value

	// cut the amount from active store
	remainValue, err := delegationValue.Minus(undelegateValue)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrDeductingDelegationAmount, ud.Tags(), err)
	}

	err := ds.WithPrefixType(activePrefix).Set(ud.Delegator, remainValue)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrSettingActiveDelegationAmount, ud.Tags(), err)
	}

	// get mature height
	delegationOptions, err := ctx.GovernanceStore.GetDelegationOption()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, ErrGettingDelegationOption, ud.Tags(), err)
	}
	matureHeight := ctx.Header.GetHeight() + delegationOptions.MaturityPeriod

	// form the pending amount key and check if there is already an entry in pending store,
	// this means same delegator at least undelegated once in this block
	if ds.WithPrefixType(pendingPrefix).Exist(ud.Delegator, matureHeight) {
		// if so, change the amount
		existingPendingAmount, err := ds.WithPrefixType(pendingPrefix).Get(ud.Delegator, matureHeight)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, ErrGettingActiveDelegationAmount, ud.Tags(), err)
		}
		newAmount := existingPendingAmount + delegationAmount
		ds.WithPrefixType(pendingPrefix).Set(newAmount)

		return true, action.Response{Events: action.GetEvent(ud.Tags(), "un_delegate_success")}
	}

	// if not, add an entry to pending store
	ds.WithPrefixType(pendingPrefix).Set(delegationAmount)


	return true, action.Response{Events: action.GetEvent(ud.Tags(), "un_delegate_success")}
}