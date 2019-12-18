package ons

import (
	"encoding/json"
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

var _ Ons = &RenewDomain{}

type RenewDomain struct {
	//Owner of sub Domain
	Owner action.Address `json:"owner"`

	//Name of New Sub domain
	Name ons.Name `json:"name"`

	//Amount Added to extend duration of subscription
	Price action.Amount `json:"price"`
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

func (r RenewDomain) Unmarshal(data []byte) error {
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
	err = action.ValidateFee(ctx.FeeOpt, signedTx.Fee)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	if renewDomain.Owner == nil || len(renewDomain.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !renewDomain.Price.IsValid(ctx.Currencies) || renewDomain.Price.Currency != "OLT" {
		return false, action.ErrInvalidFeeCurrency
	}

	//Checking if amount satisfies minimum requirement.
	coin := renewDomain.Price.ToCoin(ctx.Currencies)
	if coin.LessThanCoin(coin.Currency.NewCoinFromInt(CREATE_PRICE)) {
		return false, action.ErrNotEnoughFund
	}

	return true, nil
}

func (r RenewDomainTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	renewDomain := &RenewDomain{}
	err := renewDomain.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	//Check if domain exists, return error if it doesn't.
	if !ctx.Domains.Exists(renewDomain.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain: %s doesn't exist, cannot renew.", renewDomain.Name)}
	}

	//Check if domain is active, return error if it isn't
	domain, err := ctx.Domains.Get(renewDomain.Name)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	if !domain.ActiveFlag {
		return false, action.Response{Log: "domain is not active, cannot renew"}
	}

	//Transfer funds to the fee pool
	price := renewDomain.Price.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(renewDomain.Owner, price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	result := action.Response{Tags: renewDomain.Tags()}

	return true, result
}

func (r RenewDomainTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	renewDomain := &RenewDomain{}
	err := renewDomain.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	//Check if domain exists, return error if it doesn't.
	if !ctx.Domains.Exists(renewDomain.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain: %s doesn't exist, cannot renew.", renewDomain.Name)}
	}

	//Check if domain is active, return error if it isn't
	domain, err := ctx.Domains.Get(renewDomain.Name)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	if !domain.ActiveFlag {
		return false, action.Response{Log: "domain is not active, cannot renew"}
	}

	//Transfer funds to the fee pool
	price := renewDomain.Price.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(renewDomain.Owner, price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	//Calculate Number of blocks to extend subscription by. Then update domain store.
	//TODO: Need to define price per block, calculate block expiry height, then update all sub domains.

	return true, action.Response{Tags: renewDomain.Tags()}
}

func (r RenewDomainTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
