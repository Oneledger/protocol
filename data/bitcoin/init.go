/*git
 */

package bitcoin

import (
	"os"

	"github.com/Oneledger/protocol/utils/transition"
)

var (
	Engine transition.Engine
)

func init() {

	Engine = transition.NewEngine(
		[]transition.Status{Available, BusySigning, BusyBroadcasting, BusyFinalizing},
	)

	err := engine.Register(transition.Transition{
		Name: "makeAvailable",
		Fn:   MakeAvailable,
		From: BusyFinalizing,
		To:   Available,
	})
	if err != nil {
		os.Exit(1)
	}

	err = engine.Register(transition.Transition{
		Name: "reserveTracker",
		Fn:   ReserveTracker,
		From: Available,
		To:   BusySigning,
	})
	if err != nil {
		os.Exit(1)
	}

	err = engine.Register(transition.Transition{
		Name: "freezeForBroadcast",
		Fn:   FreezeForBroadcast,
		From: BusySigning,
		To:   BusyBroadcasting,
	})
	if err != nil {
		os.Exit(1)
	}

	err = engine.Register(transition.Transition{
		Name: "reportBroadcastSuccess",
		Fn:   ReportBroadcastSuccess,
		From: BusySigning,
		To:   BusyFinalizing,
	})
	if err != nil {
		os.Exit(1)
	}
}
