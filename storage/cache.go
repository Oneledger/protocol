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

	"github.com/Oneledger/protocol/data"
)

/* cache is a simple in-memory keyvalue store, to store binary data. This is not thread safe and
any concurrent read/write might throw panics.
*/
type cache struct {
	name  string
	store map[string][]byte
}

// cache satisfies data.Store interface
var _ data.Store = &cache{}

func NewCache(name string) *cache {
	return &cache{name, map[string][]byte{}}
}

// Get retrieves data for a key.
func (c *cache) Get(key data.StoreKey) ([]byte, error) {

	d, ok := c.store[string(key)]
	if !ok {
		return nil, ErrNotFound
	}

	return d, nil
}

// Exists checks if a key exists in the database.
func (c *cache) Exists(key data.StoreKey) bool {

	_, ok := c.store[string(key)]

	return ok
}

// Set is used to store or update some data with a key
func (c *cache) Set(key data.StoreKey, dat []byte) error {

	c.store[string(key)] = dat

	return nil
}

// Delete removes any data stored against a key
func (c *cache) Delete(key data.StoreKey) (bool, error) {

	delete(c.store, string(key))
	return true, nil
}

func (c *cache) GetIterator() *Iterator {
	items := make([]iteratorItem, 0)

	for k, v := range c.store {
		items = append(items, iteratorItem{[]byte(k), v})
	}

	return newIterator(items)
}


func (c *cache) GetRangeIterator(start, end []byte) *Iterator {
	items := make([]iteratorItem, 0)

	for k, v := range c.store {
		key := []byte(k)
		if isKeyInDomain(key, start, end) {
			items = append(items, iteratorItem{key, v})
		}
	}
	return newIterator(items)
}

/*
	CacheSafe starts here
*/

// cacheSafe is a thread safe implementation of above cache
type cacheSafe struct {
	sync.RWMutex

	name  string
	store map[string][]byte
}

// cacheSafe pointer satisfies data.Store interface
var _ data.Store = &cacheSafe{}

func NewCacheSafe(name string) *cacheSafe {
	return &cacheSafe{sync.RWMutex{}, name, map[string][]byte{}}
}

// Get retrieves data for a key.
func (c *cacheSafe) Get(key data.StoreKey) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()

	d, ok := c.store[string(key)]
	if !ok {
		return nil, ErrNotFound
	}

	return d, nil
}

// Exists checks if a key exists in the database.
func (c *cacheSafe) Exists(key data.StoreKey) bool {
	c.RLock()
	defer c.RUnlock()

	_, ok := c.store[string(key)]

	return ok
}

// Set is used to store or update some data with a key
func (c *cacheSafe) Set(key data.StoreKey, dat []byte) error {
	c.Lock()
	defer c.Unlock()

	c.store[string(key)] = dat

	return nil
}

// Delete removes any data stored against a key
func (c *cacheSafe) Delete(key data.StoreKey) (bool, error) {
	c.Lock()
	defer c.Unlock()

	delete(c.store, string(key))
	return true, nil
}

func (c *cacheSafe) GetIterator() *Iterator {
	items := make([]iteratorItem, 0)

	c.RLock()
	defer c.RUnlock()
	for k, v := range c.store {
		items = append(items, iteratorItem{[]byte(k), v})
	}

	return newIterator(items)
}


func (c *cacheSafe) GetRangeIterator(start, end []byte) *Iterator {
	items := make([]iteratorItem, 0)

	c.RLock()
	defer c.RUnlock()
	for k, v := range c.store {
		key := []byte(k)
		if isKeyInDomain(key, start, end) {
			items = append(items, iteratorItem{key, v})
		}
	}
	return newIterator(items)
}

/*
	utils
*/
func  isKeyInDomain(key, start, end []byte) bool {
	if bytes.Compare(key, start) < 0 {
		return false
	}

	if bytes.Compare(end, key) <= 0 {
		return false
	}

	return true
}