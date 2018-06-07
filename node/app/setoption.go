/*
	Copyright 2017 - 2018 OneLedger

	Handle setting any options for the node.
*/
package app

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Arguments for registration
type RegisterArguments struct {
	Identity   string
	Chain      string
	PublicKey  string
	PrivateKey string
}

func SetOption(app *Application, key string, value string) bool {
	log.Debug("Redirecting the option handling")

	switch key {

	case "Register":
		var arguments RegisterArguments
		result, err := comm.Deserialize([]byte(value), &arguments)
		if err != nil {
			log.Error("Can't set options", "err", err)
			return false
		}
		args := result.(*RegisterArguments)
		Register(app, args.Identity, args.Identity, id.ParseAccountType(args.Chain))

	default:
		return false
	}
	return true

}

// Register Identities and Accounts from the user.
func Register(app *Application, name string, scope string, chain data.ChainType) bool {

	status := false

	if !app.Identities.Exists(name) {
		log.Debug("Adding new Identity", "name", name)
		identity := id.NewIdentity(name, "Contact Info", false)
		app.Identities.Add(identity)
		status = true

	} else {
		log.Debug("Existing Identity", "name", name)
	}

	if chain == data.UNKNOWN {
		return status
	}

	accountName := name + "-" + scope

	if !app.Accounts.Exists(chain, accountName) {
		log.Debug("Adding new Account", "accountName", accountName)
		account := id.NewAccount(chain, accountName, id.PublicKey{})
		app.Accounts.Add(account)
		status = true

	} else {
		log.Debug("Existing Account", "accountName", accountName)
	}

	return status
}
