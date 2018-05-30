/*
	Copyright 2017 - 2018 OneLedger

	Handle setting any options for the node.
*/
package app

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

type RegisterArguments struct {
	Name       string
	Chain      string
	PublicKey  string
	PrivateKey string
}

func SetOption(app *Application, key string, value []byte) bool {
	log.Debug("Redirecting the option handling")
	switch key {
	case "Register":
		var arguments RegisterArguments
		result, err := comm.Deserialize(value, &arguments)
		if err != nil {
			log.Error("Can't set options", "err", err)
			return false
		}
		args := result.(*RegisterArguments)
		Register(app, args.Name, args.Name, id.ParseAccountType(args.Chain))

	default:
		return false
	}
	return true

}

// Register Identities and Accounts from the user.
func Register(app *Application, idName string, name string, chain id.AccountType) bool {

	status := false

	if !app.Identities.Exists(idName) {
		log.Debug("Adding new Identity", "idName", idName)
		identity := id.NewIdentity(idName, "Contact Info")
		app.Identities.Add(identity)
		status = true
	} else {
		log.Debug("Existing Identity", "idName", idName)
	}

	if chain == id.UNKNOWN {
		return status
	}

	accountName := idName + "-" + name
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
