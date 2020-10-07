package network_delegation

import (
	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
)

func EnableNetworkDelegation(r action.Router) error {
	err := r.AddHandler(action.ADD_NETWORK_DELEGATION, addNetworkDelegationTx{})
	if err != nil {
		return errors.Wrap(err, "AddNetworkDelegation")
	}

	return nil
}
