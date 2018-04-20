/*
	Copyright 2017-2018 OneLedger

	Encapsulate the underlying storage from our app. Currently using
	Tendermint's memdb (just an in-memory Merkle Tree)
*/
package app

import (
	"github.com/tendermint/iavl" // TODO: Double check this with cosmos-sdk
	"github.com/tendermint/tmlibs/db"
)

type DatabaseKey []byte // Database key

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
func NewDatastore(name string, dsType DatastoreType) *Datastore {
	switch dsType {

	case MEMORY:
		// TODO: No Merkle tree?
		return &Datastore{
			Name: name,
			Data: db.NewMemDB(),
		}

	case PERSISTENT:
		storage, err := db.NewGoLevelDB("OneLedger-"+name, "./")
		if err == nil {
			panic("Can't create a database")
		}

		tree := iavl.NewTree(storage, 1000) // Do I need a historic tree here?

		return &Datastore{
			Name: name,
			Tree: tree,
		}

	default:
		panic("Unknown Type")

	}
}

// Store inserts or updates a value under a key
func (store Datastore) Store(key DatabaseKey, value Message) {
	store.Data.Set(key, value)
}

// Load return the stored value
func (store Datastore) Load(key DatabaseKey) (value Message) {
	return store.Data.Get(key)
}
