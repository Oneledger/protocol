package bid_data

import (
	"errors"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ons"
)

var _ BidAsset = &DomainAsset{}

type DomainAsset struct {
	domainName ons.Name
}

func (da *DomainAsset) ToString() string {
	return string(da.domainName)
}

func (da *DomainAsset) ValidateAsset(ctx *action.Context, owner action.Address) (bool, error) {
	// check if domain is valid
	if !da.domainName.IsValid() || da.domainName.IsSub() {
		return false, ErrInvalidDomain
	}

	// check domain existence
	if !ctx.Domains.Exists(da.domainName) {
		return false, ErrDomainNotExists
	}

	domain, err := ctx.Domains.Get(da.domainName)
	if err != nil {
		return false, ErrGettingDomain
	}

	// if domain is on sale
	if domain.OnSaleFlag {
		return false, ErrDomainOnSale
	}

	// if domain is belong to the owner
	if !domain.Owner.Equal(owner) {
		return false, ErrDomainOwnerUnmatch
	}

	// if domain is expired
	if ctx.State.Version() >= domain.ExpireHeight {
		return false, ErrDomainExpired
	}
	return true, nil
}

func (da *DomainAsset) ExchangeAsset(ctx *action.Context, bidder action.Address, preOwner action.Address) (bool, error) {
	// change domain ownership
	domain, err := ctx.Domains.Get(da.domainName)
	if err != nil {
		return false, ErrGettingDomain
	}

	if !domain.IsChangeable(ctx.Header.Height) {
		return false, ErrDomainNotChangeable
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
	asset.domainName = ons.GetNameFromString(name)
	return &asset
}
