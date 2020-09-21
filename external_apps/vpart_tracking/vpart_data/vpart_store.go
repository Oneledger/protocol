package vpart_data

import (
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type VPartStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte
}

func NewVPartStore(state *storage.State, szlr serialize.Serializer, prefix string) *VPartStore {
	return &VPartStore{
		state:  state,
		szlr:   szlr,
		prefix: []byte(prefix),
	}
}

func (vs *VPartStore) GetState() *storage.State {
	return vs.state
}

func (vs *VPartStore) WithState(state *storage.State) data.ExtStore {
	vs.state = state
	return vs
}

func assembleVPartKey(vin Vin, partType string) storage.StoreKey {
	return storage.StoreKey(string(vin) + storage.DB_PREFIX + partType)
}

func (vs *VPartStore) Set(part *VPart) error {
	key := assembleVPartKey(part.VIN, part.PartType)
	prefixed := append(vs.prefix, key...)
	data, err := vs.szlr.Serialize(part)
	if err != nil {
		return ErrFailedInSerialization.Wrap(err)
	}

	err = vs.state.Set(prefixed, data)

	if err != nil {
		return ErrSettingRecord.Wrap(err)
	}
	return nil
}

func (vs *VPartStore) Get(vin Vin, partType string) (*VPart, error) {
	key := assembleVPartKey(vin, partType)
	prefixed := append(vs.prefix, key...)

	part := &VPart{}
	data, err := vs.state.Get(prefixed)
	if err != nil {
		return nil, ErrGettingRecord.Wrap(err)
	}
	err = vs.szlr.Deserialize(data, part)
	if err != nil {
		return nil, ErrFailedInDeserialization.Wrap(err)
	}
	return part, nil
}

func (vs *VPartStore) Exists(vin Vin, partType string) bool {
	key := assembleVPartKey(vin, partType)
	prefixed := append(vs.prefix, key...)
	return vs.state.Exists(prefixed)
}

func (vs *VPartStore) Delete(vin Vin, partType string) (bool, error) {
	key := assembleVPartKey(vin, partType)
	prefixed := append(vs.prefix, key...)
	res, err := vs.state.Delete(prefixed)
	if err != nil {
		return false, ErrDeletingRecord.Wrap(err)
	}
	return res, err
}