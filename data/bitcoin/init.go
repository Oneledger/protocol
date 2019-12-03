/*

 */

package bitcoin

const (
	Available TrackerState = iota
	Requested
	BusySigning
	BusyScheduleBroadcasting
	BusyBroadcasting
	BusyScheduleFinalizing
	BusyFinalizing
	Finalized
)

const (
	RESERVE              = "reserveTracker"
	FREEZE_FOR_BROADCAST = "freezeForBroadcast"
	REPORT_BROADCAST     = "reportBroadcastSuccess"
	CLEANUP              = "cleanup"
)
