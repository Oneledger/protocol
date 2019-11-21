/*

 */

package event

import (
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/utils/transition"
)

func init() {
	BtcEngine = transition.NewEngine(
		[]transition.Status{bitcoin.Available, bitcoin.Requested, bitcoin.BusySigning, bitcoin.BusyBroadcasting, bitcoin.BusyFinalizing},
	)

	err := BtcEngine.Register(transition.Transition{
		Name: bitcoin.RESERVE,
		Fn:   ReserveTracker,
		From: bitcoin.Requested,
		To:   bitcoin.BusySigning,
	})
	if err != nil {
		panic(err)
	}

	err = BtcEngine.Register(transition.Transition{
		Name: "freezeForBroadcast",
		Fn:   FreezeForBroadcast,
		From: bitcoin.BusySigning,
		To:   bitcoin.BusyBroadcasting,
	})
	if err != nil {
		panic(err)
	}

	err = BtcEngine.Register(transition.Transition{
		Name: "reportBroadcastSuccess",
		Fn:   ReportBroadcastSuccess,
		From: bitcoin.BusySigning,
		To:   bitcoin.BusyFinalizing,
	})
	if err != nil {
		panic(err)
	}
}

func MakeAvailable(ctx interface{}) error {
	return nil
}

func ReserveTracker(inp interface{}) error {
	data, ok := inp.(bitcoin.BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	t := data.Tracker
	t.State = bitcoin.BusySigning

	job := NewAddSignatureJob(t.Name)
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

		job := NewBTCBroadcastJob(t.Name)
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

	job := NewBTCCheckFinalityJob(t.Name)
	err := data.JobStore.SaveJob(job)
	if err != nil {
		return err
	}

	return nil
}
