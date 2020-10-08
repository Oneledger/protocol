package network_delegation

import (
	"github.com/Oneledger/protocol/action"
	"github.com/pkg/errors"
)

func EnableNetworkDelegation(r action.Router) error {
	err := r.AddHandler(action.ADD_NETWORK_DELEGATE, addNetworkDelegationTx{})
	if err != nil {
		return errors.Wrap(err, "AddNetworkDelegation")
	}
	err = r.AddHandler(action.NETWORK_UNDELEGATE, UndelegateTx{})
	if err != nil {
		return errors.Wrap(err, "UndelegateTx")
	}
	err = r.AddHandler(action.WITHDRAW_NETWORK_DELEGATE, delegWithdrawTx{})
	if err != nil {
		return errors.Wrap(err, "delegWithdrawRewardsTx")
	}
	return nil
}
