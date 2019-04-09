/*
	Copyright 2017 - 2018 OneLedger

	Easy Access to Persistent App Data, if the data isn't accessible stop immediately
*/
package action

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/persist"
)

func AnalyzeScript(app interface{}, request *OLVMRequest) interface{} {
	result := app.(persist.Access).AnalyzeScript(request)
	return result
}

func RunScript(app interface{}, request *OLVMRequest) (interface{}, error) {
	return app.(persist.Access).RunScript(request)
}

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

func GetRPCClient(app interface{}) comm.ClientInterface {
	rpcclient := app.(persist.Access).GetRPCClient().(comm.ClientInterface)
	return rpcclient
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

func GetExecutionContext(app interface{}) data.Datastore {
	context := app.(persist.Access).GetExecutionContext().(data.Datastore)
	if context == nil {
		log.Fatal("ExecutionContext Database Missing", "config", global.Current, "app", app)
	}
	return context
}

func GetValidators(app interface{}) *id.Validators {
	validators := app.(persist.Access).GetValidators().(*id.Validators)
	if validators == nil {
		log.Fatal("Validators list mission", "config", global.Current, "app", app)
	}
	return validators
}
