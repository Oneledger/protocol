/*
	Copywrite 2017-2018 OneLedger

	Encapsulate the underlying storage from our app. Currently using
	Tendermint's memdb (just in memory Merkle Tree)
*/
package app

import (
	"github.com/tendermint/tmlibs/db"
)

type Datastore struct {
	data *db.MemDB
}

// NewApplicationContext initializes a new application
func NewDatastore() *Datastore {
	return &Datastore{
		data: db.NewMemDB(),
	}
}

func (store Datastore) Store(key Key, value Message) {
	store.data.Set(key, value)
}

func (store Datastore) Load(key Key) (value Message) {
	return store.data.Get(key)
}
