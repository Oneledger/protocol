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

type VoteStatus struct {
	Result   VoteResult `json:"result"`
	PowerYes int64      `json:"powerYes"`
	PowerNo  int64      `json:"powerNo"`
	PowerAll int64      `json:"powerAll"`
}

func NewProposalVote(validator keys.Address, opinion VoteOpinion, power int64) *ProposalVote {
	return &ProposalVote{
		Validator: validator,
		Opinion:   opinion,
		Power:     power,
	}
}

func NewVoteStatus(result VoteResult, yesPower, noPower, allPower int64) *VoteStatus {
	return &VoteStatus{
		Result:   result,
		PowerYes: yesPower,
		PowerNo:  noPower,
		PowerAll: allPower,
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
