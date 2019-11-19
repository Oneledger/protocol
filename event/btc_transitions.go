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

	data.Tracker.State = bitcoin.BusySigning

	return nil
}

func FreezeForBroadcast(inp interface{}) error {
	data, ok := inp.(bitcoin.BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	if data.Tracker.Multisig.IsValid() {
		data.Tracker.State = bitcoin.BusyBroadcasting
	}

	return nil
}

func ReportBroadcastSuccess(inp interface{}) error {
	data, ok := inp.(bitcoin.BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	data.Tracker.State = bitcoin.BusyFinalizing

	return nil
}
