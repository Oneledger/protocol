/*
	Copyright 2017 - 2018 OneLedger

	Handle setting any options for the node.
*/
package app

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Arguments for registration
type RegisterArguments struct {
	Identity   string
	Chain      string
	PublicKey  string
	PrivateKey string
}

func HandleSet(app Application, path string, arguments map[string]string) []byte {
	log.Dump("Have Set", path, arguments)
	var result interface{}

	switch path {
	case "/account":
		result = HandleSetAccount(app, arguments)

	case "/register":
		result = HandleRegisterIdentity(app, arguments)
	}

	if result == nil {
		return nil
	}

	buffer, err := serial.Serialize(result, serial.CLIENT)
	if err != nil {
		log.Fatal("Failed to serialize query", "err", err)
	}

	return buffer
}

func GetChain(chainName string) data.ChainType {
	return data.ONELEDGER
}

// TODO: The datatype for Key, depends on Chain
func GetKeys(chain data.ChainType, name string, publicKey string, privateKey string) (id.PublicKeyED25519, id.PrivateKeyED25519) {
	//return id.NilPublicKey(), id.NilPrivateKey()

	// TODO: Need to push the passphrase back through the CLI
	priv, public := id.GenerateKeys([]byte(name + "as password"))
	return public, priv
}

// TODO: Pass in App pointer?
func HandleSetAccount(app Application, arguments map[string]string) interface{} {
	chain := GetChain(arguments["Chain"])
	accountName := arguments["Account"]

	publicKey, privateKey := GetKeys(chain, accountName, arguments["PublicKey"], arguments["PrivateKey"])

	log.Debug("#### Adding Accounts", "chain", chain, "name", accountName)

	AddAccount(&app, accountName, chain, publicKey, privateKey, false)

	account, err := app.Accounts.FindName(accountName)
	if err == status.SUCCESS {
		return account
	}
	return "Error in Setting up Account"
}

func HandleRegisterIdentity(app Application, arguments map[string]string) interface{} {
	return "Registering Identity"

	// TODO: Broadcast the transaction
}

// Handle a SetOption ABCi reqeust
func SetOption(app *Application, key string, value string) bool {
	log.Debug("Setting Application Options", "key", key, "value", value)

	switch key {

	case "Register":
		var arguments RegisterArguments
		result, err := serial.Deserialize([]byte(value), &arguments, serial.NETWORK)
		if err != nil {
			log.Error("Can't set options", "status", err)
			return false
		}
		args := result.(*RegisterArguments)
		privateKey, publicKey := id.GenerateKeys([]byte(args.Identity)) // TODO: Switch with passphrase
		RegisterLocally(app, args.Identity, "OneLedger", id.ParseAccountType(args.Chain),
			publicKey, privateKey)

	default:
		log.Warn("Unknown Option", "key", key)
		return false
	}
	return true
}
