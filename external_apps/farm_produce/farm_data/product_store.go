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

	return ErrSettingRecord.Wrap(err)
}
