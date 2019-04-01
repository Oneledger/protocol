package serialize

import (
	"github.com/tendermint/go-amino"
)

type Channel int

const (
	CLIENT Channel = iota
	PERSISTENT
	NETWORK
	JSON
)

var aminoCodec *amino.Codec
var JSONSzr Serializer

func init() {
	aminoCodec = amino.NewCodec()

	JSONSzr = GetSerializer(JSON)
}

type Serializer interface {
	Serialize(obj interface{}) ([]byte, error)
	Deserialize(d []byte, obj interface{}) error
}

// GetSerializer for a channel of standard types, default is a JSON serializer
func GetSerializer(channel Channel, args ...interface{}) Serializer {

	switch channel {

	case CLIENT:
		return &msgpackStrategy{}

	case PERSISTENT:
		return &msgpackStrategy{}

	case NETWORK:
		return &msgpackStrategy{}

	case JSON:
		return &jsonStrategy{}

	default:
		return &jsonStrategy{}
	}
}


// functions to register types
func RegisterInterface(obj interface{}) {
	aminoCodec.RegisterInterface(obj, &amino.InterfaceOptions{AlwaysDisambiguate: true})
}

func RegisterConcrete(obj interface{}, name string) {
	aminoCodec.RegisterConcrete(obj, name, nil)
}
