package event

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/ethereum"
	chaindriver "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/utils/transition"
)

func init() {

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
		Name: ethereum.VERIFYREDEEM,
		Fn:   VerifyRedeem,
		From: transition.Status(ethereum.BusyFinalizing),
		To:   transition.Status(ethereum.Finalized),
	})
	if err != nil {
		panic(err)
	}

	err = EthRedeemEngine.Register(transition.Transition{
		Name: ethereum.BURN,
		Fn:   Burn,
		From: transition.Status(ethereum.Finalized),
		To:   transition.Status(ethereum.Released),
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
	if tracker.State != ethereum.BusyBroadcasting {
		err := errors.New("Cannot start Finalizing Redeem from the current state")
		return errors.Wrap(err, string(tracker.State))
	}
	if context.Validators.IsValidator() {
		job :=NewETHVerifyRedeem(tracker.TrackerName,ethereum.BusyFinalizing)
		err := context.JobStore.SaveJob(job)
		if err != nil {
			return errors.Wrap(err,"Failed to save job")
		}
	}
	if tracker.State == ethereum.BusyFinalizing {
		signingJob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyFinalizing))
		if err != nil {
			return errors.Wrap(err, "Signing Job not found ")
		}
		if signingJob.IsDone(){
			tracker.State=ethereum.Finalized
		}
	}
	context.Tracker = tracker
	return nil
}



func Burn(ctx interface{}) error {
	return nil
}

func redeemCleanup(ctx interface{}) error {
	context, ok := ctx.(ethereum.TrackerCtx)
	if !ok {
		return errors.New("error casting tracker context")
	}
	tracker := context.Tracker
	//delete the tracker related jobs
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
	//Delete Tracker
	res, err := context.TrackerStore.Delete(tracker.TrackerName)
	if err != nil || !res {
		return errors.Wrap(err, "error deleting tracker from store")
	}
	return nil
}