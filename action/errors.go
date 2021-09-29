package action

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrMissingData        = codes.ProtocolError{codes.TxErrMissingData, "missing data in transaction"}
	ErrUnserializable     = codes.ProtocolError{codes.TxErrUnserializable, "unserializable tx"}
	ErrWrongTxType        = codes.ProtocolError{codes.TxErrWrongTxType, "wrong tx type"}
	ErrInvalidAmount      = codes.ProtocolError{codes.TxErrInvalidAmount, "invalid amount"}
	ErrInvalidPubkey      = codes.ProtocolError{codes.TxErrInvalidPubKey, "invalid pubkey"}
	ErrUnmatchSigner      = codes.ProtocolError{codes.TxErrUnmatchedSigner, "unmatch signers"}
	ErrInvalidSignature   = codes.ProtocolError{codes.TxErrInvalidSignature, "invalid signatures"}
	ErrInvalidFeeCurrency = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid fees currency"}
	ErrInvalidFeePrice    = codes.ProtocolError{codes.TxErrInvalidFeePrice, "fee price is smaller than minimal fee"}
	ErrPoolDoesNotExist   = codes.ProtocolError{codes.TxErrPoolDoesNotExist, "Pool does not exist"}
	ErrNotEnoughFund      = codes.ProtocolError{codes.TxErrInsufficientFunds, "not enough fund"}
	ErrGasOverflow        = codes.ProtocolError{codes.TxErrGasOverflow, "gas used exceed limit"}
	ErrInvalidExtTx       = codes.ProtocolError{codes.TxErrInvalidExtTx, "invalid external tx"}
	ErrInvalidVmExecution = codes.ProtocolError{codes.TxErrVMExecution, "vm execution error"}

	ErrInvalidAddress          = codes.ErrBadAddress
	ErrInvalidCurrency         = codes.ProtocolError{codes.TxErrInvalidFeeCurrency, "invalid fee currency"}
	ErrTokenNotSupported       = codes.ProtocolError{codes.ExternalErrTokenNotSuported, "Token not supported"}
	ErrTransactionNotSupported = codes.ProtocolError{codes.ExternalTransactionNotSupported, "TX not supported"}

	ErrGettingValidatorList = codes.ProtocolError{codes.GovErrGettingValidatorList, "fund proposal failed in getting validator list"}
	ErrGettingWitnessList   = codes.ProtocolError{codes.GovErrGettingWitnessList, "failed in getting witness list"}

	ErrInvalidValidatorAddr = codes.ProtocolError{codes.GovErrInvalidValidatorAddr, "invalid validator address"}
	ErrStakeAddressInUse    = codes.ProtocolError{codes.DelgErrStakeAddressInUse, "current stake address is in use"}
	ErrStakeAddressMismatch = codes.ProtocolError{codes.DelgErrStakeAddressMismatch, "stake address does not match"}
)
