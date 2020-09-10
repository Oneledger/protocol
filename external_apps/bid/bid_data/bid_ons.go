package bid_data

import (
	"errors"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ BidAsset = &DomainAsset{}

type DomainAsset struct {
	DomainName ons.Name `json:"domainName"`
}

func (da *DomainAsset) ToString() string {
	return string(da.DomainName)
}

func (da *DomainAsset) ValidateAsset(ctx *action.Context, owner action.Address) (bool, error) {
	// check if domain is valid
	if !da.DomainName.IsValid() || da.DomainName.IsSub() {
		return false, errors.New("error domain not valid")
	}

	// check domain existence
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

	// if domain is belong to the owner
	if !domain.Owner.Equal(owner) {
		return false, errors.New("domain does not belong to this address")
	}

	// if domain is expired
	if ctx.State.Version() >= domain.ExpireHeight {
		return false, errors.New("domain is expired, you can just create it")
	}
	return true, nil
}

func (da *DomainAsset) ExchangeAsset(ctx *action.Context, bidder action.Address, preOwner action.Address) (bool, error) {
	// change domain ownership
	domain, err := ctx.Domains.Get(da.DomainName)
	if err != nil {
		return false, errors.New("error getting domain")
	}

	if !domain.IsChangeable(ctx.Header.Height) {
		return false, errors.New("error domain not changeable")
	}
	domain.ResetAfterSale(bidder, bidder, 0, ctx.State.Version())
	err = ctx.Domains.DeleteAllSubdomains(domain.Name)
	if err != nil {
		return false, err
	}

	err = ctx.Domains.Set(domain)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (da *DomainAsset) NewAssetWithName(name string) BidAsset {
	asset := *da
	asset.DomainName = ons.GetNameFromString(name)
	return &asset
}
