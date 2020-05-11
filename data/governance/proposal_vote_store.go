package governance

import (
	"fmt"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type VoteOpinion int

const (
	UNKNOWN  VoteOpinion = 0
	POSITIVE VoteOpinion = 1
	NEGATIVE VoteOpinion = 2
	GIVEUP   VoteOpinion = 3
)

type ProposalVoteStore struct {
	prefix []byte
	store  *storage.State
}

type ProposalVote struct {
	ProposalID string
	Opinion    VoteOpinion
	Power      int64
}

func NewProposalVoteStore(prefix string, state *storage.State) *ProposalVoteStore {
	return &ProposalVoteStore{
		prefix: storage.Prefix(prefix),
		store:  state,
	}
}

func (pvs *ProposalVoteStore) WithState(state *storage.State) *ProposalVoteStore {
	pvs.store = state
	return pvs
}

// Setup an initial voting validator to proposalID
func (pvs *ProposalVoteStore) Setup(proposalID string, validator keys.Address, power int64) error {
	if proposalID == "" {
		return errors.New("failed to setup voting validator: empty proposalID")
	}

	key := storage.StoreKey(string(pvs.prefix) + proposalID + storage.DB_PREFIX + string(validator))
	vote := &ProposalVote{ProposalID: proposalID, Opinion: UNKNOWN, Power: power}
	value := vote.Bytes()
	err := pvs.store.Set(key, value)
	if err != nil {
		errMsg := fmt.Sprintf("failed to setup voting validator to proposalID: %v, storage failure", proposalID)
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// Update a validator's voting opinion to proposalID
func (pvs *ProposalVoteStore) Update(proposalID string, validator keys.Address, vote *ProposalVote) error {
	key := storage.StoreKey(string(pvs.prefix) + proposalID + storage.DB_PREFIX + string(validator))
	exist := pvs.store.Exists(key)
	if !exist {
		errMsg := fmt.Sprintf("failed to vote proposalID: %v, validator: %v does not exists", proposalID, string(validator))
		return errors.New(errMsg)
	}

	value := vote.Bytes()
	err := pvs.store.Set(key, value)
	if err != nil {
		errMsg := fmt.Sprintf("failed to vote proposalID: %v, validator: %v, storage failure", proposalID, string(validator))
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// Delete all voting records under a proposalID
func (pvs *ProposalVoteStore) Delete(proposalID string) error {
	succeed := true
	pvs.IterateByID(proposalID, func(key []byte, value []byte) bool {
		_, err := pvs.store.Delete(key)
		if err != nil {
			errMsg := fmt.Sprintf("failed to delete voting record under proposalID: %v, key: %v", proposalID, string(key))
			logger.Error(errMsg)
			succeed = false
		}
		return false
	})

	if !succeed {
		errMsg := fmt.Sprintf("failed to delete voting record under proposalID: %v", proposalID)
		return errors.New(errMsg)
	}
	return nil
}

// Check and see if a proposal has passed
func (pvs *ProposalVoteStore) IsPassed(proposalID string) bool {
	_, votes, err := pvs.GetVotesByID(proposalID)
	if err != nil {
		return false
	}
	voteResult := make([]int64, 4)
	for _, vote := range votes {
		voteResult[vote.Opinion] += 1
	}

	percent := 0.0
	percentPass := 0.67
	total := int64(len(votes)) - voteResult[GIVEUP]
	if total > 0 {
		percent = float64(voteResult[POSITIVE]) / float64(total)
	}
	if percent >= percentPass {
		return true
	}
	return false
}

// get voting votes by proposalID
func (pvs *ProposalVoteStore) GetVotesByID(proposalID string) ([]keys.Address, []*ProposalVote, error) {
	succeed := true
	addrs := make([]keys.Address, 0)
	votes := make([]*ProposalVote, 0)
	pvs.IterateByID(proposalID, func(key []byte, value []byte) bool {
		vote, err := (&ProposalVote{}).FromBytes(value)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to deserialize proposal vote under proposalID: %v", proposalID))
			succeed = false
			return false
		}
		votes = append(votes, vote)
		prefix_len := len(append(pvs.prefix, (proposalID + storage.DB_PREFIX)...))
		addr := key[prefix_len:]
		addrs = append(addrs, addr)
		return false
	})

	if !succeed {
		errMsg := fmt.Sprintf("failed to get voting records under proposalID: %v", proposalID)
		return nil, nil, errors.New(errMsg)
	}
	return addrs, votes, nil
}

// Iterate voting records by proposalID
func (pvs *ProposalVoteStore) IterateByID(proposalID string, fn func(key []byte, value []byte) bool) (stopped bool) {
	prefix := append(pvs.prefix, (proposalID + storage.DB_PREFIX)...)
	return pvs.store.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			return fn(key, value)
		},
	)
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
