package bidding

import "github.com/Oneledger/protocol/data/ons"

type DomainAsset struct {
	Name ons.Name `json:"name"`
}

func NewDomainAsset(name string) *DomainAsset {
	return &DomainAsset{
		Name: ons.Name(name),
	}
}

func (da *DomainAsset) ToString() string {
	return string(da.Name)
}
