package network_delegation

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/balance"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/network_delegation"
	netwkDeleg "github.com/Oneledger/protocol/data/network_delegation"
)

type Reinvest struct {
	Delegator action.Address `json:"delegator"`
	Amount    action.Amount  `json:"amount"`
}

func (ri Reinvest) Signers() []action.Address {
	return []action.Address{ri.Delegator}
}

func (ri Reinvest) Type() action.Type {
	return action.REWARDS_REINVEST_NETWORK_DELEGATE
}

func (ri Reinvest) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(ri.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.Delegator"),
		Value: ri.Delegator.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.Amount"),
		Value: []byte(ri.Amount.String()),
	}
	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (ri Reinvest) Marshal() ([]byte, error) {
	return json.Marshal(ri)
}

func (ri *Reinvest) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, ri)
}

type delegReinvestRewardsTx struct {
}

func (delegReinvestRewardsTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	invest := Reinvest{}
	err := invest.Unmarshal(signedTx.Data)
	if err != nil {
		return false, err
	}

	err = action.ValidateBasic(signedTx.RawBytes(), invest.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	currency, ok := ctx.Currencies.GetCurrencyByName("OLT")
	if !ok {
		return false, errors.Wrap(action.ErrInvalidCurrency, invest.Amount.Currency)
	}
	if currency.Name != invest.Amount.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, invest.Amount.String())
	}

	err = invest.Delegator.Err()
	if err != nil {
		return false, errors.Wrap(action.ErrInvalidAddress, err.Error())
	}
	return true, nil
}

func (delegReinvestRewardsTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runReinvest(ctx, tx)
}

func (delegReinvestRewardsTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runReinvest(ctx, tx)
}

func (delegReinvestRewardsTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

var _ action.Msg = &Reinvest{}
var _ action.Tx = &delegReinvestRewardsTx{}

func runReinvest(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	invest := Reinvest{}
	err := invest.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrUnserializable, invest.Tags(), err)
	}

	// cut rewards
	coinAmt := invest.Amount.ToCoin(ctx.Currencies)
	err = ctx.NetwkDelegators.Rewards.MinusRewardsBalance(invest.Delegator, coinAmt.Amount)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, netwkDeleg.ErrReinvestRewards, invest.Tags(), err)
	}

	//Add Delegation
	//Get Delegation Pool
	delagationPool, err := ctx.GovernanceStore.GetPoolByName(gov.POOL_DELEGATION)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrPoolDoesNotExist, invest.Tags(), err)
	}

	//Add balance to pool
	err = ctx.Balances.AddToAddress(delagationPool, coinAmt)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, invest.Tags(), err)
	}

	//Add balance to delegation
	currentDelegation, _ := ctx.NetwkDelegators.Deleg.WithPrefix(network_delegation.ActiveType).Get(invest.Delegator)
	newCoin := currentDelegation.Plus(coinAmt)
	err = ctx.NetwkDelegators.Deleg.WithPrefix(network_delegation.ActiveType).Set(invest.Delegator, &newCoin)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, balance.ErrBalanceErrorAddFailed, invest.Tags(), err)
	}

	ctx.Logger.Debugf("Successfully reinvested, delegator= %s, amount= %s", invest.Delegator.String(), coinAmt)
	return helpers.LogAndReturnTrue(ctx.Logger, invest.Tags(), "Success")
}
