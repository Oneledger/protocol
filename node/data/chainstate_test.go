/*
	Copyright 2017 - 2018 OneLedger
*/
package data

import (
	"testing"

	"github.com/Oneledger/protocol/node/log"
)

func TestPersistence(t *testing.T) {
	state := NewChainState("./SimpleTest", PERSISTENT)
	_ = state

	key := "Hello"
	value := "The Value"

	state.Delivered.Set(DatabaseKey(key), []byte(value))

	version := state.Delivered.Version64()
	index, result := state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Fetched", "index", index, "result", string(result))

	state.Commit()

	version = state.Delivered.Version64()
	index, result = state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Fetched", "index", index, "result", string(result))

	//state.Dump()
	state.Commit()

	version = state.Delivered.Version64()
	index, result = state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Fetched", "index", index, "result", string(result))
	//state.Dump()

}
