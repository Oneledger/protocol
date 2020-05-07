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

	ProposalStateActive ProposalState = iota
	ProposalStatePassed
	ProposalStateFailed
)
