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
	store   *ProposalFundStore
	cs      *storage.State
	address keys.Address
)

func init() {
	db := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("balance", db))
	store = NewProposalFundStore("test", cs)
	pub, _, _ := keys.NewKeyPairFromTendermint()
	h, _ := pub.GetHandler()
	address = h.Address()
}

func TestProposalFundStore_AddNewPropososer(t *testing.T) {
	fmt.Println("Adding New Proposer for funding")
	err := store.AddNewPropososer(0, address, 100)
	assert.NoError(t, err, "")
	err = store.AddNewPropososer(0, address, 100)
	assert.Error(t, err, "")
}
