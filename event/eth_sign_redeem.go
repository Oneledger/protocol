package event

import (
	"crypto/ecdsa"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Oneledger/protocol/chains/ethereum"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/Oneledger/protocol/storage"
)

type JobETHSignRedeem struct {
	TrackerName ethereum.TrackerName
	JobID       string
	RetryCount  int
	Status      jobs.Status
}

func NewETHSignRedeem (name ethereum.TrackerName,state trackerlib.TrackerState) *JobETHSignRedeem{
	return &JobETHSignRedeem{
		TrackerName: name,
		JobID:       name.String()+storage.DB_PREFIX+strconv.Itoa(int(state)),
		RetryCount:  0,
		Status:      0,
	}
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

	cd, err := ethereum.NewChainDriver(ethconfig, ethCtx.Logger, trackerStore.GetOption())
	if err != nil {
		ethCtx.Logger.Error("err trying to get ChainDriver : ", j.GetJobID(), err)
		return
	}
	rawTx := tracker.SignedETHTx
    tx,err := cd.DecodeTransaction(rawTx)
	if err != nil {
		ethCtx.Logger.Error("Error Decoding Bytes from RaxTX :", j.GetJobID(), err)
		return
	}

	//check if tx already broadcasted, if yest, job.Status = jobs.Completed
    req,err := cd.ParseRedeem(rawTx)
    if err !=nil{
    	ethCtx.Logger.Error("Error in Parsing amount from rawTx ",j.GetJobID(),err)
		return
	}
    redeemAmount := req.Amount
    msg,err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()))
    if err != nil {
    	ethCtx.Logger.Error("Error in decoding trasnaction as message : ",j.GetJobID(),err)
		return
	}
	validatorPublicKey := ethCtx.ETHPrivKey.Public()
	publicKeyECDSA, ok := validatorPublicKey.(*ecdsa.PublicKey)
	if !ok {
		ethCtx.Logger.Error("error casting public key to ECDSA",j.GetJobID())

	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	tx, err = cd.SignRedeem(fromAddress,redeemAmount,msg.From())
	if err != nil {
		ethCtx.Logger.Error("Error in creating signing trasanction : ", j.GetJobID(), err)
		return
	}
	unsignedTx,err := cd.PrepareUnsignedETHRedeem(fromAddress,redeemAmount)
	if err != nil {
		ethCtx.Logger.Error("Error in preparing unsigned Ethereum Transaction")
	}
	signedTx,err := types.SignTx(unsignedTx,types.NewEIP155Signer(tx.ChainId()),ethCtx.ETHPrivKey)
	txHash,err := cd.BroadcastTx(signedTx)
	if err != nil {
		ethCtx.Logger.Error("Unable to broadcast transaction :",j.GetJobID(),err)
	}
	ethCtx.Logger.Info("Redeem Transaction broadcasted to network : ",txHash)
	fmt.Println("Broadcast job completed for ", ethCtx.ValidatorAddress,"Job ID : ", j.GetJobID())
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
