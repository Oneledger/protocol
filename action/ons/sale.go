/*

 */

package ons

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainSale{}

type DomainSale struct {
	DomainName   string         `json:"domainName"`
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
	return s.DomainName
}

func (DomainSale) Type() action.Type {
	return action.DOMAIN_SELL
}

func (s DomainSale) Signers() []action.Address {
	return []action.Address{s.OwnerAddress}
}

func (s DomainSale) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)
	tag0 := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.DOMAIN_SELL.String()),
	}
	tag1 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: s.OwnerAddress,
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.domain_name"),
		Value: []byte(s.DomainName),
	}

	tags = append(tags, tag0, tag1, tag2)
	if s.CancelSale {
		tag3 := common.KVPair{
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

	err = action.ValidateFee(ctx.FeeOpt, tx.Fee)
	if err != nil {
		return false, err
	}

	if !sale.Price.IsValid(ctx.Currencies) {
		return false, errors.Wrap(action.ErrInvalidAmount, sale.Price.String())
	}

	// validate the sender and receiver are not nil
	if sale.OwnerAddress == nil || len(sale.DomainName) == 0 {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainSaleTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	sale := &DomainSale{}
	err := sale.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !sale.Price.IsValid(ctx.Currencies) {
		return false, action.Response{Log: "invalid price amount"}
	}

	// validate the sender and receiver are not nil
	if sale.OwnerAddress == nil || len(sale.DomainName) <= 0 {
		return false, action.Response{Log: "invalid data"}
	}

	domain, err := ctx.Domains.Get(sale.DomainName)
	if err != nil {
		if err == ons.ErrDomainNotFound {
			return false, action.Response{Log: "domain not found"}
		}
		return false, action.Response{Log: err.Error()}
	}

	if bytes.Compare(domain.OwnerAddress, sale.OwnerAddress) != 0 {
		return false, action.Response{Log: "not the owner"}
	}

	if !domain.IsChangeable(ctx.Header.Height) {
		log := fmt.Sprintf("domain not changeable; name: %s, lastUpdateheight %d",
			domain.Name, domain.LastUpdateHeight)
		return false, action.Response{Log: log}
	}

	// if action to cancel sale and domain is not on sale
	// fail the ProcessCheck
	if sale.CancelSale &&
		!domain.OnSaleFlag {
		return false, action.Response{Log: "cannot cancel sale of domain; domain not on sale"}
	}

	return true, action.Response{Tags: sale.Tags()}
}

func (domainSaleTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	sale := &DomainSale{}
	err := sale.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !sale.Price.IsValid(ctx.Currencies) {
		return false, action.Response{Log: "invalid price amount"}
	}

	// validate the sender and receiver are not nil
	if sale.OwnerAddress == nil || sale.DomainName == "" {
		return false, action.Response{Log: "invalid data"}
	}

	domain, err := ctx.Domains.Get(sale.DomainName)
	if err != nil {
		if err == ons.ErrDomainNotFound {
			return false, action.Response{Log: "domain not found"}
		}
		return false, action.Response{Log: "error getting domain"}
	}

	// verify the ownership
	if bytes.Compare(domain.OwnerAddress, sale.OwnerAddress) != 0 {
		return false, action.Response{Log: "not the owner"}
	}

	if !domain.IsChangeable(ctx.Header.Height) {
		log := fmt.Sprintf("domain not changeable; name: %s, lastUpdateheight %d",
			domain.Name, domain.LastUpdateHeight)
		return false, action.Response{Log: log}
	}

	if sale.CancelSale {
		domain.CancelSale()
	} else {
		domain.PutOnSale(sale.Price.ToCoin(ctx.Currencies))
	}
	domain.LastUpdateHeight = ctx.Header.Height

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: "error updating domain store"}
	}

	return true, action.Response{Tags: sale.Tags()}
}

func (domainSaleTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
