/*

 */

package ons

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type domainData struct {
	Owner            keys.Address `json:"a"`
	Beneficiary      keys.Address `json:"b"`
	Name             string       `json:"c"`
	CreationHeight   int64        `json:"d"`
	LastUpdateHeight int64        `json:"e"`
	ExpireHeight     int64        `json:"f"`
	ActiveFlag       bool         `json:"g"`
	OnSaleFlag       bool         `json:"h"`
	SalePriceData    []byte       `json:"i"`
	URI              string       `json:"k"`
}

func (d *Domain) NewDataInstance() serialize.Data {
	return &domainData{}
}

func (d *Domain) Data() serialize.Data {
	dd := &domainData{
		Owner:            d.Owner,
		Beneficiary:      d.Beneficiary,
		Name:             d.Name.String(),
		CreationHeight:   d.CreationHeight,
		LastUpdateHeight: d.LastUpdateHeight,
		ExpireHeight:     d.ExpireHeight,
		ActiveFlag:       d.ActiveFlag,
		OnSaleFlag:       d.OnSaleFlag,
		SalePriceData:    nil,
		URI:              d.URI,
	}
	if d.SalePrice != nil {
		dd.SalePriceData, _ = d.SalePrice.MarshalJSON()
	}
	return dd
}

func (d *Domain) SetData(a interface{}) error {
	cd, ok := a.(*domainData)
	if !ok {
		return errors.New("Wrong data type for domain")
	}

	d.Owner = cd.Owner
	d.Beneficiary = cd.Beneficiary
	d.Name = GetNameFromString(cd.Name)
	d.CreationHeight = cd.CreationHeight
	d.LastUpdateHeight = cd.LastUpdateHeight
	d.ExpireHeight = cd.ExpireHeight
	d.ActiveFlag = cd.ActiveFlag
	d.OnSaleFlag = cd.OnSaleFlag
	d.URI = cd.URI

	if cd.SalePriceData != nil {
		amt := &balance.Amount{}
		err := amt.UnmarshalJSON(cd.SalePriceData)
		if err != nil {
			return err
		}
		d.SalePrice = amt
	}
	return nil
}

func (ad *domainData) SerialTag() string {
	return ""
}
