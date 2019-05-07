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

var _ Storage = &KeyValue{}
var _ StorageSession = &KeyValueSession{}
var _ Store = &cacheSafe{}
var _ Store = &cache{}

// Storage wraps objects with option to start a session(db transaction)
type Storage interface {
	Get(StoreKey) ([]byte, error)
	Exists(StoreKey) (bool, error)

	Begin() StorageSession
	Close()
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

/*
	StorageSession
*/

// StorageSession defines a session-ed storage object of your choice
type StorageSession interface {
	Store
	Commit() bool
	FindAll() []StoreKey
}

// NewStorageSession creates a new SessionStorage
func NewSessionStorageDB(flavor, name string, ctx Context) Storage {

	switch flavor {
	case KEYVALUE:
		return newKeyValue(name, ctx.DbDir, ctx.ConfigDB, PERSISTENT)
	default:
		log.Error("incorrect session storage: ", flavor)
	}
	return nil
}

type Context struct {
	DbDir    string
	ConfigDB string
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

type Store interface {
	Get(StoreKey) ([]byte, error)
	Set(StoreKey, []byte) error
	Exists(StoreKey) bool
	Delete(StoreKey) (bool, error)
}
