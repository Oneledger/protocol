/*
	Copyright 2017 - 2018 OneLedger

	Test the database persistence
*/
package data

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Oneledger/protocol/node/serialize"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

const TagTesting = "data_testing"

type Testing struct {
	Value []byte
}

func init() {
	serial.Register(Testing{})
	serialize.RegisterConcrete(new(Testing), TagTesting)
}

func TestDatabase(t *testing.T) {

	store := NewDatastore("TestDatabase", PERSISTENT)

	key := []byte("TheKey")

	value := Testing{
		Value: []byte("TheValue"),
	}

	session := store.Begin()

	log.Dump("The message...", value)

	session.Set(key, value)
	session.Commit()

	log.Dump("The message...", value)
	Check("Commited", store, key, value)

	store.Dump()
	store.Close()

	store = NewDatastore("TestDatabase", PERSISTENT)

	value = *(store.Get(key).(*Testing))
	Check("Reopened", store, key, value)

	store.Dump()
}

func Check(text string, store Datastore, key []byte, testing Testing) {
	log.Debug(text + " Check")
	log.Dump("Inside of Check", testing)

	value := testing.Value

	if store.Exists(key) {
		result := store.Get(key).(*Testing)
		text := fmt.Sprintf("[%X]:\t[%X]\n", key, result)

		log.Debug("Found Data", "text", text)
		if bytes.Compare(result.Value, value) != 0 {
			log.Debug("but it differs", "value", value, "result", result)
		}

	} else {
		log.Debug("Missing Value for key", "key", key)
	}
}
