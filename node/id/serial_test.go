package id

import (
	"github.com/Oneledger/protocol/node/serialize"
	"testing"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/stretchr/testify/assert"
)

var pSzlr serialize.Serializer

func init() {
	pSzlr = serialize.GetSerializer(serialize.PERSISTENT)
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
	buffer, err := pSzlr.Serialize(identity)
	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var result = &Identity{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, identity, result, "These should be equal")
}

func TestIdentities(t *testing.T) {
	global.Current.RootDir = "./"
	identities := NewIdentities("TestIdentities")

	identity := Identity{
		Name:     "Tester",
		NodeName: "Here",
	}

	identities.Add(identity)

	result, _ := identities.FindName(identity.Name)

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
	buffer, err := pSzlr.Serialize(key)
	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var result = &KeyBase{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, result)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, key, result, "These should be equal")
}

func TestAccount(t *testing.T) {
	//global.Current.RootDir = "./"
	//accounts := NewAccounts("LocalAccounts")

	chain := data.ONELEDGER
	accountName := "Zero"
	publicKey := NilPublicKey()
	privateKey := NilPrivateKey()

	account := NewAccount(chain, accountName, publicKey, privateKey, nil)

	// Serialize the go data structure
	buffer, err := pSzlr.Serialize(account)
	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var result interface{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, &result)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}

	assert.Equal(t, account, result, "These should be equal")
}

func TestAccountArray(t *testing.T) {
	//global.Current.RootDir = "./"
	//accounts := NewAccounts("LocalAccounts")

	chain := data.ONELEDGER
	accountName := "Zero"
	publicKey := NilPublicKey()
	privateKey := NilPrivateKey()

	accounts := make([]Account, 3)
	accounts[0] = NewAccount(chain, accountName, publicKey, privateKey, nil)
	accounts[1] = NewAccount(chain, accountName, publicKey, privateKey, nil)
	accounts[2] = NewAccount(chain, accountName, publicKey, privateKey, nil)

	// Serialize the go data structure
	buffer, err := pSzlr.Serialize(accounts)
	if err != nil {
		log.Fatal("Serialized failed", "err", err)
	}

	var result interface{}

	// Deserialize back into a go data structure
	err = pSzlr.Deserialize(buffer, &result)
	if err != nil {
		log.Fatal("Deserialized failed", "err", err)
	}
	results := result.([]interface{})
	for i := range accounts {
		assert.Equal(t, accounts[i], results[i].(Account), "These should be equal")
	}
}

func TestAccountPersistence(t *testing.T) {
	global.Current.RootDir = "./"
	accounts := NewAccounts("LocalAccounts")

	chain := data.ONELEDGER
	accountName := "Zero"
	publicKey := NilPublicKey()
	privateKey := NilPrivateKey()

	account := NewAccount(chain, accountName, publicKey, privateKey)

	accounts.Add(account)

	result, status := accounts.Find(account)
	if status != err.SUCCESS {
		log.Fatal("Account Datastore Failed", "status", status)
	}

	assert.Equal(t, account, result, "These should be equal")
}
