package event

import (
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	ethereum2 "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"

	"github.com/ethereum/go-ethereum/rlp"
)

var _ jobs.Job = &JobETHCheckFinality{}

type JobETHCheckFinality struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	JobStatus   jobs.Status
}

func NewETHCheckFinality(name ethereum.TrackerName, state ethereum2.TrackerState) JobETHCheckFinality {
	return JobETHCheckFinality{
		TrackerName: name,
		JobID:       name.String() + storage.DB_PREFIX + strconv.Itoa(int(state)),
		RetryCount:  0,
		JobStatus:   0,
	}
}

func (job JobETHCheckFinality) DoMyJob(ctx interface{}) {

	// get tracker
	job.RetryCount += 1
	if job.RetryCount > jobs.Max_Retry_Count {
		job.JobStatus = jobs.Failed
	}
	if job.JobStatus == jobs.New {
		job.JobStatus = jobs.InProgress
	}
	ethCtx, _ := ctx.(*JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.Get(job.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("err trying to deserialize tracker: ", job.TrackerName, err)
		job.RetryCount += 1
		return
	}
	ethconfig := config.DefaultEthConfig()
	logger := log.NewLoggerWithPrefix(os.Stdout, "JOB_ETHCHECKFINALITY")
	cd, err := ethereum.NewEthereumChainDriver(ethconfig, logger, &ethCtx.ETHPrivKey)
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

	reportFinalityMint := eth.ReportFinalityMint{
		TrackerName:      job.TrackerName,
		Locker:           tracker.ProcessOwner,
		ValidatorAddress: ethCtx.ValidatorAddress,
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
	if err != nil {
		ethCtx.Logger.Error("error while broadcasting finality vote and mint txn ", job.GetJobID(), err)
		return
	}
	job.JobStatus = jobs.Completed
}

func (job JobETHCheckFinality) IsMyJobDone(ctx interface{}) bool {

	panic("implement me")
}

func (job JobETHCheckFinality) IsSufficient(ctx interface{}) bool {
	panic("implement me")
}

func (job JobETHCheckFinality) DoFinalize() {
	panic("implement me")
}

func (job JobETHCheckFinality) GetType() string {
	panic("implement me")
}

func (job JobETHCheckFinality) GetJobID() string {
	return job.JobID
}

func (job JobETHCheckFinality) IsDone() bool {
	if job.JobStatus == jobs.Completed {
		return true
	}
	return false
}
