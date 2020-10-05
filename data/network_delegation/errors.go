package network_delegation

import codes "github.com/Oneledger/protocol/status_codes"

var (
	ErrFinalizingDelgRewards         = codes.ProtocolError{codes.NetDelgErrFinalizingDelgRewards, "failed to finalize delegation rewards"}
	ErrAddingWithdrawAmountToBalance = codes.ProtocolError{codes.NetDelgErrAddingWithdrawAmountToBalance, "failed to add withdrawn amount to delegator's balance"}
)
