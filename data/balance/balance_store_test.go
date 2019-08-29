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

package balance

import (
	"testing"

	"github.com/tendermint/tendermint/libs/db"

	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	olt := currencies["OLT"]
	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("balance", db))
	store := NewStore("b", cs)

	bal := NewBalance()
	bal.AddCoin(olt.NewCoinFromInt(10))

	err := store.Set([]byte("asdfasdfasdfasdfasdf"), *bal)
	assert.NoError(t, err)

	bal2, err := store.Get([]byte("asdfasdfasdfasdfasdf"))
	assert.NoError(t, err)
	assert.Equal(t, bal, bal2)

	bal2, err = store.Get([]byte("asdfasdfasdfasdfhjkl"))
	assert.Error(t, err)
	assert.NotEqual(t, bal, bal2)

	//assert.True(t, store.Exists([]byte("asdfasdfasdfasdfasdf")))
	assert.False(t, store.Exists([]byte("asdfasdfasdfasdfhjkl")))

	store.State.Commit()
	cnt := 0
	store.State.Iterate(func(key, value []byte) bool {
		cnt++
		return false
	})
	assert.Equal(t, 1, cnt)
}
