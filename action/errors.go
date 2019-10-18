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
)
