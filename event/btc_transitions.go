/*

 */

package event

import (
	"github.com/Oneledger/protocol/data/bitcoin"
)

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
