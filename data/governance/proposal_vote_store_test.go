package governance

import (
	"encoding/hex"
	"testing"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
)

const (
	totalNum    = 8
	proposalID  = "id_test_proposal"
	passPercent = 67
	hex0        = "72143ADE3D941025468792311A0AB38D5085E15A"
	hex1        = "821437DF3C9410254A8792311A0A13255085E157"
	hex2        = "92143CDE3D941025468792311A0AB38D5085E151"
	hex3        = "A2143AD5793B910D9410225ADC68B38D5085E11C"
	hex4        = "B214863A4B8B910D941022556AAF23685085E11C"
	hex5        = "C314863A4B8B910963517CC0DC68B38D5085F00D"
	hex6        = "D25479AAF1C259910225ADCA01FF674D55744421"
	hex7        = "E25479AAF1C2599D9410225ADC68B399DDA249AB"
)

// test setup
func setupProposalVoteStore(t *testing.T) (*ProposalVoteStore, []keys.Address) {
	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewChainState("chainstate", db)
	pvs := NewProposalVoteStore("pvs", storage.NewState(cs))

	// participating validators
	addrs := make([]keys.Address, 8)
	addr0, _ := hex.DecodeString(hex0)
	addr1, _ := hex.DecodeString(hex1)
	addr2, _ := hex.DecodeString(hex2)
	addr3, _ := hex.DecodeString(hex3)
	addr4, _ := hex.DecodeString(hex4)
	addr5, _ := hex.DecodeString(hex5)
	addr6, _ := hex.DecodeString(hex6)
	addr7, _ := hex.DecodeString(hex7)
	addrs[0] = addr0
	addrs[1] = addr1
	addrs[2] = addr2
	addrs[3] = addr3
	addrs[4] = addr4
	addrs[5] = addr5
	addrs[6] = addr6
	addrs[7] = addr7
	// init voting validators
	for i := 0; i < totalNum; i++ {
		vote := NewProposalVote(addrs[i], OPIN_POSITIVE, 1)
		err := pvs.Setup(proposalID, vote)
		assert.Nil(t, err)
	}
	pvs.store.Commit()

	return pvs, addrs
}

func setupProposalVotes(t *testing.T, pvs *ProposalVoteStore, addrs []keys.Address, positive, negative, giveup, power int) {
	// create vote objects
	votes := make([]*ProposalVote, totalNum)
	for i := 0; i < totalNum; i++ {
		vote := NewProposalVote(addrs[i], OPIN_UNKNOWN, int64(power))
		votes[i] = vote
	}
	// setup POSITIVE votes
	curIndex := 0
	for j := 0; j < positive; j++ {
		vote := votes[curIndex+j]
		vote.Opinion = OPIN_POSITIVE
		err := pvs.Update(proposalID, vote)
		assert.Nil(t, err)
	}
	// setup NEGATIVE votes
	curIndex += positive
	for j := 0; j < negative; j++ {
		vote := votes[curIndex+j]
		vote.Opinion = OPIN_NEGATIVE
		err := pvs.Update(proposalID, vote)
		assert.Nil(t, err)
	}
	// setup GIVEUP votes
	curIndex += negative
	for j := 0; j < giveup; j++ {
		vote := votes[curIndex+j]
		vote.Opinion = OPIN_GIVEUP
		err := pvs.Update(proposalID, vote)
		assert.Nil(t, err)
	}
}

func checkProposalVotes(t *testing.T, pvs *ProposalVoteStore, addrs []keys.Address, positive, negative, giveup, power int) {
	validators, votes, err := pvs.GetVotesByID(proposalID)
	assert.Nil(t, err)
	// check POSITIVE votes
	power64 := int64(power)
	curIndex := 0
	for j := 0; j < positive; j++ {
		vote := votes[curIndex+j]
		assert.Equal(t, OPIN_POSITIVE, vote.Opinion)
		assert.Equal(t, power64, vote.Power)
		assert.Equal(t, addrs[curIndex+j], validators[curIndex+j])
	}
	// check NEGATIVE votes
	curIndex += positive
	for j := 0; j < negative; j++ {
		vote := votes[curIndex+j]
		assert.Equal(t, OPIN_NEGATIVE, vote.Opinion)
		assert.Equal(t, power64, vote.Power)
		assert.Equal(t, addrs[curIndex+j], validators[curIndex+j])
	}
	// check GIVEUP votes
	curIndex += negative
	for j := 0; j < giveup; j++ {
		vote := votes[curIndex+j]
		assert.Equal(t, OPIN_GIVEUP, vote.Opinion)
		assert.Equal(t, power64, vote.Power)
		assert.Equal(t, addrs[curIndex+j], validators[curIndex+j])
	}
	// check UNKNOWN votes
	curIndex += giveup
	unknown := totalNum - positive - negative - giveup
	for j := 0; j < unknown; j++ {
		vote := votes[curIndex+j]
		assert.Equal(t, OPIN_UNKNOWN, vote.Opinion)
		assert.Equal(t, power64, vote.Power)
		assert.Equal(t, addrs[curIndex+j], validators[curIndex+j])
	}
}

