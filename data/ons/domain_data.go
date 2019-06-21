/*

 */

package ons

import (
	"github.com/Oneledger/protocol/data/balance"
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

	SalePriceData *balance.CoinData
}

func (d *Domain) NewDataInstance() serialize.Data {
	return &domainData{}
}

func (d *Domain) Data() serialize.Data {
	dd := &domainData{OwnerAddress: d.OwnerAddress,
		AccountAddress:   d.AccountAddress,
		Name:             d.Name,
		CreationHeight:   d.CreationHeight,
		LastUpdateHeight: d.LastUpdateHeight,
		ActiveFlag:       d.ActiveFlag,
		OnSaleFlag:       d.OnSaleFlag,
	}
	if d.SalePrice.Amount != nil {
		dd.SalePriceData = d.SalePrice.Data().(*balance.CoinData)
	}

	return dd
}

func (d *Domain) SetData(a interface{}) error {
	cd, ok := a.(*domainData)
	if !ok {
		return errors.New("Wrong data type for domain")
	}

	d.OwnerAddress = cd.OwnerAddress
	d.AccountAddress = cd.AccountAddress
	d.Name = cd.Name
	d.CreationHeight = cd.CreationHeight
	d.LastUpdateHeight = cd.LastUpdateHeight
	d.ActiveFlag = cd.ActiveFlag
	d.OnSaleFlag = cd.OnSaleFlag

	if cd.SalePriceData != nil {
		err := d.SalePrice.SetData(cd.SalePriceData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ad *domainData) SerialTag() string {
	return ""
}
