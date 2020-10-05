package network_delegation

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
)

func init() {
	serialize.RegisterConcrete(new(Undelegate), "undelegate")
}

func EnableNetwkDelegation(r action.Router) error {
	err := r.AddHandler(action.NETWORK_DELEGATION_UNDELEGATE, UndelegateTx{})
	if err != nil {
		return errors.Wrap(err, "UndelegateTx")
	}
	return nil
}


