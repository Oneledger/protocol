package governance

import (
	"github.com/pkg/errors"
)

type (
	ProposalID      string
	ProposalType    int
	ProposalStatus  int
	ProposalOutcome int
	ProposalState   int
	VoteOpinion     int
	VoteResult      int
)

func (id ProposalID) Err() error {
	switch {
	case len(id) == 0:
		return errors.New("proposal id is empty")
	case len(id) != 32:
		return errors.New("proposal id length is incorrect: must be 32 hex characters")
	}
	return nil
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
