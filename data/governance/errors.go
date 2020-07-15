package governance

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	// Options Objects from store
	ErrGetProposalOptions = codes.ProtocolError{codes.GovErrGetProposalOptions, "failed to get proposal options"}
	ErrGetFeeOptions      = codes.ProtocolError{codes.TxErrGettingFeeOptions, "failed to get fee options"}
	ErrGetRewardOptions   = codes.ProtocolError{codes.TxErrGettingRewardsOptions, "failed to get Reward options"}
	ErrGetONSOptions      = codes.ProtocolError{codes.TxErrGettingONSOptions, "failed to get ONS options"}
	ErrGetEthOptions      = codes.ProtocolError{codes.TxErrGettingETHOptions, "failed to get ETH options"}
	ErrGetBtcOptions      = codes.ProtocolError{codes.TxErrGettingBTCOptions, "failed to get BTC options"}
	ErrValidateGovState   = codes.ProtocolError{codes.TxErrValidateGovState, "failed to validate gov state"}

	//Proposal
	ErrInvalidProposalId      = codes.ProtocolError{codes.GovErrInvalidProposalId, "invalid proposal id"}
	ErrInvalidProposalType    = codes.ProtocolError{codes.GovErrInvalidProposalType, "invalid proposal type"}
	ErrInvalidProposalDesc    = codes.ProtocolError{codes.GovErrInvalidProposalDesc, "invalid description of proposal"}
	ErrProposalExists         = codes.ProtocolError{codes.GovErrProposalExists, "proposal already exists"}
	ErrProposalNotExists      = codes.ProtocolError{codes.GovErrProposalNotExists, "proposal not exists"}
	ErrInvalidBeneficiaryAddr = codes.ProtocolError{codes.GovErrInvalidBeneficiaryAddr, "invalid withdraw beneficiary address"}
	ErrInvalidFundingGoal     = codes.ProtocolError{codes.GovErrInvalidFundingGoal, "invalid funding goal"}
	ErrInvalidPassPercentage  = codes.ProtocolError{codes.GovErrInvalidPassPercentage, "invalid pass percentage"}
	ErrInvalidFundingDeadline = codes.ProtocolError{codes.GovErrInvalidFundingDeadline, "invalid funding deadline"}
	ErrInvalidVotingDeadline  = codes.ProtocolError{codes.GovErrInvalidVotingDeadline, "invalid voting deadline"}

	//Funding
	ErrDeductFunding          = codes.ProtocolError{codes.GovErrDeductFunding, "failed to deduct funds from address"}
	ErrAddFunding             = codes.ProtocolError{codes.GovErrAddFunding, "failed to add funds to address"}
	ErrInvalidFunderAddr      = codes.ProtocolError{codes.GovErrInvalidFunderAddr, "invalid funder address"}
	ErrStatusNotFunding       = codes.ProtocolError{Code: codes.GovErrStatusNotFunding, Msg: "TX not in Funding status"}
	ErrGovFundUnableToAdd     = codes.ProtocolError{Code: codes.GovFundUnableToAdd, Msg: "Funding unable to add funds"}
	ErrGovFundUnableToDelete  = codes.ProtocolError{Code: codes.GovFundUnableToDelete, Msg: "Funding unable to delete Funds"}
	ErrFundingDeadlineCrossed = codes.ProtocolError{Code: codes.GovErrFundingDeadlineCrossed, Msg: "Funding deadline has been crossed"}

	//Voting
	ErrSetupVotingValidator        = codes.ProtocolError{codes.GovErrSetupVotingValidator, "failed to setup voting validator"}
	ErrStatusNotVoting             = codes.ProtocolError{Code: codes.GovErrStatusNotVoting, Msg: "TX not in Voting status"}
	ErrVotingHeightReached         = codes.ProtocolError{codes.GovErrVotingHeightReached, "voting height has already been reached"}
	ErrAddingVoteToVoteStore       = codes.ProtocolError{codes.GovErrAddingVoteToVoteStore, "failed to add vote to vote store"}
	ErrPeekingVoteResult           = codes.ProtocolError{codes.GovErrPeekingVoteResult, "failed to peek vote result"}
	ErrInvalidVoteOpinion          = codes.ProtocolError{codes.GovErrInvalidVoteOpinion, "invalid vote opinion"}
	ErrVoteDeleteVoteRecordsFailed = codes.ProtocolError{Code: codes.GovErrVoteDeleteVoteRecords, Msg: "ErrVote, failed to delete voting records"}
	ErrVoteCheckVoteResultFailed   = codes.ProtocolError{Code: codes.GovErrVoteCheckVoteResult, Msg: "ErrVote, failed to check voting result"}
	ErrStatusUnableToSetVoting     = codes.ProtocolError{Code: codes.StatusUnableToSetVoting, Msg: "Failed to set status to voting"}
	ErrUnabletoQueryVoteResult     = codes.ProtocolError{Code: codes.GovErrUnabletoQueryVoteResult, Msg: "Unable to query Votestore to get Vote result"}
	ErrVotingTBD                   = codes.ProtocolError{Code: codes.GovErrVotingTBD, Msg: "Voting Decision not achieved"}

	//Finalizing
	ErrStatusNotCompleted              = codes.ProtocolError{Code: codes.GovErrStatusNotCompleted, Msg: "TX not in completed status"}
	ErrFinalizeDistributtionFailed     = codes.ProtocolError{Code: codes.GovErrFinalizeDistributtionFailed, Msg: "Failed in distributing Funds"}
	ErrFinalizeConfigUpdateFailed      = codes.ProtocolError{Code: codes.GovErrFinalizeConfigUpdateFailed, Msg: "Failed to execute Config Update"}
	ErrStatusUnableToSetFinalized      = codes.ProtocolError{Code: codes.GovErrStatusUnableToSetFinalized, Msg: "Failed to set status to finalized"}
	ErrStatusUnableToSetFinalizeFailed = codes.ProtocolError{Code: codes.GovErrUnableToSetFinalizeFailed, Msg: "Failed to set status to finalized Failed"}
	ErrGovFundBalanceMismatch          = codes.ProtocolError{Code: codes.GovFundBalanceMismatch, Msg: "Balance Mismatch While Burning Funds"}

	//Withdraw
	ErrProposalWithdrawNotEligible = codes.ProtocolError{codes.GovErrProposalWithdrawNotEligible, "proposal does not meet withdraw requirement"}
	ErrNoSuchFunder                = codes.ProtocolError{codes.GovErrNoSuchFunder, "no such funder funded this proposal"}
	ErrUnmatchedProposer           = codes.ProtocolError{codes.GovErrUnmatchedProposer, "proposer does not match"}
	ErrInvalidVoterId              = codes.ProtocolError{codes.GovErrInvalidVoterId, "invalid voter id"}
	ErrWithdrawCheckFundsFailed    = codes.ProtocolError{Code: codes.GovErrWithdrawCheckFundsFailed, Msg: "ErrWithdraw, failed to check available funds to withdraw for this contributor"}

	//Proposal Prefix
	ErrAddingProposalToActiveStore     = codes.ProtocolError{codes.GovErrAddingProposalToActiveStore, "failed to add proposal to ACTIVE store"}
	ErrDeletingProposalFromActiveStore = codes.ProtocolError{codes.GovErrDeletingProposalFromActiveStore, "failed to delet proposal from ACTIVE store"}
	ErrAddingProposalToPassedStore     = codes.ProtocolError{codes.GovErrAddingProposalToPassedStore, "failed to add proposal to PASSED store"}
	ErrAddingProposalToFailedStore     = codes.ProtocolError{codes.GovErrAddingProposalToFailedStore, "failed to add proposal to FAILED store"}
	ErrDeletingProposalFromFailedStore = codes.ProtocolError{codes.GovErrDeletingProposalFromFailedStore, "failed to delet proposal from FAILED store"}
)
