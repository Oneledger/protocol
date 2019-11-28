package event

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
)

type JobETHSignRedeem struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
}

func (j JobETHSignRedeem) DoMyJob(ctx interface{}) {
	j.RetryCount += 1
	if j.RetryCount > jobs.Max_Retry_Count {
		j.Status = jobs.Failed
	}
	if j.Status == jobs.New {
		j.Status = jobs.InProgress
	}
	ethCtx, _ := ctx.(*JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.Get(j.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("err trying to deserialize tracker: ", j.TrackerName, err)
		return
	}
	ethconfig := ethCtx.cfg.EthChainDriver

	cd, err := ethereum.NewEthereumChainDriver(ethconfig, ethCtx.Logger, trackerStore.GetOption())
	if err != nil {
		ethCtx.Logger.Error("err trying to get ChainDriver : ", j.GetJobID(), err)
		return
	}
	rawTx := tracker.SignedETHTx
	tx := &types.Transaction{}
	err = rlp.DecodeBytes(rawTx, tx)
	if err != nil {
		ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", j.GetJobID(), err)
		return
	}
	//check if tx already broadcasted, if yest, job.Status = jobs.Completed
	_, err = cd.ValidatorSignRedeem()
	if err != nil {
		ethCtx.Logger.Error("Error in transaction broadcast : ", j.GetJobID(), err)
		return
	}
	fmt.Println("Broadcast job completed ", j.GetJobID())
	j.Status = jobs.Completed
}

func (j JobETHSignRedeem) IsDone() bool {
	return j.Status == jobs.Completed
}

func (j JobETHSignRedeem) GetType() string {
	return JobTypeETHSignRedeem
}

func (j JobETHSignRedeem) GetJobID() string {
	return j.JobID
}
