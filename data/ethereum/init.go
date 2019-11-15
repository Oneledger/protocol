package ethereum

import (
	"github.com/Oneledger/protocol/utils/transition"
	"github.com/pkg/errors"
)

var (
	ErrTrackerNotFound    = errors.New("tracker not found")
	errTrackerInvalidVote = errors.New("vote information is invalid")

	engine transition.Engine
)

type trackerCtx struct {
	tracker Tracker
}

func init() {
	engine = transition.NewEngine(
		[]transition.Status{
			transition.Status(New),
			transition.Status(BusyBroadcasting),
			transition.Status(BusyFinalizing),
			transition.Status(Finalized),
			transition.Status(Minted),
		})

	_ = engine.Register(transition.Transition{
		Name: "broadcasting",
		Fn:   Broadcast,
		From: transition.Status(New),
		To:   transition.Status(BusyBroadcasting),
	})

	_ = engine.Register(transition.Transition{
		Name: "finalize",
		Fn:   Finalize,
		From: transition.Status(BusyFinalizing),
		To:   transition.Status(Finalized),
	})
	_ = engine.Register(transition.Transition{
		Name: "mint",
		Fn:   Minting,
		From: transition.Status(Finalized),
		To:   transition.Status(Minted),
	})
}
