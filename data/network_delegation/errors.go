package network_delegation

import codes "github.com/Oneledger/protocol/status_codes"

var (
	ErrFinalizingDelgRewards         = codes.ProtocolError{codes.NetDelgErrFinalizingDelgRewards, "failed to finalize delegation rewards"}
	ErrAddingWithdrawAmountToBalance = codes.ProtocolError{codes.NetDelgErrAddingWithdrawAmountToBalance, "failed to add withdrawn amount to delegator's balance"}
	ErrGettingActiveDelgAmount   = codes.ProtocolError{codes.NetDelgErrGettingActiveDelgAmount, "failed to get active network delegation amount"}
	ErrDeductingActiveDelgAmount = codes.ProtocolError{codes.NetDelgErrDeductingActiveDelgAmount, "failed to deduct active network delegation amount"}
	ErrSettingActiveDelgAmount   = codes.ProtocolError{codes.NetDelgErrSettingActiveDelgAmount, "failed to set active network delegation amount"}
	ErrGettingDelgOption         = codes.ProtocolError{codes.NetDelgErrGettingDelgOption, "failed to get network delegation option from governance store"}
	ErrSettingPendingDelgAmount  = codes.ProtocolError{codes.NetDelgErrSettingPendingDelgAmount, "failed to set pending network delegation amount"}
	ErrGettingPendingDelgAmount  = codes.ProtocolError{codes.NetDelgErrGettingPendingDelgAmount, "failed to get pending network delegation amount"}
	ErrInitiateWithdrawal        = codes.ProtocolError{codes.NetDelgErrWithdraw, "failed to initiate rewards withdrawal"}
)
