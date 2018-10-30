/*
	Copyright 2017 - 2018 OneLedger
	Handle the basic local registeration for identities and accounts
*/

package app

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Register Identities and Accounts from the user.
func RegisterLocally(app *Application, name string, scope string, chain data.ChainType,
	publicKey id.PublicKeyED25519, privateKey id.PrivateKeyED25519) bool {

	status := false

	if chain == data.UNKNOWN {
		log.Warn("Can't register, chain UNKNOWN", "name", name)
		return status
	}

	// Accounts are relative to a chain
	// TODO: Scope is tied to chain for demo purposes?
	accountName := name + "-" + scope
	account, _ := app.Accounts.FindNameOnChain(accountName, chain)

	//var account id.Account = nil
	if account == nil {
		account = id.NewAccount(chain, accountName, publicKey, privateKey)
		app.Accounts.Add(account)

		log.Debug("Created New Account", "key", account.AccountKey(), "account", account)

		// TODO: This should add to a list
		if name != "Zero" && chain == data.ONELEDGER {
			log.Debug("Updating NodeAccount", "name", accountName)

			global.Current.NodeAccountName = accountName
			/*
				buffer, err := serial.Serialize(accountName, serial.NETWORK)
				if err != nil {
					log.Error("Failed to Serialize accountName")
				}
			*/
			log.Debug("Admin store", "data.DatabaseKey", data.DatabaseKey("NodeAccountName"),
				"accountName", accountName)

			session := app.Admin.Begin()
			session.Set(data.DatabaseKey("NodeAccountName"), accountName)
			session.Commit()
		}
		status = true
	} else {
		log.Debug("Existing Account", "accountName", accountName)
	}

	// Fill in the balance
	if chain == data.ONELEDGER && !app.Utxo.Exists(account.AccountKey()) {
		balance := data.NewBalance(0, "OLT")
		app.Utxo.Set(account.AccountKey(), balance)
		status = true
	}

	// Identities are global
	identity, _ := app.Identities.FindName(name)
	if identity.Name == "" {
		interim := id.NewIdentity(name, "Contact Info", false,
			global.Current.NodeName, account.AccountKey())
		identity = *interim

		global.Current.NodeIdentity = name
		app.Identities.Add(identity)

		log.Debug("Registered a New Identity", "name", name, "identity", identity)
		status = true
	}

	// Associate this account with the identity
	if chain != data.ONELEDGER {
		identity.SetAccount(chain, account)
		app.Identities.Add(identity)
	}
	return status
}
