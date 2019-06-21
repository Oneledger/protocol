package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/serialize"
)

func (sv *Service) ONS_GetDomainByName(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReplay) error {
	domains := sv.ons
	if len(req.Name) <= 0 {
		return rpc.InvalidParamsError("name not provided")
	}

	d, err := domains.Get(req.Name, true)
	if err != nil {
		return rpc.InternalError("domain not exist")
	}
	*reply = client.ONSGetDomainsReplay{
		Domains: []ons.Domain{*d},
	}

	return nil
}

func (sv *Service) ONS_GetDomainByOwner(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReplay) error {
	domains := sv.ons
	if req.Owner == nil {
		return rpc.InvalidParamsError("owner not provided")
	}
	ds := make([]ons.Domain, 0)

	domains.Iterate(func(key []byte, value []byte) bool {
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

	*reply = client.ONSGetDomainsReplay{
		Domains: ds,
	}

	return nil
}

func (sv *Service) ONS_GetDomainOnSale(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReplay) error {
	domains := sv.ons
	if req.OnSale == false {
		return rpc.InvalidParamsError("OnSale flag not set")
	}
	ds := make([]ons.Domain, 0)

	domains.Iterate(func(key []byte, value []byte) bool {
		d := &ons.Domain{}
		err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, d)
		if err != nil {
			return true
		}
		if d.OnSaleFlag {
			ds = append(ds, *d)
		}
		return false
	})

	*reply = client.ONSGetDomainsReplay{
		Domains: ds,
	}
	return nil
}
