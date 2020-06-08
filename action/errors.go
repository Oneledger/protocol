package action

import (
	codes "github.com/Oneledger/protocol/status_codes"
)


var (
	ErrMissingData                     = codes.ProtocolError{codes.TxErrMissingData, "missing data in transaction"}
	ErrUnserializable                  = codes.ProtocolError{codes.TxErrUnserializable, "unserializable tx"}
	ErrWrongTxType                     = codes.ProtocolError{codes.TxErrWrongTxType, "wrong tx type"}
	ErrInvalidAmount                   = codes.ProtocolError{codes.TxErrInvalidAmount, "invalid amount"}
	ErrInvalidPubkey                   = codes.ProtocolError{codes.TxErrInvalidPubKey, "invalid pubkey"}
	ErrUnmatchSigner                   = codes.ProtocolError{codes.TxErrUnmatchedSigner, "unmatch signers"}
	ErrInvalidSignature                = codes.ProtocolError{codes.TxErrInvalidSignature, "invalid signatures"}
	ErrInvalidFeeCurrency              = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid fees currency"}
	ErrInvalidFeePrice                 = codes.ProtocolError{codes.TxErrInvalidFeePrice, "fee price is smaller than minimal fee"}
	ErrNotEnoughFund                   = codes.ProtocolError{codes.TxErrInsufficientFunds, "not enough fund"}
	ErrGasOverflow                     = codes.ProtocolError{codes.TxErrGasOverflow, "gas used exceed limit"}
	ErrInvalidExtTx                    = codes.ProtocolError{codes.TxErrInvalidExtTx, "invalid external tx"}
	ErrInvalidAddress                  = codes.ErrBadAddress

	ErrInvalidCurrency                 = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid amount"}
	ErrTokenNotSupported               = codes.ProtocolError{codes.ExternalErrTokenNotSuported, "Token not supported"}

	// governance error
	ErrInvalidProposalId               = codes.ProtocolError{codes.GovErrInvalidProposalId, "invalid proposal id"}
	ErrInvalidProposalType             = codes.ProtocolError{codes.GovErrInvalidProposalType, "invalid proposal type"}
	ErrGetProposalOptions              = codes.ProtocolError{codes.GovErrGetProposalOptions, "failed to get proposal options"}
	ErrInvalidProposerAddr             = codes.ProtocolError{codes.GovErrInvalidProposerAddr, "invalid proposer address"}
	ErrInvalidProposalDesc             = codes.ProtocolError{codes.GovErrInvalidProposalDesc, "invalid description of proposal"}
	ErrProposalExists                  = codes.ProtocolError{codes.GovErrProposalExists, "proposal already exists"}
	ErrProposalNotExists               = codes.ProtocolError{codes.GovErrProposalNotExists, "proposal not exists"}
	ErrDeductFunding                   = codes.ProtocolError{codes.GovErrDeductFunding, "failed to deduct funds from address"}
	ErrAddFunding                      = codes.ProtocolError{codes.GovErrAddFunding, "failed to add funds to address"}
	ErrFundingHeightReached            = codes.ProtocolError{codes.GovErrFundingHeightReached, "funding height has already been reached"}
	ErrNotInFunding                    = codes.ProtocolError{codes.GovErrNotInFunding, "proposal not in FUNDING stage"}
	ErrGettingValidatorList            = codes.ProtocolError{codes.GovErrGettingValidatorList, "fund proposal failed in getting validator list"}
	ErrSetupVotingValidator            = codes.ProtocolError{codes.GovErrSetupVotingValidator, "failed to setup voting validator"}
	ErrAddingProposalToActiveStore     = codes.ProtocolError{codes.GovErrAddingProposalToActiveStore, "failed to add proposal to ACTIVE store"}
	ErrDeletingProposalFromActiveStore = codes.ProtocolError{codes.GovErrDeletingProposalFromActiveStore, "failed to delet proposal from ACTIVE store"}
	ErrAddingProposalToPassedStore     = codes.ProtocolError{codes.GovErrAddingProposalToPassedStore, "failed to add proposal to PASSED store"}
	ErrDeletingProposalFromPassedStore = codes.ProtocolError{codes.GovErrDeletingProposalFromPassedStore, "failed to delet proposal from PASSED store"}
	ErrAddingProposalToFailedStore     = codes.ProtocolError{codes.GovErrAddingProposalToFailedStore, "failed to add proposal to FAILED store"}
	ErrDeletingProposalFromFailedStore = codes.ProtocolError{codes.GovErrDeletingProposalFromFailedStore, "failed to delet proposal from FAILED store"}
	ErrInvalidContributorAddr          = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
	ErrProposalWithdrawNotEligible     = codes.ProtocolError{codes.GovErrProposalWithdrawNotEligible, "proposal does not meet withdraw requirement"}
	ErrNoSuchContributor               = codes.ProtocolError{codes.GovErrNoSuchContributor, "no such contributor funded this proposal"}
	ErrNotInVoting                     = codes.ProtocolError{codes.GovErrNotInVoting, "proposal not in VOTING status"}
	ErrVotingHeightReached             = codes.ProtocolError{codes.GovErrVotingHeightReached, "voting height has already been reached"}
	ErrAddingVoteToVoteStore           = codes.ProtocolError{codes.GovErrAddingVoteToVoteStore, "failed to add vote to vote store"}
	ErrPeekingVoteResult               = codes.ProtocolError{codes.GovErrPeekingVoteResult, "failed to peek vote result"}
	ErrUnmatchedProposer               = codes.ProtocolError{codes.GovErrUnmatchedProposer, "proposer does not match"}

)
