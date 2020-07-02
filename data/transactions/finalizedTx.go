package transactions

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"strings"

	"github.com/Oneledger/protocol/storage"
)

func (ts *TransactionStore) ExistsFinalized(id string) bool {
	key := storage.StoreKey(FINALIZE_KEY + storage.DB_PREFIX + string(id))
	return ts.State.Exists(key)
}
func (ts *TransactionStore) DeleteFinalized(id string) (bool, error) {
	key := storage.StoreKey(FINALIZE_KEY + storage.DB_PREFIX + string(id))
	return ts.Delete(key)
}

func (ts *TransactionStore) AddFinalized(id string, tx *abci.RequestDeliverTx) error {
	key := storage.StoreKey(FINALIZE_KEY + storage.DB_PREFIX + string(id))
	err := ts.Set(tx, key)
	if err != nil {
		return err
	}
	return nil
}
func (ts *TransactionStore) GetFinalized(id string) (*abci.RequestDeliverTx, error) {
	key := storage.StoreKey(FINALIZE_KEY + storage.DB_PREFIX + string(id))
	tx, err := ts.Get(key)
	if err != nil {
		return &abci.RequestDeliverTx{}, err
	}
	return tx, nil
}

func (ts *TransactionStore) IterateFinalized(fn func(key string, tx *abci.RequestDeliverTx) bool) bool {
	return ts.State.IterateRange(
		append(ts.prefix, FINALIZE_KEY...),
		storage.Rangefix(string(ts.prefix)+FINALIZE_KEY),
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
