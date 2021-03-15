package ons

import (
	"encoding/json"
	"math/big"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainPurchase{}

/*
		DomainPurchase

This transaction lets any buyer purchase a domain which is on sale or has expired.

For on sale domains the Offering should be more than the sale price and for expired domains the offering should be
more than the base domain price.

The expiry height is reset based on the extra OLT offered.
*/
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

func (dp DomainPurchase) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)
	tag0 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(dp.Type().String()),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: dp.Buyer,
	}
	tag2 := kv.Pair{
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

	// the currency should be OLT
	c, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if c.Name != buy.Offering.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, buy.Offering.String())
	}

	// validate the sender and receiver are not nil
	if buy.Buyer == nil || len(buy.Name) == 0 {
		return false, action.ErrMissingData
	}

	// check if name is valid
	if !buy.Name.IsValid() {
		return false, ErrInvalidDomain
	}

	return true, nil
}

func (domainPurchaseTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runPurchaseDomain(ctx, tx)
}

func (domainPurchaseTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runPurchaseDomain(ctx, tx)
}

func (domainPurchaseTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runPurchaseDomain(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	buy := &DomainPurchase{}
	err := buy.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	// domain ,ust be already created
	domain, err := ctx.Domains.Get(buy.Name)
	if err != nil {
		if err == ons.ErrDomainNotFound {
			return false, action.Response{Log: "domain not found"}
		}
		return false, action.Response{Log: "error getting domain"}
	}

	// if domain is not on sale or not expired, return error
	if !domain.OnSaleFlag && (ctx.State.Version() <= domain.ExpireHeight) {
		return false, action.Response{Log: "domain is not on sale or expired"}
	}

	// A sub domain cannot be purchased
	if domain.Name.IsSub() {
		return false, action.Response{Log: "cannot buy subdomain"}
	}

	olt, ok := ctx.Currencies.GetCurrencyByName(buy.Offering.Currency)
	if !ok {
		return false, action.Response{Log: action.ErrInvalidCurrency.Error()}
	}

	remain := buy.Offering.ToCoin(ctx.Currencies)

	opt, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, gov.ErrGetONSOptions, buy.Tags(), err)
	}
	var extend int64
	// if the domain is on sale and not expired
	if (ctx.State.Version() <= domain.ExpireHeight) && domain.OnSaleFlag {

		sale := olt.NewCoinFromAmount(*domain.SalePrice)
		// offering should be more than sale price
		if !sale.LessThanEqualCoin(olt.NewCoinFromAmount(buy.Offering.Value)) {
			return false, action.Response{Log: "offering is not enough"}
		}

		// debit the sale price from buyer balance
		err := ctx.Balances.MinusFromAddress(buy.Buyer, sale)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}

		// credit the sale price to the previous owner
		err = ctx.Balances.AddToAddress(domain.Owner, sale)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}

		// deduct the saleprice from the offering
		remain, err = remain.Minus(sale)
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}

		extend = big.NewInt(0).Div(remain.Amount.BigInt(), opt.PerBlockFees.BigInt()).Int64()

	} else {
		// calculate expiry from the buying price
		extend, err = calculateExpiry(&buy.Offering.Value, &opt.BaseDomainPrice, &opt.PerBlockFees)
		if err != nil {
			return false, action.Response{
				Log: err.Error(),
			}
		}

	}

	// calculate the number of blocks by which to extend the expiry height

	// minus the domain life charges from the buyer
	err = ctx.Balances.MinusFromAddress(buy.Buyer, remain)
	if err != nil {
		return false, action.Response{Log: "error deducting balance for purchase expired domain: " + err.Error()}
	}

	// add the remain to fee pool
	err = ctx.FeePool.AddToPool(remain)
	if err != nil {
		return false, action.Response{Log: "error adding domain purchase: " + err.Error()}
	}

	previousOwner := domain.Owner

	domain.ResetAfterSale(buy.Buyer, buy.Account, extend, ctx.State.Version())

	err = ctx.Domains.DeleteAllSubdomains(domain.Name)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to update domain").Error()}
	}
	return true, action.Response{Events: action.GetEvent(buy.Tags(), "purchase_domain"), Info: previousOwner.Humanize()}
}
