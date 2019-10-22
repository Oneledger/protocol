package governance

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

const (
	ADMIN_INITIAL_KEY    string = "initial"
	ADMIN_CURRENCY_KEY   string = "currency"
	ADMIN_FEE_OPTION_KEY string = "feeopt"
)

type Store struct {
	state  *storage.State
	prefix []byte
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

func (st *Store) Get(key []byte) ([]byte, error) {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)

	return st.state.Get(prefixKey)
}

func (st *Store) Set(key []byte, value []byte) error {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)
	err := st.state.Set(prefixKey, value)
	return err
}

func (st *Store) Exists(key []byte) bool {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)
	return st.state.Exists(prefixKey)
}

func (st *Store) GetCurrencies() (balance.Currencies, error) {
	result, err := st.Get([]byte(ADMIN_CURRENCY_KEY))
	currencies := make(balance.Currencies, 0, 10)
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(result, &currencies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the currencies")
	}
	return currencies, nil
}

func (st *Store) SetCurrencies(currencies balance.Currencies) error {

	currenciesBytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(currencies)
	if err != nil {
		return errors.Wrap(err, "failed to serialize currencies")
	}
	err = st.Set([]byte(ADMIN_CURRENCY_KEY), currenciesBytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the currencies")
	}
	return nil
}

func (st *Store) GetFeeOption() (*fees.FeeOption, error) {
	feeOpt := &fees.FeeOption{}
	bytes, err := st.Get([]byte(ADMIN_FEE_OPTION_KEY))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get FeeOption")
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, feeOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize FeeOption stored")
	}

	return feeOpt, nil
}

func (st *Store) SetFeeOption(feeOpt fees.FeeOption) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(feeOpt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize FeeOption")
	}
	err = st.Set([]byte(ADMIN_FEE_OPTION_KEY), bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the FeeOption")
	}
	return nil
}

func (st *Store) Initiated() bool {
	_ = st.Set([]byte(ADMIN_INITIAL_KEY), []byte("initialed"))
	return true
}

func (st *Store) InitialChain() bool {
	data, err := st.Get([]byte(ADMIN_INITIAL_KEY))
	if err != nil {
		return true
	}
	if data == nil {
		return true
	}
	return false
}
