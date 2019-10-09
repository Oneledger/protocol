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
	olt := Currency{
		Id:      0,
		Name:    "OLT",
		Chain:   0,
		Decimal: 18,
		Unit:    "nue",
	}
	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("balance", db))
	store := NewStore("b", cs)
	currencies := NewCurrencySet()
	err := currencies.Register(olt)
	assert.NoError(t, err)

	coin := olt.NewCoinFromInt(10)

	err = store.AddToAddress([]byte("asdfasdfasdfasdfasdf"), coin)
	assert.NoError(t, err)
	cs.Commit()

	bal, err := store.GetBalance([]byte("asdfasdfasdfasdfasdf"), currencies)
	assert.NoError(t, err)
	assert.Equal(t, coin, bal.GetCoin(olt))

	coin2, err := coin.Minus(olt.NewCoinFromInt(4))
	assert.NoError(t, err)
	err = store.CheckBalanceFromAddress([]byte("asdfasdfasdfasdfasdf"), coin2)
	assert.NoError(t, err)

	coin3 := coin.Plus(olt.NewCoinFromInt(1))
	err = store.CheckBalanceFromAddress([]byte("asdfasdfasdfasdfasdf"), coin3)
	assert.Error(t, err)

	err = store.MinusFromAddress([]byte("asdfasdfasdfasdfasdf"), coin2)
	assert.NoError(t, err)

	err = store.MinusFromAddress([]byte("asdfasdfasdfasdfasdf"), coin2)
	assert.Error(t, err)

	store.State.Commit()
	cnt := 0
	store.State.Iterate(func(key, value []byte) bool {
		cnt++
		return false
	})
	assert.Equal(t, 1, cnt)
}
