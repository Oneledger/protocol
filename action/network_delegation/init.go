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
	err = r.AddHandler(action.NETWORK_UNDELEGATE, undelegateTx{})
	if err != nil {
		return errors.Wrap(err, "undelegateTx")
	}
	err = r.AddHandler(action.WITHDRAW_NETWORK_DELEGATION, withdrawNetworkDelegationTx{})
	if err != nil {
		return errors.Wrap(err, "withdrawNetworkDelegationTx")
	}
	err = r.AddHandler(action.REWARDS_WITHDRAW_NETWORK_DELEGATE, delegWithdrawRewardsTx{})
	if err != nil {
		return errors.Wrap(err, "delegWithdrawRewardsTx")
	}
	err = r.AddHandler(action.REWARDS_FINALIZE_NETWORK_DELEGATE, finalizeWithdrawRewardsTx{})
	if err != nil {
		return errors.Wrap(err, "WithdrawRewardsTx")
	}
	err = r.AddHandler(action.REWARDS_REINVEST_NETWORK_DELEGATE, delegReinvestRewardsTx{})
	if err != nil {
		return errors.Wrap(err, "ReinvestRewardsTx")
	}
	return nil
}
