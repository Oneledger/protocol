package governance

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	store    *ProposalFundStore
	cs       *storage.State
	address  keys.Address
	address2 keys.Address
	address3 keys.Address
	ID1      ProposalID
	ID2      ProposalID
	ID3      ProposalID
)

func init() {
	fmt.Println("####### TESTING PROPOSAL FUND STORE #######")
	memDb := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("balance", memDb))
	store = NewProposalFundStore("test", cs)
	generateAddresses()
	generateIDs()
}

func generateAddresses() {
	pub, _, _ := keys.NewKeyPairFromTendermint()
	h, _ := pub.GetHandler()
	address = h.Address()

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()
	address2 = h2.Address()
	pub3, _, _ := keys.NewKeyPairFromTendermint()
	h3, _ := pub3.GetHandler()
	address3 = h3.Address()
}

func generateIDs() {
	ID1 = generateProposalID("Test")
	ID2 = generateProposalID("Test1")
	ID3 = generateProposalID("Test2")

}
func TestProposalFundStore_AddFunds(t *testing.T) {
	fmt.Println("Adding New Proposer for funding")
	err := store.AddFunds(ID1, address, balance.NewAmount(100))
	assert.NoError(t, err, "")
	cs.Commit()
	err = store.AddFunds(ID2, address, balance.NewAmount(100))
	assert.NoError(t, err, "")
	cs.Commit()
	err = store.AddFunds(ID1, address2, balance.NewAmount(1000))
	assert.NoError(t, err, "")
	err = store.AddFunds(ID2, address2, balance.NewAmount(120))
	assert.NoError(t, err, "")
	cs.Commit()

}

func TestNewProposalFundStore_Delete(t *testing.T) {
	//fmt.Println("Deleting fund record ID : ", ID1, "| address :", address)
	ok, err := store.DeleteFunds(ID1, address)
	if err != nil {
		fmt.Println("Error Deleting : ", err)
		return
	}
	cs.Commit()
	assert.True(t, ok, "")
}

func TestProposalFundStore_Iterate(t *testing.T) {
	//fmt.Println("Iterating Stores")
	IDLIST := []ProposalID{ID2, ID1}
	ok := store.iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) bool {
		//fmt.Println("ProposalID : ", proposalID, "ProposalAddress :", fundingAddr, "Proposal Amount :", amt.BigInt())
		assert.Contains(t, IDLIST, proposalID, "")
		return false
	})
	assert.True(t, ok, "")
}

//
func TestProposalFundStore_GetFundersForProposalID(t *testing.T) {
	//fmt.Println("Get Funders for ID :  ", ID1)
	funds := store.GetFundersForProposalID(ID1, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		return ProposalFund{
			id:            proposalID,
			address:       fundingAddr,
			fundingAmount: amt,
		}
	})
	//for _, fund := range funds {
	//	fund.Print()
	//}
	assert.EqualValues(t, 1, len(funds), "")
}

//
func TestProposalFundStore_GetProposalForFunder(t *testing.T) {
	//fmt.Println("Get Funders for Address :", address2)

	funds := store.GetProposalsForFunder(address2, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		return ProposalFund{
			id:            proposalID,
			address:       fundingAddr,
			fundingAmount: amt,
		}
	})
	//for _, fund := range funds {
	//	fund.Print()
	//}
	assert.EqualValues(t, 2, len(funds), "")
}

func TestProposalFund_getCurrentFunds(t *testing.T) {
	//fmt.Println("Getting Total fund for ProposalID")
	currentFunds := GetCurrentFunds(ID1, store)
	funds := currentFunds.BigInt().Int64()
	assert.EqualValues(t, int64(1000), funds, "")
}
