/*
	Copyright 2017 - 2018 OneLedger

	Easy Access to Persistent App Data, if the data isn't accessible stop immediately
*/
package id

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
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

func GetIdentities(app interface{}) *Identities {
	identities := app.(persist.Access).GetIdentities().(*Identities)
	if identities == nil {
		log.Fatal("Identity Database Missing", "config", global.Current, "app", app)
	}
	return identities
}

func GetAccounts(app interface{}) *Accounts {
	accounts := app.(persist.Access).GetAccounts().(*Accounts)
	if accounts == nil {
		log.Fatal("Accounts Database Missing", "config", global.Current, "app", app)
	}
	return accounts
}

func GetBalances(app interface{}) *data.ChainState {
	balances := app.(persist.Access).GetBalances().(*data.ChainState)
	if balances == nil {
		log.Fatal("Balances Database Missing", "config", global.Current, "app", app)
	}
	return balances
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

func GetSmartContracts(app interface{}) data.Datastore {
	smartContracts := app.(persist.Access).GetSmartContract().(data.Datastore)
	if smartContracts == nil {
		log.Fatal("SmartContract Database Missing", "config", global.Current, "app", app)
	}
	return smartContracts
}

func GetValidators(app interface{}) *Validators {
	validators := app.(persist.Access).GetValidators().(*Validators)
	//if validators == nil {
	//	log.Fatal("Validators Missing", "validators", validators, "app", app)
	//}
	return validators
}

func GetSequence(app interface{}) data.Datastore {
	sequenceDb := app.(persist.Access).GetSequence().(data.Datastore)
	if sequenceDb == nil {
		log.Fatal("Sequence Database Missing", "config", global.Current, "app", app)
	}
	return sequenceDb
}
