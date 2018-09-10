/*
	Copyright 2017-2018 OneLedger

	Encapsulate the underlying storage from our app. Currently using:
		Tendermint's memdb (just an in-memory Merkle Tree)
		Tendermint's persistent kvstore (with Merkle Trees & Proofs)
			- Can only be opened by one process...

*/
package data

import (
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type Message = []byte
type DatabaseKey = []byte // Database key

// ENUM for datastore type
type DatastoreType int

// Different types
const (
	MEMORY     DatastoreType = iota
	PERSISTENT DatastoreType = iota
)

// Wrap the underlying usage
type Datastore struct {
	Type DatastoreType

	Name string
	File string

	memory   *db.MemDB
	tree     *iavl.MutableTree
	database *db.GoLevelDB

	version int64
}

// Test to see if this exists already
func Exists(name string, dir string) bool {
	dbPath := filepath.Join(dir, name+".db")
	info, err := os.Stat(dbPath)
	if err != nil {
		return false
	}
	_ = info
	return true
}

// NewApplicationContext initializes a new application
func NewDatastore(name string, newType DatastoreType) *Datastore {
	switch newType {

	case MEMORY:
		// TODO: No Merkle tree?
		return &Datastore{
			Type:   newType,
			Name:   name,
			memory: db.NewMemDB(),
		}

	case PERSISTENT:
		fullname := "OneLedger-" + name

		if Exists(fullname, global.Current.RootDir) {
			//log.Debug("Appending to database", "name", fullname)
		} else {
			log.Info("Creating new database", "name", fullname)
		}

		storage, err := db.NewGoLevelDB(fullname, global.Current.RootDir)
		if err != nil {
			log.Error("Database create failed", "err", err)
			panic("Can't create a database " + global.Current.RootDir + "/" + fullname)
		}

		tree := iavl.NewMutableTree(storage, 100)

		// Note: the tree is empty, until at least one version is loaded
		tree.LoadVersion(0)

		return &Datastore{
			Type:     newType,
			Name:     name,
			File:     fullname,
			tree:     tree,
			database: storage,
			version:  tree.Version64(),
		}
	default:
		panic("Unknown Type")

	}
}

// Close the database
func (store Datastore) Close() {
	switch store.Type {

	case MEMORY:
		store.memory = nil

	case PERSISTENT:
		store.tree = nil
		store.database.Close()
		store.database = nil

	default:
		panic("Unknown Type")
	}
}

// Store inserts or updates a value under a key
func (store Datastore) Store(key DatabaseKey, value Message) Message {
	log.Info("Datastore Store", "key", key, "value", value)
	switch store.Type {

	case MEMORY:
		store.memory.Set(key, value)

	case PERSISTENT:
		store.tree.Set(key, value)

	default:
		panic("Unknown Type")
	}
	return value
}

func (store Datastore) Exists(key DatabaseKey) bool {
	switch store.Type {

	case MEMORY:
		return store.memory.Has(key)

	case PERSISTENT:
		version := store.tree.Version64()
		index, _ := store.tree.GetVersioned(key, version)
		if index != -1 {
			return true
		}

	default:
		panic("Unknown Type")
	}
	return false
}

// Load return the stored value
func (store Datastore) Load(key DatabaseKey) (value Message) {
	switch store.Type {

	case MEMORY:
		return store.memory.Get(key)

	case PERSISTENT:
		version := store.tree.Version64()
		_, value := store.tree.GetVersioned(key, version)
		return Message(value)

	default:
		panic("Unknown Type")
	}
	return Message(nil)
}

// Commit the changes to persistence
func (store Datastore) Commit() {
	switch store.Type {

	case PERSISTENT:
		_, version, err := store.tree.SaveVersion()
		if err != nil {
			log.Fatal("Database Error", "err", err)
		}
		store.version = version

		// Save only a few copies at a time
		//if store.version-10 > 10 {
		//		store.tree.DeleteVersion(store.version - 10)
		//	}
	}
}

// Dump out the contents of the database
func (store Datastore) Dump() {
	texts := store.database.Stats()
	for key, value := range texts {
		log.Debug("Stat", key, value)
	}

	iter := store.database.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		hash := iter.Key()
		node := iter.Value()
		log.Debug("Row", hash, node)
	}
}

// List all of the keys
func (store Datastore) List() (keys []DatabaseKey) {
	switch store.Type {

	case PERSISTENT:
		//store.tree.
		size := store.tree.Size()
		results := make([]DatabaseKey, size, size)
		for i := 0; i < store.tree.Size(); i++ {
			key, _ := store.tree.GetByIndex(i)
			results[i] = DatabaseKey(key)
		}
		log.Debug("Datastore List", "results", results)
		return results

	default:
		panic("Invalid Op")
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
