/*

 */

package bitcoin

type BTCTransitionContext struct {
	Tracker *Tracker
}

func ReserveTracker(inp interface{}) error {
	data, ok := inp.(BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	data.Tracker.State = BusySigning

	return nil
}

func FreezeForBroadcast(inp interface{}) error {
	data, ok := inp.(BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	if data.Tracker.Multisig.IsValid() {
		data.Tracker.State = BusyBroadcasting
	}

	return nil
}

func ReportBroadcastSuccess(inp interface{}) error {
	data, ok := inp.(BTCTransitionContext)
	if !ok {
		panic("wrong transition data")
	}

	data.Tracker.State = BusyFinalizing

	return nil
}
