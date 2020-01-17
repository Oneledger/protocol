package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/ons"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (sv *Service) ONS_GetDomainByName(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
	domains := sv.ons
	if len(req.Name) <= 0 {
		return codes.ErrBadName
	}

	d, err := domains.Get(ons.Name(req.Name))
	if err != nil {
		return codes.ErrDomainNotFound
	}

	ds := make([]ons.Domain, 0)

	ds = append(ds, *d)

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
	}

	return nil
}

func (sv *Service) ONS_GetDomainByOwner(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
	domains := sv.ons
	if req.Owner == nil {
		return codes.ErrBadOwner
	}
	ds := make([]ons.Domain, 0)

	domains.Iterate(func(name ons.Name, domain *ons.Domain) bool {

		if domain.OwnerAddress.Equal(req.Owner) {
			if req.OnSale && !domain.OnSaleFlag {
				return false
			}
			ds = append(ds, *domain)
		}
		return false
	})

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
	}

	return nil
}

func (sv *Service) ONS_GetDomainOnSale(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsOnSaleReply) error {
	domains := sv.ons
	if req.OnSale == false {
		return codes.ErrFlagNotSet
	}

	ds := make([]ons.Domain, 0)
	domains.Iterate(func(name ons.Name, domain *ons.Domain) bool {
		if domain.OnSaleFlag {
			ds = append(ds, *domain)
		}
		return false
	})

	*reply = client.ONSGetDomainsOnSaleReply{
		Domains: ds,
	}
	return nil
}

func (sv *Service) ONS_GetDomainByBeneficiary(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsOnSaleReply) error {
	domains := sv.ons
	if req.Beneficiary == nil {
		return codes.ErrBadAddress
	}

	ds := make([]ons.Domain, 0)
	domains.Iterate(func(name ons.Name, domain *ons.Domain) bool {
		if domain.Beneficiary.Equal(req.Beneficiary) {
			ds = append(ds, *domain)
		}
		return false
	})

	*reply = client.ONSGetDomainsOnSaleReply{
		Domains: ds,
	}
	return nil
}
