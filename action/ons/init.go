package ons

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

var (
	ErrInvalidDomain = errors.New("invalid domain name")
)

func init() {

	serialize.RegisterConcrete(new(DomainCreate), "action_dc")
	serialize.RegisterConcrete(new(CreateSubDomain), "action_csd")
	serialize.RegisterConcrete(new(DomainUpdate), "action_du")
	serialize.RegisterConcrete(new(DomainSale), "action_dsale")
	serialize.RegisterConcrete(new(DomainSend), "action_dsend")
	serialize.RegisterConcrete(new(DomainPurchase), "action_dp")
	serialize.RegisterConcrete(new(RenewDomain), "action_dr")
}

func EnableONS(r action.Router) error {
	err := r.AddHandler(action.DOMAIN_CREATE, domainCreateTx{})
	if err != nil {
		return errors.Wrap(err, "domainCreateTx")
	}
	err = r.AddHandler(action.DOMAIN_CREATE_SUB, CreateSubDomainTx{})
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
	err = r.AddHandler(action.DOMAIN_RENEW, RenewDomainTx{})
	if err != nil {
		return errors.Wrap(err, "domainCreateTx")
	}
	err = r.AddHandler(action.DOMAIN_DELETE_SUB, deleteSubTx{})
	if err != nil {
		return errors.Wrap(err, "deleteSubTx")
	}

	return nil
}

type Ons interface {
	action.Msg
	OnsName() string
}
