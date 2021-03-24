package evm

import (
	"encoding/binary"

	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

var (
	KeyPrefixCode       = []byte{0x01}
	KeyPrefixStorage    = []byte{0x02}
	KeyPrefixHeightHash = []byte{0x03}
	KeyPrefixLogs       = []byte{0x04}
)

type ContractStore struct {
	state  *storage.State
	prefix []byte
}

func NewContractStore(state *storage.State) *ContractStore {
	return &ContractStore{
		state:  state,
		prefix: storage.Prefix("contracts"),
	}
}

func (cs *ContractStore) WithState(state *storage.State) *ContractStore {
	cs.state = state
	return cs
}

func (cs *ContractStore) GetStoreKey(prefix []byte, key []byte) storage.StoreKey {
	prefixKey := append(cs.prefix, prefix...)
	prefixKey = append(prefixKey, key...)
	return storage.StoreKey(prefixKey)
}

func (cs *ContractStore) Get(prefix []byte, key []byte) ([]byte, error) {
	storeKey := cs.GetStoreKey(prefix, key)
	dat, err := cs.state.Get(storeKey)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (cs *ContractStore) Set(prefix []byte, key []byte, value []byte) error {
	storeKey := cs.GetStoreKey(prefix, key)
	err := cs.state.Set(storeKey, value)
	return err
}

func (cs *ContractStore) Delete(prefix []byte, key []byte) (bool, error) {
	storeKey := cs.GetStoreKey(prefix, key)
	return cs.state.Delete(storeKey)
}

func (cs *ContractStore) Iterate(prefix []byte, fn func(key []byte, value []byte) bool) (stop bool) {
	prefixKey := append(cs.prefix, prefix...)
	return cs.state.IterateRange(
		prefixKey,
		storage.Rangefix(string(prefixKey)),
		true,
		fn,
	)
}

// AddressStoragePrefix returns a prefix to iterate over a given account storage.
func AddressStoragePrefix(address ethcmn.Address) []byte {
	return append(KeyPrefixStorage, address.Bytes()...)
}

// HeightHashKey returns the key for the given chain epoch and height.
// The key will be composed in the following order:
//   key = prefix + bytes(height)
// This ordering facilitates the iteration by height for the EVM GetHashFn
// queries.
func HeightHashKey(height uint64) []byte {
	buf := make([]byte, 8)
	binary.PutVarint(buf, int64(height))
	return buf
}
