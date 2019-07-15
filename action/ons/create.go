package ons

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/data/ons"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
)

var _ Ons = &DomainCreate{}

type DomainCreate struct {
	Owner   action.Address
	Account action.Address
	Name    string
	Price   action.Amount
}

func (dc DomainCreate) Marshal() ([]byte, error) {
	return json.Marshal(dc)
}

func (dc *DomainCreate) Unmarshal(data []byte) error {
	return json.Unmarshal(data, dc)
}

func (dc DomainCreate) OnsName() string {
	return dc.Name
}

func (dc DomainCreate) Signers() []action.Address {
	return []action.Address{dc.Owner}
}

func (dc DomainCreate) Type() action.Type {
	return action.DOMAIN_CREATE
}

func (dc DomainCreate) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(dc.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: dc.Owner.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

var _ action.Tx = domainCreateTx{}

type domainCreateTx struct {
}

func (domainCreateTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {

	create := &DomainCreate{}
	err := create.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	ok, err := action.ValidateBasic(tx.RawBytes(), create.Signers(), tx.Signatures)
	if err != nil {
		return ok, err
	}

	if create.Owner == nil || len(create.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !create.Price.IsValid(ctx.Currencies) || create.Price.Currency != "OLT" {
		return false, action.ErrInvalidAmount
	}

	coin := create.Price.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromInt(CREATE_PRICE)) {
		return false, action.ErrNotEnoughFund
	}

	return true, nil
}

func (domainCreateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	create := &DomainCreate{}
	err := create.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	b, err := ctx.Balances.Get(create.Owner.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("failed to get balance for owner: %s", hex.EncodeToString(create.Owner))}
	}
	price := create.Price.ToCoin(ctx.Currencies)

	//just verify if balance is enough or not, don't set to db
	b, err = b.MinusCoin(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if ctx.Domains.Exists(create.Name) {
		return false, action.Response{Log: fmt.Sprintf("Domain already exist: %s", create.Name)}
	}
	result := action.Response{
		Tags: create.Tags(),
	}

	return true, result
}

func (domainCreateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	create := &DomainCreate{}
	err := create.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	b, err := ctx.Balances.Get(create.Owner.Bytes(), false)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("failed to get balance for owner: %s", hex.EncodeToString(create.Owner))}
	}
	price := create.Price.ToCoin(ctx.Currencies)

	// verify balance and set to db, the price for create domain is just burned for now.
	//todo: pay the price to fee pool that will be shared by validators at the fee distribution time.
	b, err = b.MinusCoin(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.Balances.Set(create.Owner.Bytes(), *b)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "set balance of owner").Error()}
	}

	//check domain existence and set to db
	if ctx.Domains.Exists(create.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain already exist: %s", create.Name)}
	}
	domain := ons.NewDomain(
		create.Owner,
		create.Account,
		create.Name,
		ctx.Header.Height,
	)
	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	ctx.Logger.Error(domain)
	ctx.Logger.Error(ctx.Domains.Get(domain.Name, false))
	result := action.Response{
		Tags: create.Tags(),
	}
	return true, result
}

func (domainCreateTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
