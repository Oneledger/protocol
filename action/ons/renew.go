package ons

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &RenewDomain{}

type RenewDomain struct {
	//Owner of sub Domain
	Owner action.Address `json:"owner"`

	//Name of New Sub domain
	Name ons.Name `json:"name"`

	//Amount Added to extend duration of subscription
	BuyingPrice action.Amount `json:"price"`
}

func (r RenewDomain) Signers() []action.Address {
	return []action.Address{r.Owner}
}

func (r RenewDomain) Type() action.Type {
	return action.DOMAIN_RENEW
}

func (r RenewDomain) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(r.Type().String()),
	}
	tag2 := common.KVPair{
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

	if renewDomain.Owner == nil || len(renewDomain.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !renewDomain.Name.IsValid() || renewDomain.Name.IsSub() {
		return false, ErrInvalidDomain
	}

	c, _ := ctx.Currencies.GetCurrencyById(0)
	if c.Name != renewDomain.BuyingPrice.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, renewDomain.BuyingPrice.String())
	}

	coin := renewDomain.BuyingPrice.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(ctx.Domains.GetOptions().PerBlockFees)) {
		return false, action.ErrNotEnoughFund
	}

	return true, nil
}

func (r RenewDomainTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runRenew(ctx, tx)
}

func (r RenewDomainTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runRenew(ctx, tx)
}

func (r RenewDomainTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runRenew(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	renewDomain := &RenewDomain{}
	err := renewDomain.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if renewDomain.Name.IsSub() {
		return false, action.Response{Log: "renew sub domain is not possible"}
	}
	//Check if domain is active, return error if it isn't
	domain, err := ctx.Domains.Get(renewDomain.Name)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if domain.IsExpired(ctx.State.Version()) {
		return false, action.Response{Log: "domain already expired, need to purchase again"}
	}

	if !bytes.Equal(renewDomain.Owner, domain.OwnerAddress) {
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

	opt := ctx.Domains.GetOptions()
	extend := big.NewInt(0).Div(renewDomain.BuyingPrice.Value.BigInt(), opt.PerBlockFees.BigInt()).Int64()

	domain.AddToExpire(extend)

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	ctx.Domains.IterateSubDomain(domain.Name, func(subname ons.Name, subdomain *ons.Domain) bool {
		subdomain.ExpireHeight = domain.ExpireHeight
		err := ctx.Domains.Set(subdomain)
		if err != nil {
			ctx.Logger.Error("failed to update sub domain expiry ", subdomain.Name, err)
			return false
		}
		return false
	})
	return true, action.Response{Tags: renewDomain.Tags()}
}
