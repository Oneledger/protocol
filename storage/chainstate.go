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
	"encoding/hex"
	"errors"
	"sync"

	"github.com/tendermint/iavl"
	tmdb "github.com/tendermint/tendermint/libs/db"
)

// Chainstate is a storage for balances on the chain, a snapshot of all accounts
type ChainState struct {
	Name string

	Delivered *iavl.MutableTree // Build us a new set of transactions

	// Last committed values
	LastVersion int64
	Version     int64
	LastHash    []byte
	Hash        []byte
	TreeHeight  int8

	//persistent data config

	// "recent" : latest number of version to persist
	// recent = 0 : keep last version only
	// recent = 3 : keep last 4 version
	recent int64

	// "every"  : every X number of version to persist
	// every = 0 : keep no other epoch version
	// every = 1 : keep every version
	// every > 1 : epoch number = every
	every int64

	// "cycles" : number of latest cycles for "every" to persist
	// cycles = 1 : only keep one of latest every
	// cycles = 0 : keep every "every"
	cycles int64

	sync.RWMutex
}

// NewChainState generates a new ChainState object
func NewChainState(name string, db tmdb.DB) *ChainState {

	chain := &ChainState{
		Name:    name,
		Version: 0,
	}
	chain.loadDB(db)

	return chain
}

// Setup the rotation configuration for ChainState,
// "recent" : latest number of version to persist
// "every"  : every X number of version to persist
// "cycles" : number of latest every version to persist
func (state *ChainState) SetupRotation(recent, every, cycles int64) {
	state.recent = recent
	state.every = every
	state.cycles = cycles
}

// Do this only for the Delivery side
func (state *ChainState) Set(key StoreKey, val []byte) error {
	state.Lock()
	defer state.Unlock()

	state.Delivered.Set(key, val)

	return nil
}

func (state *ChainState) GetIterator() Iteratable {
	return state.Delivered.ImmutableTree
}

func (state *ChainState) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	return state.Delivered.IterateRange(start, end, ascending, fn)
}

func (state *ChainState) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	return state.Delivered.Iterate(fn)
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
func (state *ChainState) Get(key StoreKey) ([]byte, error) {
	// get the value of currently working tree. it's temporary value that is not persistent yet.
	_, value := state.Delivered.ImmutableTree.Get(key)

	return value, nil
}

func (state *ChainState) GetLatestVersioned(key StoreKey) (int64, []byte) {
	// get the value of last commit version
	version := state.Delivered.Version()
	return state.Delivered.GetVersioned(key, version)

}

func (state *ChainState) GetVersioned(version int64, key StoreKey) (int64, []byte) {
	return state.Delivered.GetVersioned(key, version)
}

// TODO: Should be against the commit tree, not the delivered one!!!
func (state *ChainState) Exists(key StoreKey) bool {
	return state.Delivered.ImmutableTree.Has(key)
}

func (state *ChainState) Delete(key StoreKey) (bool, error) {
	_, ok := state.Delivered.Remove(key)
	if !ok {
		err := errors.New("Failed to delete the item from chainstate")
		log.Error(err.Error())
		return false, err
	}
	return ok, nil
}

// TODO: Not sure about this, it seems to be Cosmos-sdk's way of getting arround the immutable copy problem...
func (state *ChainState) Commit() ([]byte, int64) {

	state.RLock()
	// Persist the Delivered merkle tree
	hash, version, err := state.Delivered.SaveVersion()
	if err != nil {
		log.Fatal("Saving", "err", err)
	}

	state.LastVersion, state.Version = state.Version, version
	state.LastHash, state.Hash = state.Hash, hash

	release := state.LastVersion - state.recent

	if release > 0 {
		if state.every == 0 || release%state.every != 0 {
			err := state.Delivered.DeleteVersion(release)
			if err != nil {
				log.Error("Failed to delete old version of chainstate", "err", err)
			}
		}
		if state.cycles != 0 && release%state.every == 0 {
			release = release - state.cycles*state.every
			err := state.Delivered.DeleteVersion(release)
			if err != nil {
				log.Error("Failed to delete old version of chainstate", "err", err)
			}
		}

	}
	state.RUnlock()
	return hash, version
}

// Reset the chain state from persistence
func (state *ChainState) loadDB(db tmdb.DB) ([]byte, int64) {
	tree := iavl.NewMutableTree(db, CHAINSTATE_CACHE_SIZE) // Do I need a historic tree here?
	version, err := tree.Load()
	if err != nil {
		log.Error("error in loading tree version", "version", version, "err", err)
	}
	state.Delivered = tree
	// Essentially, the last commited value...
	state.LastHash = state.Hash
	state.LastVersion = state.Version

	// Essentially, the last commited value...
	state.Hash = state.Delivered.Hash()
	state.Version = state.Delivered.Version()
	state.TreeHeight = state.Delivered.Height()

	log.Debug("Reinitialized From Database", "version", state.Version, "tree_height", state.TreeHeight, "hash", hex.EncodeToString(state.Hash))
	return state.Hash, state.Version
}
