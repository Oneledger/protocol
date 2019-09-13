package fees

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type Store struct {
	state    *storage.State
	prefix   []byte
	currency balance.Currency
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		state:  state,
		prefix: storage.Prefix(prefix),
	}
}

func (st *Store) WithState(state *storage.State) *Store {
	st.state = state
	return st
}

func (st *Store) SetupCurrency(currency balance.Currency) {
	st.currency = currency
}

func (st *Store) Get(address []byte) (coin *balance.Coin, err error) {
	key := append(st.prefix, storage.StoreKey(address)...)
	dat, _ := st.state.Get(key)
	coin = &balance.Coin{}
	if len(dat) == 0 {
		*coin = st.currency.NewCoinFromInt(0)
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, coin)
	return
}

func (st *Store) Set(address keys.Address, coin balance.Coin) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(&coin)
	if err != nil {
		return err
	}

	key := append(st.prefix, storage.StoreKey(address)...)
	err = st.state.Set(key, dat)
	return err
}

func (st *Store) Exists(address keys.Address) bool {
	key := append(st.prefix, storage.StoreKey(address)...)
	return st.state.Exists(key)
}

func (st *Store) Iterate(fn func(addr keys.Address, coin balance.Coin) bool) bool {
	return st.state.IterateRange(
		st.prefix,
		storage.Rangefix(string(st.prefix[:len(st.prefix)-1])),
		true,
		func(key, value []byte) bool {
			coin := &balance.Coin{}
			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, coin)
			if err != nil {
				return false
			}
			return fn(key, *coin)
		},
	)
}
