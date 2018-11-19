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

// TODO: NodeAccount flag should not be here!!!
// Create a local account for this fullnode
func AddAccount(app *Application, name string, chain data.ChainType,
	publicKey id.PublicKeyED25519, privateKey id.PrivateKeyED25519, nodeAccount bool) {

	account := id.NewAccount(chain, name, publicKey, privateKey)
	app.Accounts.Add(account)

	// Set this account as the current node account
	if nodeAccount {
		global.Current.NodeAccountName = name
		SetNodeName(app)
	}
}

// Broadcast an Indentity to the chain
func AddIdentity(app *Application, name string, publicKey id.PublicKeyED25519) {
	// Broadcast Identity to Chain
	LoadPrivValidatorFile()
}

// Register Identities and Accounts from the user.
func XRegisterLocally(app *Application, name string, scope string, chain data.ChainType,
	publicKey id.PublicKeyED25519, privateKey id.PrivateKeyED25519) bool {

	log.Debug("Register Locally", "name", name, "chain", chain)

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
		if name != "Zero" && name != "Zero-OneLedger" && chain == data.ONELEDGER {
			log.Debug("Updating NodeAccount", "name", accountName)

			global.Current.NodeAccountName = accountName
			SetNodeName(app)

			/*
				log.Debug("Admin store", "data.DatabaseKey", data.DatabaseKey("NodeAccountName"),
					"accountName", accountName)

				parameters := AdminParameters{NodeAccountName: accountName}
				session := app.Admin.Begin()
				session.Set(data.DatabaseKey("NodeAccountName"), parameters)
				session.Commit()
			*/
		}
		status = true
	} else {
		log.Debug("Existing Account", "accountName", accountName)
	}

	// Fill in the balance
	if chain == data.ONELEDGER && !app.Balances.Exists(account.AccountKey()) {
		balance := data.NewBalance(0, "OLT")
		app.Balances.Set(account.AccountKey(), balance)
		status = true
	}

	// Identities are global
	identity, _ := app.Identities.FindName(name)

	LoadPrivValidatorFile()

	if identity.Name == "" {
		tendermintAddress := global.Current.TendermintAddress
		tendermintPubKey := global.Current.TendermintPubKey
		if name == "Zero" || name == "Payment" {
			tendermintAddress = ""
			tendermintPubKey = ""
		}
		interim := id.NewIdentity(name, "Contact Info", false,
			global.Current.NodeName, account.AccountKey(), tendermintAddress, tendermintPubKey)
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
