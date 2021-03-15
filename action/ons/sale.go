/*

 */

package ons

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainSale{}

type DomainSale struct {
	Name         ons.Name       `json:"name"`
	OwnerAddress action.Address `json:"ownerAddress"`
	Price        action.Amount  `json:"price"`
	CancelSale   bool           `json:"cancelSale"`
}

func (s DomainSale) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *DomainSale) Unmarshal(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s DomainSale) OnsName() string {
	return s.Name.String()
}

func (DomainSale) Type() action.Type {
	return action.DOMAIN_SELL
}

func (s DomainSale) Signers() []action.Address {
	return []action.Address{s.OwnerAddress}
}

func (s DomainSale) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)
	tag0 := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(action.DOMAIN_SELL.String()),
	}
	tag1 := kv.Pair{
		Key:   []byte("tx.owner"),
		Value: s.OwnerAddress,
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.domain_name"),
		Value: []byte(s.Name),
	}

	tags = append(tags, tag0, tag1, tag2)
	if s.CancelSale {
		tag3 := kv.Pair{
			Key:   []byte("tx.is_cancel"),
			Value: []byte{0xff},
		}
		tags = append(tags, tag3)
	}

	return tags
}

/*


	domainSaleTx

*/

type domainSaleTx struct {
}

var _ action.Tx = domainSaleTx{}

func (domainSaleTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {

	sale := &DomainSale{}
	err := sale.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), sale.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if !sale.Price.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, sale.Price.String())
	}

	// validate the sender and receiver are not nil
	if sale.OwnerAddress == nil || len(sale.Name) == 0 {
		return false, action.ErrMissingData
	}

	if !sale.Name.IsValid() || sale.Name.IsSub() {
		return false, ErrInvalidDomain
	}
	c, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if c.Name != sale.Price.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, sale.Price.String())
	}

	return true, nil
}

func (domainSaleTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runDomainSale(ctx, tx)
}

func (domainSaleTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runDomainSale(ctx, tx)
}

func (domainSaleTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runDomainSale(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	sale := &DomainSale{}
	err := sale.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	opt, err := ctx.GovernanceStore.GetONSOptions()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, gov.ErrGetONSOptions, sale.Tags(), err)
	}
	coin := sale.Price.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(opt.PerBlockFees)) {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrNotEnoughFund, sale.Tags(), err)
	}

	if !sale.Price.IsValid(ctx.Currencies) {
		return false, action.Response{Log: "invalid price amount"}
	}

	if sale.OwnerAddress == nil {
		return false, action.Response{Log: "invalid owner"}
	}

	if sale.Name.IsSub() {
		return false, action.Response{Log: "put sub domain on sale not allowed"}
	}

	domain, err := ctx.Domains.Get(sale.Name)
	if err != nil {
		if err == ons.ErrDomainNotFound {
			return false, action.Response{Log: "domain not found"}
		}
		return false, action.Response{Log: "error getting domain"}
	}

	// verify the ownership
	if bytes.Compare(domain.Owner, sale.OwnerAddress) != 0 {
		return false, action.Response{Log: "not the owner"}
	}

	if !domain.IsChangeable(ctx.Header.Height) {
		log := fmt.Sprintf("domain not changeable; name: %s, lastUpdateheight %d",
			domain.Name, domain.LastUpdateHeight)
		return false, action.Response{Log: log}
	}

	if domain.IsExpired(ctx.Header.Height) {
		return false, action.Response{Log: "domain expired"}
	}

	if sale.CancelSale {
		domain.CancelSale()
	} else {
		domain.PutOnSale(sale.Price.Value)
	}
	domain.LastUpdateHeight = ctx.Header.Height

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: "error updating domain store"}
	}

	return true, action.Response{Events: action.GetEvent(sale.Tags(), "domain_on_sale")}
}
