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

func (p ProposalID) String() string {
	return string(p)
}

func (state ProposalState) String() string {
	switch state {
	case ProposalStateError:
		return "Error"
	case ProposalStateActive:
		return "Active"
	case ProposalStatePassed:
		return "Passed"
	case ProposalStateFailed:
		return "Failed"
	default:
		return "Invalid state"
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
		return "To Be Determined"
	default:
		return "Invalid vote result"
	}
}
