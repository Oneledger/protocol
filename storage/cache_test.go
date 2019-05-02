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
	"testing"

	"github.com/Oneledger/protocol/data"
	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	assert.Equal(t, cache{"name", map[string][]byte{}}, NewCache("name"))
}

func TestCache_SetGet(t *testing.T) {
	c := NewCache("test")

	key := data.StoreKey("key")
	dat := []byte("data")
	assert.NoError(t, c.Set(key, dat))

	dat2, err := c.Get(key)
	assert.Equal(t, dat, dat2)
	assert.NoError(t, err)

	exists, err := c.Exists(key)
	assert.True(t, exists)
	assert.NoError(t, err)

	deleted, err := c.Delete(key)
	assert.True(t, deleted)
	assert.NoError(t, err)

	exists, err = c.Exists(key)
	assert.False(t, exists)
	assert.NoError(t, err)

	key2 := data.StoreKey("key2")
	assert.NoError(t, c.Set(key2, dat))

	exists, err = c.Exists(key)
	assert.False(t, exists)
	assert.NoError(t, err)

	exists, err = c.Exists(key2)
	assert.True(t, exists)
	assert.NoError(t, err)

}

