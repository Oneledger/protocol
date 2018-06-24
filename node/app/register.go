/*
	Copyright 2017 - 2018 OneLedger
	Handle the basic local registeration for identities and accounts
*/

package app

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Register Identities and Accounts from the user.
func RegisterLocally(app *Application, name string, scope string, chain data.ChainType,
	publicKey id.PublicKey, privateKey id.PrivateKey) bool {

	status := false

	if chain == data.UNKNOWN {
		return status
	}

	// Accounts are relative to a chain
	// TODO: Scope is tied to chain for demo purposes?
	accountName := name + "-" + scope

	account, _ := app.Accounts.FindNameOnChain(accountName, chain)
	if account == nil {
		account = id.NewAccount(chain, accountName, publicKey, privateKey)
		app.Accounts.Add(account)
		//app.Accounts.Commit()

		log.Debug("Created New Account", "key", account.AccountKey(), "account", account)

		// TODO: This should add to a list
		if name != "Zero" && chain == data.ONELEDGER {
			log.Debug("Updating NodeAccount", "name", accountName)

			global.Current.NodeAccountName = accountName
			buffer, _ := comm.Serialize(accountName)
			app.Admin.Store(data.DatabaseKey("NodeAccountName"), buffer)
			app.Admin.Commit()
		}

		status = true

	} else {
		log.Debug("Existing Account, Ignoring", "accountName", accountName)
	}

	// Fill in the balance
	if chain == data.ONELEDGER && !app.Utxo.Exists(account.AccountKey()) {
		balance := data.NewBalance(0, "OLT")
		buffer, _ := comm.Serialize(balance)

		app.Utxo.Delivered.Set(account.AccountKey(), buffer)
		app.Utxo.Commit()

		log.Debug("Adding New Utxo Entry", "key", account.AccountKey(), "balance", balance)
		status = true
	}

	// Identities are global
	identity, _ := app.Identities.FindName(name)
	if identity == nil {
		identity = id.NewIdentity(name, "Contact Info", false,
			global.Current.NodeName, account.AccountKey())
		global.Current.NodeIdentity = name
		app.Identities.Add(identity)
		//app.Identities.Commit()

		log.Debug("Registered a New Identity", "name", name, "identity", identity)
		status = true
	}

	// Associate this account with the identity
	if chain != data.ONELEDGER {
		identity.SetAccount(chain, account)
		app.Identities.Add(identity)
		//app.Identities.Commit()
	}

	return status
}