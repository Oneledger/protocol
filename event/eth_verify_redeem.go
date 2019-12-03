package event

import (
	"crypto/ecdsa"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"

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

func NewETHVerifyRedeem (name ethereum.TrackerName,state trackerlib.TrackerState) *JobETHVerifyRedeem{
	return &JobETHVerifyRedeem{
		TrackerName: name,
		JobID:       name.String()+storage.DB_PREFIX+strconv.Itoa(int(state)),
		RetryCount:  0,
		Status:      0,
	}
}

func (job *JobETHVerifyRedeem) DoMyJob(ctx interface{}) {
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
	if err!= nil {
		ethCtx.Logger.Error("Unable to get Tracker",job.JobID)
		return
	}
	cd,err := ethereum.NewChainDriver(ethCtx.cfg.EthChainDriver,ethCtx.Logger,trackerStore.GetOption())
	if err != nil {
		ethCtx.Logger.Error("Unable to get Chain Driver",job.JobID)
		return
	}
	tx,err := cd.DecodeTransaction(tracker.SignedETHTx)
	if err != nil {
		ethCtx.Logger.Error("Unable to decode transaction")
		return
	}
	msg,err := cd.GetTransactionMessage(tx)
	if err != nil {
		ethCtx.Logger.Error("Error in decoding transaction as message : ",job.GetJobID(),err)
		return
	}
	validatorPublicKey := ethCtx.ETHPrivKey.Public()
	publicKeyECDSA, ok := validatorPublicKey.(*ecdsa.PublicKey)
	if !ok {
		ethCtx.Logger.Error("error casting public key to ECDSA",job.GetJobID())

	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	success,err := cd.VerifyRedeem(fromAddress,msg.From())
	if err!= nil {
		ethCtx.Logger.Error("Error in verifying redeem :",job.GetJobID(),err)
	}
	if success {
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
	return job.GetJobID()
}
