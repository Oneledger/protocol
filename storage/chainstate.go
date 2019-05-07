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

	Keep the state of the MerkleTrees between the different stages of consensus

	We need an up to date tree to check new transactions against. We then
	need to apply them when delivered. We also need to get to the last tree.

	The difficulty comes from the underlying code not quite being thread-safe...
*/

package storage

import (
	"fmt"
	"sync"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// Number of times we initialized since starting
var count int

// Chainstate is a storage for balances on the chain, a snapshot of all accounts
type ChainState struct {
	Name string
	Type StorageType

	Delivered *iavl.MutableTree // Build us a new set of transactions
	database  db.DB

	// Last committed values
	LastVersion int64
	Version     int64
	LastHash    []byte
	Hash        []byte
	TreeHeight  int8
	configDB    string
	dbDir       string

	sync.RWMutex
}

// NewChainState generates a new ChainState object
func NewChainState(name, dbDir, configDB string, newType StorageType) *ChainState {
	count = 0

	chain := &ChainState{Name: name, Type: newType, Version: 0}
	chain.dbDir = dbDir
	chain.configDB = configDB

	chain.reset()

	return chain
}

// Do this only for the Delivery side
func (state *ChainState) Set(key StoreKey, val []byte) error {
	state.Lock()
	defer state.Unlock()

	setOk := state.Delivered.Set(key, val)
	if !setOk {
		return fmt.Errorf("%s %#v \n", "failed to set bal", val)
	}
}

// Expensive O(n) search through everything...
func (state *ChainState) FindAll() map[string][]byte {
	mapping := make(map[string][]byte, 1)

	for i := int64(0); i < state.Delivered.Size(); i++ {
		key, value := state.Delivered.GetByIndex(i)
		mapping[string(key)] = value
	}

	return mapping
}

// TODO: Should be against the commit tree, not the delivered one!!!
func (state *ChainState) Get(key StoreKey, lastCommit bool) []byte {

	// TODO: Should not be this hardcoded, but still needs protection
	if len(key) != CHAINKEY_MAXLEN {
		log.Fatal("Not a valid account key")
	}

	var value []byte
	if lastCommit {
		// get the value of last commit version
		version := state.Delivered.Version()
		_, value = state.Delivered.GetVersioned(key, version)
	} else {
		// get the value of currently working tree. it's temporary value that is not persistent yet.
		_, value = state.Delivered.ImmutableTree.Get(key)
	}

	return value
}

// TODO: Should be against the commit tree, not the delivered one!!!
func (state *ChainState) Exists(key StoreKey) bool {

	version := state.Delivered.Version()
	_, value := state.Delivered.GetVersioned([]byte(key), version)
	if value == nil {
		return false
	}

	return true
}

// TODO: Not sure about this, it seems to be Cosmos-sdk's way of getting arround the immutable copy problem...
func (state *ChainState) Commit() ([]byte, int64) {

	state.RLock()
	// Persist the Delivered merkle tree
	hash, version, err := state.Delivered.SaveVersion()
	if err != nil {
		log.Fatal("Saving", "err", err)
	}
	state.RUnlock()

	state.LastVersion, state.Version = state.Version, version
	state.LastHash, state.Hash = state.Hash, hash

	if state.LastVersion-1 > 0 {
		err := state.Delivered.DeleteVersion(state.LastVersion - 1)
		if err != nil {
			log.Fatal("Failed to delete old version of chainstate", "err", err)
		}
	}

	return hash, version
}

func (state *ChainState) Dump() {
	texts := state.database.Stats()

	for key, value := range texts {
		log.Debug("Stat", key, value)
	}

}

// Reset the chain state from persistence
func (state *ChainState) reset() ([]byte, int64) {

	state.Delivered, state.database = initializeDatabase(state.Name, state.dbDir, state.configDB, state.Type, state.Version)

	// Essentially, the last commited value...
	state.LastHash = state.Hash
	state.LastVersion = state.Version

	// Essentially, the last commited value...
	state.Hash = state.Delivered.Hash()
	state.Version = state.Delivered.Version()
	state.TreeHeight = state.Delivered.Height()

	log.Debug("Reinitialized Database", "version", state.Version, "tree_height", state.TreeHeight, "hash", state.Hash)
	return state.Hash, state.Version
}

// Create or attach to a database
func initializeDatabase(name, dbDir, configDB string, newType StorageType, version int64) (*iavl.MutableTree, db.DB) {
	// TODO: Assuming persistence for right now
	storage, err := getDatabase(name, dbDir, configDB)
	if err != nil {
		log.Error("Database create failed", "err", err, "count", count)
		panic("Can't create a database: " + dbDir + "/OneLedger-" + name)
	}

	tree := iavl.NewMutableTree(storage, CHAINSTATE_CACHE_SIZE) // Do I need a historic tree here?
	_, err = tree.LoadVersion(version)
	if err != nil {
		log.Error("error in loading tree version", "version", version, "err", err)
	}

	count = count + 1

	return tree, storage
}

func (state *ChainState) Close() {
	state.database.Close()
}
