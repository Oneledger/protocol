package event

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Oneledger/protocol/chains/ethereum"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
)

var _ jobs.Job = &JobETHSignRedeem{}

type JobETHSignRedeem struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
	TxHash      *ethereum.TransactionHash
}

func NewETHSignRedeem(name ethereum.TrackerName, state trackerlib.TrackerState) *JobETHSignRedeem {
	return &JobETHSignRedeem{
		TrackerName: name,
		JobID:       name.String() + storage.DB_PREFIX + strconv.Itoa(int(state)),
		RetryCount:  0,
		Status:      0,
	}
}

func (j *JobETHSignRedeem) DoMyJob(ctx interface{}) {
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

	cd, err := ethereum.NewChainDriver(ethconfig, ethCtx.Logger, trackerStore.GetOption())
	if err != nil {
		ethCtx.Logger.Error("err trying to get ChainDriver : ", j.GetJobID(), err)
		return
	}

	rawTx := tracker.SignedETHTx
	tx, err := cd.DecodeTransaction(rawTx)
	if err != nil {
		ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", j.GetJobID(), err)
		return
	}

	msg, err := cd.GetTransactionMessage(tx)
	if err != nil {
		ethCtx.Logger.Error("Error in decoding transaction as message : ", j.GetJobID(), err)
		return
	}

	addr := ethCtx.GetValidatorETHAddress()

	success, err := cd.HasValidatorSigned(addr, msg.From())
	if err != nil {
		ethCtx.Logger.Error("Error in verifying redeem :", j.GetJobID(), err)
	}
	if success {
		j.Status = jobs.Completed
		return
	}

	//check if tx already broadcasted, if yest, job.Status = jobs.Completed
	req, err := cd.ParseRedeem(rawTx)
	if err != nil {
		ethCtx.Logger.Error("Error in Parsing amount from rawTx ", j.GetJobID(), err)
		return
	}

	redeemAmount := req.Amount

	tx, err = cd.SignRedeem(addr, redeemAmount, msg.From())
	if err != nil {
		ethCtx.Logger.Error("Error in creating signing trasanction : ", j.GetJobID(), err)
		return
	}

	recipientAddr := common.HexToAddress(tracker.To.String())
	unsignedTx, err := cd.PrepareUnsignedETHRedeem(addr, redeemAmount, recipientAddr)
	if err != nil {
		ethCtx.Logger.Error("Error in preparing unsigned Ethereum Transaction", err)
		return
	}
	privkey := ethCtx.GetValidatorETHPrivKey()
	chainid, err := cd.ChainId()
	if err != nil {
		ethCtx.Logger.Error("Failed to get chain id ", err)
		return
	}

	signedTx, err := types.SignTx(unsignedTx, types.NewEIP155Signer(chainid), privkey)
	privkey = nil
	txHash, err := cd.BroadcastTx(signedTx)
	if err != nil {
		ethCtx.Logger.Error("Unable to broadcast transaction :", j.GetJobID(), err)
		return
	}

	ethCtx.Logger.Info("Redeem Transaction broadcasted to network : ", txHash)
	fmt.Println("Broadcast job completed for ", ethCtx.ValidatorAddress, "Job ID : ", j.GetJobID())

	j.TxHash = &txHash
	// j.Status = jobs.Completed
}

func (j *JobETHSignRedeem) IsDone() bool {
	return j.Status == jobs.Completed
}

func (j *JobETHSignRedeem) GetType() string {
	return JobTypeETHSignRedeem
}

func (j *JobETHSignRedeem) GetJobID() string {
	return j.JobID
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}
