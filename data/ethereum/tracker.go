package ethereum

import (
	//"errors"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/event"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"strconv"
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
	voted := 0
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
	return ""
}

//TODO Go back to Busy broadcasting if there is a failure in Finalizing.
func Broadcasting(ctx interface{}) error {
	context := ctx.(TrackerCtx)
	tracker := context.tracker

	if tracker.State != New {
		err := errors.New("Cannot Broadcast from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	tracker.State = BusyBroadcasting

	//create broadcasting job
	job := event.NewETHBroadcast(tracker.TrackerName, tracker.State)
	err := context.jobStore.SaveJob(job)
	if err != nil {
		return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
	}

	return nil
}

func Finalizing(ctx interface{}) error {
	context := ctx.(TrackerCtx)
	tracker := context.tracker

	if tracker.State != BusyBroadcasting {
		err := errors.New("Cannot start Finalizing from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	//Check if current Node voted
	voted := tracker.CheckIfVoted(context.currNodeAddr)

	if !voted {
		//Check Broadcasting job
		job, typ := context.jobStore.GetJob(tracker.GetJobID(BusyBroadcasting))
		broadcastJob := event.MakeJob(job, typ)

		if broadcastJob.IsDone() {

			//Create job to check finality
			job := event.NewETHCheckFinality(tracker.TrackerName, BusyFinalizing)
			err := context.jobStore.SaveJob(job)
			if err != nil {
				return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
			}
		}
	}

	numVotes := tracker.GetVotes()

	if numVotes > 0 {
		tracker.State = BusyFinalizing
	}

	return nil
}

func Finalization(ctx interface{}) error {
	context := ctx.(TrackerCtx)
	tracker := context.tracker

	if tracker.State != BusyFinalizing {
		err := errors.New("cannot finalize from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	//Check if current Node voted
	voted := tracker.CheckIfVoted(context.currNodeAddr)

	if !voted {
		//Create job to check finality
		job := event.NewETHCheckFinality(tracker.TrackerName, tracker.State)
		err := context.jobStore.SaveJob(job)
		if err != nil {
			return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
		}
	}

	if tracker.Finalized() {
		tracker.State = Finalized
	}
	return nil
}

func Minting(ctx interface{}) error {
	context := ctx.(TrackerCtx)
	tracker := context.tracker

	if tracker.State != Finalized {
		err := errors.New("Cannot Mint from the current state")
		return errors.Wrap(err, string(tracker.State))
	}
	//todo: create a job to mint

	if tracker.TaskCompleted {
		tracker.State = Minted
	}
	return nil
}

func Cleanup(ctx interface{}) error {
	context := ctx.(TrackerCtx)
	tracker := context.tracker
	//todo: delete the tracker and jobs related

	//Delete Broadcasting Job
	job, typ := context.jobStore.GetJob(tracker.GetJobID(BusyBroadcasting))
	broadcastJob := event.MakeJob(job, typ)

	err := context.jobStore.DeleteJob(broadcastJob)
	if err != nil {
		return err
	}

	//Delete CheckFinality Job
	job, typ = context.jobStore.GetJob(tracker.GetJobID(BusyFinalizing))
	checkFinJob := event.MakeJob(job, typ)

	err = context.jobStore.DeleteJob(checkFinJob)
	if err != nil {
		return err
	}

	//Delete Tracker
	res, err := context.trackerStore.Delete(tracker.TrackerName)
	if err != nil || !res {
		return err
	}

	return nil
}
