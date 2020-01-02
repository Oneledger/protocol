package ons

import (
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

type DomainCreate struct {
	Owner       action.Address `json:"owner"`
	Beneficiary action.Address `json:"account"`
	Name        ons.Name       `json:"name"`
	Uri         string         `json:"uri"`
	BuyingPrice action.Amount  `json:"buyingprice"`
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
	c, _ := ctx.Currencies.GetCurrencyById(0)
	if c.Name != create.BuyingPrice.Currency {
		fmt.Println(c, create)
		return false, errors.Wrap(action.ErrInvalidAmount, create.BuyingPrice.String())
	}

	coin := create.BuyingPrice.ToCoin(ctx.Currencies)
	if coin.LessThanEqualCoin(coin.Currency.NewCoinFromAmount(ctx.Domains.GetOptions().BaseDomainPrice)) {
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

	price := create.BuyingPrice.ToCoin(ctx.Currencies)
	err = ctx.Balances.MinusFromAddress(create.Owner.Bytes(), price)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, hex.EncodeToString(create.Owner)).Error()}
	}

	err = ctx.FeePool.AddToPool(price)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	opt := ctx.Domains.GetOptions()

	expiry, err := calculateExpiry(&create.BuyingPrice.Value, &opt.BaseDomainPrice, &opt.PerBlockFees)
	if err != nil {
		return false, action.Response{
			Log: "Unable to calculate Expiry" + err.Error(),
		}
	}

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

	domain, err := ons.NewDomain(
		create.Owner,
		create.Beneficiary,
		create.Name.String(),
		"",
		ctx.Header.Height,
		create.Uri,
		ctx.State.Version()+expiry,
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

func calculateExpiry(buyingPrice *balance.Amount, basePrice *balance.Amount, pricePerBlock *balance.Amount) (int64, error) {
	if buyingPrice.BigInt().Cmp(basePrice.BigInt()) < 0 {
		return 0, errors.New("Buying price too less")
	}
	remain := big.NewInt(0).Sub(buyingPrice.BigInt(), basePrice.BigInt())
	return big.NewInt(0).Div(remain, pricePerBlock.BigInt()).Int64(), nil
}

func verifyDomainName(name ons.Name, feeOpt *ons.Options) bool {
	return feeOpt.IsNameAllowed(name) && name.IsValid()
}
