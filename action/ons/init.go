package ons

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
)

const (
	CREATE_PRICE = 100
)

var (
	ErrInvalidDomain = errors.New("invalid domain name")
)

func init() {

	serialize.RegisterConcrete(new(DomainCreate), "action_dc")
	serialize.RegisterConcrete(new(DomainUpdate), "action_du")
	serialize.RegisterConcrete(new(DomainSale), "action_dsale")
	serialize.RegisterConcrete(new(DomainSend), "action_dsend")
	serialize.RegisterConcrete(new(DomainPurchase), "action_dp")

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
	err = r.AddHandler(action.DOMAIN_SELL, domainSaleTx{})
	if err != nil {
		return errors.Wrap(err, "domainSaleTx")
	}
	err = r.AddHandler(action.DOMAIN_PURCHASE, domainPurchaseTx{})
	if err != nil {
		return errors.Wrap(err, "domainPurchaseTx")
	}
	err = r.AddHandler(action.DOMAIN_SEND, domainSendTx{})
	if err != nil {
		return errors.Wrap(err, "domainSendTx")
	}

	return nil
}

type Ons interface {
	action.Msg
	OnsName() string
}
