package shared

import "github.com/Oneledger/protocol/node/serialize"

var clSerializer serialize.Serializer

func init() {
	clSerializer = serialize.GetSerializer(serialize.CLIENT)
}
