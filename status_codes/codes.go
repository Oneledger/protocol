/*

 */

package status_codes

const (
	InvalidParams       = 1001
	IncorrectAddress    = 100101
	DomainMissing       = 100102
	OwnerAddressMissing = 100103
	OnSaleFlagNotSet    = 100104

	IOError        = 1002
	IOErrorNodeKey = 100201

	ParseError        = 1003
	ParseErrorAddress = 100301

	ConfigurationError          = 1004
	ConfigurationErrorChainType = 100401

	ResourceNotFoundError = 1005
	AccountNotFound       = 100501
	DomainNotFound        = 100502

	InternalError                  = 1006
	InternalErrorSerialization     = 100601
	InternalErrorSigning           = 100602
	InternalErrorGeneratingKeyPair = 100603
	InternalErrorGettingBalance    = 100604
	InternalErrorListValidators    = 100605

	WalletError               = 2006
	WalletErrorAddingAccount  = 200601
	WalletErrorGettingAccount = 200602
	WalletErrorDeleteAccount  = 200603

	AccountsError                     = 2007
	AccountsErrorGeneratingNewAccount = 200701

	// Transaction statuses
	TxErrMisingData         = 300101
	TxErrUnserializable     = 300102
	TxErrWrongTxType        = 300103
	TxErrInvalidAmount      = 300104
	TxErrInvalidPubKey      = 300105
	TxErrUnmatchedSigner    = 300106
	TxErrInvalidSignature   = 300107
	TxErrInvalidFeeCurrency = 300108
	TxErrInvalidFeePrice    = 300109
	TxErrInsufficientFunds  = 300110
	TxErrGasOverflow        = 300111
)
