// +build gcc

// This file is for grabbing a leveldb instance WITH cleveldb support
// It is only loaded if the code is compiled with CGO_ENABLED=1 and the "gcc" tag added
package data

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/libs/db"
)

func init() {
	log.Info("Node running with cleveldb support...")
}

func getDatabase(name string) (db.DB, error) {
	if global.Current.Config.Node.DB == "cleveldb" {
		log.Info("Getting cleveldb...")
		return db.NewCLevelDB(name, global.Current.DatabaseDir())
	}
	log.Info("Getting goleveldb...")
	return db.NewGoLevelDB(name, global.Current.DatabaseDir())
	// panic("nogo")
}
