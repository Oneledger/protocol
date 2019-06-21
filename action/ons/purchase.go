package ons

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type DomainPurchase struct {
	Name     string
	Buyer    action.Address
	Account  action.Address
	Offering action.Amount
}

var _ Ons = DomainPurchase{}

func (dp DomainPurchase) OnsName() string {
	return dp.Name
}

func (dp DomainPurchase) Signers() []action.Address {
	return []action.Address{dp.Buyer}
}

func (dp DomainPurchase) Type() action.Type {
	return action.DOMAIN_PURCHASE
}

func (dp DomainPurchase) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)
	tag0 := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(dp.Type().String()),
	}
	tag1 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: dp.Buyer,
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.domain_name"),
		Value: []byte(dp.Name),
	}

	tags = append(tags, tag0, tag1, tag2)
	return tags
}

type domainPurchaseTx struct {
}

func (domainPurchaseTx) Validate(ctx *action.Context, msg action.Msg, fee action.Fee, memo string, signatures []action.Signature) (bool, error) {

	// validate basic signature
	ok, err := action.ValidateBasic(msg, fee, memo, signatures)
	if err != nil {
		return ok, err
	}

	buy, ok := msg.(*DomainPurchase)
	if !ok {
		return false, action.ErrWrongTxType
	}

	if !buy.Offering.IsValid(ctx) {
		return false, errors.Wrap(action.ErrInvalidAmount, buy.Offering.String())
	}

	// validate the sender and receiver are not nil
	if buy.Buyer == nil || len(buy.Name) == 0 {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainPurchaseTx) ProcessCheck(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {

	buy, ok := msg.(*DomainPurchase)
	if !ok {
		return false, action.Response{Log: "failed to cast msg"}
	}

	domain, err := ctx.Domains.Get(buy.Name, false)
	if err != nil {
		if err == ons.ErrDomainNotFound {
			return false, action.Response{Log: "domain not found"}
		}
		return false, action.Response{Log: "error getting domain"}
	}

	if !domain.OnSaleFlag {
		return false, action.Response{Log: "domain is not on sale"}
	}

	if domain.SalePrice.LessThanEqualCoin(buy.Offering.ToCoin(ctx)) {
		return false, action.Response{Log: "offering price not enough"}
	}

	buyerBalance, err := ctx.Balances.Get(buy.Buyer.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to get buyerBalance balance").Error()}
	}

	buyerBalance, err = buyerBalance.MinusCoin(buy.Offering.ToCoin(ctx))
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	return true, action.Response{Tags: buy.Tags()}

}

func (domainPurchaseTx) ProcessDeliver(ctx *action.Context, msg action.Msg, fee action.Fee) (bool, action.Response) {
	buy, ok := msg.(*DomainPurchase)
	if !ok {
		return false, action.Response{Log: "failed to cast msg"}
	}

	domain, err := ctx.Domains.Get(buy.Name, false)
	if err != nil {
		if err == ons.ErrDomainNotFound {
			return false, action.Response{Log: "domain not found"}
		}
		return false, action.Response{Log: "error getting domain"}
	}

	if !domain.OnSaleFlag {
		return false, action.Response{Log: "domain is not on sale"}
	}

	coin := buy.Offering.ToCoin(ctx)
	if domain.SalePrice.LessThanEqualCoin(coin) {
		return false, action.Response{Log: "offering price not enough"}
	}

	buyerBalance, err := ctx.Balances.Get(buy.Buyer.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to get buyerBalance balance").Error()}
	}

	buyerBalance, err = buyerBalance.MinusCoin(coin)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	salerBalance, err := ctx.Balances.Get(domain.OwnerAddress, false)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to get salerBalance balance").Error()}
	}
	salerBalance.AddCoin(coin)

	err = ctx.Balances.Set(buy.Buyer, *buyerBalance)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	err = ctx.Balances.Set(domain.OwnerAddress, *salerBalance)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	domain.OwnerAddress = buy.Buyer

	if buy.Account != nil {
		domain.SetAccountAddress(buy.Account)
	} else {
		domain.SetAccountAddress(buy.Buyer)
	}
	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to update domain").Error()}
	}
	return true, action.Response{Tags: buy.Tags()}
}

func (domainPurchaseTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
