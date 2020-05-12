package governance

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	store    *ProposalFundStore
	cs       *storage.State
	address  keys.Address
	address2 keys.Address
	ID1      ProposalID
	ID2      ProposalID
)

func init() {
	fmt.Println("####### TESTING PROPOSAL FUND STORE #######")
	db := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("balance", db))
	store = NewProposalFundStore("test", cs)
	pub, _, _ := keys.NewKeyPairFromTendermint()
	h, _ := pub.GetHandler()
	address = h.Address()

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()
	address2 = h2.Address()
	ID1 = generateProposalID("Test")
	ID2 = generateProposalID("Test")
}

func TestProposalFundStore_AddFunds(t *testing.T) {
	fmt.Println("Adding New Proposer for funding")
	err := store.AddFunds(ID1, address, NewAmount(100))
	assert.NoError(t, err, "")
	cs.Commit()
	//err = store.AddFunds("ID1", address, NewAmount(100))
	//assert.NoError(t, err, "")
	//cs.Commit()
	err = store.AddFunds(ID2, address, NewAmount(100))
	assert.NoError(t, err, "")
	cs.Commit()
	err = store.AddFunds(ID1, address2, NewAmount(1000))
	assert.NoError(t, err, "")
	err = store.AddFunds(ID2, address2, NewAmount(120))
	assert.NoError(t, err, "")
	cs.Commit()

}

func TestNewProposalFundStore_Delete(t *testing.T) {
	fmt.Println("Deleting fund record ID : ", ID1, "| address :", address)
	ok, err := store.DeleteFunds(ID1, address)
	if err != nil {
		fmt.Println("Error Deleting : ", err)
		return
	}
	fmt.Println("OK :", ok)
	cs.Commit()
	assert.True(t, ok, "")
}

func TestProposalFundStore_Iterate(t *testing.T) {
	fmt.Println("Iterating Stores")
	store.iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) bool {
		fmt.Println("ProposalID : ", proposalID, "ProposalAddress :", fundingAddr, "Proposal Amount :", amt.BigInt())
		return false
	})
}

//
func TestProposalFundStore_GetFundersForProposalID(t *testing.T) {
	fmt.Println("Get Funders for ID :  ", ID1)
	funds := store.GetFundersForProposalID(ID1, func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) ProposalFund {
		return ProposalFund{
			id:            proposalID,
			address:       fundingAddr,
			fundingAmount: *amt,
		}
	})
	for _, fund := range funds {
		fund.Print()
	}
	assert.EqualValues(t, 1, len(funds), "")
}

//
func TestProposalFundStore_GetProposalForFunder(t *testing.T) {
	fmt.Println("Get Funders for Address :", address2)

	funds := store.GetProposalsForFunder(address2, func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) ProposalFund {
		return ProposalFund{
			id:            proposalID,
			address:       fundingAddr,
			fundingAmount: *amt,
		}
	})
	for _, fund := range funds {
		fund.Print()
	}
	assert.EqualValues(t, 2, len(funds), "")
}
