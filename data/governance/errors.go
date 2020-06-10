package governance

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	//Vote Error
	ErrVoteSetupValidatorFailed    = codes.ProtocolError{Code: codes.GovErrVoteSetupValidator, Msg: "ErrVote, failed to setup voting validator"}
	ErrVoteUpdateVoteFailed        = codes.ProtocolError{Code: codes.GovErrVoteUpdateVote, Msg: "ErrVote, failed to vote proposal"}
	ErrVoteDeleteVoteRecordsFailed = codes.ProtocolError{Code: codes.GovErrVoteDeleteVoteRecords, Msg: "ErrVote, failed to delete voting records"}
	ErrVoteCheckVoteResultFailed   = codes.ProtocolError{Code: codes.GovErrVoteCheckVoteResult, Msg: "ErrVote, failed to check voting result"}
	ErrWithdrawCheckFundsFailed    = codes.ProtocolError{Code: codes.GovErrWithdrawCheckFundsFailed, Msg: "ErrWithdraw, failed to check available funds to withdraw for this funder"}
)
