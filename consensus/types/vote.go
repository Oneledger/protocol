package types

import (
	"errors"
	"fmt"
	"io"

	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/tendermint/types"
)

var (
	ErrVoteChannelNotMatch = errors.New("Invalid Channel")
)

type ErrVoteConflictingVotes struct {
	*types.DuplicateVoteEvidence
}

func (err *ErrVoteConflictingVotes) Error() string {
	return fmt.Sprintf("Conflicting votes from validator %v", err.PubKey.Address())
}

func NewConflictingVoteError(val *types.Validator, voteA, voteB *Vote) *ErrVoteConflictingVotes {
	return &ErrVoteConflictingVotes{
		&types.DuplicateVoteEvidence{
			PubKey: val.PubKey,
			VoteA:  &voteA.Vote,
			VoteB:  &voteB.Vote,
		},
	}
}

// Types of votes
// TODO Make a new type "VoteType"
const (
	VoteTypePrevote   = byte(0x01)
	VoteTypePrecommit = byte(0x02)
)

func IsVoteTypeValid(type_ byte) bool {
	switch type_ {
	case VoteTypePrevote:
		return true
	case VoteTypePrecommit:
		return true
	default:
		return false
	}
}

// Represents a prevote, precommit, or commit vote from validators for consensus.
type Vote struct {
	Vote    types.Vote `json:"vote"`
	Channel int        `json:"channel"`
}

func (vote *Vote) WriteSignBytes(chainID string, w io.Writer, n *int, err *error) {
	wire.WriteJSON(types.CanonicalJSONOnceVote{
		chainID,
		types.CanonicalVote(&vote.Vote),
	}, w, n, err)
}

func (vote *Vote) Copy() *Vote {
	voteCopy := *vote
	return &voteCopy
}

func (vote *Vote) String() string {
	return fmt.Sprintf("channel: %s, vote: %q", vote.Channel, vote.Vote.String())
}

func (vote *Vote) Verify(chainID string, pubKey crypto.PubKey, channel int) error {
	if channel != vote.Channel {
		return ErrVoteChannelNotMatch
	}
	return vote.Vote.Verify(chainID, pubKey)
}
