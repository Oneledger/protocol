package network_delegation

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrGettingOptions     = codes.ProtocolError{codes.NetwkDelgErrGetOptions, "failed to get network delegation options"}
	ErrInitiateWithdrawal = codes.ProtocolError{codes.NetwkDelgErrWithdraw, "failed to initiate rewards withdrawal"}
)
