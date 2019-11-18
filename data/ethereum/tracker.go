package ethereum

import (
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
)

type TrackerState int

const (
	New TrackerState = iota
	BusyBroadcasting
	BusyFinalizing
	Finalized
	Minted

	votesThreshold float32 = 0.6667
)

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

//TODO Go back to Busy broadcasting if there is a failure in Finalizing.
func Broadcasting(ctx interface{}) error {
	context := ctx.(trackerCtx)
	tracker := context.tracker

	if tracker.State != New {
		err := errors.New("Cannot Broadcast from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	tracker.State = BusyBroadcasting
	return nil
}

func Finalizing(ctx interface{}) error {
	context := ctx.(trackerCtx)
	tracker := context.tracker

	if tracker.State != BusyBroadcasting {
		err := errors.New("Cannot start Finalizing from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	numVotes := tracker.GetVotes()

	if numVotes > 0 {
		tracker.State = BusyFinalizing
	}

	return nil
}

func Finalization(ctx interface{}) error {
	context := ctx.(trackerCtx)
	tracker := context.tracker

	if tracker.State != BusyFinalizing {
		err := errors.New("Cannot Finalize from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	if tracker.Finalized() {
		tracker.State = Finalized
	}
	return nil
}

func Minting(ctx interface{}) error {
	context := ctx.(trackerCtx)
	tracker := context.tracker

	if tracker.State != Finalized {
		err := errors.New("Cannot Mint from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	tracker.State = Minted
	return nil
}
