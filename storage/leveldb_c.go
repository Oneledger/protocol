// +build gcc

/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


	Copyright 2017 - 2019 OneLedger

   This file is for grabbing a leveldb instance WITH cleveldb support
   It is only loaded if the code is compiled with CGO_ENABLED=1 and the "gcc" tag added

*/

package storage

import (
	"github.com/tendermint/tendermint/libs/db"
)

func init() {
	// log.Info("Node running with cleveldb support...")
}

func GetDatabase(name, dbDir, configDB string) (db.DB, error) {

	if configDB == "cleveldb" {
		log.Info("Getting cleveldb...")
		return db.NewCLevelDB(name, dbDir)
	}

	log.Info("Getting goleveldb...")
	return db.NewGoLevelDB(name, dbDir)
}
