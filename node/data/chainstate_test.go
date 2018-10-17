/*
	Copyright 2017 - 2018 OneLedger
*/
package data

import (
	"flag"
	"os"
	"testing"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
)

// Control the execution
func TestMain(m *testing.M) {
	flag.Parse()

	// Set the debug flags according to whether the -v flag is set in go test
	if testing.Verbose() {
		global.Current.Debug = true
	} else {
		global.Current.Debug = false
	}

	// Run it all.
	code := m.Run()

	os.Exit(code)
}

func TestPersistence(t *testing.T) {
	log.Debug("Create new chain state")

	global.Current.RootDir = "./"

	state := NewChainState("SimpleTest", PERSISTENT)
	_ = state
	//log.Dump("The chain state", state)

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

	state.Dump()
	state.Commit()

	version = state.Delivered.Version64()
	index, result = state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Fetched", "index", index, "result", string(result))
	state.Dump()

}
