package governance

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/Oneledger/protocol/data/balance"

	"github.com/Oneledger/protocol/data/keys"
)

const EmptyStr = ""

type ProposalOption struct {
	InitialFunding  *balance.Amount `json:"baseDomainPrice"`
	FundingGoal     *balance.Amount `json:"fundingGoal"`
	FundingDeadline int64           `json:"fundingDeadline"`
	VotingDeadline  int64           `json:"votingDeadline"`
	PassPercentage  int             `json:"passPercentage"`
}

type ProposalOptionSet struct {
	ConfigUpdate ProposalOption
	CodeChange   ProposalOption
	General      ProposalOption
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

func NewProposal(propType ProposalType, desc string, proposer keys.Address, fundingDeadline int64, fundingGoal *balance.Amount,
	votingDeadline int64, passPercentage int) *Proposal {

	return &Proposal{
		ProposalID:      generateProposalID(proposer.String()),
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

func generateProposalID(key string) ProposalID {
	uniqueKey := key + time.Now().String()
	hashHandler := md5.New()
	_, err := hashHandler.Write([]byte(uniqueKey))
	if err != nil {
		return EmptyStr
	}
	return ProposalID(hex.EncodeToString(hashHandler.Sum(nil)))
}