package ethereum

import (
	"github.com/Oneledger/protocol/utils/transition"
	"github.com/pkg/errors"
)

var (
	ErrTrackerNotFound    = errors.New("tracker not found")
	errTrackerInvalidVote = errors.New("vote information is invalid")

	Engine transition.Engine
)

const (
	New TrackerState = iota
	BusyBroadcasting
	BusyFinalizing
	Finalized
	Minted

	votesThreshold float32 = 0.6667

	BROADCASTING string = "broadcasting"
	FINALIZING   string = "finalizing"
	FINALIZE     string = "finalize"
	MINTING      string = "minting"
	CLEANUP      string = "cleanup"
)

func init() {
	Engine = transition.NewEngine(
		[]transition.Status{
			transition.Status(New),
			transition.Status(BusyBroadcasting),
			transition.Status(BusyFinalizing),
			transition.Status(Finalized),
			transition.Status(Minted),
		})

	_ = Engine.Register(transition.Transition{
		Name: BROADCASTING,
		Fn:   Broadcasting,
		From: transition.Status(New),
		To:   transition.Status(BusyBroadcasting),
	})

	_ = Engine.Register(transition.Transition{
		Name: FINALIZING,
		Fn:   Finalizing,
		From: transition.Status(BusyBroadcasting),
		To:   transition.Status(BusyFinalizing),
	})

	_ = Engine.Register(transition.Transition{
		Name: FINALIZE,
		Fn:   Finalization,
		From: transition.Status(BusyFinalizing),
		To:   transition.Status(Finalized),
	})

	_ = Engine.Register(transition.Transition{
		Name: MINTING,
		Fn:   Minting,
		From: transition.Status(Finalized),
		To:   transition.Status(Minted),
	})
	_ = Engine.Register(transition.Transition{
		Name: CLEANUP,
		Fn:   Cleanup,
		From: transition.Status(Minted),
		To:   0,
	})
}
