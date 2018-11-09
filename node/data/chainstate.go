/*
	Copyright 2017-2018 OneLedger

	Keep the state of the MerkleTrees between the different stages of consensus

	We need an up to date tree to check new transactions against. We then
	need to apply them when delivered. We also need to get to the last tree.

	The difficulty comes from the underlying code not quite being thread-safe...
*/
package data

import (
	"bytes"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

// Number of times we initialized since starting
var count int

type ChainState struct {
	Name string
	Type StorageType

	Delivered *iavl.MutableTree // Build us a new set of transactions
	database  *db.GoLevelDB

	Checked   *iavl.MutableTree // Temporary and can be Rolled Back
	Committed *iavl.MutableTree // Last Persistent Tree

	// Last committed values
	LastVersion int64
	Version     int64
	LastHash    []byte
	Hash        []byte
	TreeHeight  int
}

func NewChainState(name string, newType StorageType) *ChainState {
	count = 0
	chain := &ChainState{Name: name, Type: newType}
	chain.reset()
	return chain
}

/*
// Test this against the checked UTXO data to make sure the transaction is legit
func (state *ChainState) Test(key DatabaseKey, balance Balance) bool {
	//buffer := serial.Serialize(balance)
	//state.Checked.Set(key, buffer)
	return true
}
*/

// Do this only for the Delivery side
func (state *ChainState) Set(key DatabaseKey, balance Balance) {
	buffer, err := serial.Serialize(balance, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Failed to Deserialize balance: ", err)
	}

	// TODO: Get some error handling in here
	state.Delivered.Set(key, buffer)
}

// Do this only for the Delivery side
func (state *ChainState) Test(key DatabaseKey, balance Balance) {
	buffer, err := serial.Serialize(balance, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Failed to Deserialize balance: ", err)
	}

	// TODO: Get some error handling in here
	state.Checked.Set(key, buffer)
}

// Expensive O(n) search through everything...
func (state *ChainState) FindAll() map[string]*Balance {
	mapping := make(map[string]*Balance, 1)

	for i := int64(0); i < state.Delivered.Size64(); i++ {
		key, value := state.Delivered.GetByIndex64(i)

		var balance Balance
		result, err := serial.Deserialize(value, balance, serial.PERSISTENT)
		if err != nil {
			log.Fatal("Failed to Deserialize: FindAll", "i", i, "key", string(key))
			continue
		}

		final := result.(Balance)
		mapping[string(key)] = &final
	}
	return mapping
}

// TODO: Should be against the commit tree, not the delivered one!!!
func (state *ChainState) Get(key DatabaseKey) *Balance {

	version := state.Delivered.Version64()
	_, value := state.Delivered.GetVersioned(key, version)

	if value != nil {
		var balance Balance
		result, err := serial.Deserialize(value, balance, serial.PERSISTENT)
		if err != nil {
			log.Fatal("Failed to deserialize Balance in chainstate: ", err)
			return nil
		}
		final := result.(Balance)
		return &final
	}

	// By definition, if a balance doesn't exist, it is zero
	empty := NewBalance(0, "OLT")
	return &empty
}

// TODO: Should be against the commit tree, not the delivered one!!!
func (state *ChainState) Exists(key DatabaseKey) bool {

	version := state.Delivered.Version64()
	_, value := state.Delivered.GetVersioned([]byte(key), version)

	if value != nil {
		return true
	}
	return false
}

// TODO: Not sure about this, it seems to be Cosmos-sdk's way of getting arround the immutable copy problem...
func (state *ChainState) Commit() ([]byte, int64) {

	hash, version, err := state.Delivered.SaveVersion()
	if err != nil {
		log.Fatal("Saving", "err", err)
	}

	// TODO: Force the database to completely close, then repoen it.
	state.database.Close()
	state.database = nil

	nhash, nversion := state.reset()
	if bytes.Compare(hash, nhash) != 0 || version != nversion {
		log.Fatal("Persistence Failed, difference in hash,version",
			"version", version, "nversion", nversion, "hash", hash, "nhash", nhash)
	}

	return hash, version
}

func (state *ChainState) Dump() {
	texts := state.database.Stats()

	for key, value := range texts {
		log.Debug("Stat", key, value)
	}

	// TODO: Need a way to just list out the last changes, not all of them
	/*
		iter := state.database.Iterator(nil, nil)
		for ; iter.Valid(); iter.Next() {
			hash := iter.Key()
			node := iter.Value()
			log.Debug("ChainState", hash, node)
		}
	*/
}

// Reset the chain state from persistence
func (state *ChainState) reset() ([]byte, int64) {
	// TODO: I need three copies of the tree, only one is ultimately mutable... (checked changed rollback)
	// TODO: Close before reopen, better just update...

	//state.Checked = createDatabase(state.Name, state.Type)
	state.Delivered, state.database = initializeDatabase(state.Name, state.Type)
	//state.Committed = createDatabase(state.Name, state.Type)

	// TODO: Can I stick the delivered database into the checked tree?

	// Essentially, the last commited value...
	state.LastHash = state.Hash
	state.LastVersion = state.Version

	// Essentially, the last commited value...
	state.Hash = state.Delivered.Hash()
	state.Version = state.Delivered.Version64()
	state.TreeHeight = state.Delivered.Height()

	log.Debug("Reinitialized Database", "version", state.Version, "tree_height", state.TreeHeight, "hash", state.Hash)
	return state.Hash, state.Version
}

// Create or attach to a database
func initializeDatabase(name string, newType StorageType) (*iavl.MutableTree, *db.GoLevelDB) {
	// TODO: Assuming persistence for right now
	storage, err := db.NewGoLevelDB("OneLedger-"+name, global.Current.RootDir)
	if err != nil {
		log.Error("Database create failed", "err", err, "count", count)
		panic("Can't create a database")
	}

	// TODO: cosmos seems to be using MutableTree now????
	tree := iavl.NewMutableTree(storage, 1000) // Do I need a historic tree here?
	tree.LoadVersion(0)

	count = count + 1

	return tree, storage
}
