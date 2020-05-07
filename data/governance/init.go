package governance

const (
	ProposalTypeConfigUpdate ProposalType = iota
	ProposalTypeCodeChange
	ProposalTypeGeneral

	ProposalStateActive ProposalState = iota
	ProposalStatePassed
	ProposalStateFailed

	ProposalOutcomeInProgress ProposalOutcome = iota
	ProposalOutcomeInsufficientFunds
	ProposalOutcomeInsufficientVotes
	ProposalOutcomeCompleted
)
