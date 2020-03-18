package event

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/utils/transition"
)

func init() {
	EthLockEngine = transition.NewEngine(
		[]transition.Status{
			transition.Status(ethereum.New),
			transition.Status(ethereum.BusyBroadcasting),
			//	transition.Status(ethereum.BroadcastSuccess),
			transition.Status(ethereum.BusyFinalizing),
			transition.Status(ethereum.Finalized),
			transition.Status(ethereum.Released),
			transition.Status(ethereum.Failed),
		})

	err := EthLockEngine.Register(transition.Transition{
		Name: ethereum.BROADCASTING,
		Fn:   Broadcasting,
		From: transition.Status(ethereum.New),
		To:   transition.Status(ethereum.BusyBroadcasting),
	})
	if err != nil {
		panic(err)
	}

	err = EthLockEngine.Register(transition.Transition{
		Name: ethereum.FINALIZING,
		Fn:   Finalizing,
		From: transition.Status(ethereum.BusyBroadcasting),
		To:   transition.Status(ethereum.BusyFinalizing),
	})
	if err != nil {
		panic(err)
	}

	err = EthLockEngine.Register(transition.Transition{
		Name: ethereum.FINALIZE,
		Fn:   Finalization,
		From: transition.Status(ethereum.BusyFinalizing),
		To:   transition.Status(ethereum.Finalized),
	})
	if err != nil {
		panic(err)
	}

	err = EthLockEngine.Register(transition.Transition{
		Name: ethereum.CLEANUP,
		Fn:   Cleanup,
		From: transition.Status(ethereum.Released),
		To:   0,
	})

	if err != nil {
		panic(err)
	}
	err = EthLockEngine.Register(transition.Transition{
		Name: ethereum.CLEANUPFAILED,
		Fn:   CleanupFailed,
		From: transition.Status(ethereum.Failed),
		To:   0,
	})
	if err != nil {
		panic(err)
	}

}

//TODO Go back to Busy broadcasting if there is a failure in Finalizing.
func Broadcasting(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	if tracker.State != ethereum.New {
		err := errors.New("Cannot Broadcast from the current state")
		return errors.Wrap(err, string((*tracker).State))
	}

	tracker.State = ethereum.BusyBroadcasting
	context.Tracker = tracker

	//create broadcasting
	if context.Validators.IsValidator() {

		job := NewETHBroadcast((*tracker).TrackerName, ethereum.BusyBroadcasting)
		err := context.JobStore.SaveJob(job)
		if err != nil {
			return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
		}
	}
	return nil
}

func Finalizing(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}
	tracker := context.Tracker

	if tracker.State != ethereum.BusyBroadcasting {
		err := errors.New("Cannot start Finalizing from the current state")
		return errors.Wrap(err, tracker.State.String())
	}
	y, n := tracker.GetVotes()

	if y+n > 0 {
		tracker.State = ethereum.BusyFinalizing
	}

	context.Tracker = tracker
	if context.Validators.IsValidator() {
		_, voted := tracker.CheckIfVoted(context.CurrNodeAddr)
		if voted {
			return nil
		}

		bjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyBroadcasting))
		if err != nil {
			return errors.Wrap(err, "failed to get job")
		}

		if !bjob.IsDone() || bjob.IsFailed() {
			return nil
		}

		fjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyFinalizing))
		if fjob != nil {
			return nil
		}
		job := NewETHCheckFinality(tracker.TrackerName, ethereum.BusyFinalizing)
		err = context.JobStore.SaveJob(job)
		if err != nil {
			return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
		}
	}

	return nil
}

func Finalization(ctx interface{}) error {

	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	if tracker.State != ethereum.BusyFinalizing {
		err := errors.New("cannot finalize from the current state")
		return errors.Wrap(err, string(tracker.State))
	}

	if tracker.Finalized() {
		tracker.State = ethereum.Finalized
		return nil
	}

	if context.Validators.IsValidator() {
		//Check if current Node voted
		_, voted := tracker.CheckIfVoted(context.CurrNodeAddr)

		if !voted {
			//Create job to check finality
			fmt.Println("Creating finality job from finalization")
			job := NewETHCheckFinality(tracker.TrackerName, tracker.State)

			err := context.JobStore.SaveJob(job)
			if err != nil {
				return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
			}
			//} else {
			//	job, err := context.JobStore.GetJob(tracker.GetJobID(tracker.State))
			//	if err != nil {
			//		return errors.Wrap(errors.LockNew("job serialization failed err: "), err.Error())
			//	}
			//	job.
		}
	}
	context.Tracker = tracker
	return nil
}

func Cleanup(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	//Delete Tracker
	context.Logger.Info("Setting tracker to success (ethLock):", tracker.State.String())
	err := context.TrackerStore.WithPrefixType(ethereum.PrefixPassed).Set(tracker.Clean())
	if err != nil {
		context.Logger.Error("error saving eth tracker", err)
		return err
	}
	context.Logger.Info("Deleting tracker (ethLock):", tracker.State.String())
	res, err := context.TrackerStore.WithPrefixType(ethereum.PrefixOngoing).Delete(tracker.TrackerName)
	if err != nil || !res {
		return err
	}

	//Delete Broadcasting Job
	if context.Validators.IsValidator() {
		bjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyBroadcasting))
		if err != nil {
			return errors.Wrap(err, "Failed to get Broadcasting Job")
		}

		err = context.JobStore.DeleteJob(bjob)
		if err != nil {
			return err
		}
		//Delete CheckFinality Job
		fjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyFinalizing))
		if err != nil {
			return errors.Wrap(err, "Failed to get Finalizing Job")
		}
		err = context.JobStore.DeleteJob(fjob)
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanupFailed(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	context.Logger.Info("Setting Tracker to failed (ethLock):", tracker.State.String())
	err := context.TrackerStore.WithPrefixType(ethereum.PrefixFailed).Set(tracker.Clean())
	if err != nil {
		context.Logger.Error("error saving eth tracker", err)
		return err
	}
	res, err := context.TrackerStore.WithPrefixType(ethereum.PrefixOngoing).Delete(tracker.TrackerName)
	if err != nil || !res {
		return err
	}

	//Delete Broadcasting Job It its there
	if context.Validators.IsValidator() {
		bjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyBroadcasting))
		if err == nil {
			err = context.JobStore.DeleteJob(bjob)
			if err != nil {
				return err
			}
		}

		//Delete CheckFinality Job If its there
		fjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyFinalizing))
		if err == nil {
			err = context.JobStore.DeleteJob(fjob)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