func TestNewProposalVoteStore(t *testing.T) {
	pvs, _ := setupProposalVoteStore(t)
	assert.NotEmpty(t, pvs)
}

func TestProposalVoteStore_Setup(t *testing.T) {
	t.Run("test setup initial voting validators", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		checkProposalVotes(t, pvs, addrs, 0, 0, 0, 1)
	})
}

func TestProposalVoteStore_Update(t *testing.T) {
	t.Run("test updated vote records of a proposal, should match", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 5, 3, 0, 2)
		checkProposalVotes(t, pvs, addrs, 5, 3, 0, 2)
	})
	t.Run("test vote using invalid address, should fail", func(t *testing.T) {
		pvs, _ := setupProposalVoteStore(t)
		hexNotExist := "1111111111111111111792311A0AB38D5085E15A"
		addrInvalid, _ := hex.DecodeString(hexNotExist)
		vote := NewProposalVote(addrInvalid, OPIN_NEGATIVE, 1)
		err := pvs.Update(proposalID, vote)
		assert.NotNil(t, err)
	})
}

func TestProposalVoteStore_Delete(t *testing.T) {
	t.Run("test deleting vote records of an initial proposal", func(t *testing.T) {
		pvs, _ := setupProposalVoteStore(t)
		err := pvs.Delete(proposalID)
		assert.Nil(t, err)

		addrs, votes, err := pvs.GetVotesByID(proposalID)
		assert.Error(t, err)
		assert.Nil(t, addrs)
		assert.Nil(t, votes)
	})
	t.Run("test deleting vote records of a voted proposal", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 5, 3, 0, 2)
		err := pvs.Delete(proposalID)

		assert.Nil(t, err)
		addrs, votes, err := pvs.GetVotesByID(proposalID)
		assert.Error(t, err)
		assert.Nil(t, addrs)
		assert.Nil(t, votes)
	})
	t.Run("test deleting vote records of invalid proposal", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 5, 3, 0, 2)
		err := pvs.Delete("fake_proposal_id")
		assert.Nil(t, err)
	})
}

func TestProposalVoteStore_ProposalPassed(t *testing.T) {
	t.Run("test a proposal that passed successfully, nobody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 6, 2, 0, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_PASSED, result)
	})
	t.Run("test a proposal that passed successfully, somebody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 5, 2, 1, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_PASSED, result)
	})
	t.Run("test a proposal that passed successfully, somebody give up & somebody give no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 5, 1, 1, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_PASSED, result)
	})
	t.Run("test a proposal that passed successfully, all support", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 8, 0, 0, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_PASSED, result)
	})
	t.Run("test a proposal that does not exist, should return error", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 8, 0, 0, 2)
		_, err := pvs.ResultSoFar("fake_proposal_id", passPercent)
		assert.NotNil(t, err)
	})
}

func TestProposalVoteStore_ProposalFailed(t *testing.T) {
	t.Run("test a proposal that failed to pass, nobody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 5, 3, 0, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_FAILED, result)
	})
	t.Run("test a proposal that failed to pass, somebody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 4, 3, 1, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_FAILED, result)
	})
	t.Run("test a proposal that failed to pass, somebody give up & somebody give no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 3, 3, 1, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_FAILED, result)
	})
	t.Run("test a proposal that failed to pass, all disagree", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 0, 8, 0, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_FAILED, result)
	})
}

func TestProposalVoteStore_ProposalTBD(t *testing.T) {
	t.Run("test a proposal that to be determined, somebody give no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 4, 2, 0, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_TBD, result)
	})
	t.Run("test a proposal that to be determined, somebody give up & somebody give no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 3, 1, 1, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_TBD, result)
	})
	t.Run("test a proposal that to be determined, all give up or no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 0, 0, 2, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_TBD, result)
	})
	t.Run("test a proposal that to be determined, all no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 0, 0, 0, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_TBD, result)
	})
	t.Run("test a proposal that to be determined, all give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore(t)
		setupProposalVotes(t, pvs, addrs, 0, 0, 8, 2)
		result, _ := pvs.ResultSoFar(proposalID, passPercent)
		assert.Equal(t, VOTE_RESULT_TBD, result)
	})
}
