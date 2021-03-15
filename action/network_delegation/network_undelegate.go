package network_delegation

import (
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/balance"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	net_delg "github.com/Oneledger/protocol/data/network_delegation"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

type Undelegate struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
}

var _ action.Msg = &Undelegate{}

func (ud Undelegate) Signers() []action.Address {
	return []action.Address{ud.Delegator.Bytes()}
}

func (ud Undelegate) Type() action.Type {
	return action.NETWORK_UNDELEGATE
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

type undelegateTx struct{}

var _ action.Tx = &undelegateTx{}

func (u undelegateTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	ctx.Logger.Debug("Validate undelegateTx transaction for CheckTx", tx)
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
		return false, action.ErrInvalidAddress
	}

	return true, nil
}

func (u undelegateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck undelegateTx transaction for CheckTx", tx)
	return runUndelegate(ctx, tx)
}

func (u undelegateTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func (u undelegateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver undelegateTx transaction for DeliverTx", tx)
	return runUndelegate(ctx, tx)
}

func runUndelegate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ud := &Undelegate{}
	err := ud.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, ud.Tags(), err)
	}

	// get coin for active delegation amount and the amount to undelegate
	ds := ctx.NetwkDelegators.Deleg
	ds.WithPrefix(net_delg.ActiveType)
	delegationCoin, err := ds.Get(ud.Delegator)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrGettingActiveDelgAmount, ud.Tags(), err)
	}

	undelegateCoin := ud.Amount.ToCoin(ctx.Currencies)
	// cut the amount from active store
	remainCoin, err := delegationCoin.Minus(undelegateCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrDeductingActiveDelgAmount, ud.Tags(), err)
	}
	err = ds.Set(ud.Delegator, &remainCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrSettingActiveDelgAmount, ud.Tags(), err)
	}

	// get mature height
	delegationOptions, err := ctx.GovernanceStore.GetNetworkDelegOptions()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrGettingDelgOption, ud.Tags(), err)
	}
	matureHeight := ctx.Header.GetHeight() + delegationOptions.RewardsMaturityTime

	// check if there is already an entry in pending store with same address and height,
	// this means same delegator at least undelegated once in this block
	ds.WithPrefix(net_delg.PendingType)
	if !ds.PendingExists(ud.Delegator, matureHeight) {
		// if not, add an entry to pending store
		err := ds.SetPendingAmount(ud.Delegator, matureHeight, &undelegateCoin)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrSettingPendingDelgAmount, ud.Tags(), err)
		}
	} else {
		// if so, change the amount
		existingPendingCoin, err := ds.GetPendingAmount(ud.Delegator, matureHeight)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrGettingPendingDelgAmount, ud.Tags(), err)
		}
		newPendingCoin := existingPendingCoin.Plus(undelegateCoin)
		err = ds.SetPendingAmount(ud.Delegator, matureHeight, &newPendingCoin)
		if err != nil {
			return helpers.LogAndReturnFalse(ctx.Logger, net_delg.ErrSettingPendingDelgAmount, ud.Tags(), err)
		}
	}

	//Get Delegation Pool
	delagationPool, err := ctx.GovernanceStore.GetPoolByName(gov.POOL_DELEGATION)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrPoolDoesNotExist, ud.Tags(), err)
	}

	//cut balance from pool
	err = ctx.Balances.MinusFromAddress(delagationPool, undelegateCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, ud.Tags(), err)
	}

	return true, action.Response{Events: action.GetEvent(ud.Tags(), "undelegate_success")}
}
