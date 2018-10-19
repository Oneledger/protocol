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
	"github.com/stretchr/testify/assert"
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
	state := NewChainState("PersistentTest", PERSISTENT)

	key := "Hello"
	value := "The Value"

	state.Delivered.Set(DatabaseKey(key), []byte(value))

	version := state.Delivered.Version64()
	index, result := state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Uncommitted Fetched", "index", index, "version", version, "result", string(result))

	state.Commit()

	version = state.Delivered.Version64()
	index, result = state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Commited Fetched", "index", index, "version", version, "result", string(result))

	assert.Equal(t, []byte(value), result, "These should be equal")

}

func TestChainState(t *testing.T) {
	state := NewChainState("ChainState", PERSISTENT)
	balance := NewBalance(10000, "OLT")
	key := []byte("Ahhhhhhh")
	state.Set(key, balance)
	state.Commit()
	result := state.Find(key)

	assert.Equal(t, balance, *result, "These should be equal")
}
