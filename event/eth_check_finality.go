package event

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/chains/ethereum"
	ethtracker "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"

	"github.com/ethereum/go-ethereum/rlp"
)

var _ jobs.Job = &JobETHCheckFinality{}

type JobETHCheckFinality struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
}

func NewETHCheckFinality(name ethereum.TrackerName, state ethtracker.TrackerState) *JobETHCheckFinality {
	return &JobETHCheckFinality{
		TrackerName: name,
		JobID:       name.String() + storage.DB_PREFIX + strconv.Itoa(int(state)),
		RetryCount:  0,
		Status:      0,
	}
}

func (job *JobETHCheckFinality) DoMyJob(ctx interface{}) {

	// get tracker

	job.RetryCount += 1
	if job.RetryCount > jobs.Max_Retry_Count {
		job.Status = jobs.Failed
	}
	if job.Status == jobs.New {
		job.Status = jobs.InProgress
	}
	ethCtx, _ := ctx.(*JobsContext)
	fmt.Println("Starting check Finality JOB ", ethCtx.ValidatorAddress)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.Get(job.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("err trying to deserialize tracker: ", job.TrackerName, err)
		job.RetryCount += 1
		return
	}
	ethconfig := ethCtx.cfg.EthChainDriver
	cd, err := ethereum.NewEthereumChainDriver(ethconfig, ethCtx.Logger, trackerStore.GetOption())
	if err != nil {
		ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err)
		job.RetryCount += 1
		return
	}
	rawTx := tracker.SignedETHTx
	tx := &types.Transaction{}
	err = rlp.DecodeBytes(rawTx, tx)
	if err != nil {
		ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", job.GetJobID(), err)
		return
	}

	receipt, err := cd.CheckFinality(tx.Hash())
	if err != nil {
		ethCtx.Logger.Error("Error in Receiving TX receipt : ", job.GetJobID(), err)
		return
	}
	if receipt == nil {
		ethCtx.Logger.Info("Transaction not added to Ethereum Network yet ", job.GetJobID())
		return
	}
	index, _ := tracker.CheckIfVoted(ethCtx.ValidatorAddress)
	if index < 0 {
		return
	}
	reportFinalityMint := &eth.ReportFinalityMint{
		TrackerName:      job.TrackerName,
		Locker:           tracker.ProcessOwner,
		ValidatorAddress: ethCtx.ValidatorAddress,
		VoteIndex:        index,
	}

	fmt.Println("Creating Internal Transaction to add vote:", reportFinalityMint)
	txData, err := reportFinalityMint.Marshal()
	if err != nil {
		ethCtx.Logger.Error("Error while preparing mint txn ", job.GetJobID(), err)
		return
	}
	fmt.Println("after serialization", txData)
	internalMintTx := action.RawTx{
		Type: action.ETH_REPORT_FINALITY_MINT,
		Data: txData,
		Fee:  action.Fee{},
		Memo: job.GetJobID(),
	}

	req := InternalBroadcastRequest{
		RawTx: internalMintTx,
	}
	rep := BroadcastReply{}
	err = ethCtx.Service.InternalBroadcast(req, &rep)
	fmt.Println("Reply :", rep)
	if err != nil || !rep.OK {
		ethCtx.Logger.Error("error while broadcasting finality vote and mint txn ", job.GetJobID(), err, rep.Log)
		return
	}
	fmt.Println("Completed Check Finality JOB : ", ethCtx.ValidatorAddress)
	job.Status = jobs.Completed
}

func (job *JobETHCheckFinality) GetType() string {
	return JobTypeETHCheckfinalty
}

func (job *JobETHCheckFinality) GetJobID() string {
	return job.JobID
}

func (job *JobETHCheckFinality) IsDone() bool {
	return job.Status == jobs.Completed
}
