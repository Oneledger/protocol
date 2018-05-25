/*
	Copyright 2017 - 2018 OneLedger

	Test the database persistence
*/
package data

import (
	"testing"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
)

func TestDatabase(t *testing.T) {
	global.Current.RootDir = "./"
	ds := NewDatastore("localTestingDatabase", PERSISTENT)

	key := []byte("TheKey")
	value := []byte("TheValue")

	ds.Store(key, value)

	log.Debug("Store Check")
	Check(ds, key)

	ds.Commit()

	log.Debug("Commit Check")
	Check(ds, key)
	ds.Dump()

	ds.Close()

	ds = NewDatastore("localTestingDatabase", PERSISTENT)
	log.Debug("Reopen Check")
	Check(ds, key)
	ds.Dump()
}

func Check(ds *Datastore, key []byte) {
	if ds.Exists(key) {
		result := ds.Load(key)
		log.Debug("Found Data", "result", result)
	} else {
		log.Debug("Missing key", "key", key)
	}
}
