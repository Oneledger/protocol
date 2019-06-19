package action

import "errors"

var (
	ErrMissingData      = errors.New("missing data for tx")
	ErrUnserializable   = errors.New("unserializable tx")
	ErrWrongTxType      = errors.New("wrong tx type")
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrInvalidPubkey    = errors.New("invalid pubkey")
	ErrUnmatchSigner    = errors.New("unmatch signers")
	ErrInvalidSignature = errors.New("invalid signatures")
	ErrInvalidFee       = errors.New("invalid fees")
	ErrNotEnoughFund    = errors.New("not enough fund")

	ErrInvalidDomain = errors.New("invalid domain name")
	ErrDomainAuth    = errors.New("domain auth failed")
)
