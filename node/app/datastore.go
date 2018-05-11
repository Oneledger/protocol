/*
	Copyright 2017-2018 OneLedger

	Encapsulate the underlying storage from our app. Currently using:
		Tendermint's memdb (just an in-memory Merkle Tree)
		Tendermint's persistent kvstore (with Merkle Trees & Proofs)
			- Can only be opened by one process...

*/
package app

import (
	"github.com/Oneledger/prototype/node/log"
	"github.com/tendermint/iavl" // TODO: Double check this with cosmos-sdk
	"github.com/tendermint/tmlibs/db"
)

type DatabaseKey = []byte // Database key

// ENUM for datastore type
type DatastoreType int

const (
	MEMORY     DatastoreType = iota
	PERSISTENT DatastoreType = iota
)

// Wrap the underlying usage
type Datastore struct {
	Type DatastoreType
	Name string
	Data *db.MemDB
	Tree *iavl.Tree
}

// NewApplicationContext initializes a new application
func NewDatastore(name string, newType DatastoreType) *Datastore {
	switch newType {

	case MEMORY:
		// TODO: No Merkle tree?
		return &Datastore{
			Type: newType,
			Name: name,
			Data: db.NewMemDB(),
		}

	case PERSISTENT:
		storage, err := db.NewGoLevelDB("OneLedger-"+name, Current.RootDir)
		if err != nil {
			log.Error("Database create failed", "err", err)
			panic("Can't create a database " + Current.RootDir + "/" + "OneLedger-" + name)
		}

		tree := iavl.NewTree(storage, 1000) // Do I need a historic tree here?

		return &Datastore{
			Type: newType,
			Name: name,
			Tree: tree,
		}
	default:
		panic("Unknown Type")

	}
}

// Store inserts or updates a value under a key
func (store Datastore) Store(key DatabaseKey, value Message) Message {
	switch store.Type {

	case MEMORY:
		store.Data.Set(key, value)

	case PERSISTENT:
		store.Tree.Set(key, value)

	default:
		panic("Unknown Type")
	}
	return value
}

// Load return the stored value
func (store Datastore) Load(key DatabaseKey) (value Message) {
	switch store.Type {

	case MEMORY:
		return store.Data.Get(key)

	case PERSISTENT:
		_, value := store.Tree.Get(key)
		return Message(value)

	default:
		panic("Unknown Type")
	}
}
