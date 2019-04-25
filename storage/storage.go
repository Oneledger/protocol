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
	"github.com/Oneledger/protocol/data"
)

type Storage interface {
	Get(data.StoreKey) ([]byte, error)
	Exists(data.StoreKey) (bool, error)

	Begin() StorageSession
	Close()
}

type StorageSession interface {
	data.Store
	Commit() bool
}

type Context struct {
	DbDir string
	ConfigDB string
}

func NewSessionStorage(flavor, name string, ctx Context) Storage {

	switch flavor {
	case KEYVALUE:
		return NewKeyValue(name, ctx.DbDir, ctx.ConfigDB, PERSISTENT)
	default:
		log.Error("incorrect session storage: ", flavor)
	}
	return nil
}

func NewStorage(flavor, name string) data.Store {
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
