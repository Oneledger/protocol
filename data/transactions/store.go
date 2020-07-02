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

Description: Local store of pending transactions triggered and executed internally through the application.

*Note: Not part of the chain state

*/

package transactions

import (
	"errors"
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

const (
	FINALIZE_KEY = "FINALIZE"
	EXPIRE_KEY   = "EXPIRE"
)

type TransactionStore struct {
	State  *storage.State
	szlr   serialize.Serializer
	prefix []byte
}

func NewTransactionStore(prefix string, state *storage.State) *TransactionStore {
	return &TransactionStore{
		State:  state,
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
		prefix: storage.Prefix(prefix),
	}
}

func (ts *TransactionStore) WithState(state *storage.State) *TransactionStore {
	ts.State = state
	return ts
}

func (ts *TransactionStore) Set(tx *abci.RequestDeliverTx, key string) error {
	storeKey := append(ts.prefix, storage.StoreKey(storage.DB_PREFIX+key)...)
	data, err := ts.szlr.Serialize(tx)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("error serializing transaction")
	}
	return ts.State.Set(storeKey, data)
}

func (ts *TransactionStore) Get(key string) (tx *abci.RequestDeliverTx, err error) {
	storeKey := append(ts.prefix, storage.StoreKey(storage.DB_PREFIX+key)...)
	data, err := ts.State.Get(storeKey)
	tx = &abci.RequestDeliverTx{}
	if len(data) == 0 {

		return tx, errors.New("key doesn't exist")
	}
	err = ts.szlr.Deserialize(data, tx)
	return
}

//Iterate through all Transactions
func (ts *TransactionStore) Iterate(fn func(key string, tx *abci.RequestDeliverTx) bool) bool {
	return ts.State.IterateRange(
		append(ts.prefix),
		storage.Rangefix(string(ts.prefix)),
		true,
		func(key, value []byte) bool {
			tx := &abci.RequestDeliverTx{}

			err := ts.szlr.Deserialize(value, tx)
			if err != nil {
				return true
			}

			arr := strings.Split(string(key), storage.DB_PREFIX)
			return fn(arr[len(arr)-1], tx)
		},
	)
}

func (ts *TransactionStore) Delete(key string) (bool, error) {
	storeKey := append(ts.prefix, storage.StoreKey(storage.DB_PREFIX+key)...)
	return ts.State.Delete(storeKey)
}

func (ts *TransactionStore) Exists(key string) bool {
	storeKey := append(ts.prefix, storage.StoreKey(storage.DB_PREFIX+key)...)
	return ts.State.Exists(storeKey)
}
