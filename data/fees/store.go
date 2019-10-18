package fees

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type Store struct {
	state  *storage.State
	prefix []byte
	feeOpt *FeeOption
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

func (st *Store) SetupOpt(feeOpt *FeeOption) {
	st.feeOpt = feeOpt
}

func (st *Store) GetOpt() *FeeOption {
	return st.feeOpt
}

func (st *Store) Get(address []byte) (coin balance.Coin, err error) {
	key := append(st.prefix, storage.StoreKey(address)...)
	dat, _ := st.state.Get(key)
	a := balance.NewAmount(0)
	if len(dat) == 0 {
		return st.feeOpt.FeeCurrency.NewCoinFromInt(0), nil
	}

	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, a)
	if err != nil {
		return st.feeOpt.FeeCurrency.NewCoinFromInt(0), nil
	}

	coin = st.feeOpt.FeeCurrency.NewCoinFromAmount(*a)
	return
}

func (st *Store) Set(address keys.Address, coin balance.Coin) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(&(coin.Amount))
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

func (st *Store) Iterate(fn func(addr keys.Address, coin balance.Coin) (stop bool)) bool {
	return st.state.IterateRange(
		st.prefix,
		storage.Rangefix(string(st.prefix)),
		true,
		func(key, value []byte) bool {
			amt := &balance.Amount{}
			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, amt)
			if err != nil {
				return false
			}
			coin := st.feeOpt.FeeCurrency.NewCoinFromAmount(*amt)
			addr := key[len(st.prefix):]
			return fn(addr, coin)
		},
	)
}

func (st *Store) AddToAddress(addr keys.Address, coin balance.Coin) error {

	baseCoin, err := st.Get(addr)
	if err != nil {
		return errors.Wrapf(err, "failed to get address balance %s", addr.String())
	}

	newCoin := baseCoin.Plus(coin)

	return st.Set(addr, newCoin)
}

func (st *Store) MinusFromAddress(addr keys.Address, coin balance.Coin) error {
	baseCoin, err := st.Get(addr)
	if err != nil {
		return errors.Wrapf(err, "failed to get address balance %s", addr.String())
	}

	newCoin, err := baseCoin.Minus(coin)
	if err != nil {
		return err
	}

	return st.Set(addr, newCoin)
}

func (st *Store) AddToPool(coin balance.Coin) error {
	return st.AddToAddress(keys.Address(POOL_KEY), coin)
}

func (st *Store) MinusFromPool(coin balance.Coin) error {
	return st.MinusFromAddress(keys.Address(POOL_KEY), coin)
}

func (st *Store) GetAllowedWithdraw(addr keys.Address) balance.Coin {
	prefixed := append(st.prefix, addr...)

	data := st.state.GetPrevious(FEE_LOCK_BLOCKS, prefixed)
	amt := balance.Amount{}
	coin := balance.Coin{}
	if len(data) == 0 {
		coin = st.feeOpt.FeeCurrency.NewCoinFromInt(0)
		return coin
	}
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(data, &amt)
	if err != nil {
		coin = st.feeOpt.FeeCurrency.NewCoinFromInt(0)
	}
	return st.feeOpt.FeeCurrency.NewCoinFromAmount(amt)
}
