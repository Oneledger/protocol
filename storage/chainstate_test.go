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
	"github.com/Oneledger/protocol/config"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/db"
	"math"
	"strconv"
	"testing"
)

var cacheDB db.DB

func init() {
	cacheDB = db.NewDB("test", db.MemDBBackend, "")
}

func TestPersistence(t *testing.T) {
	log.Debug("Create new chain state")

	state := NewChainState("PersistentTest", cacheDB)

	key := "Hello"
	value := "The Value"

	state.Delivered.Set(StoreKey(key), []byte(value))

	version := state.Delivered.Version()
	index, result := state.Delivered.GetVersioned(StoreKey(key), version)
	log.Debug("Uncommitted Fetched", "index", index, "version", version, "result", string(result))

	state.Commit()

	version = state.Delivered.Version()
	index, result = state.Delivered.GetVersioned(StoreKey(key), version)
	log.Debug("Commited Fetched", "index", index, "version", version, "result", string(result))

	assert.Equal(t, []byte(value), result, "These should be equal")

}

func TestChainState(t *testing.T) {
	state := NewChainState("ChainState", cacheDB)
	key := make([]byte, 20)
	key[0] = 0xaf
	value := []byte("value1")

	err := state.Set(key, value)
	assert.Nil(t, err)
	state.Commit()
	result, err := state.Get(key)
	assert.NoError(t, err, "get")
	assert.Equal(t, value, result, "These should be equal")
}

func TestChainStateContinueUpdate(t *testing.T) {

	state := NewChainState("Continue", cacheDB)

	key := make([]byte, 20)
	key[0] = 0x03

	value := []byte("value1")

	err := state.Set(key, value)
	assert.Nil(t, err)

	state.Commit()

	//_, b1 := state.Delivered.ImmutableTree.Get(key)
	r1, err := state.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, r1, "value not eqaul after commit")

	value2 := []byte("value2")

	err = state.Set(key, value2)
	assert.Nil(t, err)

	r2, err := state.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value2, r2, "value not eqaul without commit")

}

func TestChainState_Commit(t *testing.T) {

	state := NewChainState("commit", cacheDB)

	key := make([]byte, 20)
	key[0] = 0x05
	value := []byte("value1")

	state.Set(key, value)

	hash, version := state.Commit()

	nhash, nversion := state.Commit()

	// Check the reset
	assert.Equal(t, 0, bytes.Compare(hash, nhash), "hash of persistent after commit not match")

	assert.Equal(t, version+1, nversion, "version of persistent after commit not match")
}

func TestChainState_Rotation(t *testing.T) {
	//set up test round
	testRound := 10000
	//the last round is special
	calculatedRound := testRound - 1

	//set up recent, every, cycles
	rotation := config.ChainStateRotationCfg{
		Recent : int64(10),
		Every : int64(100),
		Cycles : int64(10),
	}

	//generate multiple versions
	state := NewChainState("RotationTest", cacheDB)
	state.SetupRotation(rotation)


	//version start from 1
	for i := 1; i <= testRound; i++ {

		key := "Hello " + strconv.Itoa(i)
		value := "Value " + strconv.Itoa(i)

		state.Delivered.Set(StoreKey(key), []byte(value))

		state.Commit()

		version := state.Delivered.Version()
		index, result := state.Delivered.GetVersioned(StoreKey(key), version)
		log.Debug("Commited Fetched", "index", index, "version", version, "result", string(result))
	}

	//looking for each version
	counter := int64(0)
	for i := 1; i <= testRound; i++ {
		if state.Delivered.VersionExists(int64(i)) {
			log.Debug("remaining version ", i)
			counter++
		}

	}

	var correct int64
	recent := rotation.Recent
	every := rotation.Every
	cycles := rotation.Cycles

	// first figure out how many 'every' is reached
	// here reachedEvery only consider the ones outside 'recent'
	var reachedEvery int64
	if every != 0 {
		if cycles != 0{
			if (int64(testRound) - recent) % every == 0 { //'every' overlaps with the oldest in recent
				reachedEvery = int64(math.Min(float64((int64(testRound) - recent)/every - 1), float64(cycles)))
			} else {
				reachedEvery = int64(math.Min(float64((int64(testRound) - recent) / every), float64(cycles)))
			}
		} else {
			if (int64(testRound) - recent) % every == 0 { //'every' overlaps with the oldest in recent
				reachedEvery = (int64(testRound) - recent)/every - 1
			} else {
				reachedEvery = (int64(testRound) - recent)/every
			}
		}

	} else {
		reachedEvery = 0
	}


	if int64(calculatedRound) <= recent {
		correct = int64(calculatedRound)
	} else { // if 'every' and last round overlaps
		correct = recent + reachedEvery
	}

	correct += 1// add the last round

	assert.Equal(t, correct, counter, "These should be equal")

}
