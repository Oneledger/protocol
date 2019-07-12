package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/rpc"
	"github.com/Oneledger/protocol/serialize"
)

func (sv *Service) ONS_GetDomainByName(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
	domains := sv.ons
	if len(req.Name) <= 0 {
		return rpc.InvalidParamsError("name not provided")
	}

	d, err := domains.Get(req.Name, true)
	if err != nil {
		return rpc.InternalError("domain not exist")
	}

	ds := make([]client.DomainData, 0)
	dd := &client.DomainData{
		d.Name,
		d.SalePrice.Humanize(),
		d.OwnerAddress,
		d.AccountAddress,
		d.CreationHeight,
		d.LastUpdateHeight,
		d.ActiveFlag,
		d.OnSaleFlag,
	}
	ds = append(ds, *dd)

	*reply = client.ONSGetDomainsReply{
		Domains: ds,
	}

	return nil
}

func (sv *Service) ONS_GetDomainByOwner(req client.ONSGetDomainsRequest, reply *client.ONSGetDomainsReply) error {
	domains := sv.ons
	if req.Owner == nil {
		return rpc.InvalidParamsError("owner not provided")
	}
	ds := make([]client.DomainData, 0)

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
			d := &client.DomainData{
				d.Name,
				d.SalePrice.Humanize(),
				d.OwnerAddress,
				d.AccountAddress,
				d.CreationHeight,
				d.LastUpdateHeight,
				d.ActiveFlag,
				d.OnSaleFlag,
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
		return rpc.InvalidParamsError("OnSale flag not set")
	}

	dds := []client.DomainData{}
	domains.Iterate(func(key []byte, value []byte) bool {
		d := &ons.Domain{}
		err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, d)
		if err != nil {
			return true
		}
		if d.OnSaleFlag {
			dd := &client.DomainData{
				d.Name,
				d.SalePrice.Humanize(),
				d.OwnerAddress,
				d.AccountAddress,
				d.CreationHeight,
				d.LastUpdateHeight,
				d.ActiveFlag,
				d.OnSaleFlag,
			}
			dds = append(dds, *dd)
		}
		return false
	})

	*reply = client.ONSGetDomainsOnSaleReply{
		Domains: dds,
	}
	return nil
}
