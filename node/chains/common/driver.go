package common

import (
	"github.com/Oneledger/protocol/node/data"
)

type Contract interface {
	Chain() data.ChainType
	ToBytes() []byte
	ToKey() []byte
	FromBytes([]byte)
}
