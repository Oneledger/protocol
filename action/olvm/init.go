package olvm

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {

	serialize.RegisterConcrete(new(Transaction), "olvm")

}

func EnableOLVM(r action.Router) error {

	err := r.AddHandler(action.OLVM, olvmTx{})
	if err != nil {
		return errors.Wrap(err, "olvmTx")
	}
	return nil
}
