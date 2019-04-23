package action

import (
	"flag"
	"github.com/Oneledger/protocol/node/serialize"
	"os"
	"testing"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

var pSzlr serialize.Serializer

func init() {
	pSzlr = serialize.GetSerializer(serialize.PERSISTENT)
}

// Control the execution
func TestMain(m *testing.M) {
	flag.Parse()

	// Set the debug flags according to whether the -v flag is set in go test
	if testing.Verbose() {
		log.Debug("DEBUG TURNED ON")
		global.Current.Debug = true
	} else {
		log.Debug("DEBUG TURNED OFF")
		global.Current.Debug = false
	}

	// Run it all.
	code := m.Run()

	os.Exit(code)
}

func TestType(t *testing.T) {
	var value Type

	name := serial.GetBaseTypeString(value)
	log.Debug("String type", "name", name)

	entry := serial.GetTypeEntry(name, 1)
	if entry.Category != serial.UNKNOWN {
		log.Dump("Type entry is", name, entry)
	} else {
		log.Fatal("Missing Type Information")
	}
}

type Wrapper struct {
	Type Type
}

func init() {
	serial.Register(Wrapper{})
}

func TestTypeSerial(t *testing.T) {
	value := Wrapper{SEND}
	buffer, err := serial.Serialize(value, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Have Error", "err", err)
	}
	var proto Wrapper
	result, err := serial.Deserialize(buffer, proto, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Have Error", "err", err)
	}
	log.Dump("Final", result)
}

func TestIdentity(t *testing.T) {
	chain := map[data.ChainType]id.AccountKey{
		data.ONELEDGER: []byte("xxxx"),
	}
	value := id.Identity{
		Chain: chain,
	}

	buffer, err := pSzlr.Serialize(value)
	if err != nil {
		log.Fatal("Have Error", "err", err)
	}
	var result = &id.Identity{}
	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		log.Fatal("Have Error", "err", err)
	}
	log.Dump("Final", result)
}
