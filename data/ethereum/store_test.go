package ethereum

import (
	"fmt"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/db"
	"testing"
)
var store *TrackerStore
var cs *storage.State

func init() {
	db := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("balance", db))
	store = NewTrackerStore("test", cs)
}

func TestTrackerStore_Get(t *testing.T) {


	h := &common.Hash{}
	h.SetBytes([]byte("testhash"))
	tracker := NewTracker([]byte("tracker1addr"),[]byte("signedeth"), ethereum.TrackerName(*h), []keys.Address{keys.Address("s")})
	err := store.Set(*tracker)
	assert.NoError(t, err, "set")

	trackerNew, err := store.Get(tracker.TrackerName)
	assert.NoError(t, err, "get")
	assert.Equal(t, tracker, trackerNew, "equal")

}


func TestTrackerStore_Iterate(t *testing.T) {
	cs.Commit()
	store.Iterate(func(name *ethereum.TrackerName, tracker *Tracker) bool {
		fmt.Println(name, tracker)
		return false
	})
}