/*
	Copyright 2017 - 2018 OneLedger

	Easy Access to Persistent App Data
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/persist"
)

func GetAdmin(app interface{}) *data.Datastore {
	admin := app.(persist.Access).GetAdmin().(*data.Datastore)
	return admin
}

func GetStatus(app interface{}) *data.Datastore {
	status := app.(persist.Access).GetStatus().(*data.Datastore)
	return status
}

func GetIdentities(app interface{}) *id.Identities {
	identities := app.(persist.Access).GetIdentities().(*id.Identities)
	return identities
}

func GetAccounts(app interface{}) *id.Accounts {
	identities := app.(persist.Access).GetAccounts().(*id.Accounts)
	return identities
}

func GetUtxo(app interface{}) *data.ChainState {
	chain := app.(persist.Access).GetUtxo().(*data.ChainState)
	return chain
}
