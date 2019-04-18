package serialize

import (
	"github.com/tendermint/go-amino"
	"github.com/vmihailenco/msgpack"
	"sync"
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
var registeredConcretes = []string{}
var lockRegisteredConcretes sync.Mutex

func init() {
	aminoCodec = amino.NewCodec()
	JSONSzr = GetSerializer(JSON)

	RegisterConcrete(new(string), "std_string")
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

	msgpackRegConc(obj, name)
	aminoCodec.RegisterConcrete(obj, name, nil)

}

func msgpackRegConc(obj interface{}, name string) {
	lockRegisteredConcretes.Lock()

	registeredConcretes = append(registeredConcretes, name)
	lRegisteredConcretes := len(registeredConcretes)
	if lRegisteredConcretes > 128 {
		panic("can't lock more than 128 struct types for serialization")
	}
	msgpack.RegisterExt(int8(lRegisteredConcretes-1), obj)

	lockRegisteredConcretes.Unlock()
}