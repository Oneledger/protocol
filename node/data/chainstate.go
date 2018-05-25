/*
	Copyright 2017-2018 OneLedger

	Keep the state of the MerkleTrees between the different stages of consensus

	We need an up to date tree to check new transactions against. We then
	need to apply them when delivered. We also need to get to the last tree.

	The difficulty comes from the underlying code not quite being thread-safe...
*/
package data

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tmlibs/db"
)

var count int

type ChainState struct {
	Name string
	Type DatastoreType

	// TODO: This doesn't work (can connect the roots to the underlying db)

	Checked   *iavl.VersionedTree
	Delivered *iavl.VersionedTree
	Committed *iavl.VersionedTree
}

func NewChainState(name string, newType DatastoreType) *ChainState {
	count = 0
	chain := &ChainState{Name: name}
	chain.reset()
	return chain
}

func createDatabase(name string, newType DatastoreType) *iavl.VersionedTree {
	// TODO: Assuming persistence for right now
	storage, err := db.NewGoLevelDB("OneLedger-"+name, global.Current.RootDir)
	if err != nil {
		log.Error("Database create failed", "err", err, "count", count)
		panic("Can't create a database")
	}

	// TODO: cosmos seems to be using VersionedTree now????
	tree := iavl.NewVersionedTree(storage, 100) // Do I need a historic tree here?

	count = count + 1

	return tree
}

// TODO: Not sure about this, it seems to be Cosmos-sdk's way of getting arround the immutable copy problem...
func (state *ChainState) Commit() {

	state.Delivered.SaveVersion() // TODO: This does not seem to be updating the database
	state.reset()
}

func (state *ChainState) reset() {
	//state.Committed = createDatabase(state.Name, state.Type)
	state.Delivered = createDatabase(state.Name, state.Type)
	//state.Checked = createDatabase(state.Name, state.Type)
}
