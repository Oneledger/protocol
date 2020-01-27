package event

import (
	"strconv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/chains/ethereum"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
)

var _ jobs.Job = &JobETHCheckFinality{}

type JobETHCheckFinality struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
}

func NewETHCheckFinality(name ethereum.TrackerName, state trackerlib.TrackerState) *JobETHCheckFinality {
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

	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.Get(job.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("err trying to deserialize tracker: ", job.TrackerName, err)
		job.RetryCount += 1
		return
	}

	ethconfig := ethCtx.cfg.EthChainDriver
	ethoptions := trackerStore.GetOption()
	cd := new(ethereum.ETHChainDriver)
	if tracker.Type == trackerlib.ProcessTypeLock {
		cd, err = ethereum.NewChainDriver(ethconfig, ethCtx.Logger, ethoptions.ContractAddress, ethoptions.ContractABI, ethereum.ETH)
		if err != nil {
			ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err, tracker.Type)
			return
		}
	} else if tracker.Type == trackerlib.ProcessTypeLockERC {
		cd, err = ethereum.NewChainDriver(ethconfig, ethCtx.Logger, ethoptions.ERCContractAddress, ethoptions.ERCContractABI, ethereum.ERC)
		if err != nil {
			ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err, tracker.Type)
			return
		}
	}

	rawTx := tracker.SignedETHTx
	tx, err := cd.DecodeTransaction(rawTx)
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
	reportFinalityMint := &eth.ReportFinality{
		TrackerName:      job.TrackerName,
		Locker:           tracker.ProcessOwner,
		ValidatorAddress: ethCtx.ValidatorAddress,
		VoteIndex:        index,
		Refund:           false,
	}

	txData, err := reportFinalityMint.Marshal()
	if err != nil {
		ethCtx.Logger.Error("Error while preparing mint txn ", job.GetJobID(), err)
		return
	}

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

	if err != nil || !rep.OK {
		ethCtx.Logger.Error("error while broadcasting finality vote and mint txn ", job.GetJobID(), err, rep.Log)
		return
	}

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
