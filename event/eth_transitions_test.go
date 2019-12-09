package event

import (
	"fmt"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils/transition"
	"github.com/ethereum/go-ethereum/common"
	"github.com/magiconair/properties/assert"
	"strconv"
	"testing"
)

var (
	numberOfPrivKey = 16
	threshold       = 11
	tracker         ethereum.Tracker

	privKeys    = make([]keys.PrivateKey, numberOfPrivKey)
	addresses   = make([]keys.Address, numberOfPrivKey)
	chainState  storage.ChainState
	btcTrackers bitcoin.TrackerStore
	ethTrackers ethereum.TrackerStore
	jobStore    jobs.JobStore

	testCases map[int]Case
)

type Case struct {
	//Input Variables
	InState      ethereum.TrackerState
	CurrNodeAddr keys.Address
	NumVotes     int
	CompleteTask bool

	//Output Variables
	OutState ethereum.TrackerState
	Err      error
}

// global setup
func setup() {
	fmt.Println("Initializing Database For Transition Test...")
	db, err := storage.GetDatabase("chainstate", "test_dbpath", "db")
	if err != nil {
		fmt.Println("error initializing database")
	}

	chainState = *storage.NewChainState("chainstate", db)
	btcTrackers = *bitcoin.NewTrackerStore("btct", storage.NewState(&chainState))
	ethTrackers = *ethereum.NewTrackerStore("etht", storage.NewState(&chainState))
	jobStore = *jobs.NewJobStore(*config.DefaultServerConfig(), "test_dbpath")
}

func init() {

	for i := 0; i < numberOfPrivKey; i++ {
		pub, priv, _ := keys.NewKeyPairFromTendermint()
		privKeys[i] = priv

		h, _ := pub.GetHandler()
		addresses[i] = h.Address()
	}

	testCases = make(map[int]Case)

	//Successful Transitions Forward
	testCases[0] = Case{ethereum.New, addresses[0], 0, false, ethereum.BusyBroadcasting, nil}
	testCases[1] = Case{ethereum.BusyBroadcasting, addresses[0], 0, false, ethereum.BusyBroadcasting, nil}
	testCases[2] = Case{ethereum.BusyBroadcasting, addresses[0], 1, false, ethereum.BusyFinalizing, nil}
	testCases[3] = Case{ethereum.BusyFinalizing, addresses[0], threshold - 1, false, ethereum.BusyFinalizing, nil}
	testCases[4] = Case{ethereum.BusyFinalizing, addresses[0], threshold, false, ethereum.Finalized, nil}
	testCases[5] = Case{ethereum.Finalized, addresses[0], threshold, true, ethereum.Released, nil}

	setup()
}

func TestTransitions(t *testing.T) {
	fmt.Println("*** RUNNING ETH TRANSITION TEST ***")

	for i, testCase := range testCases {
		t.Run("Testing case "+strconv.Itoa(i), func(t *testing.T) {

			fmt.Println("***** Test Case: ", i, " *****")
			validatorIndex := 0
			h := &common.Hash{}
			h.SetBytes([]byte("test"))
			tracker = *ethereum.NewTracker(ethereum.ProcessTypeLock, testCase.CurrNodeAddr, []byte("test"), *h, addresses)

			//Add Votes
			for validatorIndex < testCase.NumVotes {
				//Add Vote to tracker
				err := tracker.AddVote(addresses[validatorIndex], int64(validatorIndex), true)
				if err != nil {
					t.Errorf(err.Error())
				}
				fmt.Println("Vote Added at index: ", validatorIndex)
				validatorIndex++
			}

			//Set Status
			tracker.State = testCase.InState

			ctx := ethereum.TrackerCtx{
				Tracker:      &tracker,
				TrackerStore: &ethTrackers,
				JobStore:     &jobStore,
				CurrNodeAddr: testCase.CurrNodeAddr,
				Validators:   nil,
			}

			fmt.Println("In State:", transition.Status(tracker.State))

			_, err := EthLockEngine.Process(tracker.NextStep(), ctx, transition.Status(tracker.State))
			if err != nil {
				t.Errorf(err.Error())
			}

			//Validate Output
			assert.Equal(t, transition.Status(testCase.OutState), transition.Status(tracker.State))
			assert.Equal(t, testCase.Err, err)
		})
	}
}
