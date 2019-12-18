/*

 */

package ons

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

const HEIGHT_INTERVAL = 1

type Domain struct {
	// addresses of the owner and the account the domain points to
	OwnerAddress keys.Address `json:"ownerAddress"`
	Beneficiary  keys.Address `json:"Beneficiary"`

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
	OnSaleFlag bool `json:"onSaleFlag"`
    URI string `json:"uri"`
	// the asking price in OLT set by the owner
	SalePrice balance.Coin `json:"salePrice"`

	// parent domain name
	Parent Name `json:"parent"`

	Expiry int64 `json:"expiry"`

}

func NewDomain(ownerAddress, accountAddress keys.Address,
	name string, parent string,
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
		return nil, ErrDomainNameNotValid
	}
	var p Name
	if len(parent) > 0 {
		p := GetNameFromString(parent)
		if !p.IsValid() {
			return nil, ErrDomainNameNotValid
		}
	}
	return &Domain{
		OwnerAddress:     ownerAddress,
		Beneficiary:      ownerAddress,
		Name:             n,
		CreationHeight:   height, // height of current txn
		LastUpdateHeight: height, // height of current txn
		ExpireHeight:     height, // height of current txn
		ActiveFlag:       true,   // Active by default

		SalePrice:  balance.Coin{},
		OnSaleFlag: false,
		Parent: p,
		URI:    uri,
		Expiry:expiry,
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
	d.OwnerAddress = addr
}

func (d *Domain) PutOnSale(price balance.Coin) {

	d.OnSaleFlag = true
	d.SalePrice = price
}

func (d *Domain) IsChangeable(currentHeight int64) bool {

	if currentHeight >= d.LastUpdateHeight+HEIGHT_INTERVAL {
		return true
	}

	return false
}

func (d *Domain) CancelSale() {
	d.OnSaleFlag = false
	d.SalePrice = balance.Coin{}
}

func (d *Domain) AddToExpire(h int64) {
	d.ExpireHeight = d.ExpireHeight + h
}

func (d Domain) IsActive(height int64) bool {
	return d.ActiveFlag && d.ExpireHeight > height
}

func (d Domain) GetParent() Name {
	return d.Parent
}

func (d *Domain) ResetAfterSale(buyer keys.Address, nBlocks, currentHeight int64) {
	d.Beneficiary = nil
	d.ExpireHeight = currentHeight + nBlocks
	d.OwnerAddress = buyer
	d.SalePrice = balance.Coin{}
	d.LastUpdateHeight = currentHeight
	d.ActiveFlag = true
	d.Parent = Name("")
	d.URI = ""
}
