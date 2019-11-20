package event

import (
	"os"

	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/utils/transition"
)

var (
	EthEngine transition.Engine
	BtcEngine transition.Engine
)

const (
	JobTypeAddSignature     = "addSignature"
	JobTypeBTCBroadcast     = "btcBroadcast"
	JobTypeBTCCheckFinality = "btcCheckFinality"

	MaxJobRetries = 10
)

func init() {
	EthEngine = transition.NewEngine(
		[]transition.Status{
			transition.Status(ethereum.New),
			transition.Status(ethereum.BusyBroadcasting),
			transition.Status(ethereum.BusyFinalizing),
			transition.Status(ethereum.Finalized),
			transition.Status(ethereum.Minted),
		})

	_ = EthEngine.Register(transition.Transition{
		Name: ethereum.BROADCASTING,
		Fn:   Broadcasting,
		From: transition.Status(ethereum.New),
		To:   transition.Status(ethereum.BusyBroadcasting),
	})

	_ = EthEngine.Register(transition.Transition{
		Name: ethereum.FINALIZING,
		Fn:   Finalizing,
		From: transition.Status(ethereum.BusyBroadcasting),
		To:   transition.Status(ethereum.BusyFinalizing),
	})

	_ = EthEngine.Register(transition.Transition{
		Name: ethereum.FINALIZE,
		Fn:   Finalization,
		From: transition.Status(ethereum.BusyFinalizing),
		To:   transition.Status(ethereum.Finalized),
	})

	_ = EthEngine.Register(transition.Transition{
		Name: ethereum.MINTING,
		Fn:   Minting,
		From: transition.Status(ethereum.Finalized),
		To:   transition.Status(ethereum.Minted),
	})
	_ = EthEngine.Register(transition.Transition{
		Name: ethereum.CLEANUP,
		Fn:   Cleanup,
		From: transition.Status(ethereum.Minted),
		To:   0,
	})

	BtcEngine = transition.NewEngine(
		[]transition.Status{bitcoin.Available, bitcoin.BusySigning, bitcoin.BusyBroadcasting, bitcoin.BusyFinalizing},
	)

	/*
		err := BtcEngine.Register(transition.Transition{
			Name: "makeAvailable",
			Fn:   MakeAvailable,
			From: bitcoin.BusyFinalizing,
			To:   bitcoin.Available,
		})
		if err != nil {
			os.Exit(1)
		}
	*/

	err := BtcEngine.Register(transition.Transition{
		Name: bitcoin.RESERVE,
		Fn:   ReserveTracker,
		From: bitcoin.Requested,
		To:   bitcoin.BusySigning,
	})
	if err != nil {
		os.Exit(1)
	}

	err = BtcEngine.Register(transition.Transition{
		Name: "freezeForBroadcast",
		Fn:   FreezeForBroadcast,
		From: bitcoin.BusySigning,
		To:   bitcoin.BusyBroadcasting,
	})
	if err != nil {
		os.Exit(1)
	}

	err = BtcEngine.Register(transition.Transition{
		Name: "reportBroadcastSuccess",
		Fn:   ReportBroadcastSuccess,
		From: bitcoin.BusySigning,
		To:   bitcoin.BusyFinalizing,
	})
	if err != nil {
		os.Exit(1)
	}

}
