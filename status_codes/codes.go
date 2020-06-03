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

	ParseError           = 1003
	ParseErrorAddress    = 100301
	ParseErrorBadBTCTxn  = 100302
	ParseErrorBTCAddress = 100303

	ConfigurationError          = 1004
	ConfigurationErrorChainType = 100401

	ResourceNotFoundError = 1005
	AccountNotFound       = 100501
	DomainNotFound        = 100502
	CurrencyNotFound      = 100503
	TxNotFound            = 100504

	InternalError                           = 1006
	InternalErrorSerialization              = 100601
	InternalErrorSigning                    = 100602
	InternalErrorGeneratingKeyPair          = 100603
	InternalErrorGettingBalance             = 100604
	InternalErrorListValidators             = 100605
	InternalErrorGettingTracker             = 100606
	InternalErrorTrackerNotFound            = 100607
	InternalErrorTrackerBusy                = 100608
	InternalErrorTrackerInsufficientBalance = 100609
	InternalErrorListWitnesses              = 100610
	InternalErrorGettingProposal            = 100611

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
	TxErrInvalidExtTx       = 300112

	ExternalErr                        = 400100
	ExternalErrBitcoinTxNotFound       = 400101
	ExternalErrGettingBTCTxn           = 400102
	ExternalErrNotEnoughConfirmations  = 400103
	ExternalErrNotSpendable            = 400104
	ExternalErrUnableToCreateEthTX     = 400105
	ExternalErrUnableToCreateOLTLockTX = 400106
	ErrUnmarshalingRedeem              = 400107

	//ERC20
	ExternalErrUnableToCreateErc20OLTLockTX = 500100
	ExternalErrTokenNotSuported             = 500101

	//Tracker
	ETHTrackerNotFoundFailed  = 600100
	ETHTrackerNotFoundSuccess = 600101
	ETHTrackerNotFoundOngoing = 600102

	//Governance
	GovErrVoteSetupValidator       = 700100
	GovErrVoteUpdateVote           = 700101
	GovErrVoteDeleteVoteRecords    = 700102
	GovErrVoteCheckVoteResult      = 700103
	GovErrWithdrawCheckFundsFailed = 700104
	ProposalNotFound               = 700105
	UnauthorizedCall               = 700106
	StatusNotCompleted             = 700107
	StatusNotVoting                = 700108
	StatusNotFunding               = 700109
	VotingTBD                      = 700110
	FinalizeDistributtionFailed    = 700111
	FinalizeConfigUpdateFailed     = 700112
	StatusUnableToSetFinalized     = 700113
	UnabletoQueryVoteResult        = 700114
)
