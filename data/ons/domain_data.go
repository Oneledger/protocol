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
	OwnerAddress     keys.Address `json:"a"`
	Beneficiary      keys.Address `json:"b"`
	Name             string       `json:"c"`
	CreationHeight   int64        `json:"d"`
	LastUpdateHeight int64        `json:"e"`
	ExpireHeight     int64        `json:"f"`
	ActiveFlag       bool         `json:"g"`
	OnSaleFlag       bool         `json:"h"`

	SalePriceData *balance.CoinData `json:"i"`
	Parent        string            `json:"j"`
	URI           string            `json:"k"`
}

func (d *Domain) NewDataInstance() serialize.Data {
	return &domainData{}
}

func (d *Domain) Data() serialize.Data {
	dd := &domainData{
		OwnerAddress:     d.OwnerAddress,
		Beneficiary:      d.Beneficiary,
		Name:             d.Name.String(),
		CreationHeight:   d.CreationHeight,
		LastUpdateHeight: d.LastUpdateHeight,
		ExpireHeight:     d.ExpireHeight,
		ActiveFlag:       d.ActiveFlag,
		OnSaleFlag:       d.OnSaleFlag,
		SalePriceData:    nil,
		Parent:           d.Parent.String(),
		URI:              d.URI,
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
	d.Beneficiary = cd.Beneficiary
	d.Name = GetNameFromString(cd.Name)
	d.CreationHeight = cd.CreationHeight
	d.LastUpdateHeight = cd.LastUpdateHeight
	d.ExpireHeight = cd.ExpireHeight
	d.ActiveFlag = cd.ActiveFlag
	d.OnSaleFlag = cd.OnSaleFlag
	d.Parent = GetNameFromString(cd.Parent)
	d.URI = cd.URI

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
