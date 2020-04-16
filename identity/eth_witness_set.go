package identity

import (
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

var isETHWitness bool

type WitnessStore struct {
	prefix []byte
	store  *storage.State
}

func NewWitnessStore(prefix string, state *storage.State) *WitnessStore {
	return &WitnessStore{
		prefix: storage.Prefix(prefix),
		store:  state,
	}
}

func (ws *WitnessStore) WithState(state *storage.State) *WitnessStore {
	ws.store = state
	return ws
}

func (ws *WitnessStore) Init(chain chain.Type, nodeValidatorAddress keys.Address) {
	isETHWitness = ws.Exists(chain, nodeValidatorAddress)
}

func (ws *WitnessStore) Get(chain chain.Type, addr keys.Address) (*Witness, error) {
	key := append(ws.prefix, addr...)
	value, _ := ws.store.Get(key)
	if value == nil {
		return nil, errors.New("failed to get witness from store")
	}
	witness := &Witness{}
	witness, err := witness.FromBytes(value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize witness")
	}

	if witness.Chain != chain {
		return nil, errors.Wrap(err, "witness does not match")
	}

	return witness, nil
}

func (ws *WitnessStore) Exists(chain chain.Type, addr keys.Address) bool {
	_, err := ws.Get(chain, addr)
	return err == nil
}

func (ws *WitnessStore) Iterate(fn func(addr keys.Address, witness *Witness) bool) (stopped bool) {
	return ws.store.IterateRange(
		ws.prefix,
		storage.Rangefix(string(ws.prefix)),
		true,
		func(key, value []byte) bool {
			witness, err := (&Witness{}).FromBytes(value)
			if err != nil {
				logger.Error("failed to deserialize witness")
				return false
			}
			addr := key[len(ws.prefix):]
			return fn(addr, witness)
		},
	)
}

// Get witness addresses
func (ws *WitnessStore) GetWitnessAddresses(chain chain.Type) ([]keys.Address, error) {
	witnessList := make([]keys.Address, 0)
	ws.Iterate(func(addr keys.Address, witness *Witness) bool {
		if witness.Chain == chain {
			witnessList = append(witnessList, addr)
		}
		return false
	})
	return witnessList, nil
}

// This node is a ethereum witness or not
func (ws *WitnessStore) IsETHWitness() bool {
	return isETHWitness
}

func (ws *WitnessStore) IsWitnessAddress(chain chain.Type, addr keys.Address) bool {
	return ws.Exists(chain, addr)
}

// Add a witness to store
func (ws *WitnessStore) AddWitness(chain chain.Type, apply Stake) error {
	if ws.Exists(chain, apply.ValidatorAddress) {
		return nil
	}

	witness := &Witness{
		Address:     apply.ValidatorAddress,
		PubKey:      apply.Pubkey,
		ECDSAPubKey: apply.ECDSAPubKey,
		Name:        apply.Name,
		Chain:       chain,
	}

	value := witness.Bytes()
	vkey := append(ws.prefix, witness.Address.Bytes()...)
	err := ws.store.Set(vkey, value)
	if err != nil {
		return errors.Wrap(err, "failed to add witness")
	}

	return nil
}
