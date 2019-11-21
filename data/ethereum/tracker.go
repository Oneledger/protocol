package ethereum

import (
	//"errors"

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
	TaskCompleted bool
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
	}
}

func (t *Tracker) IsTaskCompleted() bool {
	return t.TaskCompleted
}

func (t *Tracker) CompleteTask() {
	if !t.TaskCompleted && t.Finalized() {
		t.TaskCompleted = true
	}
}

func (t *Tracker) AddVote(vote keys.Address, index int64) error {
	if len(t.Validators) <= int(index) {
		return errTrackerInvalidVote
	}
	if t.Validators[index].Equal(vote) {
		t.FinalityVotes += 1 << index
		return nil
	}
	return errTrackerInvalidVote
}

func (t *Tracker) GetJobID(state TrackerState) string {
	return t.TrackerName.String() + storage.DB_PREFIX + strconv.Itoa(int(state))
}

func (t *Tracker) GetVotes() int {
	v := t.FinalityVotes

	cnt := 0
	for v >= 1 {
		if v%2 == 1 {
			v = v >> 1
			cnt++
		}
	}

	return cnt
}

func (t *Tracker) CheckIfVoted(node keys.Address) bool {
	index := 0
	voted := int64(0)
	for i, addr := range t.Validators {
		if addr.Equal(node) {
			index = i
			break
		}
	}

	if index < len(t.Validators) {
		voted = (t.FinalityVotes >> index) % 2
	}

	return voted > 0
}

func (t *Tracker) Finalized() bool {
	l := len(t.Validators)
	num := int(float32(l)*votesThreshold) + 1
	v := t.FinalityVotes

	cnt := 0
	for v >= 1 {
		if v%2 == 1 {
			v = v >> 1
			cnt++
		}
	}
	return cnt >= num
}

func (t Tracker) NextStep() string {

	switch t.State {
	case New:
		return BROADCASTING
	case BusyBroadcasting:
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
