package governance

import (
	"fmt"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"testing"
)

const (
	numPrivateKeys = 5
	numProposals   = 10

	codeChange   = 2
	configUpdate = 3
	general      = 4
)

var (
	addrList    []keys.Address
	proposals   []*Proposal
	proposalOpt ProposalOptions

	govStore      *Store
	proposalStore *ProposalStore
)

func init() {
	fmt.Println("####### TESTING PROPOSAL STORE #######")

	//Generate key pairs for proposers
	for i := 0; i < numPrivateKeys; i++ {
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()
		addrList = append(addrList, h.Address())
	}

	//Create new proposal options
	proposalOpt.CodeChange = options{
		InitialFunding:  codeChange,
		FundingDeadline: codeChange,
		FundingGoal:     codeChange,
		VotingDeadline:  codeChange,
	}

	proposalOpt.ConfigUpdate = options{
		InitialFunding:  configUpdate,
		FundingDeadline: configUpdate,
		FundingGoal:     configUpdate,
		VotingDeadline:  configUpdate,
	}

	proposalOpt.General = options{
		InitialFunding:  general,
		FundingDeadline: general,
		FundingGoal:     general,
		VotingDeadline:  general,
	}

	//Create new proposals
	for i := 0; i < numProposals; i++ {
		j := i / 2      //address list ranges from 0 - 4
		k := i/4 + 0x20 //proposal type ranges from 0x20 - 0x22

		proposer := addrList[j]

		var opt options
		switch ProposalType(k) {
		case ProposalTypeConfigUpdate:
			opt = proposalOpt.ConfigUpdate

		case ProposalTypeCodeChange:
			opt = proposalOpt.CodeChange

		case ProposalTypeGeneral:
			opt = proposalOpt.General
		}

		proposals = append(proposals, NewProposal(ProposalType(k), "Test Proposal", proposer,
			opt.FundingDeadline, opt.FundingGoal, opt.VotingDeadline))
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	//Create Governance store
	govStore = NewStore("g", cs)

	//Create Proposal store
	proposalStore = NewProposalStore("p_active", "p_passed", "p_failed", cs)
}

func TestProposalStore_Set(t *testing.T) {
	err := proposalStore.Set(proposals[0])
	assert.Equal(t, nil, err)

	proposal, err := proposalStore.Get(proposals[0].ProposalID)
	assert.Equal(t, nil, err)

	assert.Equal(t, proposal.ProposalID, proposals[0].ProposalID)
}

func TestProposalStore_Exists(t *testing.T) {
	exists := proposalStore.Exists(proposals[0].ProposalID)
	assert.Equal(t, true, exists)

	exists = proposalStore.Exists(proposals[1].ProposalID)
	assert.Equal(t, false, exists)
}

func TestProposalStore_Delete(t *testing.T) {
	_, err := proposalStore.Get(proposals[0].ProposalID)
	assert.Equal(t, nil, err)

	res, err := proposalStore.Delete(proposals[0].ProposalID)
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	_, err = proposalStore.Get(proposals[0].ProposalID)
	assert.NotEqual(t, nil, err)
}

func TestProposalStore_Iterate(t *testing.T) {
	for _, val := range proposals {
		_ = proposalStore.Set(val)
	}
	proposalStore.state.Commit()

	proposalCount := 0
	proposalStore.Iterate(func(id ProposalID, proposal *Proposal) bool {
		proposalCount++
		return false
	})

	assert.Equal(t, numProposals, proposalCount)
}

func TestProposalStore_IterateProposer(t *testing.T) {
	for _, val := range addrList {
		proposer := val

		proposalCount := 0
		proposalStore.IterateProposer(func(id ProposalID, proposal *Proposal) bool {
			proposalCount++
			return false
		}, proposer)
		assert.Equal(t, 2, proposalCount)
	}
}

func TestProposalStore_IterateProposalType(t *testing.T) {
	proposalCount := 0
	proposalStore.IterateProposalType(func(id ProposalID, proposal *Proposal) bool {
		proposalCount++
		return false
	}, ProposalTypeCodeChange)
	assert.Equal(t, 4, proposalCount)

	proposalCount = 0
	proposalStore.IterateProposalType(func(id ProposalID, proposal *Proposal) bool {
		proposalCount++
		return false
	}, ProposalTypeGeneral)
	assert.Equal(t, 2, proposalCount)
}

func TestProposalStore_SetOptions(t *testing.T) {
	err := govStore.SetProposalOptions(proposalOpt)
	assert.Equal(t, nil, err)

	propOpt, err := govStore.GetProposalOptions()
	assert.Exactly(t, &proposalOpt, propOpt)
}
