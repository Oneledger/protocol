package ethereum

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"fmt"
	"testing"
)

var (
	numberOfPrivKey = 16
	threshold       = 11
	tracker         Tracker

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
	tracker = *NewTracker(addresses[0], []byte("test"), *h, addresses)
}

func TestTracker_AddVote(t *testing.T) {
	for i, addr := range addresses {
		if i >= threshold {
			continue
		}
		err := tracker.AddVote(addr, int64(i))
		assert.NoError(t, err, "add vote error")
	}
}

func TestTracker_Finalized(t *testing.T) {
	ok := tracker.Finalized()
	assert.True(t, ok)
}

func TestTracker_TestStateMachine(t *testing.T) {
	index := 0
	validatorIndex := 0

	for index < 100 {
		if index%3 == 0 {
			if validatorIndex < len(addresses) {
				//Add Vote to tracker
				tracker.AddVote(addresses[validatorIndex], int64(validatorIndex))
				validatorIndex++
			}
		}

		//Transition Engine Process
		//Engine.Process(tracker.TrackerName, tracker, tracker.State)

		index++
	}
}

func TestTracker_CheckIfVoted(t *testing.T) {

	h := &common.Hash{}
	h.SetBytes([]byte("test"))
	tracker = *NewTracker(addresses[0], []byte("test"), *h, addresses)
	tracker.AddVote(addresses[1],1)
	index0,ok := tracker.CheckIfVoted(addresses[0])
	assert.False(t, ok)
	index1,ok := tracker.CheckIfVoted(addresses[1])
	assert.True(t, ok)

	assert.Equal(t, index0,int64(0))
	assert.Equal(t,index1,int64(1))
	assert.Equal(t, ok,true)
}


func TestTracker_GetVotes(t *testing.T) {
	//tracker.FinalityVotes = int64(0)
	tracker.AddVote(addresses[0],0)
	tracker.AddVote(addresses[1],1)
	fmt.Println("VOTES: " , tracker.GetVotes())
	assert.Equal(t,2,tracker.GetVotes())
}
