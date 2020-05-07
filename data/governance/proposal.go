package governance

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/Oneledger/protocol/data/keys"
)

const EmptyStr = ""

type options struct {
	FundingDeadline int64
	FundingGoal     int64
	VotingDeadline  int64
}

type ProposalOptions struct {
	ConfigUpdate options
	CodeChange   options
	General      options
}

type ProposerInfo struct {
	address keys.Address
	name    string
}

type Proposal struct {
	ProposalID      ProposalID
	Type            ProposalType
	Status          ProposalStatus
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
		ProposalID:      GenerateProposalID(proposer.address.String() + desc),
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

func GenerateProposalID(key string) ProposalID {
	hashHandler := md5.New()
	_, err := hashHandler.Write([]byte(key))
	if err != nil {
		return EmptyStr
	}
	return ProposalID(hex.EncodeToString(hashHandler.Sum(nil)))
}
