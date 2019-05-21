// +build gcc

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

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/db"
)

// TestCLevelDB is just a basic test to see if cleveldb is being used properly
func TestCLevelDB(t *testing.T) {
	storage, err := db.NewCLevelDB("test123", dbDir())
	if err != nil {
		panic(err)
	}
	want := []byte("world")
	storage.Set([]byte("hello"), want)

	got := storage.Get([]byte("hello"))

	assert.Equal(t, got, want, "Get and set should be the same")
}
