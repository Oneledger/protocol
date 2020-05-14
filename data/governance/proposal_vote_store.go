package governance

import (
	"encoding/hex"
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

func (opinion VoteOpinion) String() string {
	switch opinion {
	case UNKNOWN:
		return "Unknown"
	case POSITIVE:
		return "Positive"
	case NEGATIVE:
		return "Negative"
	case GIVEUP:
		return "Giveup"
	default:
		return "Invalid opinion"
	}
}

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
	info := fmt.Sprintf("Vote Setup: proposalID= %v, validator= %v, power= %v", proposalID, hex.EncodeToString(validator), power)

	if proposalID == "" {
		logger.Errorf("%v, empty proposalID", info)
		return ErrVoteSetupValidatorFailed
	}

	key := storage.StoreKey(string(pvs.prefix) + proposalID + storage.DB_PREFIX + string(validator))
	vote := &ProposalVote{ProposalID: proposalID, Opinion: UNKNOWN, Power: power}
	value := vote.Bytes()
	err := pvs.store.Set(key, value)
	if err != nil {
		logger.Errorf("%v, storage failure", info)
		return ErrVoteSetupValidatorFailed
	}
	logger.Info(info)

	return nil
}

// Update a validator's voting opinion to proposalID
func (pvs *ProposalVoteStore) Update(proposalID string, validator keys.Address, vote *ProposalVote) error {
	info := fmt.Sprintf("Vote Update: %v, validator= %v", vote.String(), hex.EncodeToString(validator))

	key := storage.StoreKey(string(pvs.prefix) + proposalID + storage.DB_PREFIX + string(validator))
	exist := pvs.store.Exists(key)
	if !exist {
		logger.Errorf("%v, can't participate in voting", info)
		return ErrVoteUpdateVoteFailed
	}

	value := vote.Bytes()
	err := pvs.store.Set(key, value)
	if err != nil {
		logger.Errorf("%v, storage failure", info)
		return ErrVoteUpdateVoteFailed
	}
	logger.Info(info)

	return nil
}

// Delete all voting records under a proposalID
func (pvs *ProposalVoteStore) Delete(proposalID string) error {
	info := fmt.Sprintf("Vote Delete: proposalID= %v", proposalID)

	if !pvs.Exists(proposalID) {
		logger.Errorf("%v, proposal does not exist", info)
		return ErrVoteDeleteVoteRecordsFailed
	}

	succeed := true
	pvs.IterateByID(proposalID, func(key []byte, value []byte) bool {
		_, err := pvs.store.Delete(key)
		if err != nil {
			logger.Errorf("%v, failed to delete vote, key= %v", info, string(key))
			succeed = false
		}
		return false
	})

	if !succeed {
		logger.Errorf("%v, delete failed", info)
		return ErrVoteDeleteVoteRecordsFailed
	}
	logger.Info(info)

	return nil
}

// Check and see if a proposal has passed
func (pvs *ProposalVoteStore) IsPassed(proposalID string, passPercent int64) (bool, error) {
	info := fmt.Sprintf("Vote IsPassed: proposalID= %v", proposalID)

	_, votes, err := pvs.GetVotesByID(proposalID)
	if err != nil {
		logger.Errorf("%v, getVotesByID failed", info)
		return false, ErrVoteCheckVoteResultFailed
	}

	// Currently each validator is treated equally in voting power
	voteResult := make([]int64, 4)
	for _, vote := range votes {
		voteResult[vote.Opinion] += 1
	}

	// Excludes validators that give up voting in percent calculation
	percent := 0.0
	passed := false
	total := int64(len(votes)) - voteResult[GIVEUP]
	if total > 0 {
		percent = float64(voteResult[POSITIVE]) / float64(total)
	}
	if percent >= float64(passPercent)/100.0 {
		passed = true
	}
	logger.Infof("%v, passed= %v", info, passed)

	return passed, nil
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

// get voting votes by proposalID
func (pvs *ProposalVoteStore) GetVotesByID(proposalID string) ([]keys.Address, []*ProposalVote, error) {
	info := fmt.Sprintf("Vote getVotesByID: proposalID= %v", proposalID)

	if !pvs.Exists(proposalID) {
		errMsg := fmt.Sprintf("%v, proposal does not exist", info)
		logger.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	succeed := true
	addrs := make([]keys.Address, 0)
	votes := make([]*ProposalVote, 0)
	pvs.IterateByID(proposalID, func(key []byte, value []byte) bool {
		vote, err := (&ProposalVote{}).FromBytes(value)
		if err != nil {
			logger.Errorf("%v, key= %v, deserialize proposal vote failed", info, key)
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
		errMsg := fmt.Sprintf("%v, operation failed", info)
		logger.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}
	return addrs, votes, nil
}

// check existance of proposal
// a proposal does not exist if there is no vote records
func (pvs *ProposalVoteStore) Exists(proposalID string) bool {
	prefix := append(pvs.prefix, (proposalID + storage.DB_PREFIX)...)
	exist := false
	pvs.store.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			exist = true
			return true
		},
	)
	return exist
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

func (pv *ProposalVote) String() string {
	return fmt.Sprintf("proposalID= %v, Opinion= %v, Power= %v", pv.ProposalID, pv.Opinion.String(), pv.Power)
}
