package event

import (
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	ethereum2 "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

var _ jobs.Job = &JobETHBroadcast{}

type JobETHBroadcast struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
}

func NewETHBroadcast(name ethereum.TrackerName, state ethereum2.TrackerState) JobETHCheckFinality {
	return JobETHCheckFinality{
		TrackerName: name,
		JobID:       name.String() + storage.DB_PREFIX + strconv.Itoa(int(state)),
		RetryCount:  0,
		JobStatus:   0,
	}
}

func (job JobETHBroadcast) DoMyJob(ctx interface{}) {

	// get tracker
	job.RetryCount += 1
	if job.RetryCount > jobs.Max_Retry_Count {
		job.Status = jobs.Failed
	}
	if job.Status == jobs.New {
		job.Status = jobs.InProgress
	}
	ethCtx, _ := ctx.(*JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.Get(job.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("err trying to deserialize tracker: ", job.TrackerName, err)
		return
	}
	ethconfig := config.DefaultEthConfig()
	logger := log.NewLoggerWithPrefix(os.Stdout, "JOB_ETHBROADCAST")
	cd, err := ethereum.NewEthereumChainDriver(ethconfig, logger, &ethCtx.ETHPrivKey)
	if err != nil {
		ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err)

		return
	}
	rawTx := tracker.SignedETHTx
	tx := &types.Transaction{}
	err = rlp.DecodeBytes(rawTx, tx)
	if err != nil {
		ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", job.GetJobID(), err)
		return
	}
	_, err = cd.BroadcastTx(tx)
	if err != nil {
		ethCtx.Logger.Error("Error in transaction broadcast : ", job.GetJobID(), err)
		return
	}
	job.Status = jobs.Completed
}

func (job JobETHBroadcast) IsMyJobDone(ctx interface{}) bool {

	return job.Status == jobs.Completed
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
	return job.JobID
}

func (job JobETHBroadcast) IsDone() bool {
	if job.Status == jobs.Completed {
		return true
	}
	return false
}
