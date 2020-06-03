package governance

import (
	"encoding/hex"
	"fmt"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type ProposalVote struct {
	Validator keys.Address
	Opinion   VoteOpinion
	Power     int64
}

func NewProposalVote(validator keys.Address, opinion VoteOpinion, power int64) *ProposalVote {
	return &ProposalVote{
		Validator: validator,
		Opinion:   opinion,
		Power:     power,
	}
}

func (vote *ProposalVote) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(vote)
	if err != nil {
		logger.Error("proposal vote not serializable", err)
		return []byte{}
	}
	return value
}

func (vote *ProposalVote) FromBytes(msg []byte) (*ProposalVote, error) {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, vote)
	if err != nil {
		logger.Error("failed to deserialize a proposal vote from bytes", err)
		return nil, err
	}
	return vote, nil
}

func (vote *ProposalVote) String() string {
	return fmt.Sprintf("validator= %v, opinion= %v, power= %v",
		hex.EncodeToString(vote.Validator), vote.Opinion.String(), vote.Power)
}
