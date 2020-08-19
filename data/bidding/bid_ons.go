package bidding

import (
	"errors"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

type DomainAsset struct {
	DomainName ons.Name `json:"name"`
}

func NewDomainAsset(name string) *DomainAsset {
	return &DomainAsset{
		DomainName: ons.Name(name),
	}
}

func (da *DomainAsset) ToString() string {
	return string(da.DomainName)
}

func (da *DomainAsset) ValidateAsset(ctx *action.Context) (bool, error) {
	// check domain existence and set to db
	if !ctx.Domains.Exists(da.DomainName) {
		return false, errors.New("domain does not exist, you can just create it")
	}

	domain, err := ctx.Domains.Get(da.DomainName)
	if err != nil {
		return false, errors.New("error getting domain")
	}

	// if domain is on sale
	if domain.OnSaleFlag {
		return false, errors.New("domain is on sale")
	}

	// if domain is expired
	if ctx.State.Version() >= domain.ExpireHeight {
		return false, errors.New("domain is expired, you can just create it")
	}
	return true, nil
}
