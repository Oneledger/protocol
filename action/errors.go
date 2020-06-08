package action

import (
	"encoding/json"
	codes "github.com/Oneledger/protocol/status_codes"
)

// this struct is to build a unified structure for wrapped and unwrapped errors
// so that it would be easy for sdk/client side to handles errors
type ErrorToReturn struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
}

func ErrorMarshal(code int, msg string) string {
	errStruct := ErrorToReturn{
		Code: code,
		Msg: msg,
	}
	errInByte, _ := json.Marshal(errStruct)
	return string(errInByte)
}

var (
	ErrMissingData                  = codes.ProtocolError{codes.TxErrMisingData, "missing data in transaction"}
	ErrUnserializable               = codes.ProtocolError{codes.TxErrUnserializable, "unserializable tx"}
	ErrWrongTxType                  = codes.ProtocolError{codes.TxErrWrongTxType, "wrong tx type"}
	ErrInvalidAmount                = codes.ProtocolError{codes.TxErrInvalidAmount, "invalid amount"}
	ErrInvalidPubkey                = codes.ProtocolError{codes.TxErrInvalidPubKey, "invalid pubkey"}
	ErrUnmatchSigner                = codes.ProtocolError{codes.TxErrUnmatchedSigner, "unmatch signers"}
	ErrInvalidSignature             = codes.ProtocolError{codes.TxErrInvalidSignature, "invalid signatures"}
	ErrInvalidFeeCurrency           = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid fees currency"}
	ErrInvalidFeePrice              = codes.ProtocolError{codes.TxErrInvalidFeePrice, "fee price is smaller than minimal fee"}
	ErrNotEnoughFund                = codes.ProtocolError{codes.TxErrInsufficientFunds, "not enough fund"}
	ErrGasOverflow                  = codes.ProtocolError{codes.TxErrGasOverflow, "gas used exceed limit"}
	ErrInvalidExtTx                 = codes.ProtocolError{codes.TxErrInvalidExtTx, "invalid external tx"}
	ErrInvalidAddress               = codes.ErrBadAddress

	ErrInvalidCurrency              = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid amount"}
	ErrTokenNotSupported            = codes.ProtocolError{codes.ExternalErrTokenNotSuported, "Token not supported"}

	ErrInvalidProposalId            = codes.ProtocolError{codes.GovErrInvalidProposalId, "invalid proposal id"}
	ErrInvalidProposalType          = codes.ProtocolError{codes.GovErrInvalidProposalType, "invalid proposal type"}
	ErrGetProposalOptions           = codes.ProtocolError{codes.GovErrGetProposalOptions, "failed to get proposal options"}
	ErrInvalidProposerAddr          = codes.ProtocolError{codes.GovErrInvalidProposerAddr, "invalid proposer address"}
	ErrInvalidProposalDesc          = codes.ProtocolError{codes.GovErrInvalidProposalDesc, "invalid description of proposal"}
	ErrProposalExists               = codes.ProtocolError{codes.GovErrProposalExists, "proposal already exists"}
	ErrProposalUnmarshal            = codes.ProtocolError{codes.GovErrProposalUnmarshal, "failed to unmarshal proposal"}
	ErrDeductFunding                = codes.ProtocolError{codes.GovErrDeductFunding, "failed to deduct funds from address"}
	ErrAddFunding                   = codes.ProtocolError{codes.GovErrAddFunding, "failed to add funds to address"}
	ErrFundingHeightReached         = codes.ProtocolError{codes.GovErrFundingHeightReached, "funding height has already been reached"}
	ErrNotInFunding                 = codes.ProtocolError{codes.GovErrNotInFunding, "proposal not in FUNDING stage"}
	ErrGettingValidatorList         = codes.ProtocolError{codes.GovErrGettingValidatorList, "fund proposal failed in getting validator list"}
	ErrSetupVotingValidator         = codes.ProtocolError{codes.GovErrSetupVotingValidator, "failed to setup voting validator"}
	ErrAddingProposalToActiveDB     = codes.ProtocolError{codes.GovErrAddingProposalToActiveDB, "failed to add proposal to ACTIVE store"}
	ErrDeletingProposalFromActiveDB = codes.ProtocolError{codes.GovErrDeletingProposalFromActiveDB, "failed to delet proposal from ACTIVE db"}
	ErrAddingProposalToPassedDB     = codes.ProtocolError{codes.GovErrAddingProposalToPassedDB, "failed to add proposal to PASSED db"}
	ErrDeletingProposalFromPassedDB = codes.ProtocolError{codes.GovErrDeletingProposalFromPassedDB, "failed to delet proposal from PASSED db"}
	ErrAddingProposalToFailedDB     = codes.ProtocolError{codes.GovErrAddingProposalToFailedDB, "failed to add proposal to FAILED db"}
	ErrDeletingProposalFromFailedDB = codes.ProtocolError{codes.GovErrDeletingProposalFromFailedDB, "failed to delet proposal from FAILED db"}
	ErrInvalidContributorAddr       = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
	ErrProposalWithdrawNotEligible  = codes.ProtocolError{codes.GovErrProposalWithdrawNotEligible, "proposal does not meet withdraw requirement"}
	ErrNoSuchContributor            = codes.ProtocolError{codes.GovErrNoSuchContributor, "this contributor has not funded this proposal"}
	ErrNotInVoting                  = codes.ProtocolError{codes.GovErrNotInVoting, "proposal not in VOTING status"}
	ErrVotingHeightReached          = codes.ProtocolError{codes.GovErrVotingHeightReached, "voting height has already been reached"}
	ErrAddingVoteToVoteStore        = codes.ProtocolError{codes.GovErrAddingVoteToVoteStore, "failed to add vote to vote store"}
	ErrPeekingVoteResult            = codes.ProtocolError{codes.GovErrPeekingVoteResult, "failed to peek vote result"}

)
