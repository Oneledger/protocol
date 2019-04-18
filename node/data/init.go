package data

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/serialize"
)

func init() {
	serial.Register(Balance{})

	serialize.RegisterConcrete(new(Balance), TagBalance)
}
