package governance

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/data/rewards"
)

type GovernanceState struct {
	FeeOption      fees.FeeOption             `json:"feeOption"`
	ETHCDOption    ethchain.ChainDriverOption `json:"ethchaindriverOption"`
	BTCCDOption    bitcoin.ChainDriverOption  `json:"bitcoinChainDriverOption"`
	ONSOptions     ons.Options                `json:"onsOptions"`
	PropOptions    ProposalOptionSet          `json:"propOptions"`
	StakingOptions delegation.Options         `json:"stakingOptions"`
	RewardOptions  rewards.Options            `json:"rewardOptions"`
}
type (
	ProposalID      string
	ProposalType    int
	ProposalStatus  int
	ProposalOutcome int
	ProposalState   int
	VoteOpinion     int
	VoteResult      int
)

func (id ProposalID) String() string {
	return string(id)
}

func (id ProposalID) Err() error {
	switch {
	case len(id) == 0:
		return errors.New("proposal id is empty")
	case len(id) != SHA256LENGTH:
		return errors.New("proposal id length is incorrect")
	}
	return nil
}

func NewProposalState(prefix string) ProposalState {
	prefix = strings.ToLower(prefix)
	switch prefix {
	case "active":
		return ProposalStateActive
	case "passed":
		return ProposalStatePassed
	case "failed":
		return ProposalStateFailed
	case "finalized":
		return ProposalStateFinalized
	case "finalizeFailed":
		return ProposalStateFinalizeFailed
	default:
		return ProposalStateInvalid
	}
}

func (state ProposalState) String() string {
	switch state {
	case ProposalStateInvalid:
		return "Invalid"
	case ProposalStateActive:
		return "Active"
	case ProposalStatePassed:
		return "Passed"
	case ProposalStateFailed:
		return "Failed"
	case ProposalStateFinalized:
		return "Finalized"
	case ProposalStateFinalizeFailed:
		return "FinalizeFailed"
	default:
		return "Invalid state"
	}
}

func NewProposalType(propType string) ProposalType {
	propType = strings.ToLower(propType)
	switch propType {
	case "codechange":
		return ProposalTypeCodeChange
	case "configupdate":
		return ProposalTypeConfigUpdate
	case "general":
		return ProposalTypeGeneral
	default:
		return ProposalTypeInvalid
	}
}

func (propType ProposalType) String() string {
	switch propType {
	case ProposalTypeCodeChange:
		return "Code change"
	case ProposalTypeConfigUpdate:
		return "Config update"
	case ProposalTypeGeneral:
		return "General"
	default:
		return "Invalid type"
	}
}

func (status ProposalStatus) String() string {
	switch status {
	case ProposalStatusFunding:
		return "Funding"
	case ProposalStatusVoting:
		return "Voting"
	case ProposalStatusCompleted:
		return "Completed"
	default:
		return "Invalid status"
	}
}

func (outCome ProposalOutcome) String() string {
	switch outCome {
	case ProposalOutcomeInProgress:
		return "In progress"
	case ProposalOutcomeInsufficientFunds:
		return "Failed [insufficient funds]"
	case ProposalOutcomeInsufficientVotes:
		return "Failed [insufficient votes]"
	case ProposalOutcomeCancelled:
		return "Failed [cancelled]"
	case ProposalOutcomeCompletedYes:
		return "PassedYes"
	case ProposalOutcomeCompletedNo:
		return "PassedNo"
	default:
		return "Invalid outcome"
	}
}

func NewVoteOpinion(opin string) VoteOpinion {
	opin = strings.ToUpper(opin)
	switch opin {
	case "YES":
		return OPIN_POSITIVE
	case "NO":
		return OPIN_NEGATIVE
	case "GIVEUP":
		return OPIN_GIVEUP
	default:
		return OPIN_UNKNOWN
	}
}

func (opinion VoteOpinion) String() string {
	switch opinion {
	case OPIN_UNKNOWN:
		return "Unknown"
	case OPIN_POSITIVE:
		return "Positive"
	case OPIN_NEGATIVE:
		return "Negative"
	case OPIN_GIVEUP:
		return "Giveup"
	default:
		return "Invalid opinion"
	}
}

func (opinion VoteOpinion) Err() error {
	opName := opinion.String()
	if opName == "" {
		return errors.New("vote opinion must be one of [UNKNOWN, POSITIVE, NEGATIVE, GIVEUP]")
	}
	return nil
}

func (opinion VoteResult) String() string {
	switch opinion {
	case VOTE_RESULT_PASSED:
		return "Passed"
	case VOTE_RESULT_FAILED:
		return "Failed"
	case VOTE_RESULT_TBD:
		return "To be determined"
	default:
		return "Invalid vote result"
	}
}
