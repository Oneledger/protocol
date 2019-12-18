package ons

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

var _ Ons = &CreateSubDomain{}

type CreateSubDomain struct {
	//Owner of sub Domain
	Owner action.Address `json:"owner"`

	//Beneficiary Address containing balance
	Beneficiary action.Address `json:"account"`

	//Name of New Sub domain
	Name ons.Name `json:"name"`

	//Name of Parent domain
	Parent ons.Name `json:"parent"`

	//Amount paid for the Sub Domain
	Price action.Amount `json:"price"`
}

func (c CreateSubDomain) Signers() []action.Address {
	return []action.Address{c.Owner}
}

func (c CreateSubDomain) Type() action.Type {
	return action.DOMAIN_CREATE_SUB
}

func (c CreateSubDomain) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(c.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: c.Owner.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (c CreateSubDomain) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c CreateSubDomain) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c CreateSubDomain) OnsName() string {
	return c.Name.String()
}

var _ action.Tx = CreateSubDomainTx{}

type CreateSubDomainTx struct {
}

func makeSubDomainKey(parent ons.Name, subDomain ons.Name) ons.Name {
	return parent + "_" + subDomain
}

func (crSubD CreateSubDomainTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	createSub := &CreateSubDomain{}
	err := createSub.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	//Validate whether signers match those of the transaction and verify the signed transaction.
	err = action.ValidateBasic(signedTx.RawBytes(), createSub.Signers(), signedTx.Signatures)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	//Verify fee currency is valid and the amount exceeds the minimum.
	err = action.ValidateFee(ctx.FeeOpt, signedTx.Fee)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	if createSub.Owner == nil || len(createSub.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !createSub.Price.IsValid(ctx.Currencies) || createSub.Price.Currency != "OLT" {
		return false, action.ErrInvalidFeeCurrency
	}

	coin := createSub.Price.ToCoin(ctx.Currencies)
	if coin.LessThanCoin(coin.Currency.NewCoinFromInt(CREATE_SUB_PRICE)) {
		return false, action.ErrNotEnoughFund
	}

	return true, nil
}

func (crSubD CreateSubDomainTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	createSub := &CreateSubDomain{}
	err := createSub.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	//Check if Parent exists in Domain Store
	if !ctx.Domains.Exists(createSub.Parent) {
		return false, action.Response{Log: "Parent domain doesn't exist, cannot create sub domain!"}
	}

	//Create Sub Domain Key for Domain Store. Then Check to see if it already exists.
	if ctx.Domains.Exists(makeSubDomainKey(createSub.Parent, createSub.Name)) {
		return false, action.Response{Log: fmt.Sprintf("sub domain already exist: %s", createSub.Name)}
	}

	//Deduct Price from Owner
	price := createSub.Price.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(createSub.Owner, price)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, hex.EncodeToString(createSub.Owner)).Error()}
	}

	//Add Price to the Fee Pool
	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	result := action.Response{Tags: createSub.Tags()}

	return true, result
}

func (crSubD CreateSubDomainTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	createSub := &CreateSubDomain{}
	err := createSub.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	//Check if Parent exists in Domain Store
	if !ctx.Domains.Exists(createSub.Parent) {
		return false, action.Response{Log: "Parent domain doesn't exist, cannot create sub domain!"}
	}

	//Create Sub Domain Key for Domain Store. Then Check to see if it already exists.
	subDomain := makeSubDomainKey(createSub.Parent, createSub.Name)
	if ctx.Domains.Exists(subDomain) {
		return false, action.Response{Log: fmt.Sprintf("sub domain already exist: %s", createSub.Name)}
	}

	//Deduct Price from Owner
	price := createSub.Price.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(createSub.Owner, price)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, hex.EncodeToString(createSub.Owner)).Error()}
	}

	//Add Price to the Fee Pool
	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	parent, err := ctx.Domains.Get(createSub.Parent)
	if err != nil {
		return false, action.Response{
			Log: "Unable to get parent domain",
		}
	}
	//TODO: Need to get fields from Parent domain to populate the sub domain. URI, Expiry height, etc.
	domain, err := ons.NewDomain(
		createSub.Owner,
		createSub.Beneficiary,
		subDomain.String(),
		createSub.Parent.String(),
		ctx.Header.Height,
		"",
		parent.Expiry,
	)

	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	//Save New Sub Domain to Domain Store
	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	ctx.Logger.Error(domain)
	ctx.Logger.Error(ctx.Domains.Get(domain.Name))
	result := action.Response{
		Tags: createSub.Tags(),
	}

	return true, result
}

func (crSubD CreateSubDomainTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
