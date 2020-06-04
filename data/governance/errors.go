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
	ErrWithdrawCheckFundsFailed    = codes.ProtocolError{Code: codes.GovErrWithdrawCheckFundsFailed, Msg: "ErrWithdraw, failed to check available funds to withdraw for this contributor"}
	ErrProposalNotFound            = codes.ProtocolError{Code: codes.ProposalNotFound, Msg: "Proposal not found in proposal Store"}
	ErrUnauthorizedCall            = codes.ProtocolError{Code: codes.UnauthorizedCall, Msg: "Caller not authorized to execute this TX"}
	ErrStatusNotCompleted          = codes.ProtocolError{Code: codes.StatusNotCompleted, Msg: "TX not in completed status"}
	ErrStatusNotVoting             = codes.ProtocolError{Code: codes.StatusNotVoting, Msg: "TX not in Voting status"}
	ErrStatusNotFunding            = codes.ProtocolError{Code: codes.StatusNotFunding, Msg: "TX not in Funding status"}
	ErrVotingTBD                   = codes.ProtocolError{Code: codes.VotingTBD, Msg: "Voting Decision not achieved"}
	ErrFinalizeDistributtionFailed = codes.ProtocolError{Code: codes.FinalizeDistributtionFailed, Msg: "Failed in distributing Funds"}
	ErrFinalizeConfigUpdateFailed  = codes.ProtocolError{Code: codes.FinalizeConfigUpdateFailed, Msg: "Failed to execute Config Update"}
	ErrStatusUnableToSetFinalized  = codes.ProtocolError{Code: codes.StatusUnableToSetFinalized, Msg: "Failed to set status to finalized"}
	ErrStatusUnableToSetVoting     = codes.ProtocolError{Code: codes.StatusUnableToSetVoting, Msg: "Failed to set status to voting"}
	ErrUnabletoQueryVoteResult     = codes.ProtocolError{Code: codes.UnabletoQueryVoteResult, Msg: "Unable to query Votestore to get Vote result"}
	ErrGovFundUnableToAdd          = codes.ProtocolError{Code: codes.GovFundUnableToAdd, Msg: "Funding unable to add funds"}
	ErrGovFundUnableToDelete       = codes.ProtocolError{Code: codes.GovFundUnableToDelete, Msg: "Funding unable to delete Funds"}
	ErrFundingDeadlineCrossed      = codes.ProtocolError{Code: codes.FundingDeadlineCrossed, Msg: "Funding deadline has been crossed"}
)
