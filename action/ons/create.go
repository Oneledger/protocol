package ons

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainCreate{}

/*
		DomainCreate

	This transaction is used to create a domain or sub-domain. This transaction is supposed to be sent by creator of an
original domain or the owner of the parent of the sub-domain. The Beneficiary Address can be a valid Oneledger Address,
which will receive funds in any currency sent to this domain. The domain name has to be of the type domain.ol or sub.domain.ol,
always terminating in the .ol top level domain. The URI can be a valid RFC-3986 compliant string of characters, which
is expected to be used as a meta resource location.

The BuyingPrice which must be in OLT, for domain is the sum of the creation price and the domain per block fees.

*/
type DomainCreate struct {
	Owner       action.Address `json:"owner"`
	Beneficiary action.Address `json:"account"`
	Name        ons.Name       `json:"name"`
	Uri         string         `json:"uri"`
	BuyingPrice action.Amount  `json:"buyingPrice"`
}

func (dc DomainCreate) Marshal() ([]byte, error) {
	return json.Marshal(dc)
}

func (dc *DomainCreate) Unmarshal(data []byte) error {
	return json.Unmarshal(data, dc)
}

func (dc DomainCreate) OnsName() string {
	return dc.Name.String()
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

	//Validate whether signers match those of the transaction and verify the signed transaction.
	err = action.ValidateBasic(tx.RawBytes(), create.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	//Verify fee currency is valid and the amount exceeds the minimum.
	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if create.Owner == nil || len(create.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !create.Name.IsValid() {
		return false, ErrInvalidDomain
	}

	// the currency should be OLT
	c, ok := ctx.Currencies.GetCurrencyById(0)
	if !ok {
		panic("no default currency available in the network")
	}
	if c.Name != create.BuyingPrice.Currency {
		return false, errors.Wrap(action.ErrInvalidAmount, create.BuyingPrice.String())
	}

	// the buying amount should be more than base creation price for domains
	coin := create.BuyingPrice.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(ctx.Domains.GetOptions().BaseDomainPrice)) {
		return false, action.ErrInvalidAmount
	}

	return true, nil
}

func (domainCreateTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runCreate(ctx, tx)
}

func (domainCreateTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runCreate(ctx, tx)
}

func (domainCreateTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runCreate(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	create := &DomainCreate{}
	err := create.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	isSub := create.Name.IsSub()

	// check domain existence and set to db
	if ctx.Domains.Exists(create.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain already exist: %s", create.Name)}
	}

	// we debit the buying price from Sender
	price := create.BuyingPrice.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(create.Owner.Bytes(), price)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, hex.EncodeToString(create.Owner)).Error()}
	}

	// add domain price to feepool
	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	opt := ctx.Domains.GetOptions()

	ok := verifyDomainName(create.Name, opt)
	if !ok {
		return false, action.Response{
			Log: ErrInvalidDomain.Error(),
		}
	}

	if len(create.Uri) > 0 {

		ok = opt.IsValidURI(create.Uri)
		if !ok {
			return false, action.Response{
				Log: "invalid uri provided",
			}
		}
	}

	var expiry int64

	if isSub {
		// if sub-domain, check if parent exists on our blockchain
		parentName, err := create.Name.GetParentName()
		if err != nil {
			return false, action.Response{Log: err.Error()}
		}
		parent, err := ctx.Domains.Get(parentName)
		//Check if Parent exists in Domain Store
		if err != nil {
			return false, action.Response{Log: "Parent domain doesn't exist, cannot create sub domain!"}
		}

		// check if sender is owner of Parent Domain
		if !bytes.Equal(parent.Owner, create.Owner) {
			return false, action.Response{Log: "parent domain not owned"}
		}

		// buying price should be more than base domain price
		if create.BuyingPrice.Value.BigInt().Cmp(opt.BaseDomainPrice.BigInt()) < 0 {
			return false, action.Response{Log: action.ErrInvalidAmount.Error()}
		}

		// set subdomain expiry height same as parent expiry height
		expiry = parent.ExpireHeight

	} else {

		// calculate expiry from the buying price
		extend, err := calculateExpiry(&create.BuyingPrice.Value, &opt.BaseDomainPrice, &opt.PerBlockFees)
		if err != nil {
			return false, action.Response{
				Log: err.Error(),
			}
		}

		// set expiry
		expiry = ctx.State.Version() + extend
	}

	domain, err := ons.NewDomain(
		create.Owner,
		create.Beneficiary,
		create.Name.String(),
		ctx.Header.Height,
		create.Uri,
		expiry,
	)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	result := action.Response{
		Tags: create.Tags(),
	}
	return true, result

}

func calculateExpiry(buyingPrice *balance.Amount, basePrice *balance.Amount, pricePerBlock *balance.Amount) (int64, error) {

	if buyingPrice.BigInt().Cmp(basePrice.BigInt()) < 0 {
		return 0, errors.New("Buying price too less")
	}

	remain := big.NewInt(0).Sub(buyingPrice.BigInt(), basePrice.BigInt())

	return big.NewInt(0).Div(remain, pricePerBlock.BigInt()).Int64(), nil
}

func calculateRenewal(buyingPrice *balance.Amount, pricePerBlock *balance.Amount) (int64, error) {

	if buyingPrice.BigInt().Cmp(pricePerBlock.BigInt()) < 0 {
		return 0, errors.New("Buying price too less")
	}

	return big.NewInt(0).Div(buyingPrice.BigInt(), pricePerBlock.BigInt()).Int64(), nil
}

func verifyDomainName(name ons.Name, feeOpt *ons.Options) bool {

	return feeOpt.IsNameAllowed(name) && name.IsValid()
}
