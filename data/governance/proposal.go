package governance

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/Oneledger/protocol/data/keys"
)

const EmptyStr = ""

type options struct {
	InitialFunding  int64
	FundingDeadline int64
	FundingGoal     int64
	VotingDeadline  int64
	PassPercentage  float64
}

type ProposalOptions struct {
	ConfigUpdate options
	CodeChange   options
	General      options
}

type Proposal struct {
	ProposalID      ProposalID
	Type            ProposalType
	Status          ProposalStatus
	Outcome         ProposalOutcome
	Description     string
	Proposer        keys.Address
	FundingDeadline int64
	FundingGoal     int64
	VotingDeadline  int64
}

func NewProposal(propType ProposalType, desc string, proposer keys.Address, fundingDeadline int64, fundingGoal int64,
	votingDeadline int64) *Proposal {

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
