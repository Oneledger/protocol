/*

 */

package storage

/* sessionCache is a simple in-memory keyvalue store, to store binary  This is not thread safe and
any concurrent read/write might throw panics.
*/
type sessionCache struct {
	name  string
	store map[string][]byte
	keys  []string
}

// sessionCache satisfies SessionedDirectStorage interface
var _ SessionedDirectStorage = &sessionCache{}
var _ Iteratable = &sessionCache{}
var _ SessionedDirectStorage = &sessionCache{}

func NewSessionCache(name string) *sessionCache {
	return &sessionCache{
		name:  name,
		store: make(map[string][]byte),
		keys:  make([]string, 0, 100),
	}
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
	c.keys = append(c.keys, string(key))
	return nil
}

// Delete removes any store stored against a key
func (c *sessionCache) Delete(key StoreKey) (bool, error) {

	tombstoneBytes := []byte(TOMBSTONE)
	c.store[string(key)] = tombstoneBytes

	c.keys = append(c.keys, string(key))
	return true, nil
}

func (c *sessionCache) GetIterator() Iteratable {
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
	return nil
}

func (c *cacheSession) Exists(key StoreKey) bool {

	_, ok := c.store[string(key)]

	return ok
}

func (c *cacheSession) Delete(key StoreKey) (bool, error) {

	tombstoneBytes := []byte(TOMBSTONE)
	c.store[string(key)] = tombstoneBytes

	return true, nil
}

func (c *cacheSession) GetIterator() Iteratable {
	return c
}

func (c *cacheSession) Commit() bool {

	var err error

	for k, v := range c.store {
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
