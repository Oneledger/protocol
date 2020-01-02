package ons

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &CreateSubDomain{}

type CreateSubDomain struct {
	//Owner of sub Domain
	Owner action.Address `json:"owner"`

	//Beneficiary Address containing balance
	Beneficiary action.Address `json:"account"`

	//Name of New Sub domain
	Name ons.Name `json:"name"`

	//Amount paid for the Sub Domain
	BuyingPrice action.Amount `json:"price"`

	Uri string `json:"uri"`
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
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, errors.Wrap(err, err.Error())
	}

	if createSub.Owner == nil || len(createSub.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !createSub.Name.IsValid() {
		return false, ErrInvalidDomain
	}

	c, _ := ctx.Currencies.GetCurrencyById(0)
	if c.Name != createSub.BuyingPrice.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, createSub.BuyingPrice.String())
	}

	coin := createSub.BuyingPrice.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(ctx.Domains.GetOptions().BaseDomainPrice)) {
		return false, action.ErrNotEnoughFund
	}
	return true, nil
}

func (crSubD CreateSubDomainTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runCreateSub(ctx, tx)
}

func (crSubD CreateSubDomainTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runCreateSub(ctx, tx)
}

func (crSubD CreateSubDomainTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runCreateSub(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	createSub := &CreateSubDomain{}
	err := createSub.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	//Deduct Price from Owner
	price := createSub.BuyingPrice.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(createSub.Owner, price)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, hex.EncodeToString(createSub.Owner)).Error()}
	}

	//Add Price to the Fee Pool
	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if createSub.Name.IsSub() {
		return false, action.Response{Log: "invalid sub domain name"}
	}
	parentName, err := createSub.Name.GetParentName()
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}
	parent, err := ctx.Domains.Get(parentName)
	//Check if Parent exists in Domain Store
	if err != nil {
		return false, action.Response{Log: "Parent domain doesn't exist, cannot create sub domain!"}
	}
	if !bytes.Equal(parent.OwnerAddress, createSub.Owner) {
		return false, action.Response{Log: "parent domain not owned"}
	}

	//Create Sub Domain Key for Domain Store. Then Check to see if it already exists.
	if ctx.Domains.Exists(createSub.Name) {
		return false, action.Response{Log: fmt.Sprintf("sub domain already exist: %s", createSub.Name)}
	}

	opt := ctx.Domains.GetOptions()

	if createSub.BuyingPrice.Value.BigInt().Cmp(opt.BaseDomainPrice.BigInt()) < 0 {
		return false, action.Response{Log: action.ErrInvalidAmount.Error()}
	}

	ok := verifyDomainName(createSub.Name, opt)
	if !ok {
		return false, action.Response{
			Log: ErrInvalidDomain.Error(),
		}
	}
	if len(createSub.Uri) > 0 {
		ok = opt.IsValidURI(createSub.Uri)
		if !ok {
			return false, action.Response{
				Log: "invalid uri provided",
			}
		}
	}
	//TODO: Need to get fields from Parent domain to populate the sub domain. URI, Expiry height, etc.
	domain, err := ons.NewDomain(
		createSub.Owner,
		createSub.Beneficiary,
		createSub.Name.String(),
		parentName.String(),
		ctx.Header.Height,
		createSub.Uri,
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

	result := action.Response{
		Tags: createSub.Tags(),
	}
	return true, result
}
