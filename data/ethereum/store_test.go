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

var (
	store *TrackerStore
	cs    *storage.State
)

func init() {

	for i := 0; i < numberOfPrivKey; i++ {
		pub, priv, _ := keys.NewKeyPairFromTendermint()
		privKeys[i] = priv

		h, _ := pub.GetHandler()
		addresses[i] = h.Address()
	}

	db := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("balance", db))
	store = NewTrackerStore("test", cs)
}

func TestTrackerStore_Get(t *testing.T) {
	fmt.Println("***Test Tracker Store Get***")

	h := &common.Hash{}
	h.SetBytes([]byte("testhash"))
	tracker := NewTracker(ProcessTypeLock, []byte("tracker1addr"), []byte("signedeth"), ethereum.TrackerName(*h), addresses)

	i := 0
	for i < 4 {
		fmt.Println("address: ", addresses[i], "index: ", i)
		tracker.AddVote(addresses[i], int64(i), true)

		fmt.Println(fmt.Sprintf("Finality Votes from cache: %b", tracker.FinalityVotes))

		yes, _ := tracker.GetVotes()
		fmt.Println("tracker.GetVotes() returns: ", yes)

		fmt.Println("Saving tracker to cache")
		err := store.Set(tracker)
		assert.NoError(t, err, "")

		fmt.Println("Getting tracker from cache")
		trackerNew, err := store.Get(tracker.TrackerName)
		assert.NoError(t, err, "")

		fmt.Println(fmt.Sprintf("Finality Votes from cache: %b", trackerNew.FinalityVotes))

		yesNew, _ := trackerNew.GetVotes()
		fmt.Println("tracker.GetVotes() returns: ", yesNew)

		assert.Equal(t, yes, yesNew)

		i++
	}
}

func TestTrackerStore_Iterate(t *testing.T) {
	fmt.Println("***Test Tracker Store Iterate***")
	cs.Commit()
	store.Iterate(func(name *ethereum.TrackerName, tracker *Tracker) bool {
		fmt.Println("Name: ", name, "Tracker: ", tracker)
		return false
	})
}
