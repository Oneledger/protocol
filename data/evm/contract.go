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

func (cs *ContractStore) Get(key []byte) ([]byte, error) {
	prefixKey := append(cs.prefix, key...)

	dat, err := cs.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (cs *ContractStore) Set(key []byte, value []byte) error {
	prefixKey := append(cs.prefix, key...)
	err := cs.state.Set(storage.StoreKey(prefixKey), value)
	return err
}

func (cs *ContractStore) Delete(key []byte) (bool, error) {
	prefixed := append(cs.prefix, key...)
	return cs.state.Delete(prefixed)
}

// AddressStoragePrefix returns a prefix to iterate over a given account storage.
func AddressStoragePrefix(address ethcmn.Address) []byte {
	return append(KeyPrefixStorage, address.Bytes()...)
}

// CodeStoragePrefix returns a prefix to iterate over a given contract storage.
func CodeStoragePrefix(code []byte) []byte {
	return append(KeyPrefixCode, code...)
}

// HeightHashKey returns the key for the given chain epoch and height.
// The key will be composed in the following order:
//   key = prefix + bytes(height)
// This ordering facilitates the iteration by height for the EVM GetHashFn
// queries.
func HeightHashKey(height uint64) []byte {
	buf := make([]byte, 8)
	binary.PutVarint(buf, int64(height))
	return append(KeyPrefixHeightHash, buf...)
}
