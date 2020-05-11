package governance

const (
	ProposalTypeConfigUpdate ProposalType = iota
	ProposalTypeCodeChange
	ProposalTypeGeneral

	ProposalStatusFunding ProposalStatus = iota
	ProposalStatusVoting
	ProposalStatusCompleted

	ProposalOutcomeInProgress ProposalOutcome = iota
	ProposalOutcomeInsufficientFunds
	ProposalOutcomeInsufficientVotes
	ProposalOutcomeCancelled
	ProposalOutcomeCompleted

	ProposalStateError  ProposalState = -1
	ProposalStateActive ProposalState = iota
	ProposalStatePassed
	ProposalStateFailed

	//Error Codes
	errorSerialization   = "321"
	errorDeSerialization = "322"
	errorSettingRecord   = "323"
	errorGettingRecord   = "324"
	errorDeletingRecord  = "325"
)
