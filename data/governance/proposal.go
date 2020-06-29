package governance

import (
	"github.com/Oneledger/protocol/data/balance"

	"github.com/Oneledger/protocol/data/keys"
)

const EmptyStr = ""

type ProposalOption struct {
	InitialFunding         *balance.Amount          `json:"initialFunding"`
	FundingGoal            *balance.Amount          `json:"fundingGoal"`
	FundingDeadline        int64                    `json:"fundingDeadline"`
	VotingDeadline         int64                    `json:"votingDeadline"`
	PassPercentage         int                      `json:"passPercentage"`
	PassedFundDistribution ProposalFundDistribution `json:"passedFundDistribution"`
	FailedFundDistribution ProposalFundDistribution `json:"failedFundDistribution"`
	ProposalExecutionCost  string                   `json:"proposalExecutionCost"`
}

type ProposalFundDistribution struct {
	Validators     float64 `json:"validators"`
	FeePool        float64 `json:"feePool"`
	Burn           float64 `json:"burn"`
	ExecutionCost  float64 `json:"executionCost"`
	BountyPool     float64 `json:"bountyPool"`
	ProposerReward float64 `json:"proposerReward"`
}

type ProposalOptionSet struct {
	ConfigUpdate      ProposalOption `json:"configUpdate"`
	CodeChange        ProposalOption `json:"codeChange"`
	General           ProposalOption `json:"general"`
	BountyProgramAddr string         `json:"bountyProgramAddr"`
}

type Proposal struct {
	ProposalID      ProposalID      `json:"proposalId"`
	Type            ProposalType    `json:"proposalType"`
	Status          ProposalStatus  `json:"status"`
	Outcome         ProposalOutcome `json:"outcome"`
	Headline        string          `json:"headline"`
	Description     string          `json:"descr"`
	Proposer        keys.Address    `json:"proposer"`
	FundingDeadline int64           `json:"fundingDeadline"`
	FundingGoal     *balance.Amount `json:"fundingGoal"`
	VotingDeadline  int64           `json:"votingDeadline"`
	PassPercentage  int             `json:"passPercent"`
}

func NewProposal(proposalID ProposalID, propType ProposalType, desc string, headline string, proposer keys.Address, fundingDeadline int64, fundingGoal *balance.Amount,
	votingDeadline int64, passPercentage int) *Proposal {

	return &Proposal{
		ProposalID:      proposalID,
		Type:            propType,
		Status:          ProposalStatusFunding,
		Outcome:         ProposalOutcomeInProgress,
		Description:     desc,
		Headline:        headline,
		Proposer:        proposer,
		FundingDeadline: fundingDeadline,
		FundingGoal:     fundingGoal,
		VotingDeadline:  votingDeadline,
		PassPercentage:  passPercentage,
	}
}
