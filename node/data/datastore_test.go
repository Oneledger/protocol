/*
	Copyright 2017 - 2018 OneLedger

	Test the database persistence
*/
package data

import (
	"bytes"
	"fmt"
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
	Check("Storage", ds, key, value)

	ds.Commit()
	Check("Commited", ds, key, value)

	ds.Dump()
	ds.Close()

	ds = NewDatastore("localTestingDatabase", PERSISTENT)
	Check("Reopened", ds, key, value)

	ds.Dump()
}

func Check(text string, ds *Datastore, key []byte, value []byte) {
	log.Debug(text + " Check")

	if ds.Exists(key) {
		result := ds.Load(key)
		text := fmt.Sprintf("[%X]:\t[%X]\n", key, result)
		log.Debug("Found Data", "text", text)
		if bytes.Compare(result, value) != 0 {
			log.Debug("but it differs", "value", value, "result", result)
		}

	} else {
		log.Debug("Missing key", "key", key)
	}
}
