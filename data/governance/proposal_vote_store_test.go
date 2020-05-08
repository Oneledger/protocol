package governance

import (
	"encoding/hex"
	"testing"

	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/storage"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/data/keys"
)

const (
	totalNum   = 8
	proposalID = "id_test_proposal"
	hex0       = "72143ADE3D941025468792311A0AB38D5085E15A"
	hex1       = "821437DF3C9410254A8792311A0A13255085E157"
	hex2       = "92143CDE3D941025468792311A0AB38D5085E151"
	hex3       = "A2143AD5793B910D9410225ADC68B38D5085E11C"
	hex4       = "B214863A4B8B910D941022556AAF23685085E11C"
	hex5       = "C314863A4B8B910963517CC0DC68B38D5085F00D"
	hex6       = "D25479AAF1C259910225ADCA01FF674D55744421"
	hex7       = "E25479AAF1C2599D9410225ADC68B399DDA249AB"
)

// test setup
func setupProposalVoteStore() (*ProposalVoteStore, []keys.Address) {
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
		pvs.Setup(proposalID, addrs[i], 1)
	}
	return pvs, addrs
}

func setupProposalVotes(pvs *ProposalVoteStore, addrs []keys.Address, positive, negative, giveup, power int) {
	// create vote objects
	votes := make([]*ProposalVote, totalNum)
	for i := 0; i < totalNum; i++ {
		vote := &ProposalVote{
			ProposalID: proposalID,
			Opinion:    UNKNOWN,
			Power:      int64(power),
		}
		votes[i] = vote
	}
	// setup positive votes
	curIndex := 0
	for j := 0; j < positive; j, curIndex = j+1, curIndex+1 {
		vote := votes[curIndex+j]
		vote.Opinion = POSITIVE
		pvs.Update(proposalID, addrs[j], vote)
	}
	// setup negative votes
	for j := 0; j < negative; j, curIndex = j+1, curIndex+1 {
		vote := votes[curIndex+j]
		vote.Opinion = NEGATIVE
		pvs.Update(proposalID, addrs[j], vote)
	}
	// setup giveup votes
	for j := 0; j < giveup; j, curIndex = j+1, curIndex+1 {
		vote := votes[curIndex+j]
		vote.Opinion = GIVEUP
		pvs.Update(proposalID, addrs[j], vote)
	}
}

func checkProposalVotes(t *testing.T, pvs *ProposalVoteStore, addrs []keys.Address, positive, negative, giveup, power int) {
	validators, votes, err := pvs.GetVotesByID(proposalID)
	assert.Nil(t, err)
	// check positive votes
	power64 := int64(power)
	curIndex := 0
	for j := 0; j < positive; j, curIndex = j+1, curIndex+1 {
		vote := votes[curIndex+j]
		assert.Equal(t, POSITIVE, vote.Opinion)
		assert.Equal(t, power64, vote.Power)
		assert.Equal(t, addrs[j], validators[j])
	}
	// check negative votes
	for j := 0; j < negative; j, curIndex = j+1, curIndex+1 {
		vote := votes[curIndex+j]
		assert.Equal(t, NEGATIVE, vote.Opinion)
		assert.Equal(t, power64, vote.Power)
		assert.Equal(t, addrs[j], validators[j])
	}
	// check giveup votes
	for j := 0; j < giveup; j, curIndex = j+1, curIndex+1 {
		vote := votes[curIndex+j]
		assert.Equal(t, GIVEUP, vote.Opinion)
		assert.Equal(t, power64, vote.Power)
		assert.Equal(t, addrs[j], validators[j])
	}
}

func TestNewProposalVoteStore(t *testing.T) {
	pvs, _ := setupProposalVoteStore()
	assert.NotEmpty(t, pvs)
}

func TestProposalVoteStore_Setup(t *testing.T) {
	t.Run("test setup initial voting validators", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		addrsActual, _, err := pvs.GetVotesByID(proposalID)
		assert.Nil(t, err)
		assert.EqualValues(t, addrs, addrsActual)
		checkProposalVotes(t, pvs, addrs, 0, 0, 0, 1)
	})
}

func TestProposalVoteStore_Update(t *testing.T) {
	t.Run("test updated vote records of a proposal, should match", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 5, 3, 0, 2)
		checkProposalVotes(t, pvs, addrs, 5, 3, 0, 2)
	})
}

func TestProposalVoteStore_Delete(t *testing.T) {
	t.Run("test deleting vote records of an initial proposal", func(t *testing.T) {
		pvs, _ := setupProposalVoteStore()
		err := pvs.Delete(proposalID)

		assert.Nil(t, err)
		addrs, votes, err := pvs.GetVotesByID(proposalID)
		assert.Nil(t, err)
		assert.Len(t, addrs, 0)
		assert.Len(t, votes, 0)
	})
	t.Run("test deleting vote records of a voted proposal", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 5, 3, 0, 2)
		err := pvs.Delete(proposalID)

		assert.Nil(t, err)
		addrs, votes, err := pvs.GetVotesByID(proposalID)
		assert.Nil(t, err)
		assert.Len(t, addrs, 0)
		assert.Len(t, votes, 0)
	})
}

func TestProposalVoteStore_ProposalPassed(t *testing.T) {
	t.Run("test a proposal that passed successfully, nobody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 6, 2, 0, 2)
		passed := pvs.IsPassed(proposalID)
		assert.True(t, passed)
	})
	t.Run("test a proposal that passed successfully, somebody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 5, 2, 1, 2)
		passed := pvs.IsPassed(proposalID)
		assert.True(t, passed)
	})
	t.Run("test a proposal that passed successfully, somebody give up & somebody give no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 5, 1, 1, 2)
		passed := pvs.IsPassed(proposalID)
		assert.True(t, passed)
	})
	t.Run("test a proposal that passed successfully, all support", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 8, 0, 0, 2)
		passed := pvs.IsPassed(proposalID)
		assert.True(t, passed)
	})
}

func TestProposalVoteStore_ProposalNotPassed(t *testing.T) {
	t.Run("test a proposal that failed to pass, nobody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 5, 3, 0, 2)
		passed := pvs.IsPassed(proposalID)
		assert.False(t, passed)
	})
	t.Run("test a proposal that failed to pass, somebody give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 4, 3, 1, 2)
		passed := pvs.IsPassed(proposalID)
		assert.False(t, passed)
	})
	t.Run("test a proposal that failed to pass, somebody give up & somebody give no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 3, 3, 1, 2)
		passed := pvs.IsPassed(proposalID)
		assert.False(t, passed)
	})
	t.Run("test a proposal that failed to pass, all give or no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 0, 0, 2, 2)
		passed := pvs.IsPassed(proposalID)
		assert.False(t, passed)
	})
	t.Run("test a proposal that failed to pass, all give up", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 0, 0, 8, 2)
		passed := pvs.IsPassed(proposalID)
		assert.False(t, passed)
	})
	t.Run("test a proposal that failed to pass, all no response", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 0, 0, 0, 2)
		passed := pvs.IsPassed(proposalID)
		assert.False(t, passed)
	})
	t.Run("test a proposal that failed to pass, all disagree", func(t *testing.T) {
		pvs, addrs := setupProposalVoteStore()
		setupProposalVotes(pvs, addrs, 8, 0, 0, 2)
		passed := pvs.IsPassed(proposalID)
		assert.False(t, passed)
	})
}
