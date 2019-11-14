package eth

import (
	"context"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type JobETHBroadcast struct {

	TrackerName ethereum.TrackerName


}

func (j JobETHBroadcast) DoMyJob(ctx interface{}) {

	// get tracker
	ethCtx, _ := ctx.(*action.JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker,_ := trackerStore.Get(j.TrackerName)
	client, _ := ethclient.Dial(ethCtx.ETHConnection)
	rawTx := tracker.SignedETHTx
	tx := &types.Transaction{}
	_ := rlp.DecodeBytes(rawTx,tx)
	txHash := client.SendTransaction(context.Background(),tx)
	// Put tx hash into tracker   ?

}

func (j JobETHBroadcast) IsMyJobDone(ctx interface{}) bool {}
	panic("implement me")


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
