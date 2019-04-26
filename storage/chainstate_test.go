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
*/

package storage


import (
	"bytes"
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/balance"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersistence(t *testing.T) {
	log.Debug("Create new chain state")

	state := NewChainState("PersistentTest", "goleveldb", PERSISTENT)

	key := "Hello"
	value := "The Value"

	state.Delivered.Set(data.StoreKey(key), []byte(value))

	version := state.Delivered.Version()
	index, result := state.Delivered.GetVersioned(data.StoreKey(key), version)
	log.Debug("Uncommitted Fetched", "index", index, "version", version, "result", string(result))

	state.Commit()

	version = state.Delivered.Version()
	index, result = state.Delivered.GetVersioned(data.StoreKey(key), version)
	log.Debug("Commited Fetched", "index", index, "version", version, "result", string(result))

	assert.Equal(t, []byte(value), result, "These should be equal")

}

func TestChainState(t *testing.T) {
	state := NewChainState("ChainState", "goleveldb", PERSISTENT)
	bal := balance.NewBalanceFromInt(10000, "OLT")
	key := make([]byte, 20)
	key[0] = 0xaf

	state.Set(key, bal)
	state.Commit()
	result := state.Get(key, false)

	assert.Equal(t, bal, result, "These should be equal")
}

func TestChainStateContinueUpdate(t *testing.T) {

	state := NewChainState("Continue","goleveldb", PERSISTENT)

	key := make([]byte, 20)
	key[0] = 0x03

	bal := balance.NewBalanceFromInt(10, "OLT")

	state.Set(key, bal)

	state.Commit()

	//_, b1 := state.Delivered.ImmutableTree.Get(key)
	b1 := state.Get(key, false)
	log.Debug("with commit", "balance", b1)

	assert.Equal(t, true, bal.GetCoinByName("OLT").Equals(b1.GetCoinByName("OLT")), "balance not eqaul after commit", b1)

	newbalance := balance.NewBalanceFromInt(14, "OLT")

	state.Set(key, newbalance)

	b2 := state.Get(key, false)
	//_, b2 := state.Delivered.ImmutableTree.Get(key)
	log.Debug("without commit", "balance", b2)

	assert.Equal(t, true, newbalance.GetCoinByName("OLT").Equals(b2.GetCoinByName("OLT")), "balance not eqaul without commit", b2)

}

func TestChainState_Commit(t *testing.T) {

	state := NewChainState("commit", "goleveldb", PERSISTENT)

	key := make([]byte, 20)
	key[0] = 0x05
	bal := balance.NewBalanceFromInt(10, "OLT")

	state.Set(key, bal)

	hash, version := state.Commit()

	// Force the database to completely close, then repoen it.
	state.database.Close()
	state.database = nil

	// Update all of the chain parameters
	nhash, nversion := state.reset()

	// Check the reset
	assert.Equal(t, 0, bytes.Compare(hash, nhash), "hash of persistent after commit not match")

	assert.Equal(t, version, nversion, "version of persistent after commit not match")
}

