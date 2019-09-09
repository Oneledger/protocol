package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/service/codes"
)

func (sv *Service) ONS_GetDomainByName(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
	domains := sv.ons
	if len(req.Name) <= 0 {
		return codes.ErrBadName
	}

	d, err := domains.Get(req.Name)
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

	domains.State.GetIterator().Iterate(func(key []byte, value []byte) bool {
		d := &ons.Domain{}
		err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, d)
		if err != nil {
			return true
		}
		if d.OwnerAddress.Equal(req.Owner) {
			if req.OnSale && !d.OnSaleFlag {
				return false
			}
			ds = append(ds, *d)
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

	dds := make([]ons.Domain, 0)
	domains.State.GetIterator().Iterate(func(key []byte, value []byte) bool {
		d := &ons.Domain{}
		err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, d)
		if err != nil {
			return true
		}
		if d.OnSaleFlag {
			dds = append(dds, *d)
		}
		return false
	})

	*reply = client.ONSGetDomainsOnSaleReply{
		Domains: dds,
	}
	return nil
}
