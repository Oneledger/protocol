/*
	Copyright 2017-2018 OneLedger

	Keep the state of the MerkleTrees between the different stages of consensus

	We need an up to date tree to check new transactions against. We then
	need to apply them when delivered. We also need to get to the last tree.

	The difficulty comes from the underlying code not quite being thread-safe...
*/
package data

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tmlibs/db"
)

// Number of times we initialized since starting
var count int

type ChainState struct {
	Name string
	Type DatastoreType

	Delivered *iavl.VersionedTree // Build us a new set of transactions
	database  *db.GoLevelDB

	Checked   *iavl.VersionedTree // Temporary and can be Rolled Back
	Committed *iavl.VersionedTree // Last Persistent Tree

	// Last committed values
	Version int64
	Height  int
	Hash    []byte
}

func NewChainState(name string, newType DatastoreType) *ChainState {
	count = 0
	chain := &ChainState{Name: name}
	chain.reset()
	return chain
}

// Test this against the checked UTXO data to make sure the transaction is legit
func (state *ChainState) Test(key DatabaseKey, balance Balance) bool {
	//buffer := comm.Serialize(balance)
	//state.Checked.Set(key, buffer)
	return true
}

// Do this for the Delivery side
func (state *ChainState) Set(key DatabaseKey, balance Balance) {
	buffer, err := comm.Serialize(balance)
	if err != nil {
		log.Error("Failed to Deserialize balance: ", err)
	}

	// TODO: Get some error handling in here
	state.Delivered.Set(key, buffer)
}

func (state *ChainState) FindAll() map[string]*Balance {
	mapping := make(map[string]*Balance, 1)

	for i := int64(0); i < state.Delivered.Size64(); i++ {
		key, value := state.Delivered.GetByIndex64(i)

		var balance Balance
		result, err := comm.Deserialize(value, &balance)
		if err != nil {
			log.Error("Failed to Deserialize: FindAll", "i", i, "key", string(key))
			continue
		}

		log.Debug("FindAll", "i", i, "key", string(key), "value", value, "result", result)
		mapping[string(key)] = result.(*Balance)
	}
	return mapping
}

// TODO: Should be against the commit tree, not the delivered one!!!
func (state *ChainState) Find(key DatabaseKey) *Balance {

	version := state.Delivered.Version64()
	_, value := state.Delivered.GetVersioned(key, version)

	if value != nil {
		var balance Balance
		result, err := comm.Deserialize(value, &balance)
		if err != nil {
			log.Error("Failed to deserialize Balance in chainstate: ", err)
			return nil
		}
		return result.(*Balance)
	}
	return nil
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

	// Force the database to completely close, then repoen it.
	state.database.Close()
	state.database = nil

	state.reset()

	return hash, version
}

func (state *ChainState) Dump() {
	texts := state.database.Stats()

	for key, value := range texts {
		log.Debug("Stat", key, value)
	}

	iter := state.database.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		hash := iter.Key()
		node := iter.Value()
		log.Debug("ChainState", hash, node)
	}
}

func (state *ChainState) reset() {
	// TODO: I need three copies of the tree, only one is ultimately mutable... (checked changed rollback)
	// TODO: Close before reopen, better just update...

	//state.Checked = createDatabase(state.Name, state.Type)
	state.Delivered, state.database = initializeDatabase(state.Name, state.Type)
	//state.Committed = createDatabase(state.Name, state.Type)

	// TODO: Can I stick the delivered database into the checked tree?

	// Essentially, the last commited value...
	state.Hash = state.Delivered.Hash()
	state.Version = state.Delivered.Version64()
	state.Height = state.Delivered.Height()
}

func initializeDatabase(name string, newType DatastoreType) (*iavl.VersionedTree, *db.GoLevelDB) {
	// TODO: Assuming persistence for right now
	storage, err := db.NewGoLevelDB("OneLedger-"+name, global.Current.RootDir)
	if err != nil {
		log.Error("Database create failed", "err", err, "count", count)
		panic("Can't create a database")
	}

	// TODO: cosmos seems to be using VersionedTree now????
	tree := iavl.NewVersionedTree(storage, 1000) // Do I need a historic tree here?
	tree.LoadVersion(0)

	count = count + 1

	return tree, storage
}
