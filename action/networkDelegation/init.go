package networkDelegation

import (
	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
)

func EnableNetworkDelegation(r action.Router) error {
	err := r.AddHandler(action.NETWORKDELEGATE, networkDelegateTx{})
	if err != nil {
		return errors.Wrap(err, "NetworkDelegate")
	}

	return nil
}
