package governance

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/Oneledger/protocol/data/balance"

	"github.com/Oneledger/protocol/data/keys"
)

const EmptyStr = ""

type ProposalOption struct {
	InitialFunding         *balance.Amount          `json:"baseDomainPrice"`
	FundingGoal            *balance.Amount          `json:"fundingGoal"`
	FundingDeadline        int64                    `json:"fundingDeadline"`
	VotingDeadline         int64                    `json:"votingDeadline"`
	PassPercentage         int                      `json:"passPercentage"`
	PassedFundDistribution ProposalFundDistribution `json:"passed_fund_distribution"`
	FailedFundDistribution ProposalFundDistribution `json:"failed_fund_distribution"`
	ProposalExecutionCost  string                   `json:"proposal_execution_cost"`
}

type ProposalFundDistribution struct {
	Validators     float64 `json:"validators"`
	FeePool        float64 `json:"fee_pool"`
	Burn           float64 `json:"burn"`
	ExecutionCost  float64 `json:"execution_cost"`
	BountyPool     float64 `json:"bounty_pool"`
	ProposerReward float64 `json:"proposer_reward"`
}

type ProposalOptionSet struct {
	ConfigUpdate      ProposalOption `json:"config_update"`
	CodeChange        ProposalOption `json:"code_change"`
	General           ProposalOption `json:"general"`
	BountyProgramAddr string         `json:"bounty_program_addr"`
}

type Proposal struct {
	ProposalID      ProposalID
	Type            ProposalType
	Status          ProposalStatus
	Outcome         ProposalOutcome
	Description     string
	Proposer        keys.Address
	FundingDeadline int64
	FundingGoal     *balance.Amount
	VotingDeadline  int64
	PassPercentage  int
}

func NewProposal(proposalID ProposalID, propType ProposalType, desc string, proposer keys.Address, fundingDeadline int64, fundingGoal *balance.Amount,
	votingDeadline int64, passPercentage int) *Proposal {

	return &Proposal{
		ProposalID:      generateProposalID(proposalID),
		Type:            propType,
		Status:          ProposalStatusFunding,
		Outcome:         ProposalOutcomeInProgress,
		Description:     desc,
		Proposer:        proposer,
		FundingDeadline: fundingDeadline,
		FundingGoal:     fundingGoal,
		VotingDeadline:  votingDeadline,
		PassPercentage:  passPercentage,
	}
}

func generateProposalID(key ProposalID) ProposalID {
	hashHandler := md5.New()
	_, err := hashHandler.Write([]byte(key))
	if err != nil {
		return EmptyStr
	}
	return ProposalID(hex.EncodeToString(hashHandler.Sum(nil)))
}
