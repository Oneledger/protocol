package balance

import (
	"encoding/binary"
	"sync"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

const (
	INM = iota + 1
	OUTM
)

type AccountWithAlg struct {
	Address   keys.Address   `json:"address"`
	Algorithm keys.Algorithm `json:"algorithm"`
}

type AccountJoiner struct {
	Legacy  AccountWithAlg `json:"legacy"`
	New     AccountWithAlg `json:"new"`
	Enabled bool           `json:"enabled"`
}

type AccountMapper struct {
	state  *storage.State
	prefix []byte
	mux    sync.Mutex
}

func NewAccountMapper(state *storage.State) *AccountMapper {
	return &AccountMapper{
		state:  state,
		prefix: storage.Prefix("mapper"),
	}
}

func (am *AccountMapper) WithState(state *storage.State) *AccountMapper {
	am.state = state
	return am
}

func (am *AccountMapper) getPrefix(addr keys.Address, algorithm keys.Algorithm) storage.StoreKey {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(algorithm))
	prefixKey := append(am.prefix, b...)
	prefixKey = append(am.prefix, addr.Bytes()...)
	return prefixKey
}

func (am *AccountMapper) Get(addr keys.Address, algorithm keys.Algorithm) (*AccountJoiner, error) {
	dat, err := am.state.Get(am.getPrefix(addr, algorithm))
	if err != nil {
		return nil, err
	}

	aj := &AccountJoiner{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, aj)
	if err != nil {
		return nil, err
	}
	return aj, nil
}

func (am *AccountMapper) Set(aj *AccountJoiner) error {
	am.mux.Lock()
	defer am.mux.Unlock()

	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(aj)
	if err != nil {
		return errors.Errorf("Failed to serialize: %s", err)
	}
	err = am.state.Set(am.getPrefix(aj.Legacy.Address, aj.Legacy.Algorithm), dat)
	if err != nil {
		return errors.Errorf("Failed to update storage for legacy mapper: %s", err)
	}
	err = am.state.Set(am.getPrefix(aj.New.Address, aj.New.Algorithm), dat)
	if err != nil {
		return errors.Errorf("Failed to update storage for new mapper: %s", err)
	}
	return nil
}

func (am *AccountMapper) GetOrCreateED25519ToETHECDSA(addr0 keys.Address, addr1 keys.Address) (*AccountJoiner, error) {
	var aj *AccountJoiner

	dat, err := am.state.Get(am.getPrefix(addr0, OUTM))
	if err == nil {
		if len(dat) == 0 {
			awa0 := AccountWithAlg{
				Address:   addr0,
				Algorithm: keys.ED25519,
			}
			awa1 := AccountWithAlg{
				Address:   addr1,
				Algorithm: keys.ETHSECP,
			}
			aj = &AccountJoiner{
				Legacy:  awa0,
				New:     awa1,
				Enabled: true,
			}
			err = am.Set(aj)
			if err != nil {
				return nil, err
			}
			return aj, nil
		}
		aj = &AccountJoiner{}
		err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, aj)
		if err != nil {
			return nil, err
		}
		return aj, nil
	}
	return nil, err
}
