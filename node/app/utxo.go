/*
	Copyright 2017-2018 OneLedger

	Keep the state of the MerkleTree between the different stages of consensus
*/
package app

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tmlibs/db"
)

type ChainState struct {
	Check     *iavl.Tree
	Deliver   *iavl.Tree
	Committed *iavl.Tree
}

func NewChainState(name string, newType DatastoreType) *ChainState {

	// TODO: Assuming persistence for right now
	storage, err := db.NewGoLevelDB("OneLedger-"+name, Current.RootDir)
	if err != nil {
		Log.Error("Database create failed", "err", err)
		panic("Can't create a database")
	}

	tree := iavl.NewTree(storage, 1000) // Do I need a historic tree here?

	// TODO: Get the chain state from persistence
	return &ChainState{
		Committed: tree,
		//Deliver:   tree.Copy(),
		//Check:     tree.Copy(),
	}
}
