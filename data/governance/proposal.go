package governance

import "github.com/Oneledger/protocol/data/keys"

type (
	ProposalType    int
	ProposalState   int
	ProposalOutcome int
)

type ProposerInfo struct {
	address keys.Address
	name    string
}

type Proposal struct {
	Type            ProposalType
	State           ProposalState
	Outcome         ProposalOutcome
	Description     string
	Proposer        ProposerInfo
	FundingDeadline int64
	FundingGoal     int64
	VotingDeadline  int64
}

func NewProposal(propType ProposalType, desc string, proposer ProposerInfo, fundingDeadline int64, fundingGoal int64,
	votingDeadline int64) *Proposal {

	return &Proposal{
		Type:            propType,
		State:           ProposalStateActive,
		Outcome:         ProposalOutcomeInProgress,
		Description:     desc,
		Proposer:        proposer,
		FundingDeadline: fundingDeadline,
		FundingGoal:     fundingGoal,
		VotingDeadline:  votingDeadline,
	}
}
