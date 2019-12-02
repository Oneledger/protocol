package event

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/ethereum"
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
		Name: ethereum.FINALIZESIGNING,
		Fn:   FinalizeSigning,
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
		Fn:   redeemcleanup,
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

func FinalizeSigning(ctx interface{}) error {
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
		_, voted := tracker.CheckIfVoted(context.CurrNodeAddr)
		if !voted {
			signingjob, err := context.JobStore.GetJob(tracker.GetJobID(ethereum.BusyBroadcasting))
			if err != nil {
				return errors.Wrap(err, "failed to get job")
			}
			if signingjob.IsDone() {
				fmt.Println("Signing job is done but Vote not received ")
				job := NewETHCheckFinality(tracker.TrackerName, ethereum.BusyFinalizing)
				err := context.JobStore.SaveJob(job)
				if err != nil {
					return errors.Wrap(errors.New("job serialization failed err: "), err.Error())
				}
			}
		}
	}
	yesVotes, _ := tracker.GetVotes()
	if yesVotes > 0 {
		fmt.Println("Setting Status to BusyFinalizing")
		tracker.State = ethereum.BusyFinalizing
	}

	context.Tracker = tracker
	return nil
}

func VerifyRedeem(ctx interface{}) error {
	return nil
}

func Burn(ctx interface{}) error {
	return nil
}

func redeemcleanup(ctx interface{}) error {
	return nil
}
