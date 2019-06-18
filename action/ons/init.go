package ons

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
)

const (
	CREATE_PRICE = 100
)

func init() {

	serialize.RegisterConcrete(new(DomainCreate), "action_dc")
	serialize.RegisterConcrete(new(DomainUpdate), "action_du")
}

func EnableONS(r action.Router) error {
	err := r.AddHandler(action.DOMAIN_CREATE, domainCreateTx{})
	if err != nil {
		return errors.Wrap(err, "domainCreateTx")
	}
	err = r.AddHandler(action.DOMAIN_UPDATE, domainUpdateTx{})
	if err != nil {
		return errors.Wrap(err, "domainUpdateTx")
	}

	return nil
}
