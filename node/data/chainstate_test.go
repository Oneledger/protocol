/*
	Copyright 2017 - 2018 OneLedger
*/
package data

import (
	"testing"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/stretchr/testify/assert"
)

func TestPersistence(t *testing.T) {
	log.Debug("Create new chain state")

	global.Current.RootDir = "./"
	state := NewChainState("PersistentTest", PERSISTENT)

	key := "Hello"
	value := "The Value"

	state.Delivered.Set(DatabaseKey(key), []byte(value))

	version := state.Delivered.Version()
	index, result := state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Uncommitted Fetched", "index", index, "version", version, "result", string(result))

	state.Commit()

	version = state.Delivered.Version()
	index, result = state.Delivered.GetVersioned(DatabaseKey(key), version)
	log.Debug("Commited Fetched", "index", index, "version", version, "result", string(result))

	assert.Equal(t, []byte(value), result, "These should be equal")

}

func TestChainState(t *testing.T) {
	state := NewChainState("ChainState", PERSISTENT)
	balance := NewBalanceFromInt(10000, "OLT")
	key := make([]byte, 20)
	key[0] = 0xaf

	state.Set(key, balance)
	state.Commit()
	result := state.Get(key, false)

	assert.Equal(t, balance, result, "These should be equal")
}

func TestChainStateContinueUpdate(t *testing.T) {

	state := NewChainState("Continue", PERSISTENT)

	key := make([]byte, 20)
	key[0] = 0x03

	balance := NewBalanceFromInt(10, "OLT")

	state.Set(key, balance)

	state.Commit()

	//_, b1 := state.Delivered.ImmutableTree.Get(key)
	b1 := state.Get(key, false)
	log.Debug("with commit", "balance", b1)

	assert.Equal(t, true, balance.GetAmountByName("OLT").Equals(b1.GetAmountByName("OLT")), "balance not eqaul after commit", b1)

	newbalance := NewBalanceFromInt(14, "OLT")

	state.Set(key, newbalance)

	b2 := state.Get(key, false)
	//_, b2 := state.Delivered.ImmutableTree.Get(key)
	log.Debug("without commit", "balance", b2)

	assert.Equal(t, true, newbalance.GetAmountByName("OLT").Equals(b2.GetAmountByName("OLT")), "balance not eqaul without commit", b2)

}
