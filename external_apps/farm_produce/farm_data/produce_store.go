package farm_data

import (
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

var _ data.ExtStore = &ProduceStore{}

type ProduceStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte
}

func NewProduceStore(state *storage.State, prefix string) *ProduceStore {
	return &ProduceStore{
		state:  state,
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
		prefix: []byte(prefix),
	}
}

func (ps *ProduceStore) GetState() *storage.State {
	return ps.state
}

func (ps *ProduceStore) WithState(state *storage.State) data.ExtStore {
	ps.state = state
	return ps
}

func (ps *ProduceStore) Exists(key BatchID) bool {
	prefix := append(ps.prefix, key...)
	return ps.state.Exists(prefix)
}

func (ps *ProduceStore) Set(product *Produce) error {
	prefixed := append(ps.prefix, product.BatchID...)
	data, err := ps.szlr.Serialize(product)
	if err != nil {
		return ErrFailedInSerialization.Wrap(err)
	}

	err = ps.state.Set(prefixed, data)

	if err != nil {
		return ErrSettingRecord.Wrap(err)
	}
	return nil
}

func (ps *ProduceStore) Get(batchId BatchID) (*Produce, error) {
	product := &Produce{}
	prefixed := append(ps.prefix, []byte(batchId)...)
	data, err := ps.state.Get(prefixed)
	if err != nil {
		return nil, ErrGettingRecord.Wrap(err)
	}
	err = ps.szlr.Deserialize(data, product)
	if err != nil {
		return nil, ErrFailedInDeserialization.Wrap(err)
	}

	return product, nil
}

func (ps *ProduceStore) Delete(batchId BatchID) (bool, error) {
	prefixed := append(ps.prefix, batchId...)
	res, err := ps.state.Delete(prefixed)
	if err != nil {
		return false, ErrDeletingRecord.Wrap(err)
	}
	return res, err
}
