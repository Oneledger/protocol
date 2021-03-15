package ons

import (
	"bytes"
	"encoding/json"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &RenewDomain{}

/*
		RenewDomain
This transaction is used to renew
*/
type RenewDomain struct {
	//Owner of sub Domain
	Owner action.Address `json:"owner"`

	//Name of New Sub domain
	Name ons.Name `json:"name"`

	//Amount Added to extend duration of subscription
	BuyingPrice action.Amount `json:"buyingPrice"`
}

func (r RenewDomain) Signers() []action.Address {
	return []action.Address{r.Owner}
}

func (r RenewDomain) Type() action.Type {
	return action.DOMAIN_RENEW
}

func (r RenewDomain) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(r.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: r.Owner.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (r RenewDomain) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RenewDomain) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r RenewDomain) OnsName() string {
	return r.Name.String()
}

var _ action.Tx = &RenewDomainTx{}

type RenewDomainTx struct {
}

func (r RenewDomainTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	renewDomain := &RenewDomain{}
	err := renewDomain.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	//Validate whether signers match those of the transaction and verify the signed transaction.
	err = action.ValidateBasic(signedTx.RawBytes(), renewDomain.Signers(), signedTx.Signatures)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	//Verify fee currency is valid and the amount exceeds the minimum.

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	// basic validation
	if renewDomain.Owner == nil || len(renewDomain.Name) <= 0 {
		return false, action.ErrMissingData
	}

	// check if Name is Valid and not a sub domain
	if !renewDomain.Name.IsValid() || renewDomain.Name.IsSub() {
		return false, ErrInvalidDomain
	}

	// the buying currency must be OLT
	c, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if c.Name != renewDomain.BuyingPrice.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, renewDomain.BuyingPrice.String())
	}

	return true, nil
}

func (r RenewDomainTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runRenew(ctx, tx)
}

func (r RenewDomainTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runRenew(ctx, tx)
}

func (r RenewDomainTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runRenew(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	renewDomain := &RenewDomain{}
	err := renewDomain.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	opt, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, gov.ErrGetONSOptions, renewDomain.Tags(), err)
	}
	coin := renewDomain.BuyingPrice.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(opt.PerBlockFees)) {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrNotEnoughFund, renewDomain.Tags(), errors.New("Less than per block fees"))
	}

	// domain should not be a sub domain
	if renewDomain.Name.IsSub() {
		return false, action.Response{Log: "renew sub domain is not possible"}
	}

	// Check if domain is active, return error if it isn't
	domain, err := ctx.Domains.Get(renewDomain.Name)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !domain.IsChangeable(ctx.Header.Height) {
		return false, action.Response{Log: "domain is not changeable"}
	}

	// if domain is expired it can't be renewed
	if domain.IsExpired(ctx.State.Version()) {
		return false, action.Response{Log: "domain already expired, need to purchase again"}
	}

	// the sender must be the owner of the domain
	if !bytes.Equal(renewDomain.Owner, domain.Owner) {
		return false, action.Response{Log: "only domain owner can renew a domain"}
	}

	//Transfer funds to the fee pool
	price := renewDomain.BuyingPrice.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(renewDomain.Owner, price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	// calculate the blocks

	extend, err := calculateRenewal(&renewDomain.BuyingPrice.Value, &opt.PerBlockFees)
	if err != nil {
		return false, action.Response{
			Log: err.Error(),
		}
	}

	// increase the expiry height & save domain
	domain.AddToExpire(extend)
	domain.SetLastUpdatedHeight(ctx.Header.Height)

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	// set expiry of all subdomains
	ctx.Domains.IterateSubDomain(domain.Name, func(subname ons.Name, subdomain *ons.Domain) bool {

		subdomain.ExpireHeight = domain.ExpireHeight
		err := ctx.Domains.Set(subdomain)
		if err != nil {
			ctx.Logger.Error("failed to update sub domain expiry ", subdomain.Name, err)
			return false
		}

		domain.SetLastUpdatedHeight(ctx.Header.Height)
		return false
	})

	return true, action.Response{Events: action.GetEvent(renewDomain.Tags(), "renew_domain")}
}
