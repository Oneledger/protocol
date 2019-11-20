/*

 */

package bitcoin

import "github.com/Oneledger/protocol/utils/transition"

const (
	Available transition.Status = iota
	Requested
	BusySigning
	BusyBroadcasting
	BusyFinalizing
)

const (
	RESERVE              = "reserveTracker"
	FREEZE_FOR_BROADCAST = "freezeForBroadcast"
	REPORT_BROADCAST     = "reportBroadcastSuccess"
)
