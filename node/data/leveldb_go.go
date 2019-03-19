// +build !gcc

// This file is for grabbing a leveldb instance without cleveldb support
package data

import (
	"errors"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/libs/db"
)

func init() {
	log.Info("Compiled without GCC, no cleveldb support...")
}

func getDatabase(name string) (db.DB, error) {
	if global.Current.Config.Node.DB == "cleveldb" {
		return nil, errors.New("Binary compiled without cleveldb support. Failed because \"cleveldb\" was specified in config")
	}
	return db.NewGoLevelDB("OneLedger-"+name, global.Current.DatabaseDir())
}
