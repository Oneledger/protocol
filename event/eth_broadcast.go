package event

import (
	"fmt"
	"strconv"

	"github.com/Oneledger/protocol/chains/ethereum"
	ethereum2 "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
)

var _ jobs.Job = &JobETHBroadcast{}

type JobETHBroadcast struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
}

func NewETHBroadcast(name ethereum.TrackerName, state ethereum2.TrackerState) *JobETHBroadcast {
	fmt.Println("CREATING NEW JOB FOR NewETHBroadcast")
	return &JobETHBroadcast{
		TrackerName: name,
		JobID:       name.String() + storage.DB_PREFIX + strconv.Itoa(int(state)),
		RetryCount:  0,
		Status:      0,
	}
}

func (job *JobETHBroadcast) DoMyJob(ctx interface{}) {

	fmt.Println("Do job for broadcast")
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
	ethconfig := ethCtx.cfg.EthChainDriver

	//logger := log.NewLoggerWithPrefix(os.Stdout, "JOB_ETHBROADCAST")
	cd, err := ethereum.NewChainDriver(ethconfig, ethCtx.Logger, trackerStore.GetOption())
	if err != nil {
		ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err)

		return
	}
	rawTx := tracker.SignedETHTx
	tx,err := cd.DecodeTransaction(rawTx)
	if err != nil {
		ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", job.GetJobID(), err)
		return
	}
	//check if tx already broadcasted, if yest, job.Status = jobs.Completed
	_, err = cd.BroadcastTx(tx)
	if err != nil {
		ethCtx.Logger.Error("Error in transaction broadcast : ", job.GetJobID(), err)
		return
	}
	fmt.Println("Broadcast job completed ", job.GetJobID())
	job.Status = jobs.Completed
}

func (job *JobETHBroadcast) GetType() string {
	return JobTypeETHBroadcast
}

func (job *JobETHBroadcast) GetJobID() string {
	return job.JobID
}

func (job *JobETHBroadcast) IsDone() bool {
	return job.Status == jobs.Completed
}
