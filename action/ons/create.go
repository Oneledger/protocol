package ons

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainCreate{}

type DomainCreate struct {
	Owner       action.Address `json:"owner"`
	Beneficiary action.Address `json:"account"`
	Name        ons.Name       `json:"name"`
	Price       action.Amount  `json:"price"`
	Uri         string         `json:"uri"`
	BuyingPrice int64          `json:"buying_price"`
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

	err = action.ValidateBasic(tx.RawBytes(), create.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeeOpt, tx.Fee)
	if err != nil {
		return false, err
	}

	if create.Owner == nil || len(create.Name) <= 0 {
		return false, action.ErrMissingData
	}

	if !create.Price.IsValid(ctx.Currencies) || create.Price.Currency != "olt" {
		return false, action.ErrInvalidAmount
	}

	coin := create.Price.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromInt(CREATE_PRICE)) {
		return false, action.ErrNotEnoughFund
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

	//check domain existence and set to db
	if ctx.Domains.Exists(create.Name) {
		return false, action.Response{Log: fmt.Sprintf("domain already exist: %s", create.Name)}
	}

	price := create.Price.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(create.Owner.Bytes(), price)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, hex.EncodeToString(create.Owner)).Error()}
	}

	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	expiry, err := calculateExpiry(create.BuyingPrice, ctx.OnsOptions.BaseDomainPrice.Amount.BigInt().Int64(), ctx.OnsOptions.PerBlockFees)
	if err != nil {
		return false, action.Response{
			Log: "Unable to calculate Expiry" + err.Error(),
		}
	}

	err = verifyDomainName(create.Name, ctx.OnsOptions.FirstLevelDomain)
	if err != nil {
		return false,action.Response{
			Log: err.Error(),
		}
	}
	err = parseUrl(create.Uri)
	if err != nil {
		return false, action.Response{
			Log: err.Error(),
		}
	}

	domain, err := ons.NewDomain(
		create.Owner,
		create.Beneficiary,
		create.Name.String(),
		"",
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
	//ctx.Logger.Info("Domain :" ,domain)

	//ctx.Logger.Info(ctx.Domains.Get(domain.Name))
	result := action.Response{
		Tags: create.Tags(),
	}
	return true, result

}

func calculateExpiry(buyingPrice int64, basePrice int64, pricePerBlock int64) (int64, error) {
	if buyingPrice < basePrice {
		return 0, errors.New("Buying price too less")
	}
	return (buyingPrice - basePrice) / pricePerBlock, nil
}

func verifyDomainName(name ons.Name, firstlevelDomain []string) error {

	firstlevelMap := make(map[string]bool)
	for i := 0; i < len(firstlevelDomain); i ++ {
		firstlevelMap[firstlevelDomain[i]] = true
	}
	nameAr:= strings.Split(name.String(),".")
	if !firstlevelMap[nameAr[len(nameAr)-1]] {
		return errors.New("First level domain does is not allowed")
	} 
	return nil
}

func parseUrl(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	validProtocol := map[string]bool{"http": true, "ipfs": true, "ftp": true}
	if !validProtocol[u.Scheme] {
		return errors.New("Protocol not recognised")
	}
	return nil
}
