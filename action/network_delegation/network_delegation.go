package network_delegation

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/balance"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &NetworkDelegate{}

type NetworkDelegate struct {
	DelegationAddress keys.Address
	Amount            action.Amount
}

func (n NetworkDelegate) Signers() []action.Address {
	return []action.Address{n.DelegationAddress}
}

func (n NetworkDelegate) Type() action.Type {
	return action.NETWORKDELEGATE
}

func (n NetworkDelegate) Tags() kv.Pairs {
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

func (n NetworkDelegate) Marshal() ([]byte, error) {
	return json.Marshal(n)
}

func (n *NetworkDelegate) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, n)
}

var _ action.Tx = networkDelegateTx{}

type networkDelegateTx struct{}

func (n networkDelegateTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	delegate := &NetworkDelegate{}
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

func (n networkDelegateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runNetworkDelegate(ctx, tx)
}

func (n networkDelegateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runNetworkDelegate(ctx, tx)
}

func (n networkDelegateTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.Logger.Detail("Processing Delegate Transaction for ProcessFee", signedTx)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runNetworkDelegate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	delegate := NetworkDelegate{}
	err := delegate.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, delegate.Tags(), err)
	}

	// Check if delegation address has funds
	coin := delegate.Amount.ToCoinWithBase(ctx.Currencies)
	if !coin.IsValid() {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidAmount, delegate.Tags(), errors.New("Coin is not valid"))
	}
	if coin.Currency.Name != "OLT" {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidCurrency, delegate.Tags(), errors.New("currency is not OLT"))
	}
	currencyOlt, ok := ctx.Currencies.GetCurrencyByName("OLT")
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

	//Get Delegation Pool
	poolList, err := ctx.GovernanceStore.GetPoolList()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, gov.ErrPoolList, delegate.Tags(), err)
	}
	if _, ok := poolList["DelegationPool"]; !ok {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrPoolDoesNotExist, delegate.Tags(), err)
	}
	delagationPool := poolList["DelegationPool"]

	//Add balance to pool
	err = ctx.Balances.AddToAddress(delagationPool, coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, delegate.Tags(), err)
	}

	//Store the entries in network_delegation_store
	err = ctx.NetwkDelegators.Deleg.Set(delegate.DelegationAddress, &coin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, delegate.Tags(), err)
	}

	oldBalance, err := ctx.Balances.GetBalance(delagationPool, ctx.Currencies)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidCurrency, delegate.Tags(), errors.Wrap(err, "Pool is not Funded by OLT"))
	}
	updatedBalance := oldBalance.GetCoin(currencyOlt).Plus(coin)
	ctx.Logger.Infof("Delegation Pool has been updated , New Network Delegation amount: %s ", updatedBalance.String())

	return helpers.LogAndReturnTrue(ctx.Logger, delegate.Tags(), "Success")
}
