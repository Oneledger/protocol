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
)

func init() {
	db := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("balance", db))
	store = NewProposalFundStore("test", cs)
	pub, _, _ := keys.NewKeyPairFromTendermint()
	h, _ := pub.GetHandler()
	address = h.Address()

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()
	address2 = h2.Address()
}

func TestProposalFundStore_AddFunds(t *testing.T) {
	fmt.Println("Adding New Proposer for funding")
	err := store.AddFunds("ID1", address, NewAmount(10))
	assert.NoError(t, err, "")
	cs.Commit()
	err = store.AddFunds("ID1", address, NewAmount(100))
	assert.NoError(t, err, "")
	cs.Commit()
	err = store.AddFunds("ID2", address2, NewAmount(120))
	assert.NoError(t, err, "")
	cs.Commit()

}

func TestProposalFundStore_Iterate(t *testing.T) {
	fmt.Println("Iterating Stores")
	store.Iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) bool {
		fmt.Println("ProposalID : ", proposalID, "ProposalAddress :", fundingAddr, "Proposal Amount :", amt.BigInt())
		return false
	})
}

func TestProposalFundStore_GetFundersForProposalID(t *testing.T) {
	fmt.Println("Get Funders for ID")

	ok := GetFundersForProposalID(store, "ID1", func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) Funder {
		return Funder{
			id:            proposalID,
			address:       fundingAddr.String(),
			fundingAmount: amt.String(),
		}
	})
	fmt.Println("Found Proposals :", foundProposals)
	assert.True(t, ok, "")
}
