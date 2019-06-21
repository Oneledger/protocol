/*

 */

package ons

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
)

type domainData struct {
	OwnerAddress     keys.Address
	AccountAddress   keys.Address
	Name             string
	CreationHeight   int64
	LastUpdateHeight int64
	ActiveFlag       bool
	OnSaleFlag       bool

	SalePriceData serialize.Data
}

func (d *Domain) NewDataInstance() serialize.Data {
	return &domainData{}
}

func (d *Domain) Data() serialize.Data {
	return &domainData{d.OwnerAddress,
		d.AccountAddress,
		d.Name,
		d.CreationHeight,
		d.LastUpdateHeight,
		d.ActiveFlag,
		d.OnSaleFlag,
		d.SalePrice.Data(),
	}
}

func (d *Domain) SetData(a interface{}) error {
	cd, ok := a.(*domainData)
	if !ok {
		return errors.New("Wrong data")
	}

	d.OwnerAddress = cd.OwnerAddress
	d.AccountAddress = cd.AccountAddress
	d.Name = cd.Name
	d.CreationHeight = cd.CreationHeight
	d.LastUpdateHeight = cd.LastUpdateHeight
	d.ActiveFlag = cd.ActiveFlag
	d.OnSaleFlag = cd.OnSaleFlag

	err := d.SalePrice.SetData(cd.SalePriceData)
	if err != nil {
		return err
	}

	return nil
}

func (ad *domainData) SerialTag() string {
	return ""
}
