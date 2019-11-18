package event

import (
	"context"
	"time"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type JobETHBroadcast struct {
	TrackerName         ethereum.TrackerName
	RetryCount          int8
	Done                bool
	BroadcastSuccessful bool
	BroadcastedHash     ethereum.TransactionHash
}

func (job JobETHBroadcast) DoMyJob(ctx interface{}) {

	// get tracker
	ethCtx, _ := ctx.(*JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.Get(job.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("err trying to deserialize tracker: ", job.TrackerName, err)
		job.RetryCount += 1
		return
	}
	client, err := ethclient.Dial(ethCtx.ETHConnection)
	if err != nil {
		ethCtx.Logger.Error("Unable to create Ethereum connection for the connection string :,", ethCtx.ETHConnection)
		return
	}
	if !job.BroadcastSuccessful {
		rawTx := tracker.SignedETHTx
		tx := &types.Transaction{}
		err = rlp.DecodeBytes(rawTx, tx)
		if err != nil {
			ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", job.TrackerName)
			return
		}
		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			ethCtx.Logger.Error("Error in tranascation broadcast : ", job.TrackerName)
			return
		}
		job.BroadcastSuccessful = true
		job.BroadcastedHash = tx.Hash()
	} else {

	}

}

func (job JobETHBroadcast) IsMyJobDone(ctx interface{}) bool {

	panic("implement me")
}

func (job JobETHBroadcast) IsSufficient(ctx interface{}) bool {
	panic("implement me")
}

func (job JobETHBroadcast) DoFinalize() {
	panic("implement me")
}

func (job JobETHBroadcast) GetType() string {
	panic("implement me")
}

func (job JobETHBroadcast) GetJobID() string {
	panic("implement me")
}

func (job JobETHBroadcast) IsDone() bool {
	panic("implement me")
}

func CheckTxForSuccess(client *ethclient.Client, tx *types.Transaction, maxWait time.Duration, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				result, err := client.TransactionReceipt(context.Background(), tx.Hash())
				if err == nil {
					if result.Status == types.ReceiptStatusSuccessful {
						ticker.Stop()
					}
				}
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	time.Sleep(maxWait)
	close(stop)
	return
}
