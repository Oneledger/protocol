/*
	Copyright 2017-2018 OneLedger

	Encapsulate the underlying storage from our app. Currently using
	Tendermint's memdb (just in memory Merkle Tree)
*/
package app

import (
	"github.com/tendermint/tmlibs/db"
)

type DatastoreType int

const (
	MEMORY     DatastoreType = iota
	PERSISTENT DatastoreType = iota
)

type Datastore struct {
	name  string
	ttype DatastoreType
	data  *db.MemDB
}

// NewApplicationContext initializes a new application
func NewDatastore(name string, dsType DatastoreType) *Datastore {
	switch dsType {

	case MEMORY:
		return &Datastore{
			name: name,
			data: db.NewMemDB(),
		}

	case PERSISTENT:
		panic("Not yet implemented")

	default:
		panic("Unknown Type")

	}
}

// Store inserts or updates a value under a key
func (store Datastore) Store(key DatabaseKey, value Message) {
	store.data.Set(key, value)
}

// Load return the stored value
func (store Datastore) Load(key DatabaseKey) (value Message) {
	return store.data.Get(key)
}
