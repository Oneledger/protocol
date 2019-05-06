package action

import "errors"

var (
	ErrMissingData      = errors.New("missing data for send tx")
	ErrWrongTxType      = errors.New("wrong tx type")
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrInvalidPubkey    = errors.New("invalid pubkey")
	ErrUnmatchSigner    = errors.New("unmatch signers")
	ErrInvalidSignature = errors.New("invalid signatures")
	ErrInvalidFee       = errors.New("invalid fees")
)
