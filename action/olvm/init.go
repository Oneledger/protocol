package olvm

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

const HandlerName = "olvm"

func init() {

	serialize.RegisterConcrete(new(Transaction), HandlerName)

}

func EnableOLVM(r action.Router) error {

	err := r.AddHandler(action.OLVM, olvmTx{})
	if err != nil {
		return errors.Wrap(err, "olvmTx")
	}
	return nil
}
