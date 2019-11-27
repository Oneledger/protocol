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
	"bytes"
	"sync"
)

/* cache is a simple in-memory keyvalue store, to store binary  This is not thread safe and
any concurrent read/write might throw panics.
*/
type cache struct {
	name  string
	store map[string][]byte
	keys  []string
}

// cache satisfies Store interface
var _ Store = &cache{}
var _ Iteratable = &cache{}

func NewCache(name string) *cache {
	return &cache{
		name:  name,
		store: make(map[string][]byte),
		keys:  make([]string, 0, 100),
	}
}

// Get retrieves data for a key.
func (c *cache) Get(key StoreKey) ([]byte, error) {

	d, ok := c.store[string(key)]
	if !ok {
		return nil, ErrNotFound
	}

	return d, nil
}

// Exists checks if a key exists in the database.
func (c *cache) Exists(key StoreKey) bool {

	_, ok := c.store[string(key)]

	return ok
}

// Set is used to store or update some data with a key
func (c *cache) Set(key StoreKey, dat []byte) error {

	c.store[string(key)] = dat
	c.keys = append(c.keys, string(key))
	return nil
}

// Delete removes any data stored against a key
func (c *cache) Delete(key StoreKey) (bool, error) {

	delete(c.store, string(key))
	return true, nil
}

func (c *cache) GetIterator() Iteratable {
	return c
}

func (c *cache) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	for _, k := range c.keys {
		v, ok := c.store[k]
		if !ok {
			continue
		}
		if fn([]byte(k), v) {
			return true
		}
	}
	return true
}

//
func (c *cache) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	panic("IterateRange not implemented for cache kv")
}

/*
	CacheSafe starts here
*/

// cacheSafe is a thread safe implementation of above cache
type cacheSafe struct {
	sync.RWMutex

	cache
}

// cacheSafe pointer satisfies Store interface
var _ Store = &cacheSafe{}

func NewCacheSafe(name string) *cacheSafe {
	return &cacheSafe{sync.RWMutex{}, *NewCache(name)}
}

// Get retrieves data for a key.
func (c *cacheSafe) Get(key StoreKey) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()

	return c.cache.Get(key)
}

// Exists checks if a key exists in the database.
func (c *cacheSafe) Exists(key StoreKey) bool {
	c.RLock()
	defer c.RUnlock()

	return c.Exists(key)
}

// Set is used to store or update some data with a key
func (c *cacheSafe) Set(key StoreKey, dat []byte) error {
	c.Lock()
	defer c.Unlock()

	return c.cache.Set(key, dat)
}

// Delete removes any data stored against a key
func (c *cacheSafe) Delete(key StoreKey) (bool, error) {
	c.Lock()
	defer c.Unlock()

	return c.cache.Delete(key)
}

func (c *cacheSafe) GetIterator() Iteratable {
	c.RLock()
	defer c.RUnlock()

	return c.cache.GetIterator()
}

func (c *cacheSafe) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	c.RLock()
	defer c.RUnlock()

	return c.cache.Iterate(fn)
}

//
func (c *cacheSafe) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	panic("IterateRange not implemented for cache kv")
}

/*
	utils
*/
func isKeyInDomain(key, start, end []byte) bool {
	if bytes.Compare(key, start) < 0 {
		return false
	}

	if bytes.Compare(end, key) <= 0 {
		return false
	}

	return true
}
