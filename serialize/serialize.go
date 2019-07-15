package serialize

import (
	"os"
	"sync"

	"github.com/Oneledger/protocol/log"

	"github.com/google/uuid"

	"github.com/tendermint/go-amino"
	"github.com/vmihailenco/msgpack"
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
var registeredConcretes = make([]string, 0, 100)
var lockRegisteredConcretes sync.Mutex
var logger *log.Logger

func init() {
	aminoCodec = amino.NewCodec()
	JSONSzr = GetSerializer(JSON)
	logger = log.NewLoggerWithPrefix(os.Stdout, "serialize")

	RegisterConcrete(new(string), "std_string")
	RegisterConcrete(new([]byte), "std_byte_arr")
	RegisterConcrete(new(uuid.UUID), "google_uuid_UUID")
}

type Serializer interface {
	Serialize(obj interface{}) ([]byte, error)
	Deserialize(d []byte, obj interface{}) error
}

// GetSerializer for a channel of standard types, default is a JSON serializer
func GetSerializer(channel Channel, args ...interface{}) Serializer {

	switch channel {

	case CLIENT:
		return &jsonStrategy{}

	case PERSISTENT:
		return &jsonStrategy{}

	case NETWORK:
		return &jsonStrategy{}

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
