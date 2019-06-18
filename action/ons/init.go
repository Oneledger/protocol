package ons

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {

	serialize.RegisterConcrete(new(DomainCreate), "action_dc")

}

func EnableONS(r action.Router) error {

}
