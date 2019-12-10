package ethereum

import (
	"fmt"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/utils/transition"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	numberOfPrivKey = 16
	threshold       = 11
	trackerLock     Tracker
	trackerRedeem   Tracker

	privKeys  []keys.PrivateKey = make([]keys.PrivateKey, numberOfPrivKey)
	addresses []keys.Address    = make([]keys.Address, numberOfPrivKey)
)

func init() {
	for i := 0; i < numberOfPrivKey; i++ {
		pub, priv, _ := keys.NewKeyPairFromTendermint()
		privKeys[i] = priv

		h, _ := pub.GetHandler()
		addresses[i] = h.Address()
	}
	h := &common.Hash{}
	h.SetBytes([]byte("test"))
	trackerLock = *NewTracker(ProcessTypeLock, addresses[0], []byte("test"), *h, addresses)
	trackerRedeem = *NewTracker(ProcessTypeRedeem, addresses[0], []byte("test"), *h, addresses)
}

func TestTracker_AddVote(t *testing.T) {
	fmt.Println("***Test AddVote()***")
	for i, addr := range addresses {
		if i >= threshold {
			continue
		}
		err := trackerLock.AddVote(addr, int64(i), true)
		assert.NoError(t, err, "add vote error")
	}
	yesVotes, _ := trackerLock.GetVotes()
	assert.Equal(t, yesVotes, 11)
}

func TestTracker_Finalized(t *testing.T) {
	fmt.Println("***Test Finalized()***")
	ok := trackerLock.Finalized()
	assert.True(t, ok)
}

func TestTracker_CheckIfVoted(t *testing.T) {
	fmt.Println("***Test CheckIfVoted()***")
	addr := addresses[12]
	err := trackerLock.AddVote(addr, int64(12), true)
	if err != nil {
		t.Errorf(err.Error())
	}

	_, voted := trackerLock.CheckIfVoted(addr)
	assert.True(t, voted)
}

func TestTracker_GetVotes(t *testing.T) {
	fmt.Println("***Test GetVotes()***")
	yesVotes, _ := trackerLock.GetVotes()
	assert.Equal(t, yesVotes, threshold+1) //Vote was added in TestTracker_CheckIfVoted
}

type trackerCase struct {
	Input  TrackerState
	Output string
}

func TestTracker_NextStep(t *testing.T) {
	fmt.Println("***Test NextStep()***")
	testCases := make(map[int]trackerCase)
	testCases[0] = trackerCase{New, BROADCASTING}
	testCases[1] = trackerCase{BusyBroadcasting, FINALIZING}
	testCases[2] = trackerCase{BusyFinalizing, FINALIZE}
	testCases[3] = trackerCase{Finalized, MINTING}
	testCases[4] = trackerCase{Released, CLEANUP}
	testCases[5] = trackerCase{-1, transition.NOOP}

	//Test Lock
	for _, test := range testCases {
		trackerLock.State = test.Input
		assert.Equal(t, trackerLock.NextStep(), test.Output)
	}

	testCases[0] = trackerCase{New, SIGNING}
	testCases[1] = trackerCase{BusyBroadcasting, VERIFYREDEEM}
	testCases[2] = trackerCase{BusyFinalizing, REDEEMCONFIRM}
	testCases[3] = trackerCase{Finalized, BURN}

	//Test Redeem
	for _, test := range testCases {
		trackerRedeem.State = test.Input
		assert.Equal(t, trackerRedeem.NextStep(), test.Output)
	}

}
