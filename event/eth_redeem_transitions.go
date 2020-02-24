package event

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/utils/transition"
)

func init() {

	EthRedeemEngine = transition.NewEngine(
		[]transition.Status{
			transition.Status(ethereum.New),
			transition.Status(ethereum.BusyBroadcasting),
			transition.Status(ethereum.BusyFinalizing),
			transition.Status(ethereum.Finalized),
			transition.Status(ethereum.Released),
			transition.Status(ethereum.Failed),
		})

	err := EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.SIGNING,
		Fn:   Signing,
		From: transition.Status(ethereum.New),
		To:   transition.Status(ethereum.BusyBroadcasting),
	})
	if err != nil {

		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.VERIFYREDEEM,
		Fn:   VerifyRedeem,
		From: transition.Status(ethereum.BusyBroadcasting),
		To:   transition.Status(ethereum.BusyFinalizing),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.REDEEMCONFIRM,
		Fn:   RedeemConfirmed,
		From: transition.Status(ethereum.BusyFinalizing),
		To:   transition.Status(ethereum.Finalized),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.CLEANUP,
		Fn:   redeemCleanup,
		From: transition.Status(ethereum.Released),
		To:   transition.Status(0),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.CLEANUPFAILED,
		Fn:   redeemCleanupFailed,
		From: transition.Status(ethereum.Failed),
		To:   0,
	})
	if err != nil {
		panic(err)
	}
}

func Signing(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	tracker := context.Tracker

	if tracker.State != ethereum.New {
		err := errors.New("Cannot Start Sign and Broadcast from Current State")
		return errors.Wrap(err, string((*tracker).State))
	}

	tracker.State = ethereum.BusyBroadcasting

	if context.Validators.IsValidator() {

		job := NewETHSignRedeem(tracker.TrackerName, tracker.State)
		err := context.JobStore.SaveJob(job)
		if err != nil {

			return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
		}

	}
	context.Tracker = tracker
	return nil
}

func VerifyRedeem(ctx interface{}) error {

	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}
	tracker := context.Tracker

	// create verify job for the first time from the state of broadcasting
	if tracker.State == ethereum.BusyBroadcasting && context.Validators.IsValidator() {
		job := NewETHVerifyRedeem(tracker.TrackerName, ethereum.BusyFinalizing)
		err := context.JobStore.SaveJob(job)
		if err != nil {
			return errors.Wrap(err, "Failed to save job")
		}
		tracker.State = ethereum.BusyFinalizing
	}

	context.Tracker = tracker
	return nil
}

func RedeemConfirmed(ctx interface{}) error {

	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}
	tracker := context.Tracker

	if tracker.State == ethereum.BusyFinalizing {
		if tracker.Finalized() {
			tracker.State = ethereum.Finalized
		}
	}
	context.Tracker = tracker
	return nil
}

func redeemCleanup(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}
	tracker := context.Tracker
	//delete the tracker related jobs
	if context.Validators.IsValidator() {
		for state := ethereum.BusyBroadcasting; state <= ethereum.Released; state++ {
			job, err := context.JobStore.GetJob(tracker.GetJobID(state))
			if err != nil {
				//fmt.Println(err, "failed to get job from state: ", state)
				continue
			}
			err = context.JobStore.DeleteJob(job)
			if err != nil {
				return errors.Wrap(err, "error deleting job from store")
			}
		}
	}
	//Delete Tracker

	fmt.Println("Deleting tracker (ethRedeem):",ethereum.GetTrackerStateString(tracker.State))
	res, err := context.TrackerStore.Delete(tracker.TrackerName)
	if err != nil || !res {
		return errors.Wrap(err, "error deleting tracker from store")
	}
	return nil
}

func redeemCleanupFailed(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}
	tracker := context.Tracker
	//delete the tracker related jobs
	if context.Validators.IsValidator() {
		for state := ethereum.BusyBroadcasting; state <= ethereum.Failed; state++ {
			job, err := context.JobStore.GetJob(tracker.GetJobID(state))
			if err != nil {
				//fmt.Println(err, "failed to get job from state: ", state)
				continue
			}
			if job != nil {
				err = context.JobStore.DeleteJob(job)
				if err != nil {
					return errors.Wrap(err, "error deleting job from store")
				}
			}
		}
	}
	//Delete Tracker
	fmt.Println("Deleting tracker (ethRedeem):",ethereum.GetTrackerStateString(tracker.State))

	res, err := context.TrackerStore.Delete(tracker.TrackerName)
	if err != nil || !res {
		return errors.Wrap(err, "error deleting tracker from store")
	}
	return nil
}
