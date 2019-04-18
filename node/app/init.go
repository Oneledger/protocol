package app

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/serialize"
)

var clSerializer serialize.Serializer

func init() {
	// TODO: Should be driven from config
	ChainId = "OneLedger-Root"

	clSerializer = serialize.GetSerializer(serialize.CLIENT)

	serial.Register(AdminParameters{})
	serialize.RegisterConcrete(new(AdminParameters), TagAdminParameters)
}
