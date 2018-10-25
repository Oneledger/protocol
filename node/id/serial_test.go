package id

import (
	"flag"
	"os"
	"testing"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/stretchr/testify/assert"
)

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

func TestKeyType(t *testing.T) {
	var key AccountKey

	name := serial.GetBaseTypeString(key)
	log.Debug("String type", "name", name)

	entry := serial.GetTypeEntry(name, 1)
	if entry.Category != serial.UNKNOWN {
		log.Dump("AccountKey entry is", name, entry)
	} else {
		log.Fatal("Missing Type Information")
	}
}

func TestIdentity(t *testing.T) {
	var identity Identity

	// Serialize the go data structure
	buffer, err := serial.Serialize(identity, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var opp2 Identity

	// Deserialize back into a go data structure
	result, err := serial.Deserialize(buffer, opp2, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, identity, result, "These should be equal")
}

type KeyBase struct {
	Key PublicKeyED25519
}

func init() {
	serial.Register(KeyBase{})
}

func TestPublicKey(t *testing.T) {
	var key KeyBase

	// Serialize the go data structure
	buffer, err := serial.Serialize(key, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var opp2 KeyBase

	// Deserialize back into a go data structure
	result, err := serial.Deserialize(buffer, opp2, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, key, result, "These should be equal")
}

func TestAccount(t *testing.T) {
	//global.Current.RootDir = "./"
	//accounts := NewAccounts("LocalAccounts")

	chain := data.ONELEDGER
	accountName := "BaseAccount"
	publicKey := NilPublicKey()
	privateKey := NilPrivateKey()

	account := NewAccount(chain, accountName, publicKey, privateKey)

	// Serialize the go data structure
	buffer, err := serial.Serialize(account, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var opp2 Account

	// Deserialize back into a go data structure
	result, err := serial.Deserialize(buffer, opp2, serial.PERSISTENT)

	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, account, result, "These should be equal")
}
