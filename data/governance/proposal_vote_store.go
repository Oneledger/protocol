package governance

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
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
	info := fmt.Sprintf("Vote Setup: proposalID= %v, %v", proposalID, vote.String())

	if proposalID == "" {
		logger.Errorf("%v, empty proposalID", info)
		return ErrSetupVotingValidator
	}

	vote.Opinion = OPIN_UNKNOWN // initialize as OPIN_UNKNOWN
	key := GetKey(pvs.prefix, proposalID, vote)
	value := vote.Bytes()
	err := pvs.store.Set(key, value)
	if err != nil {
		logger.Errorf("%v, storage failure", info)
		return ErrSetupVotingValidator
	}
	logger.Detail(info)

	return nil
}

// Update a validator's voting opinion to proposalID
func (pvs *ProposalVoteStore) Update(proposalID ProposalID, vote *ProposalVote) error {
	info := fmt.Sprintf("Vote Update: proposalID= %v, %v", proposalID, vote.String())

	// Get this vote record
	key := GetKey(pvs.prefix, proposalID, vote)
	msg, err := pvs.store.Get(key)
	if err != nil {
		logger.Errorf("%v, can't participate in voting", info)
		return ErrSetupVotingValidator
	}

	// Deserialize it
	pv, err := (&ProposalVote{}).FromBytes(msg)
	if err != nil {
		logger.Errorf("%v, deserialize proposal vote failed", info)
		return ErrSetupVotingValidator
	}

	// Update opinion field only
	pv.Opinion = vote.Opinion
	value := pv.Bytes()
	err = pvs.store.Set(key, value)
	if err != nil {
		logger.Errorf("%v, storage failure", info)
		return ErrSetupVotingValidator
	}
	logger.Detail(info)

	return nil
}

// Delete all voting records under a proposalID
func (pvs *ProposalVoteStore) Delete(proposalID ProposalID) error {
	info := fmt.Sprintf("Vote Delete: proposalID= %v", proposalID)

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
	logger.Detail(info)

	return nil
}

//ResultSoFar check and see if a proposal has already passed or failed
//Proposal passed if passPercent already achieved
//Proposal never pass if received enough NEGATIVE votes
func (pvs *ProposalVoteStore) ResultSoFar(proposalID ProposalID, passPercent int) (*VoteStatus, error) {
	info := fmt.Sprintf("Vote IsPassed: proposalID= %v", proposalID)

	_, votes, err := pvs.GetVotesByID(proposalID)
	if err != nil {
		logger.Errorf("%v, getVotesByID failed", info)
		stat := NewVoteStatus(VOTE_RESULT_TBD, 0, 0, 0)
		return stat, ErrVoteCheckVoteResultFailed
	}

	// Accumulates power of each opinion
	allPower := int64(0)
	eachPower := make([]int64, 4)
	for _, vote := range votes {
		allPower += vote.Power
		eachPower[vote.Opinion] += vote.Power
	}

	// Excludes validators that give up voting in percent calculation
	totalPower := allPower - eachPower[OPIN_GIVEUP]

	// Calculate actual percentage
	yesPower := eachPower[OPIN_POSITIVE]
	noPower := eachPower[OPIN_NEGATIVE]
	yesPercentage := 0.0
	noPercentage := 0.0
	passPercentage := float64(passPercent) / 100.0
	if totalPower > 0 {
		yesPercentage = float64(yesPower) / float64(totalPower)
		noPercentage = float64(noPower) / float64(totalPower)
	}

	// Proposal passed if received enough votes of YES
	if yesPercentage >= passPercentage {
		logger.Detailf("%v, passed, YES percentage= %v", info, yesPercentage)
		stat := NewVoteStatus(VOTE_RESULT_PASSED, yesPower, noPower, allPower)
		return stat, nil
	}
	// Proposal failed if received enough votes of NO
	if (1.0 - noPercentage) < passPercentage {
		logger.Detailf("%v, failed, NO percentage= %v", info, noPercentage)
		stat := NewVoteStatus(VOTE_RESULT_FAILED, yesPower, noPower, allPower)
		return stat, nil
	}

	// Result to be dertermined
	stat := NewVoteStatus(VOTE_RESULT_TBD, yesPower, noPower, allPower)
	return stat, nil
}

// Iterate voting records by proposalID
func (pvs *ProposalVoteStore) IterateByID(proposalID ProposalID, fn func(key []byte, value []byte) bool) (stopped bool) {
	prefix := append(pvs.prefix, proposalID+storage.DB_PREFIX...)
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
	info := fmt.Sprintf("Vote getVotesByID: proposalID= %v", proposalID)

	succeed := true
	addrs := make([]keys.Address, 0)
	votes := make([]*ProposalVote, 0)
	pvs.IterateByID(proposalID, func(key []byte, value []byte) bool {
		vote, err := (&ProposalVote{}).FromBytes(value)
		if err != nil {
			logger.Errorf("%v, key= %v, deserialize proposal vote failed", info, key)
			succeed = false
			return true
		}
		votes = append(votes, vote)
		prefix_len := len(append(pvs.prefix, proposalID+storage.DB_PREFIX...))
		addr := key[prefix_len:]
		addrs = append(addrs, addr)
		return false
	})

	if !succeed {
		errMsg := fmt.Sprintf("%v, operation failed", info)
		logger.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}
	// Caused by invalid/deleted proposalID
	if len(votes) == 0 {
		errMsg := fmt.Sprintf("%v, no votes records found", info)
		logger.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	return addrs, votes, nil
}

func GetKey(prefix []byte, proposalID ProposalID, vote *ProposalVote) []byte {
	return storage.StoreKey(string(prefix) + string(proposalID) + storage.DB_PREFIX + string(vote.Validator))
}
