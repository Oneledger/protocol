/*
	Copyright 2017 - 2018 OneLedger

	Basic datatypes
*/
package data

// Database is a consistent interface for all underlying persistent data stores.
type Datastore interface {
	// Contruction/Destruction
	//Initialize(name string)
	Close()
	Reopen()

	// Some of the databases require commits, some are persisted right way
	Begin() Session

	// Readonly, not in a session
	FindAll() []DatabaseKey
	Exists(key DatabaseKey) bool
	Get(key DatabaseKey) interface{}

	// Give out a listof errors
	Errors() string
	Dump()
}

// Open a session (transaction in the database sense, not blockchain).
type Session interface {
	// Primary operations
	FindAll() []DatabaseKey
	Exists(key DatabaseKey) bool
	Get(key DatabaseKey) interface{}
	Set(key DatabaseKey, value interface{}) bool
	Delete(key DatabaseKey) bool

	// Finish the sessions
	Commit() bool
	Rollback() bool

	Errors() string
	Dump()
}

func NewDatastore(name string, stype StorageType) Datastore {
	return NewKeyValue(name, stype)
}
