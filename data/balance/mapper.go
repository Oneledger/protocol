package balance

import (
	"encoding/binary"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

const (
	INM = iota + 1
	OUTM
)

type AccountJoiner struct {
	Address   keys.Address   `json:"address"`
	Algorithm keys.Algorithm `json:"algorithm"`
	Enabled   bool           `json:"enabled"`
}

type AccountMapper struct {
	state  *storage.State
	prefix []byte
}

func NewAccountMapper(state *storage.State) *AccountMapper {
	return &AccountMapper{
		state:  state,
		prefix: storage.Prefix("mapper"),
	}
}

func (ac *AccountMapper) WithState(state *storage.State) *AccountMapper {
	ac.state = state
	return ac
}

func (ac *AccountMapper) Get(addr keys.Address, direction uint64) (*AccountJoiner, error) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, direction)
	prefixKey := append(ac.prefix, b...)
	prefixKey = append(ac.prefix, addr.Bytes()...)

	dat, err := ac.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil, err
	}

	am := &AccountJoiner{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, am)
	if err != nil {
		return nil, err
	}
	return am, nil
}
