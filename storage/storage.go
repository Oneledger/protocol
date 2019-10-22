/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/

	Copyright 2017 - 2019 OneLedger

*/

package storage

import (
	"encoding/hex"
)

var _ SessionedStorage = &KeyValue{}
var _ Session = &KeyValueSession{}
var _ Store = &cacheSafe{}
var _ Store = &cache{}

// SessionedStorage wraps objects with option to start a session(db transaction)
type SessionedStorage interface {
	Get(StoreKey) ([]byte, error)
	Exists(StoreKey) (bool, error)

	BeginSession() Session
	Close()
}

// Session defines a session-ed storage object of your choice
type Session interface {
	Store
	Commit() bool
}

// NewStorageSession creates a new SessionStorage
func NewStorageDB(flavor, name string, DBDir, DBType string) SessionedStorage {

	switch flavor {
	case KEYVALUE:
		return newKeyValue(name, DBDir, DBType, PERSISTENT)
	default:
		log.Error("incorrect session storage: ", flavor)
	}
	return nil
}

/*
	Base interfaces
*/

type StoreKey []byte

func (sk StoreKey) Bytes() []byte {
	return sk
}

func (sk StoreKey) String() string {
	return hex.EncodeToString(sk)
}

// store wraps object with option to use a cache type db
type Store interface {
	Get(StoreKey) ([]byte, error)
	Set(StoreKey, []byte) error
	Exists(StoreKey) bool
	Delete(StoreKey) (bool, error)
	GetIterator() Iteratable
}

// The iteratable interface include the function for iteration
// Iteratable function only be implemented for persistent data, doesn't guaranteed in the cache storage
type Iteratable interface {
	Iterate(fn func(key, value []byte) bool) (stop bool)
	IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool)
}

// NewStorage initializes a non sessioned storage
func NewStorage(flavor, name string) Store {
	switch flavor {
	case CACHE:
		return NewCache(name)
	case CACHE_SAFE:
		return NewCacheSafe(name)
	default:
		log.Error("incorrect storage: ", flavor)
	}
	return nil
}
