package identity

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

var isETHWitness bool

type EthWitnessStore struct {
	prefix []byte
	store  *storage.State
}

func NewEthWitnessStore(prefix string, state *storage.State) *EthWitnessStore {
	return &EthWitnessStore{
		prefix: storage.Prefix(prefix),
		store:  state,
	}
}

func (ws *EthWitnessStore) WithState(state *storage.State) *EthWitnessStore {
	ws.store = state
	return ws
}

func (ws *EthWitnessStore) Init(nodeValidatorAddress keys.Address) {
	isETHWitness = ws.Exists(nodeValidatorAddress)
}

func (ws *EthWitnessStore) Get(addr keys.Address) (*EthWitness, error) {
	key := append(ws.prefix, addr...)
	value, _ := ws.store.Get(key)
	if value == nil {
		return nil, errors.New("failed to get ethereum witness from store")
	}
	witness := &EthWitness{}
	witness, err := witness.FromBytes(value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize ethereum witness")
	}
	return witness, nil
}

func (ws *EthWitnessStore) Exists(addr keys.Address) bool {
	key := append(ws.prefix, addr...)
	return ws.store.Exists(key)
}

func (ws *EthWitnessStore) Iterate(fn func(addr keys.Address, witness *EthWitness) bool) (stopped bool) {
	return ws.store.IterateRange(
		ws.prefix,
		storage.Rangefix(string(ws.prefix)),
		true,
		func(key, value []byte) bool {
			witness, err := (&EthWitness{}).FromBytes(value)
			if err != nil {
				logger.Error("failed to deserialize ethereum witness")
				return false
			}
			addr := key[len(ws.prefix):]
			return fn(addr, witness)
		},
	)
}

// Get Ethereum witness addresses
func (ws *EthWitnessStore) GetETHWitnessAddresses() ([]keys.Address, error) {
	witnessList := make([]keys.Address, 0)
	ws.Iterate(func(addr keys.Address, witness *EthWitness) bool {
		witnessList = append(witnessList, addr)
		return false
	})
	return witnessList, nil
}

// This node is a ethereum witness or not
func (ws *EthWitnessStore) IsETHWitness() bool {
	return isETHWitness
}

func (ws *EthWitnessStore) IsETHWitnessAddress(addr keys.Address) bool {
	return ws.Exists(addr)
}

// Add a ethereum witness to store
func (ws *EthWitnessStore) AddWitness(apply Stake) error {
	if ws.Exists(apply.ValidatorAddress) {
		return nil
	}

	witness := &EthWitness{
		Address:     apply.ValidatorAddress,
		PubKey:      apply.Pubkey,
		ECDSAPubKey: apply.ECDSAPubKey,
		Name:        apply.Name,
	}

	value := witness.Bytes()
	vkey := append(ws.prefix, witness.Address.Bytes()...)
	err := ws.store.Set(vkey, value)
	if err != nil {
		return errors.Wrap(err, "failed to add ethereum witness")
	}

	return nil
}
