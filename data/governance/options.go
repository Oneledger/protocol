package governance

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/serialize"
)

type Options interface {
	SetOptions(*interface{}, string) error
	GetOptions(string) (*interface{}, error)
	ValidateOptions(string)
}

var OptionList map[string]interface{}
var Validation map[string]func(interface{}) bool

func IntializeOptions() {
	OptionList[ADMIN_STAKING_OPTION] = delegation.Options{}
	Validation[ADMIN_STAKING_OPTION] = validateStalking
}

func (st *Store) GetOptions(key string) (interface{}, error) {

	bytes, err := st.Get(key)
	if err != nil {
		return nil, err
	}
	optionsType := OptionList[key]
	r := &optionsType
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize options")
	}

	return r, nil
}

func (st *Store) SetOptions(opt *interface{}, key string) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(opt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize options")
	}

	err = st.Set(key, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set options")
	}

	return nil
}

func (st *Store) ValidateOptions(key string, options *interface{}) {
	validationFunction := Validation[key]
	validationFunction(*options)

}

func validateStalking(opt interface{}) bool {
	stakingOptions, ok := opt.(delegation.Options)
	fmt.Println(stakingOptions)
	fmt.Println(ok)
	return true
}
