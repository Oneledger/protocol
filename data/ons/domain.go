/*

 */

package ons

import (
	"fmt"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

const HEIGHT_INTERVAL = 1

type Domain struct {
	// addresses of the owner and the account the domain points to
	Owner       keys.Address `json:"owner"`
	Beneficiary keys.Address `json:"beneficiary"`

	// the domain name; this is als a unique identifier of
	// the domain object over the chain
	Name Name `json:"name"`

	// block heights at which the domain was first created and updated
	CreationHeight   int64 `json:"creationHeight"`
	LastUpdateHeight int64 `json:"lastUpdateHeight"`

	// expire block height
	ExpireHeight int64 `json:"expireHeight"`

	// flag to denote whether send2Domain is active on this domain
	ActiveFlag bool `json:"activeFlag"`

	// denotes whether the domain is for sale
	OnSaleFlag bool   `json:"onSaleFlag"`
	URI        string `json:"uri"`
	// the asking price in OLT set by the owner
	SalePrice *balance.Amount `json:"salePrice"`
}

func NewDomain(ownerAddress, accountAddress keys.Address,
	name string,
	height int64,
	uri string,
	expiry int64,
) (*Domain, error) {

	if accountAddress == nil ||
		len(accountAddress) == 0 {
		accountAddress = ownerAddress
	}

	n := GetNameFromString(name)
	if !n.IsValid() {
		fmt.Println("Inavlid Name :", n.String())
		return nil, ErrDomainNameNotValid
	}

	return &Domain{
		Owner:            ownerAddress,
		Beneficiary:      accountAddress,
		Name:             n,
		CreationHeight:   height, // height of current txn
		LastUpdateHeight: height, // height of current txn
		ExpireHeight:     expiry, // height of expiry
		ActiveFlag:       true,   // Active by default

		SalePrice:  nil,
		OnSaleFlag: false,
		URI:        uri,
	}, nil
}

func (d *Domain) SetAccountAddress(addr keys.Address) {
	d.Beneficiary = addr
}

func (d *Domain) Activate() {
	d.ActiveFlag = true
}

func (d *Domain) Deactivate() {
	d.ActiveFlag = false
}

func (d *Domain) SetLastUpdatedHeight(height int64) {
	d.LastUpdateHeight = height
}

func (d *Domain) ChangeOwner(addr keys.Address) {
	d.Owner = addr
}

func (d *Domain) PutOnSale(price balance.Amount) {
	d.ActiveFlag = false
	d.OnSaleFlag = true
	d.SalePrice = &price
}

func (d *Domain) IsChangeable(currentHeight int64) bool {

	if currentHeight >= d.LastUpdateHeight+HEIGHT_INTERVAL {
		return true
	}

	return false
}

func (d *Domain) CancelSale() {
	d.OnSaleFlag = false
	d.SalePrice = nil
}

func (d *Domain) AddToExpire(h int64) {
	d.ExpireHeight = d.ExpireHeight + h
}

func (d Domain) IsActive(height int64) bool {
	return d.ActiveFlag && d.ExpireHeight > height
}

func (d Domain) IsExpired(height int64) bool {
	return d.ExpireHeight < height
}

func (d *Domain) ResetAfterSale(buyer, account keys.Address, nBlocks, currentHeight int64) {
	d.Beneficiary = account
	d.ExpireHeight = currentHeight + nBlocks
	d.Owner = buyer
	d.SalePrice = nil
	d.LastUpdateHeight = currentHeight
	d.ActiveFlag = true
	d.URI = ""
	d.OnSaleFlag = false
}
