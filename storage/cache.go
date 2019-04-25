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
	"sync"
)


type cache struct {
	name string
	store map[string][]byte
}

func (c *cache) Get(key data.StoreKey) ([]byte, error) {

	d, ok := c.store[string(key)]
	if !ok {
		return nil, ErrNotFound
	}

	return d, nil
}

func (c *cache) Exists(key data.StoreKey) (bool, error) {

	_, ok := c.store[string(key)]

	return ok, nil
}

func (c *cache) Set(key data.StoreKey, dat []byte) (error) {

	c.store[string(key)] = dat

	return nil
}

func (c *cache) Delete(key data.StoreKey) (bool, error) {

	delete(c.store, string(key))
	return true, nil
}


/*
	CacheSafe starts here
 */


type cacheSafe struct {
	name string
	store map[string][]byte
	lock sync.RWMutex
}

func (c *cacheSafe) Get(key data.StoreKey) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	d, ok := c.store[string(key)]
	if !ok {
		return nil, ErrNotFound
	}

	return d, nil
}

func (c *cacheSafe) Exists(key data.StoreKey) (bool, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.store[string(key)]

	return ok, nil
}

func (c *cacheSafe) Set(key data.StoreKey, dat []byte) (error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.store[string(key)] = dat

	return nil
}

func (c *cacheSafe) Delete(key data.StoreKey) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.store, string(key))
	return true, nil
}
