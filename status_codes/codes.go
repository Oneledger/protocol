/*

 */

package status_codes

const (
	GeneralErr       = 999 // all errors without error code
	InvalidParams    = 1001
	IncorrectAddress = 100101

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
	InternalErrorGettingBidConv             = 100612

	ONSError                        = 1007
	ONSErrDomainMissing             = 100701
	ONSErrOwnerAddressMissing       = 100702
	ONSErrOnSaleFlagNotSet          = 100703
	ONSErrDomainExists              = 100704
	ONSErrDebitingFromAddress       = 100705
	ONSErrAddingToFeePool           = 100706
	ONSErrInvalidUri                = 100707
	ONSErrGettingParentName         = 100708
	ONSErrParentDoesNotExist        = 100709
	ONSErrParentNotOwned            = 100710
	ONSErrFailedToCalculateExpiry   = 100711
	ONSErrFailedToCreateDomain      = 100712
	ONSErrFailedAddingDomainToStore = 100713
	ONSErrInvalidDomainName         = 100714

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
	TxErrPoolDoesNotExist   = 300113
	TxErrInsufficientFunds  = 300110
	TxErrGasOverflow        = 300111
	TxErrInvalidExtTx       = 300112
	TxErrMaliciousValidator = 300113

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
	ExternalTransactionNotSupported         = 500102

	//ValidatorStore
	ValidatorsUnableGetList = 410100

	//Tracker
	ETHTrackerNotFoundFailed  = 600100
	ETHTrackerNotFoundSuccess = 600101
	ETHTrackerNotFoundOngoing = 600102
	ETHTrackerExists          = 600103
	ETHTrackerUnabletoSet     = 600104

	// Staking
	DelgErr                     = 6003
	DelgErrStakeAddressInUse    = 600301
	DelgErrStakeAddressMismatch = 600302

	NetDelgErr                              = 6005
	NetDelgErrGettingActiveDelgAmount       = 600501
	NetDelgErrDeductingActiveDelgAmount     = 600502
	NetDelgErrSettingActiveDelgAmount       = 600503
	NetDelgErrGettingDelgOption             = 600504
	NetDelgErrSettingPendingDelgAmount      = 600505
	NetDelgErrGettingPendingDelgAmount      = 600506
	NetDelgErrWithdraw                      = 600507
	NetDelgErrFinalizingDelgRewards         = 600508
	NetDelgErrAddingWithdrawAmountToBalance = 600509
	NetDelgErrReinvest                      = 600510

	GovErr                                = 7001
	GovErrGetProposalOptions              = 700101
	GovErrInvalidProposalId               = 700102
	GovErrInvalidProposalType             = 700103
	GovErrInvalidProposerAddr             = 700104
	GovErrInvalidProposalDesc             = 700105
	GovErrProposalExists                  = 700106
	GovErrProposalNotExists               = 700107
	GovErrAddingProposalToActiveStore     = 700108
	GovErrDeletingProposalFromActiveStore = 700109
	GovErrAddingProposalToPassedStore     = 700110
	GovErrAddingProposalToFailedStore     = 700111
	GovErrDeletingProposalFromFailedStore = 700112
	GovErrDeductFunding                   = 700113
	GovErrAddFunding                      = 700114
	GovErrInvalidFunderAddr               = 700115
	GovErrInvalidBeneficiaryAddr          = 700116
	GovErrGettingValidatorList            = 700117
	GovErrGettingWitnessList              = 700217
	GovErrSetupVotingValidator            = 700118
	GovErrProposalWithdrawNotEligible     = 700119
	GovErrNoSuchFunder                    = 700120
	GovErrVotingHeightReached             = 700121
	GovErrAddingVoteToVoteStore           = 700122
	GovErrPeekingVoteResult               = 700123
	GovErrUnmatchedProposer               = 700124
	GovErrInvalidVoterId                  = 700125
	GovErrInvalidValidatorAddr            = 700126
	GovErrInvalidVoteOpinion              = 700127
	GovErrVoteUpdate                      = 700128
	GovErrVoteDeleteVoteRecords           = 700129
	GovErrVoteCheckVoteResult             = 700130
	GovErrWithdrawCheckFundsFailed        = 700131
	GovErrUnauthorizedCall                = 700132
	GovErrStatusNotCompleted              = 700133
	GovErrStatusNotVoting                 = 700134
	GovErrStatusNotFunding                = 700135
	GovErrVotingTBD                       = 700136
	GovErrFinalizeDistributtionFailed     = 700137
	GovErrFinalizeConfigUpdateFailed      = 700138
	GovErrStatusUnableToSetFinalized      = 700139
	GovErrUnabletoQueryVoteResult         = 700140
	GovErrFundingDeadlineCrossed          = 700141
	StatusUnableToSetVoting               = 700142
	GovFundUnableToAdd                    = 700143
	GovFundUnableToDelete                 = 700144
	GovFundBalanceMismatch                = 700145
	GovErrUnableToSetFinalizeFailed       = 700146
	GovErrInvalidFundingGoal              = 700147
	GovErrInvalidPassPercentage           = 700148
	GovErrInvalidFundingDeadline          = 700149
	GovErrInvalidVotingDeadline           = 700150
	TxErrGettingETHOptions                = 700151
	TxErrGettingONSOptions                = 700152
	TxErrGettingBTCOptions                = 700153
	TxErrGettingRewardsOptions            = 700154
	TxErrGettingFeeOptions                = 700155
	TxErrValidateGovState                 = 700156
	TxErrGetStakingOptions                = 700158
	TxErrInvalidOptions                   = 700159
	TxErrEvidenceError                    = 700160

	//Rewards Error
	RewardsUnableToGetMaturedAmount = 800001
	RewardsUnableToWithdraw         = 800002
	RewardsYearRewardsMissing       = 800003

	TxErrVMExecution  = 900001
	TxErrIntrinsicGas = 900002

	//Ethereum Errors

	//6 digit error codes starting with 99 are reserved for external apps
)
