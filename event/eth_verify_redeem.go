package event

import (
	"fmt"
	"strconv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/chains/ethereum"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
)

var _ jobs.Job = &JobETHVerifyRedeem{}

type JobETHVerifyRedeem struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
}

func NewETHVerifyRedeem(name ethereum.TrackerName, state trackerlib.TrackerState) *JobETHVerifyRedeem {
	return &JobETHVerifyRedeem{
		TrackerName: name,
		JobID:       name.String() + storage.DB_PREFIX + strconv.Itoa(int(state)),
		RetryCount:  0,
		Status:      0,
	}
}

func (job *JobETHVerifyRedeem) DoMyJob(ctx interface{}) {
	job.RetryCount += 1
	if job.Status == jobs.New {
		job.Status = jobs.InProgress
	}
	ethCtx, _ := ctx.(*JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.WithPrefixType(trackerlib.PrefixOngoing).Get(job.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("Unable to get Tracker", job.JobID)
		return
	}
	ethconfig := ethCtx.cfg.EthChainDriver
	ethoptions := trackerStore.GetOption()
	cd := new(ethereum.ETHChainDriver)
	if tracker.Type == trackerlib.ProcessTypeRedeem {
		cd, err = ethereum.NewChainDriver(ethconfig, ethCtx.Logger, ethoptions.ContractAddress, ethoptions.ContractABI, ethereum.ETH)
		if err != nil {
			ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err, tracker.Type)
			return
		}
	} else if tracker.Type == trackerlib.ProcessTypeRedeemERC {
		cd, err = ethereum.NewChainDriver(ethconfig, ethCtx.Logger, ethoptions.ERCContractAddress, ethoptions.ERCContractABI, ethereum.ERC)
		if err != nil {
			ethCtx.Logger.Error("err trying to get ChainDriver : ", job.GetJobID(), err, tracker.Type)
			return
		}
	}

	tx, err := cd.DecodeTransaction(tracker.SignedETHTx)
	if err != nil {
		ethCtx.Logger.Error("Unable to decode transaction")
		return
	}

	msg, err := cd.GetTransactionMessage(tx)
	if err != nil {
		ethCtx.Logger.Error("Error in decoding transaction as message : ", job.GetJobID(), err)
		return
	}

	addr := ethCtx.GetValidatorETHAddress()
	status, err := cd.VerifyRedeem(addr, msg.From())
	if err != nil {
		ethCtx.Logger.Error("Error in verifying redeem :", job.GetJobID(), err)
	}
	if err == ethereum.ErrRedeemExpired {
		fmt.Println("Failing from verify : Redeem Expired")
		job.Status = jobs.Failed
		BroadcastReportFinalityETHTx(ctx.(*JobsContext), job.TrackerName, job.JobID, false)
	}
	if status == 0 {
		ethCtx.Logger.Info("Waiting for RedeemTx to be Completed ,67 % Signing Votes")
	}
	// create internal check finality to report that the redeem is done on ethereum chain
	if status == 1 {
		index, _ := tracker.CheckIfVoted(ethCtx.ValidatorAddress)
		if index < 0 {
			return
		}
		//TODO : Replace with BroadcastReportFinalityETHTx(ethCtx,job.TrackerName,job.JobID,true)
		cf := &eth.ReportFinality{
			TrackerName:      tracker.TrackerName,
			Locker:           tracker.ProcessOwner,
			ValidatorAddress: ethCtx.ValidatorAddress,
			VoteIndex:        index,
			Success:          true,
		}

		txData, err := cf.Marshal()
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
		//TODO END
		job.Status = jobs.Completed
	}
}

func (job *JobETHVerifyRedeem) IsDone() bool {
	return job.Status == jobs.Completed
}

func (job *JobETHVerifyRedeem) GetType() string {
	return JobTypeETHVerifyRedeem
}

func (job *JobETHVerifyRedeem) GetJobID() string {
	return job.JobID
}

func (job *JobETHVerifyRedeem) IsFailed() bool {
	return job.Status == jobs.Failed
}
