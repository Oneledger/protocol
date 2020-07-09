package rewards

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	UnableToGetMaturedAmount = codes.ProtocolError{codes.RewardsUnableToGetMaturedAmount, "failed to get matured balance"}
	UnableToWithdraw         = codes.ProtocolError{codes.RewardsUnableToWithdraw, "Unable to withdraw"}
)
