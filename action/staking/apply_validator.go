package staking

import (
	"github.com/Oneledger/protocol/action"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/tendermint/tendermint/libs/common"
)

var _ action.Msg = ApplyValidator{}

type ApplyValidator struct {
	Address          action.Address
	Stake            action.Amount
	NodeName         string
	ValidatorAddress action.Address
	ValidatorPubKey  keys.PublicKey
	Purge            bool
}

func (apply ApplyValidator) Signers() []action.Address {
	return []action.Address{apply.Address.Bytes()}
}

func (apply ApplyValidator) Type() action.Type {
	return action.APPLYVALIDATOR
}

var _ action.Tx = applyTx{}

type applyTx struct {
}

func (applyTx) Validate(ctx *action.Context, msg action.Msg, fee action.Fee, memo string, signatures []action.Signature) (bool, error) {
	ok, err := action.ValidateBasic(msg, fee, memo, signatures)
	if err != nil {
		return ok, err
	}

	apply, ok := msg.(*ApplyValidator)
	if !ok {
		return false, action.ErrWrongTxType
	}

	if msg == nil || len(apply.Address) == 0 {
		return false, action.ErrMissingData
	}

	if apply.ValidatorAddress == nil {
		return false, action.ErrMissingData
	}
	_, err = apply.ValidatorPubKey.GetHandler()
	if err != nil {
		return false, action.ErrInvalidPubkey
	}

	coin := apply.Stake.ToCoin(ctx)
	if coin.LessThanCoin(coin.Currency.NewCoinFromInt(0)) {
		return false, action.ErrInvalidAmount
	}

	if coin.Currency.Name != "VT" {
		return false, action.ErrInvalidAmount
	}

	return true, nil
}

func (a applyTx) ProcessCheck(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	apply, ok := msg.(*ApplyValidator)
	if !ok {
		return false, action.Response{Log: "Apply validator cast failed"}
	}
	result, err := checkBalances(ctx, apply.Address, apply.Stake)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	return result, action.Response{Tags: apply.Tags()}
}

func checkBalances(ctx *action.Context, address action.Address, stake action.Amount) (bool, error) {

	balances := ctx.Balances

	// check identity's VT is equal to the stake
	balance, err := balances.Get(address, false)
	if err != nil {
		return false, action.ErrNotEnoughFund
	}
	c, ok := ctx.Currencies.GetCurrencyByName("VT")
	if !ok {
		return false, action.ErrInvalidAmount
	}
	if balance.GetCoin(c).LessThanCoin(stake.ToCoin(ctx)) {
		return false, action.ErrNotEnoughFund
	}
	return true, nil
}

func (applyTx) ProcessDeliver(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	apply, ok := msg.(*ApplyValidator)
	if !ok {
		return false, action.Response{Log: "Apply validator cast failed"}
	}
	_, err := checkBalances(ctx, apply.Address, apply.Stake)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	validators := ctx.Validators

	balances := ctx.Balances
	balance, err := balances.Get(apply.Address.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	b, err := balance.MinusCoin(apply.Stake.ToCoin(ctx))
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = balances.Set(apply.Address.Bytes(), *b)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !apply.Purge {
		stake := identity.Stake{
			ValidatorAddress: apply.ValidatorAddress,
			StakeAddress:     apply.Address,
			Pubkey:           apply.ValidatorPubKey,
			Name:             apply.NodeName,
			Amount:           apply.Stake.ToCoin(ctx),
		}
		err = validators.HandleStake(stake)
	} else {
		unstake := identity.Unstake{
			Address: apply.ValidatorAddress,
			Amount:  apply.Stake.ToCoin(ctx),
		}
		err = validators.HandleUnstake(unstake)
	}
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	return true, action.Response{Tags: apply.Tags()}
}

func (applyTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	// TODO: implement fee charge for apply
	return true, action.Response{Info: "Unimplemented"}
}

func (apply ApplyValidator) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(apply.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: apply.Address.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}
