/*

 */

package storage

import (
	"fmt"
	"sort"
)

/* sessionCache is a simple in-memory keyvalue store, to store binary  This is not thread safe and
any concurrent read/write might throw panics.
*/
type sessionCache struct {
	name  string
	store map[string][]byte
	keys  []string
	done  map[string]bool
}

// sessionCache satisfies SessionedDirectStorage interface
var _ SessionedDirectStorage = &sessionCache{}
var _ Iterable = &sessionCache{}
var _ SessionedDirectStorage = &sessionCache{}

func NewSessionCache(name string) *sessionCache {
	return &sessionCache{
		name:  name,
		store: make(map[string][]byte),
		keys:  make([]string, 0, 100),
		done:  make(map[string]bool),
	}
}

func (c *sessionCache) DumpState() {
	keys1 := make([]string, 0)
	for _, key := range c.keys {
		keys1 = append(keys1, key)
	}
	sort.Strings(keys1)

	keys2 := make([]string, 0)
	for key, _ := range c.store {
		keys2 = append(keys2, key)
	}
	sort.Strings(keys2)

	fmt.Println("--- Start dump cache keys ---")
	for _, key := range keys1 {
		fmt.Println("key", key)
	}
	fmt.Println("--- End dump cache keys ---")
	fmt.Println("--- Start dump store data keys ---")
	for _, key := range keys2 {
		value := c.store[key]

		var val []byte
		if len(value) > 24 {
			val = value[:24]
		} else {
			val = value
		}
		fmt.Println("key", key, "value", val)
	}
	fmt.Println("--- End dump store data keys ---")
}

// Get retrieves store for a key.
func (c *sessionCache) Get(key StoreKey) ([]byte, error) {

	d, ok := c.store[string(key)]
	if !ok {
		return nil, ErrNotFound
	}

	return d, nil
}

// Exists checks if a key exists in the database.
func (c *sessionCache) Exists(key StoreKey) bool {

	_, ok := c.store[string(key)]

	return ok
}

// Set is used to store or update some store with a key
func (c *sessionCache) Set(key StoreKey, dat []byte) error {
	c.store[string(key)] = dat
	if d, ok := c.done[string(key)]; !ok || !d {
		c.keys = append(c.keys, string(key))
		c.done[string(key)] = true
	}
	return nil
}

// Delete removes any store stored against a key
func (c *sessionCache) Delete(key StoreKey) (bool, error) {

	tombstoneBytes := []byte(TOMBSTONE)
	c.store[string(key)] = tombstoneBytes
	if d, ok := c.done[string(key)]; !ok || !d {
		c.keys = append(c.keys, string(key))
		c.done[string(key)] = true
	}
	return true, nil
}

func (c *sessionCache) GetIterable() Iterable {
	return c
}

func (c *sessionCache) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
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
func (c *sessionCache) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	panic("IterateRange not implemented for sessionCache kv")
}

func (c *sessionCache) BeginSession() Session {
	return &cacheSession{
		parent: c,
		store:  map[string][]byte{},
		keys:   make([]string, 0, 10),
		done:   map[string]bool{},
	}
}

func (c *sessionCache) Close() {
	// pass
	// no resources to close
}

/*
	cacheSession
*/
type cacheSession struct {
	parent *sessionCache
	store  map[string][]byte
	keys   []string
	done   map[string]bool
}

func (c *cacheSession) Get(key StoreKey) ([]byte, error) {
	d, ok := c.store[string(key)]
	if !ok {
		return nil, ErrNotFound
	}

	return d, nil
}

func (c *cacheSession) Set(key StoreKey, dat []byte) error {

	c.store[string(key)] = dat
	if d, ok := c.done[string(key)]; !ok || !d {
		c.keys = append(c.keys, string(key))
		c.done[string(key)] = true
	}
	return nil
}

func (c *cacheSession) Exists(key StoreKey) bool {

	_, ok := c.store[string(key)]

	return ok
}

func (c *cacheSession) Delete(key StoreKey) (bool, error) {

	tombstoneBytes := []byte(TOMBSTONE)
	c.store[string(key)] = tombstoneBytes
	if d, ok := c.done[string(key)]; !ok || !d {
		c.keys = append(c.keys, string(key))
		c.done[string(key)] = true
	}
	return true, nil
}

func (c *cacheSession) GetIterable() Iterable {
	return c
}

func (c *cacheSession) Commit() bool {

	var err error

	for _, k := range c.keys {
		v, ok := c.store[k]
		if !ok {
			continue
		}
		err = c.parent.Set(StoreKey(k), v)
		if err != nil {
			return false
		}
	}

	return true
}

func (c *cacheSession) Iterate(fn func(key []byte, value []byte) bool) (stopped bool) {
	for k, v := range c.store {

		if fn([]byte(k), v) {
			return true
		}
	}
	return true
}

//
func (c *cacheSession) IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) (stop bool) {
	panic("IterateRange not implemented for cacheSession kv")
}
