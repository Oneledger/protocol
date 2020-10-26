package evidence

import (
	"fmt"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type ValidatorStatus struct {
	Address  keys.Address `json:"address"`
	IsActive bool         `json:"isActive"`
	Height   int64        `json:"height"`
}

func (vs *ValidatorStatus) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(vs)
	if err != nil {
		fmt.Println("validator status not serializable", err)
		return []byte{}
	}
	return value
}

func (vs *ValidatorStatus) FromBytes(msg []byte) (*ValidatorStatus, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, vs)
	if err != nil {
		fmt.Println("failed to deserialize validator status from bytes", err)
		return nil, err
	}
	return vs, nil
}

func (es *EvidenceStore) getValidatorStatusKey(validatorAddress keys.Address) []byte {
	key := []byte(fmt.Sprintf("_vss_%s", validatorAddress))
	return key
}

func (es *EvidenceStore) GetValidatorMap() map[string]bool {
	vMap := make(map[string]bool)
	es.IterateValidatorStatuses(func(vs *ValidatorStatus) bool {
		if vs.IsActive {
			vMap[vs.Address.String()] = vs.IsActive
		}
		return false
	})
	return vMap
}

func (es *EvidenceStore) IterateValidatorStatuses(fn func(vs *ValidatorStatus) bool) (stopped bool) {
	key := []byte("_vss_")
	prefixKey := append(es.prefix, key...)
	return es.state.IterateRange(
		prefixKey,
		storage.Rangefix(string(prefixKey)),
		true,
		func(key, value []byte) bool {
			vs, err := (&ValidatorStatus{}).FromBytes(value)
			if err != nil {
				fmt.Println("failed to deserialize validator status")
				return false
			}
			return fn(vs)
		},
	)
}

func (es *EvidenceStore) IsActiveValidator(addr keys.Address) bool {
	vs, err := es.GetValidatorStatus(addr)
	if err != nil {
		return false
	}
	return vs.IsActive
}

func (es *EvidenceStore) GetValidatorStatus(addr keys.Address) (*ValidatorStatus, error) {
	key := es.getValidatorStatusKey(addr)
	data, err := es.Get(key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("failed to get validator status from store")
	}
	valStatus := &ValidatorStatus{}
	valStatus, err = valStatus.FromBytes(data)
	if err != nil {
		return nil, errors.Wrap(err, "error deserialize validator")
	}
	return valStatus, nil
}

func (es *EvidenceStore) SetValidatorStatus(addr keys.Address, isActive bool, height int64) error {
	key := es.getValidatorStatusKey(addr)
	validatorStatus := ValidatorStatus{
		Address:  addr,
		IsActive: isActive,
		Height:   height,
	}
	value := (validatorStatus).Bytes()
	err := es.Set(key, value)
	if err != nil {
		return errors.Wrap(err, "failed to set validator status")
	}
	return nil
}
