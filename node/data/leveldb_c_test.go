// +build gcc

package data

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/db"

	"testing"
)

// TestCLevelDB is just a basic test to see if cleveldb is being used properly
func TestCLevelDB(t *testing.T) {
	storage, err := db.NewCLevelDB("test123", global.Current.DatabaseDir())
	if err != nil {
		panic(err)
	}
	want := []byte("world")
	storage.Set([]byte("hello"), want)

	got := storage.Get([]byte("hello"))

	assert.Equal(t, got, want, "Get and set should be the same")
}
