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
	ErrMissingData            = codes.ProtocolError{codes.TxErrMisingData, "missing data in transaction"}
	ErrUnserializable         = codes.ProtocolError{codes.TxErrUnserializable, "unserializable tx"}
	ErrWrongTxType            = codes.ProtocolError{codes.TxErrWrongTxType, "wrong tx type"}
	ErrInvalidAmount          = codes.ProtocolError{codes.TxErrInvalidAmount, "invalid amount"}
	ErrInvalidPubkey          = codes.ProtocolError{codes.TxErrInvalidPubKey, "invalid pubkey"}
	ErrUnmatchSigner          = codes.ProtocolError{codes.TxErrUnmatchedSigner, "unmatch signers"}
	ErrInvalidSignature       = codes.ProtocolError{codes.TxErrInvalidSignature, "invalid signatures"}
	ErrInvalidFeeCurrency     = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid fees currency"}
	ErrInvalidFeePrice        = codes.ProtocolError{codes.TxErrInvalidFeePrice, "fee price is smaller than minimal fee"}
	ErrNotEnoughFund          = codes.ProtocolError{codes.TxErrInsufficientFunds, "not enough fund"}
	ErrGasOverflow            = codes.ProtocolError{codes.TxErrGasOverflow, "gas used exceed limit"}
	ErrInvalidExtTx           = codes.ProtocolError{codes.TxErrInvalidExtTx, "invalid external tx"}
	ErrInvalidAddress         = codes.ErrBadAddress

	ErrInvalidCurrency        = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid amount"}
	ErrTokenNotSupported      = codes.ProtocolError{codes.ExternalErrTokenNotSuported, "Token not supported"}

	ErrInvalidProposalId      = codes.ProtocolError{codes.GovErrInvalidProposalId, "invalid proposal id"}
	ErrInvalidProposalType    = codes.ProtocolError{codes.GovErrInvalidProposalType, "invalid proposal type"}
	ErrGetProposalOptions     = codes.ProtocolError{codes.GovErrGetProposalOptions, "failed to get proposal options"}
	ErrInvalidProposerAddr    = codes.ProtocolError{codes.GovErrInvalidProposerAddr, "invalid proposer address"}
	ErrInvalidProposalDesc    = codes.ProtocolError{codes.GovErrInvalidProposalDesc, "invalid description of proposal"}
	ErrProposalExists         = codes.ProtocolError{codes.GovErrProposalExists, "proposal already exists"}
	ErrAddingProposalToDB     = codes.ProtocolError{codes.GovErrAddingProposalToDB, "failed to add proposal to db"}
	ErrProposalUnmarshal      = codes.ProtocolError{codes.GovErrProposalUnmarshal, "failed to unmarshal proposal"}
	ErrDeductFunding          = codes.ProtocolError{codes.GovErrDeductFunding, "failed to deduct funds from address"}
	ErrAddFunding             = codes.ProtocolError{codes.GovErrAddFunding, "failed to add funds to address"}
	ErrFundingHeightReached   = codes.ProtocolError{codes.GovErrFundingHeightReached, "funding Height has already been reached"}
	ErrInvalidContributorAddr = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
	ErrInvalidContributorAddr = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
	ErrInvalidContributorAddr = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
	ErrInvalidContributorAddr = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
	ErrInvalidContributorAddr = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
	ErrInvalidContributorAddr = codes.ProtocolError{codes.GovErrInvalidContributorAddr, "invalid contributor address"}
)
