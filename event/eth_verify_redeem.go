package event

import (
	"fmt"
	"strconv"

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
	status := cd.VerifyRedeem(addr, msg.From())
	if status == ethereum.ErrorConnecting {
		ethCtx.Logger.Error("Error connecting to HasValidatorSigned function in Smart Contract  :", job.JobID, err)
		panic("Error connecting  to Smart Contract ") //TODO : Possible panic
	}
	if status == ethereum.Expired {
		job.Status = jobs.Failed
		err := BroadcastReportFinalityETHTx(ctx.(*JobsContext), job.TrackerName, job.JobID, false)
		if err != nil {
			panic(fmt.Sprintf("Unable to broadcast failed TX for : %s ", job.JobID))
		}

		return
	}
	if status == ethereum.Ongoing {
		ethCtx.Logger.Info("Waiting for RedeemTx to be Completed ,67 % Signing Votes")
	}
	// create internal check finality to report that the redeem is done on ethereum chain
	if status == ethereum.Success {
		index, _ := tracker.CheckIfVoted(ethCtx.ValidatorAddress)
		if index < 0 {
			return
		}
		err := BroadcastReportFinalityETHTx(ethCtx, job.TrackerName, job.JobID, true)
		if err != nil {
			panic(fmt.Sprintf("Unable to broadcast success TX for : %s ", job.JobID))
		}
		job.Status = jobs.Completed
		return
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
