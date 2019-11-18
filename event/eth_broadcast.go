package event

import (
	"context"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	"os"
	"time"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

type JobETHBroadcast struct {
	TrackerName         ethereum.TrackerName
	RetryCount          int8
	Finalized           bool
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


	ethconfig := config.DefaultEthConfig()
	logger := log.NewLoggerWithPrefix(os.Stdout,"JOB_ETHBROADCAST")
	cd,err := ethereum.NewEthereumChainDriver(ethconfig,logger,&ethCtx.ETHPrivKey)
	if err != nil {
		ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err)
		job.RetryCount += 1
		return
	}


	if !job.BroadcastSuccessful {
		rawTx := tracker.SignedETHTx
		tx := &types.Transaction{}
		err = rlp.DecodeBytes(rawTx, tx)
		if err != nil {
			ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", job.GetJobID(),err)
			return
		}
		// get chain driver here
		txhash,err := cd.BroadcastTx(tx)
		if err != nil {
			ethCtx.Logger.Error("Error in transaction broadcast : ", job.GetJobID(),err)
			return
		}
		job.BroadcastSuccessful = true
		job.BroadcastedHash = txhash
	} else {

		receipt,err := cd.CheckFinality(job.BroadcastedHash)
		if err != nil {
			ethCtx.Logger.Error("Error in Receiving TX receipt : ", job.GetJobID(),err)
			return
		}
		if receipt == nil {
			ethCtx.Logger.Info("Transaction not added to Ethereum Network yet ",job.GetJobID())
		}
		job.Finalized =true
	}
	if job.BroadcastSuccessful&&job.Finalized{
		// Create internal  TX
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
	return "We should make a Job ID"
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
