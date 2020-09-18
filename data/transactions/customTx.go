package transactions

import (
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Oneledger/protocol/storage"
)

func (ts *TransactionStore) ExistsCustom(customKey string, id string) bool {
	key := storage.StoreKey(customKey + storage.DB_PREFIX + string(id))
	return ts.State.Exists(key)
}
func (ts *TransactionStore) DeleteCustom(customKey string, id string) (bool, error) {
	key := storage.StoreKey(customKey + storage.DB_PREFIX + string(id))
	return ts.Delete(key)
}

func (ts *TransactionStore) AddCustom(customKey string, id string, tx *abci.RequestDeliverTx) error {
	key := storage.StoreKey(customKey + storage.DB_PREFIX + string(id))
	err := ts.Set(tx, key)
	if err != nil {
		return err
	}
	return nil
}
func (ts *TransactionStore) GetCustom(customKey string, id string) (*abci.RequestDeliverTx, error) {
	key := storage.StoreKey(customKey + storage.DB_PREFIX + string(id))
	tx, err := ts.Get(key)
	if err != nil {
		return &abci.RequestDeliverTx{}, err
	}
	return tx, nil
}

func (ts *TransactionStore) IterateCustom(customKey string, fn func(key string, tx *abci.RequestDeliverTx) bool) bool {
	return ts.State.IterateRange(
		append(ts.prefix, customKey...),
		storage.Rangefix(string(ts.prefix)+customKey),
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
