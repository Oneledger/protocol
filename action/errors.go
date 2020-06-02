package action

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrMissingData        = codes.ProtocolError{codes.TxErrMisingData, "missing data in transaction"}
	ErrUnserializable     = codes.ProtocolError{codes.TxErrUnserializable, "unserializable tx"}
	ErrWrongTxType        = codes.ProtocolError{codes.TxErrWrongTxType, "wrong tx type"}
	ErrInvalidAmount      = codes.ProtocolError{codes.TxErrInvalidAmount, "invalid amount"}
	ErrInvalidPubkey      = codes.ProtocolError{codes.TxErrInvalidPubKey, "invalid pubkey"}
	ErrUnmatchSigner      = codes.ProtocolError{codes.TxErrUnmatchedSigner, "unmatch signers"}
	ErrInvalidSignature   = codes.ProtocolError{codes.TxErrInvalidSignature, "invalid signatures"}
	ErrInvalidFeeCurrency = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid fees currency"}
	ErrInvalidFeePrice    = codes.ProtocolError{codes.TxErrInvalidFeePrice, "fee price is smaller than minimal fee"}
	ErrNotEnoughFund      = codes.ProtocolError{codes.TxErrInsufficientFunds, "not enough fund"}
	ErrGasOverflow        = codes.ProtocolError{codes.TxErrGasOverflow, "gas used exceed limit"}
	ErrInvalidExtTx       = codes.ProtocolError{codes.TxErrInvalidExtTx, "invalid external tx"}

	ErrInvalidAddress = codes.ErrBadAddress

	ErrInvalidCurrency   = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid amount"}
	ErrTokenNotSupported = codes.ProtocolError{codes.ExternalErrTokenNotSuported, "Token not supported"}

	ErrProposalNotFound            = codes.ProtocolError{Code: codes.ProposalNotFound, Msg: "Proposal not found in proposal Store"}
	ErrUnauthorizedCall            = codes.ProtocolError{Code: codes.UnauthorizedCall, Msg: "Caller not authorized to execute this TX"}
	ErrStatusNotCompleted          = codes.ProtocolError{Code: codes.StatusNotCompleted, Msg: "TX not in completed status"}
	ErrStatusNotVoting             = codes.ProtocolError{Code: codes.StatusNotVoting, Msg: "TX not in Voting status"}
	ErrStatusNotFunding            = codes.ProtocolError{Code: codes.StatusNotFunding, Msg: "TX not in Funding status"}
	ErrVotingTBD                   = codes.ProtocolError{Code: codes.VotingTBD, Msg: "Voting Decision not achieved"}
	ErrFinalizeDistributtionFailed = codes.ProtocolError{Code: codes.FinalizeDistributtionFailed, Msg: "Failed in distributing Funds"}
	ErrFinalizeConfigUpdateFailed  = codes.ProtocolError{Code: codes.FinalizeConfigUpdateFailed, Msg: "Failed to execute Config Update"}
	ErrStatusUnableToSetFinalized  = codes.ProtocolError{Code: codes.StatusUnableToSetFinalized, Msg: "Failed to set status to finalized"}
	ErrUnabletoQueryVoteResult     = codes.ProtocolError{Code: codes.UnabletoQueryVoteResult, Msg: "Unable to query Votestore to get Vote result"}
)
