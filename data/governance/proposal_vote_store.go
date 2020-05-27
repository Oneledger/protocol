package governance

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type ProposalVoteStore struct {
	prefix []byte
	store  *storage.State
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
func (pvs *ProposalVoteStore) Setup(proposalID ProposalID, vote *ProposalVote) error {
	if proposalID == "" {
		return ErrVoteSetupValidatorFailed
	}

	vote.Opinion = OPIN_UNKNOWN // initialize as OPIN_UNKNOWN
	key := GetKey(pvs.prefix, proposalID, vote)
	value := vote.Bytes()
	err := pvs.store.Set(key, value)
	if err != nil {
		return ErrVoteSetupValidatorFailed
	}

	return nil
}

// Update a validator's voting opinion to proposalID
func (pvs *ProposalVoteStore) Update(proposalID ProposalID, vote *ProposalVote) error {
	key := GetKey(pvs.prefix, proposalID, vote)
	exist := pvs.store.Exists(key)
	if !exist {
		return ErrVoteUpdateVoteFailed
	}

	value := vote.Bytes()
	err := pvs.store.Set(key, value)
	if err != nil {
		return ErrVoteUpdateVoteFailed
	}

	return nil
}

// Delete all voting records under a proposalID
func (pvs *ProposalVoteStore) Delete(proposalID ProposalID) error {
	succeed := true
	pvs.IterateByID(proposalID, func(key []byte, value []byte) bool {
		_, err := pvs.store.Delete(key)
		if err != nil {
			succeed = false
		}
		return false
	})

	if !succeed {
		return ErrVoteDeleteVoteRecordsFailed
	}

	return nil
}

// Check and see if a proposal has passed
func (pvs *ProposalVoteStore) IsPassed(proposalID ProposalID, passPercent int64) (bool, error) {
	_, votes, err := pvs.GetVotesByID(proposalID)
	if err != nil {
		return false, ErrVoteCheckVoteResultFailed
	}

	// Accumulates power of each opinion
	totalPower := int64(0)
	eachPower := make([]int64, 4)
	for _, vote := range votes {
		totalPower += vote.Power
		eachPower[vote.Opinion] += vote.Power
	}

	// Excludes validators that give up voting in percent calculation
	totalPower -= eachPower[OPIN_GIVEUP]

	// Calculate actual percentage
	percent := 0.0
	passed := false
	if totalPower > 0 {
		percent = float64(eachPower[OPIN_POSITIVE]) / float64(totalPower)
	}
	if percent >= float64(passPercent)/100.0 {
		passed = true
	}

	return passed, nil
}

// Iterate voting records by proposalID
func (pvs *ProposalVoteStore) IterateByID(proposalID ProposalID, fn func(key []byte, value []byte) bool) (stopped bool) {
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
func (pvs *ProposalVoteStore) GetVotesByID(proposalID ProposalID) ([]keys.Address, []*ProposalVote, error) {
	succeed := true
	addrs := make([]keys.Address, 0)
	votes := make([]*ProposalVote, 0)
	pvs.IterateByID(proposalID, func(key []byte, value []byte) bool {
		vote, err := (&ProposalVote{}).FromBytes(value)
		if err != nil {
			succeed = false
			return true
		}
		votes = append(votes, vote)
		prefix_len := len(append(pvs.prefix, (proposalID + storage.DB_PREFIX)...))
		addr := key[prefix_len:]
		addrs = append(addrs, addr)
		return false
	})

	if !succeed {
		return nil, nil, errors.New("")
	}
	// Caused by invalid/deleted proposalID
	if len(votes) == 0 {
		return nil, nil, errors.New("")
	}

	return addrs, votes, nil
}

func GetKey(prefix []byte, proposalID ProposalID, vote *ProposalVote) []byte {
	return storage.StoreKey(string(prefix) + string(proposalID) + storage.DB_PREFIX + string(vote.Validator))
}
