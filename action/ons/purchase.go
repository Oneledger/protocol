package ons

import (
	"encoding/json"
	"math/big"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainPurchase{}

type DomainPurchase struct {
	Name     ons.Name       `json:"name"`
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
	return dp.Name.String()
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
	err = action.ValidateBasic(tx.RawBytes(), buy.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	c, _ := ctx.Currencies.GetCurrencyById(0)
	if c.Name != buy.Offering.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, buy.Offering.String())
	}

	// validate the sender and receiver are not nil
	if buy.Buyer == nil || len(buy.Name) == 0 {
		return false, action.ErrMissingData
	}

	if !buy.Name.IsValid() {
		return false, ErrInvalidDomain
	}

	if c.Name != buy.Offering.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, buy.Offering.String())
	}

	coin := buy.Offering.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(ctx.Domains.GetOptions().PerBlockFees)) {
		return false, action.ErrNotEnoughFund
	}

	return true, nil
}

func (domainPurchaseTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runPurchaseDomain(ctx, tx)
}

func (domainPurchaseTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runPurchaseDomain(ctx, tx)
}

func (domainPurchaseTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runPurchaseDomain(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	buy := &DomainPurchase{}
	err := buy.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	domain, err := ctx.Domains.Get(buy.Name)
	if err != nil {
		if err == ons.ErrDomainNotFound {
			return false, action.Response{Log: "domain not found"}
		}
		return false, action.Response{Log: "error getting domain"}
	}

	if !domain.OnSaleFlag && (ctx.State.Version() <= domain.ExpireHeight) {
		return false, action.Response{Log: "domain is not on sale or expired"}
	}

	olt, _ := ctx.Currencies.GetCurrencyByName(buy.Offering.Currency)

	remain := buy.Offering.ToCoin(ctx.Currencies)

	if (ctx.State.Version() <= domain.ExpireHeight) && domain.OnSaleFlag {
		if domain.SalePrice.LessThanEqualCoin(olt.NewCoinFromAmount(buy.Offering.Value)) {
			return false, action.Response{Log: "offering is not enough"}
		}

		err := ctx.Balances.MinusFromAddress(buy.Buyer, domain.SalePrice)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}

		err = ctx.Balances.AddToAddress(domain.Beneficiary, domain.SalePrice)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}
		remain, err = remain.Minus(domain.SalePrice)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}
	}

	opt := ctx.Domains.GetOptions()
	extend := big.NewInt(0).Div(remain.Amount.BigInt(), opt.PerBlockFees.BigInt()).Int64()

	err = ctx.Balances.MinusFromAddress(buy.Buyer, remain)
	if err != nil {
		return false, action.Response{Log: "error deducting balance for purchase expired domain: " + err.Error()}
	}

	err = ctx.FeePool.AddToPool(remain)
	if err != nil {
		return false, action.Response{Log: "error adding domain purchase: " + err.Error()}
	}

	domain.ResetAfterSale(buy.Buyer, buy.Buyer, extend, ctx.State.Version())

	err = ctx.Domains.DeleteAllSubdomains(domain.Name)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to update domain").Error()}
	}
	return true, action.Response{Tags: buy.Tags()}
}
