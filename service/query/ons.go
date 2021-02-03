package query

/* 	Disabling file since domains will not be available through RPC
	TODO: Delete file in the future

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
		Height:  sv.ons.State.Version(),
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

		if domain.Owner.Equal(req.Owner) {
			if req.OnSale && !domain.OnSaleFlag {
				return false
			}
			ds = append(ds, *domain)
		}
		return false
	})

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
		Height:  sv.ons.State.Version(),
	}

	return nil
}

func (sv *Service) ONS_GetParentDomainByOwner(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
	domains := sv.ons
	if req.Owner == nil {
		return codes.ErrBadOwner
	}
	ds := make([]ons.Domain, 0)

	domains.Iterate(func(name ons.Name, domain *ons.Domain) bool {

		if domain.Owner.Equal(req.Owner) {
			if req.OnSale && !domain.OnSaleFlag {
				return false
			}
			if domain.Name.IsSub() {
				return false
			}
			ds = append(ds, *domain)
		}
		return false
	})

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
		Height:  sv.ons.State.Version(),
	}

	return nil
}

func (sv *Service) ONS_GetSubDomainByName(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
	domains := sv.ons
	if len(req.Name) <= 0 {
		return codes.ErrBadName
	}

	_, err := domains.Get(ons.Name(req.Name))
	if err != nil {
		return codes.ErrDomainNotFound
	}

	reqName := ons.GetNameFromString(req.Name)

	ds := make([]ons.Domain, 0)

	domains.Iterate(func(name ons.Name, domain *ons.Domain) bool {

		if !domain.Name.IsSub() {
			return false
		}
		if !domain.Name.EqualTo(reqName) && domain.Name.IsSubTo(reqName) {
			ds = append(ds, *domain)
		}

		return false
	})

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
		Height:  sv.ons.State.Version(),
	}

	return nil
}


func (sv *Service) ONS_GetDomainOnSale(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
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

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
		Height:  sv.ons.State.Version(),
	}
	return nil
}

func (sv *Service) ONS_GetDomainByBeneficiary(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
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

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
		Height:  sv.ons.State.Version(),
	}
	return nil
}

func (svc *Service) ONS_GetOptions(_ struct{}, reply *client.ONSGetOptionsReply) error {

	*reply = client.ONSGetOptionsReply{
		Options: *svc.ons.GetOptions(),
	}
	return nil
}
*/
