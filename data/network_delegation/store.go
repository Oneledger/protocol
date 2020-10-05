package network_delegation

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"strconv"
	"strings"
	"sync"
)

const (
	MatureKey  = "m"
	PendingKey = "p"
	ActiveKey  = "a"
)

var (
	matureType  DelegationPrefixType = 0x103
	pendingType DelegationPrefixType = 0x102
	activeType  DelegationPrefixType = 0x101
)

type Store struct {
	state         *storage.State
	szlr          serialize.Serializer
	prefix        []byte
	currentPrefix []byte
	mux           sync.Mutex
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		state:         state,
		prefix:        storage.StoreKey(prefix),
		currentPrefix: storage.StoreKey(prefix + storage.DB_PREFIX + ActiveKey),
		szlr:          serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (st *Store) WithState(state *storage.State) *Store {
	st.state = state
	return st
}

func (st *Store) Exists(addr *keys.Address) bool {
	key := append(st.currentPrefix, addr.String()...)
	return st.state.Exists(key)
}

//Set coin to specific key
func (st *Store) set(key []byte, coin *balance.Coin) (err error) {
	dat, err := st.szlr.Serialize(coin)
	if err != nil {
		return
	}
	err = st.state.Set(storage.StoreKey(key), dat)
	return
}

//get coin from specific key
func (st *Store) get(key []byte) (coin *balance.Coin, err error) {
	coin = &balance.Coin{}
	dat, err := st.state.Get(storage.StoreKey(key))
	if err != nil {
		return
	}
	err = st.szlr.Deserialize(dat, coin)
	return
}

//Set coin to specific address
func (st *Store) Set(address keys.Address, coin *balance.Coin) (err error) {
	prefixKey := append(st.currentPrefix, address.String()...)
	err = st.set(storage.StoreKey(prefixKey), coin)
	return
}

//Get coin from address
func (st *Store) Get(address keys.Address) (coin *balance.Coin, err error) {
	prefixKey := append(st.currentPrefix, address.String()...)
	coin, err = st.get(storage.StoreKey(prefixKey))
	return
}

func (st *Store) WithPrefix(prefix DelegationPrefixType) *Store {
	switch prefix {
	case matureType:
		st.currentPrefix = st.buildMatureKey()
	case pendingType:
		st.currentPrefix = st.buildPendingKey()
	case activeType:
		st.currentPrefix = st.buildActiveKey()
	}
	return st
}

func (st *Store) iterate(prefix storage.StoreKey, fn func(key []byte, coin *balance.Coin) bool) bool {
	return st.state.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			coin := &balance.Coin{}
			err := st.szlr.Deserialize(value, coin)
			if err != nil {
				return true
			}
			return fn(key, coin)
		},
	)
}

func (st *Store) iterateAddresses(prefix storage.StoreKey, fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
	return st.iterate(prefix, func(key []byte, coin *balance.Coin) bool {
		arr := strings.Split(string(key), storage.DB_PREFIX)
		addr := &keys.Address{}
		err := addr.UnmarshalText([]byte(arr[len(arr)-1]))
		if err != nil {
			return true
		}
		return fn(addr, coin)
	})
}

//------------------------------- Active key store -------------------------------
//build Active key
func (st *Store) buildActiveKey() storage.StoreKey {
	return storage.Prefix(string(st.prefix) + storage.DB_PREFIX + ActiveKey)
}

func (st *Store) IterateActiveAmounts(fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
	prefix := st.buildActiveKey()
	return st.iterateAddresses(prefix, func(addr *keys.Address, coin *balance.Coin) bool {
		return fn(addr, coin)
	})
}

//------------------------------- Pending key store -------------------------------
//build pending key
func (st *Store) buildPendingKey() storage.StoreKey {
	return storage.Prefix(string(st.prefix) + storage.DB_PREFIX + PendingKey)
}

//Set pending amount with height and address
func (st *Store) SetPendingAmount(addr keys.Address, height int64, coin *balance.Coin) error {
	prefix := st.buildPendingKey()
	pendingKey := strconv.FormatInt(height, 10) + storage.DB_PREFIX + addr.String()
	return st.set(append(prefix, pendingKey...), coin)
}

//iterate addresses for height
func (st *Store) IteratePendingAmounts(height int64, fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
	prefix := append(st.buildPendingKey(), strconv.FormatInt(height, 10)...)
	return st.iterateAddresses(prefix, func(addr *keys.Address, coin *balance.Coin) bool {
		return fn(addr, coin)
	})
}

//iterate all pending amounts
func (st *Store) IterateAllPendingAmounts(fn func(height int64, addr *keys.Address, coin *balance.Coin) bool) bool {
	prefix := st.buildPendingKey()
	return st.iterate(prefix, func(key []byte, coin *balance.Coin) bool {
		arr := strings.Split(string(key), storage.DB_PREFIX)
		addr := &keys.Address{}
		err := addr.UnmarshalText([]byte(arr[len(arr)-1]))
		if err != nil {
			return true
		}
		height, err := strconv.Atoi(arr[len(arr)-2])
		if err != nil {
			return true
		}
		return fn(int64(height), addr, coin)
	})
}

//------------------------------- Mature key store -------------------------------

//build mature key
func (st *Store) buildMatureKey() storage.StoreKey {
	return storage.Prefix(string(st.prefix) + storage.DB_PREFIX + MatureKey)
}

//iterate all matured amounts
func (st *Store) IterateMatureAmounts(fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
	prefix := st.buildMatureKey()
	return st.iterateAddresses(prefix, func(addr *keys.Address, coin *balance.Coin) bool {
		return fn(addr, coin)
	})
}
