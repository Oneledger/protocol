/*

 */

package event

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/utils/transition"
)

func init() {
	BtcEngine = transition.NewEngine(
		[]transition.Status{
			transition.Status(bitcoin.Available),
			transition.Status(bitcoin.Requested),
			transition.Status(bitcoin.BusySigning),
			transition.Status(bitcoin.BusyScheduleBroadcasting),
			transition.Status(bitcoin.BusyBroadcasting),
			transition.Status(bitcoin.BusyScheduleFinalizing),
			transition.Status(bitcoin.BusyFinalizing),
			transition.Status(bitcoin.Finalized),
		},
	)

	err := BtcEngine.Register(transition.Transition{
		Name: bitcoin.CLEANUP,
		Fn:   MakeAvailable,
		From: transition.Status(bitcoin.Finalized),
		To:   transition.Status(bitcoin.Available),
	})
	if err != nil {
		panic(err)
	}

	err = BtcEngine.Register(transition.Transition{
		Name: bitcoin.RESERVE,
		Fn:   ReserveTracker,
		From: transition.Status(bitcoin.Requested),
		To:   transition.Status(bitcoin.BusySigning),
	})
	if err != nil {
		panic(err)
	}

	err = BtcEngine.Register(transition.Transition{
		Name: bitcoin.FREEZE_FOR_BROADCAST,
		Fn:   FreezeForBroadcast,
		From: transition.Status(bitcoin.BusyScheduleBroadcasting),
		To:   transition.Status(bitcoin.BusyBroadcasting),
	})
	if err != nil {
		panic(err)
	}

	err = BtcEngine.Register(transition.Transition{
		Name: bitcoin.REPORT_BROADCAST,
		Fn:   ReportBroadcastSuccess,
		From: transition.Status(bitcoin.BusyScheduleFinalizing),
		To:   transition.Status(bitcoin.BusyFinalizing),
	})
	if err != nil {
		panic(err)
	}
}

func MakeAvailable(input interface{}) error {
	data, ok := input.(bitcoin.BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	t := data.Tracker

	states := []bitcoin.TrackerState{bitcoin.BusySigning, bitcoin.BusyBroadcasting, bitcoin.BusyFinalizing}
	for i := range states {

		//Delete Add Signature Job
		fjob, err := data.JobStore.GetJob(t.GetJobID(states[i]))
		if err != nil {
			continue
		}
		err = data.JobStore.DeleteJob(fjob)
		if err != nil {
			return errors.Wrap(err, "failed to delete job")
		}
	}

	t.State = bitcoin.Available
	return nil
}

func ReserveTracker(inp interface{}) error {
	data, ok := inp.(bitcoin.BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	t := data.Tracker
	t.State = bitcoin.BusySigning

	job := NewAddSignatureJob(t.Name, t.GetJobID(t.State))
	err := data.JobStore.SaveJob(job)
	if err != nil {
		return err
	}

	return nil
}

func FreezeForBroadcast(inp interface{}) error {
	data, ok := inp.(bitcoin.BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	t := data.Tracker
	if t.Multisig.IsValid() {

		data.Tracker.State = bitcoin.BusyBroadcasting

		job := NewBTCBroadcastJob(t.Name, t.GetJobID(t.State))
		err := data.JobStore.SaveJob(job)
		if err != nil {
			return err
		}
	} else if t.Multisig.IsCancel() {

	}

	return nil
}

func ReportBroadcastSuccess(inp interface{}) error {
	data, ok := inp.(bitcoin.BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	t := data.Tracker
	t.State = bitcoin.BusyFinalizing

	job := NewBTCCheckFinalityJob(t.Name, t.GetJobID(t.State))
	err := data.JobStore.SaveJob(job)
	if err != nil {
		return err
	}

	return nil
}
