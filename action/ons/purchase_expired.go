/*

 */

package ons

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ Ons = &DomainExpiredPurchase{}

/*
	DomainExpiredPurchase

*/
/*

	This is the message data for running a domain-expired purchase transaction.
*/
type DomainExpiredPurchase struct {
	// Name of the expired domain to purchase
	Name ons.Name `json:"name"`

	// The buyer address, who is also the signer.
	Buyer action.Address `json:"buyer"`

	// The number of blocks for which to buy this domain.
	Blocks int64
}

func (dp DomainExpiredPurchase) Marshal() ([]byte, error) {
	return json.Marshal(dp)
}

func (dp *DomainExpiredPurchase) Unmarshal(data []byte) error {
	return json.Unmarshal(data, dp)
}

func (dp DomainExpiredPurchase) OnsName() string {
	return dp.Name.String()
}

func (dp DomainExpiredPurchase) Signers() []action.Address {
	return []action.Address{dp.Buyer}
}

func (dp DomainExpiredPurchase) Type() action.Type {
	return action.DOMAIN_EXPIRED_PURCHASE
}

func (dp DomainExpiredPurchase) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)
	tag0 := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(dp.Type().String()),
	}
	tag1 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: dp.Buyer,
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.domain_name"),
		Value: []byte(dp.Name),
	}

	tags = append(tags, tag0, tag1, tag2)
	return tags
}

/*
	domainExpiredPurchaseTx
*/
type domainExpiredPurchaseTx struct {
}

func (domainExpiredPurchaseTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	buy := &DomainExpiredPurchase{}
	err := buy.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), buy.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeeOpt, tx.Fee)
	if err != nil {
		return false, err
	}

	if buy.Blocks <= 0 {
		return false, errors.New("invalid number of blocks for domain purchase")
	}

	// validate the sender and receiver are not nil
	if buy.Buyer == nil || len(buy.Name) == 0 {
		return false, action.ErrMissingData
	}

	return true, nil
}

func (domainExpiredPurchaseTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return processDomainExpiredPurchase(ctx, tx)
}

func (domainExpiredPurchaseTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return processDomainExpiredPurchase(ctx, tx)
}

func (domainExpiredPurchaseTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func processDomainExpiredPurchase(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	buy := &DomainExpiredPurchase{}
	err := buy.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "error unmarshalling tx data: " + "DOMAIN_EXPIRED_PURCHASE"}
	}

	domain, err := ctx.Domains.Get(buy.Name)
	if err != nil {
		return false, action.Response{Log: "cannot find domain: " + buy.Name.String()}
	}

	if ctx.State.Version() <= domain.ExpireHeight {
		return false, action.Response{Log: "domain not expired yet: " + buy.Name.String()}
	}

	// TODO get per block charge from governance db
	olt, _ := ctx.Currencies.GetCurrencyByName("OLT")
	blockCharges := olt.NewCoinFromUnit(buy.Blocks * 10000000)

	err = ctx.Balances.MinusFromAddress(buy.Buyer, blockCharges)
	if err != nil {
		return false, action.Response{Log: "error deducting balance for purchase expired domain: " + err.Error()}
	}

	err = ctx.FeePool.AddToPool(blockCharges)
	if err != nil {
		return false, action.Response{Log: "error adding domain purchase: " + err.Error()}
	}

	domain.ResetAfterSale(buy.Buyer, buy.Blocks, ctx.State.Version())

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, action.Response{Log: "error saving domain: " + err.Error()}
	}

	err = ctx.Domains.DeleteSubdomains(domain.Name)
	if err != nil {
		return false, action.Response{Log: "error deleting subdomains in domain purchase" + err.Error()}
	}

	return true, action.Response{Tags: buy.Tags()}

}
