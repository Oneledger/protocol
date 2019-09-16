package action

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrMissingData      = codes.ProtocolError{codes.TxErrMisingData, "missing data in transaction"}
	ErrUnserializable   = codes.ProtocolError{codes.TxErrUnserializable, "unserializable tx"}
	ErrWrongTxType      = codes.ProtocolError{codes.TxErrWrongTxType, "wrong tx type"}
	ErrInvalidAmount    = codes.ProtocolError{codes.TxErrInvalidAmount, "invalid amount"}
	ErrInvalidPubkey    = codes.ProtocolError{codes.TxErrInvalidPubKey, "invalid pubkey"}
	ErrUnmatchSigner    = codes.ProtocolError{codes.TxErrUnmatchedSigner, "unmatch signers"}
	ErrInvalidSignature = codes.ProtocolError{codes.TxErrInvalidSignature, "invalid signatures"}
	ErrInvalidFee       = codes.ProtocolError{codes.TxErrInvalidFee, "invalid fees"}
	ErrNotEnoughFund    = codes.ProtocolError{codes.TxErrInsufficientFunds, "not enough fund"}
)
