package ethereum

import (
	//"errors"

	"fmt"
	"math/bits"
	"strconv"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils/transition"
)

type TrackerState int

// Tracker
type Tracker struct {
	// State tracks the current state of the tracker, Also used for locking distributed access
	State         TrackerState `json:"state"`
	SignedETHTx   []byte
	Validators    []keys.Address
	ProcessOwner  keys.Address
	FinalityVotes int64
	TrackerName   ethereum.TrackerName
}

//number of validator should be smaller than 64
func NewTracker(owner keys.Address, signedEthTx []byte, name ethereum.TrackerName, validators []keys.Address) *Tracker {

	return &Tracker{
		State:        New,
		TrackerName:  name,
		ProcessOwner: owner,
		SignedETHTx:  signedEthTx,
		Validators:   validators,
		FinalityVotes: int64(0),
	}
}

func (t *Tracker) AddVote(addr keys.Address, index int64) error {

	if len(t.Validators) <= int(index) {
		return errTrackerInvalidVote
	}

	_, voted := t.CheckIfVoted(addr)
	if !voted{
		if t.Validators[index].Equal(addr) {
			t.FinalityVotes = (t.FinalityVotes | (1 << index))
			return nil
		}
	}


	return errTrackerInvalidVote
}

func (t *Tracker) GetJobID(state TrackerState) string {
	return t.TrackerName.String() + storage.DB_PREFIX + strconv.Itoa(int(state))
}

func (t *Tracker) GetVotes() int {
	v := t.FinalityVotes
	/*
    fmt.Println(fmt.Sprintf("%b",t.FinalityVotes))
	cnt := 0
	for v >= 1 {
		if v%2 == 1 {

			cnt++
		}
		v = v >> 1
	}*/

	cnt := bits.OnesCount64(uint64(v))

	return cnt
}

func (t *Tracker) CheckIfVoted(node keys.Address) (index int64, voted bool) {
	index = int64(-1)
	v := false
	for i, addr := range t.Validators {
		if addr.Equal(node) {
			index = int64(i)
			break
		}
	}
	if index < int64(len(t.Validators)) && index >= 0 {
		var mybit int64 = 1 << index
		and := t.FinalityVotes & mybit
		v = and == mybit
	}

	return index, v
}

func (t *Tracker) Finalized() bool {
	l := len(t.Validators)
	num := int(float32(l)*votesThreshold) + 1
	v := t.FinalityVotes
	cnt := 0
	for v >= 1 {
		if v%2 == 1 {
			cnt++
		}
		v = v >> 1
	}
	return cnt >= num
}

func (t Tracker) NextStep() string {

	switch t.State {
	case New:
		fmt.Println("Chanjing state from NEW to Broadcasting")
		return BROADCASTING
	case BusyBroadcasting:
		fmt.Println("Changing state from BusyBroadcasting to Finalizing")
		return FINALIZING
	case BusyFinalizing:
		return FINALIZE
	case Finalized:
		return MINTING
	case Minted:
		return CLEANUP
	}
	return transition.NOOP
}
