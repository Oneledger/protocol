package ons

import (
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

var _ Ons = &DomainPurchase{}

type DomainPurchase struct {
	Name     string         `json:"name"`
	Buyer    action.Address `json:"buyer"`
	Account  action.Address `json:"account"`
	Offering action.Amount  `json:"offering"`
}

func (dp DomainPurchase) Marshal() ([]byte, error) {
	return json.Marshal(dp)
}

func (dp *DomainPurchase) Unmarshal(data []byte) error {
	return json.Unmarshal(data, dp)
}

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

func (domainPurchaseTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	buy := &DomainPurchase{}
	err := buy.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic signature
	ok, err := action.ValidateBasic(tx.RawBytes(), buy.Signers(), tx.Signatures)
	if err != nil {
		return ok, err
	}

	if !buy.Offering.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, buy.Offering.String())
	}

	// validate the sender and receiver are not nil
	if buy.Buyer == nil || len(buy.Name) == 0 {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainPurchaseTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	buy := &DomainPurchase{}
	err := buy.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
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

	if !domain.SalePrice.LessThanEqualCoin(buy.Offering.ToCoin(ctx.Currencies)) {
		return false, action.Response{Log: "offering price not enough"}
	}

	buyerBalance, err := ctx.Balances.Get(buy.Buyer.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to get buyer balance").Error()}
	}

	buyerBalance, err = buyerBalance.MinusCoin(buy.Offering.ToCoin(ctx.Currencies))
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	return true, action.Response{Tags: buy.Tags()}

}

func (domainPurchaseTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	buy := &DomainPurchase{}
	err := buy.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
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

	coin := buy.Offering.ToCoin(ctx.Currencies)
	if !domain.SalePrice.LessThanEqualCoin(coin) {
		return false, action.Response{Log: "offering price not enough"}
	}

	buyerBalance, err := ctx.Balances.Get(buy.Buyer.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to get buyer balance").Error()}
	}

	buyerBalance, err = buyerBalance.MinusCoin(coin)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	salerBalance, err := ctx.Balances.Get(domain.OwnerAddress, false)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to get saler balance").Error()}
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

	domain.CancelSale()
	domain.Activate()

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to update domain").Error()}
	}
	return true, action.Response{Tags: buy.Tags()}
}

func (domainPurchaseTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
