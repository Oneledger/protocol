package farm_data

import (
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

var _ data.ExtStore = &ProductStore{}

type ProductStore struct {
	state *storage.State
	szlr  serialize.Serializer
	prefix []byte
}

func NewProductStore(state *storage.State, prefix string) *ProductStore {
	return &ProductStore{
		state: state,
		szlr: serialize.GetSerializer(serialize.PERSISTENT),
		prefix: []byte(prefix),
	}
}

func (ps *ProductStore) GetState() *storage.State {
	return ps.state
}

func (ps *ProductStore) WithState(state *storage.State) data.ExtStore {
	ps.state = state
	return ps
}

func (ps *ProductStore) Set(product *Product) error {
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

func (ps *ProductStore) Get(batchId BatchID) (*Product, error) {
	product := &Product{}
	prefixed := append(ps.prefix, batchId...)
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

func (ps *ProductStore) Delete(batchId BatchID) (bool, error) {
	prefixed := append(ps.prefix, batchId...)
	res, err := ps.state.Delete(prefixed)
	if err != nil {
		return false, ErrDeletingRecord.Wrap(err)
	}
	return res, err
}
