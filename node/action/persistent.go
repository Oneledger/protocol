/*
	Copyright 2017 - 2018 OneLedger

	Easy Access to Persistent App Data, if the persistent data isn't avoid stop immediately
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/persist"
)

func GetAdmin(app interface{}) data.Datastore {
	admin := app.(persist.Access).GetAdmin().(data.Datastore)
	if admin == nil {
		log.Fatal("Admin Database Missing", "config", global.Current, "app", app)
	}
	return admin
}

func GetStatus(app interface{}) data.Datastore {
	status := app.(persist.Access).GetStatus()
	result := status.(data.Datastore)
	if status == nil {
		log.Fatal("Status Database Missing", "config", global.Current, "app", app)
	}
	return result
}

func GetIdentities(app interface{}) *id.Identities {
	identities := app.(persist.Access).GetIdentities().(*id.Identities)
	if identities == nil {
		log.Fatal("Identity Database Missing", "config", global.Current, "app", app)
	}
	return identities
}

func GetAccounts(app interface{}) *id.Accounts {
	accounts := app.(persist.Access).GetAccounts().(*id.Accounts)
	if accounts == nil {
		log.Fatal("Account Database Missing", "config", global.Current, "app", app)
	}
	return accounts
}

func GetUtxo(app interface{}) *data.ChainState {
	chain := app.(persist.Access).GetUtxo().(*data.ChainState)
	return chain
}

func GetChainID(app interface{}) string {
	id := app.(persist.Access).GetChainID().(string)
	return id
}

func GetEvent(app interface{}) data.Datastore {
	event := app.(persist.Access).GetEvent().(data.Datastore)
	if event == nil {
		log.Fatal("Event Database Missing", "config", global.Current, "app", app)
	}
	return event
}

func GetContracts(app interface{}) data.Datastore {
	htlcs := app.(persist.Access).GetContract().(data.Datastore)
	if htlcs == nil {
		log.Fatal("Htlc Database Missing", "config", global.Current, "app", app)
	}
	return htlcs
}
