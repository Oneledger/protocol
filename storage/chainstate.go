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
	"strconv"
	"sync"

	"github.com/Oneledger/protocol/config"

	"github.com/pkg/errors"
	"github.com/tendermint/iavl"
	tmdb "github.com/tendermint/tm-db"
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

	ChainStateRotation ChainStateRotationSetting

	sync.RWMutex
}

type ChainStateRotationSetting struct {
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
}

// NewChainState generates a new ChainState object
func NewChainState(name string, db tmdb.DB) *ChainState {

	chain := &ChainState{
		Name:    name,
		Version: 0,
	}
	log.Detail("Chain state:", name)
	chain.loadDB(db)

	return chain
}

// Setup the rotation configuration for ChainState,
// "recent" : latest number of version to persist
// "every"  : every X number of version to persist
// "cycles" : number of latest every version to persist
func (state *ChainState) SetupRotation(chainStateRotationCfg config.ChainStateRotationCfg) error {

	isValid := chainStateRotationCfg.Recent >= 0 && chainStateRotationCfg.Every >= 0 && chainStateRotationCfg.Cycles >= 0
	if isValid != true {
		err := errors.New("found negative value in chain state rotation config")
		return err
	}
	state.ChainStateRotation.recent = chainStateRotationCfg.Recent
	state.ChainStateRotation.every = chainStateRotationCfg.Every
	state.ChainStateRotation.cycles = chainStateRotationCfg.Cycles
	return nil

}

// Do this only for the Delivery side
func (state *ChainState) Set(key StoreKey, val []byte) error {
	state.Lock()
	defer state.Unlock()

	state.Delivered.Set(key, val)

	return nil
}

func (state *ChainState) GetIterable() Iterable {
	return state.Delivered.ImmutableTree
}

func (state *ChainState) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	return state.Delivered.IterateRange(start, end, ascending, fn)
}

func (state *ChainState) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	state.RLock()
	defer state.RUnlock()
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
	state.RLock()
	defer state.RUnlock()
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
	state.RLock()
	defer state.RUnlock()
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

	state.Lock()
	defer state.Unlock()
	// Persist the Delivered merkle tree
	hash, version, err := state.Delivered.SaveVersion()
	if err != nil {
		panic(errors.Wrap(err, "failed to commit, version: "+strconv.FormatInt(version, 10)))
	}

	state.LastVersion, state.Version = state.Version, version
	state.LastHash, state.Hash = state.Hash, hash

	release := state.LastVersion - state.ChainStateRotation.recent

	if release > 0 {
		if state.ChainStateRotation.every == 0 || release%state.ChainStateRotation.every != 0 {
			err := state.Delivered.DeleteVersion(release)
			if err != nil {
				log.Error("Failed to delete old version of chainstate", "err:", err, "version:", release)
			}
		}
		if state.ChainStateRotation.cycles != 0 && state.ChainStateRotation.every != 0 && release%state.ChainStateRotation.every == 0 {
			release = release - state.ChainStateRotation.cycles*state.ChainStateRotation.every
			err := state.Delivered.DeleteVersion(release)
			if err != nil {
				log.Error("Failed to delete old version of chainstate", "err", err, "version:", release)
			}
		}

	}

	return hash, version
}

func (state *ChainState) LoadVersion(version int64) (int64, error) {
	return state.Delivered.LoadVersion(version)
}

// Reset the chain state from persistence
func (state *ChainState) loadDB(db tmdb.DB) ([]byte, int64) {
	tree, _ := iavl.NewMutableTree(db, CHAINSTATE_CACHE_SIZE) // Do I need a historic tree here?
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

	log.Info("Reinitialized From Database", "version", state.Version, "tree_height", state.TreeHeight, "hash", hex.EncodeToString(state.Hash))
	return state.Hash, state.Version
}

func (state *ChainState) ClearFrom(version int64) error {
	version, err := state.Delivered.LoadVersionForOverwriting(version)
	log.Info("cleared version after: ", version)
	return err
}
