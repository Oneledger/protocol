package eth

import (
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type JobETHBroadcast struct {
	TrackerName common.Hash
}

func (j JobETHBroadcast) DoMyJob(ctx interface{}) {

	// get tracker
	tracker, err := ctx.ETHTrackerStore.Get(j.TrackerName)

	client := ethclient.Dial("")

	//
	tx := tracker.SignedTx

	client.SendTransaction()
	ethereum.TransactOpts{}
}

func (j JobETHBroadcast) IsMyJobDone(ctx interface{}) bool {
	panic("implement me")
}

func (j JobETHBroadcast) IsSufficient(ctx interface{}) bool {
	panic("implement me")
}

func (j JobETHBroadcast) DoFinalize() {
	panic("implement me")
}

func (j JobETHBroadcast) GetType() string {
	panic("implement me")
}

func (j JobETHBroadcast) GetJobID() string {
	panic("implement me")
}

func (j JobETHBroadcast) IsDone() bool {
	panic("implement me")
}
