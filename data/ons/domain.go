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
	OwnerAddress   keys.Address `json:"ownerAddress"`
	AccountAddress keys.Address `json:"accountAddress"`

	// the domain name; this is als a unique identifier of
	// the domain object over the chain
	Name string `json:"name"`

	// block heights at which the domain was first created and updated
	CreationHeight   int64 `json:"creationHeight"`
	LastUpdateHeight int64 `json:"lastUpdateHeight"`

	// flag to denote whether send2Domain is active on this domain
	ActiveFlag bool `json:"activeFlag"`

	// denotes whether the domain is for sale
	OnSaleFlag bool `json:"onSaleFlag"`

	// the asking price in OLT set by the owner
	SalePrice balance.Coin `json:"salePrice"`
}

func NewDomain(ownerAddress, accountAddress keys.Address,
	name string, height int64) *Domain {

	if accountAddress == nil ||
		len(accountAddress) == 0 {
		accountAddress = ownerAddress
	}

	return &Domain{
		OwnerAddress:     ownerAddress,
		AccountAddress:   accountAddress,
		Name:             name,
		CreationHeight:   height, // height of current txn
		LastUpdateHeight: height, // height of current txn
		ActiveFlag:       true,   // Active by default

		SalePrice:  balance.Coin{},
		OnSaleFlag: false,
	}
}

func (d *Domain) SetAccountAddress(addr keys.Address) {
	d.AccountAddress = addr
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
