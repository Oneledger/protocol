package network_delegation

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
)

func EnableNetworkDelegation(r action.Router) error {
	err := r.AddHandler(action.WITHDRAW_DELEG_REWARD, delegWithdrawTx{})
	if err != nil {
		return errors.Wrap(err, "delegWithdrawRewardsTx")
	}
	return nil
}
