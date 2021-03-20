package smart_contract

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {

	serialize.RegisterConcrete(new(Execute), "smart_contract_execute")

}

func EnableSmartContract(r action.Router) error {

	err := r.AddHandler(action.SC_EXECUTE, scExecuteTx{})
	if err != nil {
		return errors.Wrap(err, "scExecuteTx")
	}
	return nil
}
