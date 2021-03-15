package network_delegation

import (
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/balance"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/network_delegation"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &AddNetworkDelegation{}

type AddNetworkDelegation struct {
	DelegationAddress keys.Address  `json:"delegationAddress"`
	Amount            action.Amount `json:"amount"`
}

func (n AddNetworkDelegation) Signers() []action.Address {
	return []action.Address{n.DelegationAddress}
}

func (n AddNetworkDelegation) Type() action.Type {
	return action.ADD_NETWORK_DELEGATE
}

func (n AddNetworkDelegation) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(n.Type().String()),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.delegationAddress"),
		Value: n.DelegationAddress.Bytes(),
	}

	tags = append(tags, tag, tag3)
	return tags
}

func (n AddNetworkDelegation) Marshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *AddNetworkDelegation) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, n)
}

var _ action.Tx = addNetworkDelegationTx{}

type addNetworkDelegationTx struct{}

func (n addNetworkDelegationTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	delegate := &AddNetworkDelegation{}
	err := delegate.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	err = action.ValidateBasic(tx.RawBytes(), delegate.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if err := delegate.DelegationAddress.Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (n addNetworkDelegationTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runNetworkDelegate(ctx, tx)
}

func (n addNetworkDelegationTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runNetworkDelegate(ctx, tx)
}

func (n addNetworkDelegationTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.Logger.Detail("Processing Delegate Transaction for ProcessFee", signedTx)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runNetworkDelegate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	delegate := AddNetworkDelegation{}
	err := delegate.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, delegate.Tags(), err)
	}

	// Check if delegation address has funds
	coin := delegate.Amount.ToCoin(ctx.Currencies)
	if !coin.IsValid() {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidAmount, delegate.Tags(), errors.New("Coin is not valid"))
	}
	if coin.Currency.Name != "OLT" {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidCurrency, delegate.Tags(), errors.New("currency is not OLT"))
	}
	_, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidCurrency, delegate.Tags(), errors.New("currency OLT does not exist in system"))
	}
	err = ctx.Balances.CheckBalanceFromAddress(delegate.DelegationAddress, coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrNotEnoughFund, delegate.Tags(), err)
	}

	//Deduct Delegation Amount
	err = ctx.Balances.MinusFromAddress(delegate.DelegationAddress, coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorMinusFailed, delegate.Tags(), err)
	}

	//Add Delegation
	//Get Delegation Pool
	delagationPool, err := ctx.GovernanceStore.GetPoolByName(gov.POOL_DELEGATION)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrPoolDoesNotExist, delegate.Tags(), err)
	}

	//Add balance to pool
	err = ctx.Balances.AddToAddress(delagationPool, coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, delegate.Tags(), err)
	}

	//Add balance to delegation
	currentDelegation, _ := ctx.NetwkDelegators.Deleg.WithPrefix(network_delegation.ActiveType).Get(delegate.DelegationAddress)
	newCoin := currentDelegation.Plus(coin)
	err = ctx.NetwkDelegators.Deleg.WithPrefix(network_delegation.ActiveType).Set(delegate.DelegationAddress, &newCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, delegate.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, delegate.Tags(), "Success")
}
