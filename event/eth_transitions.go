package event

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/ethereum"
)

//TODO Go back to Busy broadcasting if there is a failure in Finalizing.
func Broadcasting(ctx interface{}) error {
	context, ok := ctx.(ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	if tracker.State != ethereum.New {
		err := errors.New("Cannot Broadcast from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	tracker.State = ethereum.BusyBroadcasting

	//create broadcasting job
	job := NewETHBroadcast(tracker.TrackerName, tracker.State)
	err := context.JobStore.SaveJob(job)
	if err != nil {
		return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
	}

	return nil
}

func Finalizing(ctx interface{}) error {
	context, ok := ctx.(ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	if tracker.State != ethereum.BusyBroadcasting {
		err := errors.New("Cannot start Finalizing from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	//Check if current Node voted
	voted := tracker.CheckIfVoted(context.CurrNodeAddr)

	if !voted {
		//Check Broadcasting job
		job, typ := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyBroadcasting))
		broadcastJob := MakeJob(job, typ)

		if broadcastJob.IsDone() {

			//Create job to check finality
			job := NewETHCheckFinality(tracker.TrackerName, ethereum.BusyFinalizing)
			err := context.JobStore.SaveJob(job)
			if err != nil {
				return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
			}
		}
	}

	numVotes := tracker.GetVotes()

	if numVotes > 0 {
		tracker.State = ethereum.BusyFinalizing
	}

	return nil
}

func Finalization(ctx interface{}) error {
	context, ok := ctx.(ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	if tracker.State != ethereum.BusyFinalizing {
		err := errors.New("cannot finalize from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	//Check if current Node voted
	voted := tracker.CheckIfVoted(context.CurrNodeAddr)

	if !voted {
		//Create job to check finality
		job := NewETHCheckFinality(tracker.TrackerName, tracker.State)
		err := context.JobStore.SaveJob(job)
		if err != nil {
			return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
		}
	}

	if tracker.Finalized() {
		tracker.State = ethereum.Finalized
	}
	return nil
}

func Minting(ctx interface{}) error {
	context, ok := ctx.(ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	if tracker.State != ethereum.Finalized {
		err := errors.New("Cannot Mint from the current state")
		return errors.Wrap(err, string(tracker.State))
	}
	//todo: create a job to mint

	if tracker.TaskCompleted {
		tracker.State = ethereum.Minted
	}
	return nil
}

func Cleanup(ctx interface{}) error {
	context, ok := ctx.(ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker
	//todo: delete the tracker and jobs related

	//Delete Broadcasting Job
	job, typ := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyBroadcasting))
	broadcastJob := MakeJob(job, typ)

	err := context.JobStore.DeleteJob(broadcastJob)
	if err != nil {
		return err
	}

	//Delete CheckFinality Job
	job, typ = context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyFinalizing))
	checkFinJob := MakeJob(job, typ)

	err = context.JobStore.DeleteJob(checkFinJob)
	if err != nil {
		return err
	}

	//Delete Tracker
	res, err := context.TrackerStore.Delete(tracker.TrackerName)
	if err != nil || !res {
		return err
	}

	return nil
}
