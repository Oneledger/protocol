/*
	Copyright 2017 - 2018 OneLedger

	Handle setting any options for the node.
*/
package app

import (
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

// Arguments for registration
type RegisterArguments struct {
	Identity   string
	Chain      string
	PublicKey  string
	PrivateKey string
}

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
		privateKey, publicKey := id.GenerateKeys([]byte(args.Identity), true) // TODO: Switch with passphrase
		RegisterLocally(app, args.Identity, "OneLedger", id.ParseAccountType(args.Chain),
			publicKey, privateKey)

	default:
		log.Warn("Unknown Option", "key", key)
		return false
	}
	return true
}
