/*

 */

package status_codes

const (
	GeneralErr          = 999 // all errors without error code
	InvalidParams       = 1001
	IncorrectAddress    = 100101

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

	ONSError                                = 1007
	ONSErrDomainMissing                     = 100701
	ONSErrOwnerAddressMissing               = 100702
	ONSErrOnSaleFlagNotSet                  = 100703
	ONSErrDomainExists                      = 100704
	ONSErrDebitingFromAddress               = 100705
	ONSErrAddingToFeePool                   = 100706
	ONSErrInvalidUri		                = 100707
	ONSErrGettingParentName                 = 100708
	ONSErrParentDoesNotExist                = 100709
	ONSErrParentNotOwned                    = 100710
	ONSErrFailedToCalculateExpiry           = 100711
	ONSErrFailedToCreateDomain              = 100712
	ONSErrFailedAddingDomainToStore         = 100713
	ONSErrInvalidDomainName                 = 100714


	WalletError               = 2006
	WalletErrorAddingAccount  = 200601
	WalletErrorGettingAccount = 200602
	WalletErrorDeleteAccount  = 200603

	BalanceError
	BalanceErrorAddFailed   = 210600
	BalanceErrorMinusFailed = 210601

	AccountsError                     = 2007
	AccountsErrorGeneratingNewAccount = 200701

	// Transaction statuses
	TxErrMissingData        = 300101
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

	//ValidatorStore
	ValidatorsUnableGetList = 410100

	//Tracker
	ETHTrackerNotFoundFailed  = 600100
	ETHTrackerNotFoundSuccess = 600101
	ETHTrackerNotFoundOngoing = 600102

	GovErr								  = 7001
	GovErrVoteSetupValidator              = 700100
	GovErrVoteUpdateVote                  = 700101
	GovErrVoteDeleteVoteRecords           = 700102
	GovErrVoteCheckVoteResult             = 700103
	GovErrWithdrawCheckFundsFailed        = 700104
	GovErrGetProposalOptions              = 700105
	GovErrInvalidProposalId               = 700106
	GovErrInvalidProposalType             = 700107
	GovErrInvalidProposerAddr             = 700108
	GovErrInvalidProposalDesc             = 700109
	GovErrProposalExists                  = 700110
	GovErrProposalNotExists               = 700111
	GovErrAddingProposalToActiveStore     = 700112
	GovErrDeletingProposalFromActiveStore = 700113
	GovErrAddingProposalToPassedStore     = 700114
	GovErrAddingProposalToFailedStore     = 700115
	GovErrDeletingProposalFromFailedStore = 700116
	GovErrDeductFunding                   = 700117
	GovErrAddFunding                      = 700118
	GovErrInvalidFunderAddr               = 700119
	GovErrInvalidBeneficiaryAddr          = 700120
	GovErrFundingHeightReached            = 700121
	GovErrNotInFunding                    = 700122
	GovErrGettingValidatorList            = 700123
	GovErrSetupVotingValidator            = 700124
	GovErrProposalWithdrawNotEligible     = 700125
	GovErrNoSuchFunder                    = 700126
	GovErrNotInVoting                     = 700127
	GovErrVotingHeightReached             = 700128
	GovErrAddingVoteToVoteStore           = 700129
	GovErrPeekingVoteResult               = 700130
	GovErrUnmatchedProposer               = 700131
	GovErrInvalidVoterId                  = 700132
	GovErrInvalidValidatorAddr            = 700133
	GovErrInvalidVoteOpinion              = 700134
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
	FundingDeadlineCrossed         = 700115
	StatusUnableToSetVoting        = 700116
	GovFundUnableToAdd             = 700117
	GovFundUnableToDelete          = 700118
	GovFundBalanceMismatch         = 700119
)
