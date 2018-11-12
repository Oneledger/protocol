/*
	Copyright 2017-2018 OneLedger
*/
package app

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
)

// The args write Node name directly to the context, but we can check to see if that
// is persistent, or matches the last entry
func SetNodeName(app interface{}) {
	admin := action.GetAdmin(app)

	var existing AdminParameters

	interim := admin.Get(data.DatabaseKey("AdminParameters"))
	if interim == nil {
		log.Debug("Admin Parameters not found")
		existing = AdminParameters{}
		existing.NodeName = global.Current.NodeName
		existing.NodeAccountName = global.Current.NodeAccountName
	} else {
		existing = interim.(AdminParameters)

		// TODO: Denormalized, should need to check against a default to see if it changed...
		if global.Current.NodeAccountName != "" {
			log.Debug("Getting NodeAccountName from cmd args")
			existing.NodeAccountName = global.Current.NodeAccountName
		} else {
			global.Current.NodeAccountName = existing.NodeAccountName
		}
		log.Debug("Admin Parameters found", "existing", existing)
	}

	session := admin.Begin()
	session.Set(data.DatabaseKey("AdminParameters"), existing)
	session.Commit()
}
