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
		})

	fmt.Println("EthRedeemEngine Register SIGNING")
	err := EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.SIGNING,
		Fn:   Signing,
		From: transition.Status(ethereum.New),
		To:   transition.Status(ethereum.BusyBroadcasting),
	})
	if err != nil {

		fmt.Println("EthRedeemEngine Register SIGNING", err)
		panic(err)
	}

	fmt.Println("EthRedeemEngine Register VERIFYREDEEM")
	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.VERIFYREDEEM,
		Fn:   VerifyRedeem,
		From: transition.Status(ethereum.BusyBroadcasting),
		To:   transition.Status(ethereum.BusyFinalizing),
	})
	if err != nil {
		fmt.Println("EthRedeemEngine Register VERIFYREDEEM", err)
		panic(err)
	}

	fmt.Println("EthRedeemEngine Register REDEEMCONFIRM")
	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.REDEEMCONFIRM,
		Fn:   RedeemConfirmed,
		From: transition.Status(ethereum.BusyFinalizing),
		To:   transition.Status(ethereum.Finalized),
	})
	if err != nil {
		fmt.Println("EthRedeemEngine Register REDEEMCONFIRM", err)
		panic(err)
	}

	//err = EthRedeemEngine.Register(transition.Transition{
	//	Name: ethereum.BURN,
	//	Fn:   Burn,
	//	From: transition.Status(ethereum.Finalized),
	//	To:   transition.Status(ethereum.Released),
	//})
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("EthRedeemEngine Register CLEANUP")
	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.CLEANUP,
		Fn:   redeemCleanup,
		From: transition.Status(ethereum.Released),
		To:   transition.Status(0),
	})
	if err != nil {
		fmt.Println("EthRedeemEngine Register CLEANUP", err)
		panic(err)
	}
}

func Signing(ctx interface{}) error {
	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}

	fmt.Println("Starting Redeem Signing")
	tracker := context.Tracker

	if tracker.State != ethereum.New {
		err := errors.New("Cannot Start Sign and Broadcast from Current State")
		return errors.Wrap(err, string((*tracker).State))
	}

	tracker.State = ethereum.BusyBroadcasting
	fmt.Println("STATE CHANGED TO :", tracker.State)

	if context.Validators.IsValidator() {
		fmt.Println("Created Signing JOB")
		job := NewETHSignRedeem(tracker.TrackerName, tracker.State)
		err := context.JobStore.SaveJob(job)
		if err != nil {
			fmt.Println("ERROR SAVING")
			return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
		}
		fmt.Println("saved job:", job)
	}
	context.Tracker = tracker
	return nil
}

func VerifyRedeem(ctx interface{}) error {
	fmt.Println("Starting Finalizing of Redeem Signing")

	context, ok := ctx.(*ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}
	tracker := context.Tracker

	bjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyBroadcasting))
	if err != nil {
		return errors.Wrap(err, "failed to get job")
	}
	fmt.Println("Checking Broadcasting job status", bjob.IsDone(), bjob.GetType())
	if !bjob.IsDone() {
		return errors.New("broadcast not done")
	}

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
				fmt.Println(err, "failed to get job from state: ", state)
				continue
			}
			err = context.JobStore.DeleteJob(job)
			if err != nil {
				return errors.Wrap(err, "error deleting job from store")
			}
		}
	}
	//Delete Tracker
	res, err := context.TrackerStore.Delete(tracker.TrackerName)
	if err != nil || !res {
		return errors.Wrap(err, "error deleting tracker from store")
	}
	return nil
}
