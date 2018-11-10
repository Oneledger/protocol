/*
	Copyright 2017-2018 OneLedger

	Encapsulate the underlying storage from our app. Currently using:
		Tendermint's memdb (just an in-memory Merkle Tree)
		Tendermint's persistent kvstore (with Merkle Trees & Proofs)
			- Can only be opened by one process...

	Only one connection can occur to LevelDB at a time...

*/
package data

import (
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/libs/db"
)

type DatabaseKey = []byte // Database key
type Message = []byte

// ENUM for datastore type
type StorageType int

// Different types
const (
	MEMORY StorageType = iota
	PERSISTENT
)

// Wrap the underlying usage
type KeyValue struct {
	Type StorageType

	Name string
	File string

	memory   *db.MemDB
	tree     *iavl.MutableTree
	database *db.GoLevelDB

	version int64
}

type KeyValueSession struct {
	store *KeyValue
}

// Test to see if this exists already
// TODO: Should be moved to some common/shared/utils directory
func FileExists(name string, dir string) bool {
	dbPath := filepath.Join(dir, name+".db")
	info, err := os.Stat(dbPath)
	if err != nil {
		return false
	}
	_ = info
	return true
}

func convertData(data interface{}) []byte {
	buffer, err := serial.Serialize(data, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Persistent Serialization Failed", "err", err, "data", data)
	}
	return buffer
}

func unconvertData(data []byte) interface{} {
	if data == nil || string(data) == "" {
		return nil
	}

	var proto interface{}
	result, err := serial.DumpDeserialize(data, proto, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Persistent Deserialization Failed", "err", err, "data", data)
	}
	return result
}

// NewApplicationContext initializes a new application
func NewKeyValue(name string, newType StorageType) *KeyValue {
	switch newType {

	case MEMORY:
		// TODO: No Merkle tree?
		return &KeyValue{
			Type:   newType,
			Name:   name,
			memory: db.NewMemDB(),
		}

	case PERSISTENT:
		fullname := "OneLedger-" + name

		if FileExists(fullname, global.Current.RootDir) {
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

		return &KeyValue{
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
	return nil
}

func (store KeyValue) Begin() Session {
	return NewKeyValueSession(&store)
}

func (store KeyValue) Dump() {
	// TODO: Dump out debugging information here
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

func (store KeyValue) Errors() string {
	return ""
}

// Close the database
func (store KeyValue) Close() {
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

func (store KeyValue) Reopen() {
}

func (store KeyValue) FindAll() []DatabaseKey {
	return store.list()
}

func (store KeyValue) Exists(key DatabaseKey) bool {
	version := store.tree.Version64()
	index, _ := store.tree.GetVersioned(key, version)
	if index == -1 {
		return false
	}
	return true
}

func (store KeyValue) Get(key DatabaseKey) interface{} {
	version := store.tree.Version64()
	index, value := store.tree.GetVersioned(key, version)
	if index == -1 {
		return nil
	}
	return unconvertData(value)
}

func NewKeyValueSession(store *KeyValue) Session {
	return &KeyValueSession{store: store}
}

func (session KeyValueSession) FindAll() []DatabaseKey {
	return session.store.list()
}

// Store inserts or updates a value under a key
func (session KeyValueSession) Set(key DatabaseKey, value interface{}) bool {
	log.Debug("KV Set", "key", key, "value", value)

	buffer := convertData(value)
	session.store.tree.Set(key, buffer)

	return true
}

func (session KeyValueSession) Exists(key DatabaseKey) bool {
	version := session.store.tree.Version64()
	index, _ := session.store.tree.GetVersioned(key, version)
	if index == -1 {
		return false
	}
	return true
}

// Load return the stored value
func (session KeyValueSession) Get(key DatabaseKey) interface{} {
	version := session.store.tree.Version64()
	index, value := session.store.tree.GetVersioned(key, version)
	if index == -1 {
		return nil
	}
	return unconvertData(value)
}

func (session KeyValueSession) Delete(key DatabaseKey) bool {
	return true
}

func (session KeyValueSession) Errors() string {
	return ""
}

// Commit the changes to persistence
func (session KeyValueSession) Commit() bool {
	_, version, err := session.store.tree.SaveVersion()
	if err != nil {
		log.Fatal("Database Error", "err", err)
	}
	session.store.version = version

	return true
}

func (session KeyValueSession) Rollback() bool {
	return false
}

// Dump out the contents of the database
func (session KeyValueSession) Dump() {
	texts := session.store.database.Stats()
	for key, value := range texts {
		log.Debug("Stat", key, value)
	}

	iter := session.store.database.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		hash := iter.Key()
		node := iter.Value()
		log.Debug("Row", hash, node)
	}
}

// List all of the keys
func (store KeyValue) list() (keys []DatabaseKey) {
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
func (store KeyValue) empty() {
	switch store.Type {
	case MEMORY:
	case PERSISTENT:
	default:
		panic("Unknown Type")
	}
}
