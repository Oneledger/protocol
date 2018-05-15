/*
	Copyright 2017-2018 OneLedger

	Encapsulate the underlying storage from our app. Currently using:
		Tendermint's memdb (just an in-memory Merkle Tree)
		Tendermint's persistent kvstore (with Merkle Trees & Proofs)
			- Can only be opened by one process...

*/
package data

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/iavl" // TODO: Double check this with cosmos-sdk
	"github.com/tendermint/tmlibs/db"
)

type Message = []byte
type DatabaseKey = []byte // Database key

// ENUM for datastore type
type DatastoreType int

const (
	MEMORY     DatastoreType = iota
	PERSISTENT DatastoreType = iota
)

// Wrap the underlying usage
type Datastore struct {
	Type    DatastoreType
	Name    string
	data    *db.MemDB
	tree    *iavl.VersionedTree
	version int64
}

// NewApplicationContext initializes a new application
func NewDatastore(name string, newType DatastoreType) *Datastore {
	switch newType {

	case MEMORY:
		// TODO: No Merkle tree?
		return &Datastore{
			Type: newType,
			Name: name,
			data: db.NewMemDB(),
		}

	case PERSISTENT:
		storage, err := db.NewGoLevelDB("OneLedger-"+name, global.Current.RootDir)
		if err != nil {
			log.Error("Database create failed", "err", err)
			panic("Can't create a database " + global.Current.RootDir + "/" + "OneLedger-" + name)
		}

		tree := iavl.NewVersionedTree(storage, 1000) // Do I need a historic tree here?

		return &Datastore{
			Type:    newType,
			Name:    name,
			tree:    tree,
			version: tree.Version64(),
		}
	default:
		panic("Unknown Type")

	}
}

// Store inserts or updates a value under a key
func (store Datastore) Store(key DatabaseKey, value Message) Message {
	switch store.Type {

	case MEMORY:
		store.data.Set(key, value)

	case PERSISTENT:
		store.tree.Set(key, value)

	default:
		panic("Unknown Type")
	}
	return value
}

// Load return the stored value
func (store Datastore) Load(key DatabaseKey) (value Message) {
	switch store.Type {

	case MEMORY:
		return store.data.Get(key)

	case PERSISTENT:
		_, value := store.tree.Get(key)
		return Message(value)

	default:
		panic("Unknown Type")
	}
}

// Commit the changes to persistence
func (store Datastore) Commit() {
	switch store.Type {
	case PERSISTENT:
		_, store.version, _ = store.tree.SaveVersion()

		// Save only one copy at a time
		if store.version-1 > 1 {
			store.tree.DeleteVersion(store.version - 1)
		}
	}
}

// Empty out all rows from the database
func (store Datastore) Empty() {
	switch store.Type {
	case MEMORY:
	case PERSISTENT:
	default:
		panic("Unknown Type")
	}
}
