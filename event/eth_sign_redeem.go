package event

import (
	"math/big"
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

	if j.Status == jobs.Completed {
		return
	}
	//if j.RetryCount > jobs.Max_Retry_Count {
	//	j.Status = jobs.Failed
	//	//BroadcastReportFinalityETHTx(ctx.(*JobsContext), j.TrackerName, j.JobID, false)
	//}
	if j.Status == jobs.New {
		j.Status = jobs.InProgress
	}

	ethCtx, _ := ctx.(*JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.WithPrefixType(trackerlib.PrefixOngoing).Get(j.TrackerName)
	if err != nil {
		ethCtx.Logger.Error("err trying to deserialize tracker: ", j.TrackerName, err)
		return
	}
	ethconfig := ethCtx.cfg.EthChainDriver
	ethoptions := trackerStore.GetOption()
	cd := new(ethereum.ETHChainDriver)
	redeemAmount := new(big.Int)
	if tracker.Type == trackerlib.ProcessTypeRedeem {
		cd, err = ethereum.NewChainDriver(ethconfig, ethCtx.Logger, ethoptions.ContractAddress, ethoptions.ContractABI, ethereum.ETH)
		if err != nil {
			ethCtx.Logger.Error("err trying to get ChainDriver : ", j.GetJobID(), err, tracker.Type)
			return
		}
		reqParams, err := cd.ParseRedeem(tracker.SignedETHTx, ethoptions.ContractABI)
		if err != nil {
			ethCtx.Logger.Error("Error in Parsing amount from rawTx (Ether Redeem)", j.GetJobID(), err)
			return
		}
		redeemAmount = reqParams.Amount

	} else if tracker.Type == trackerlib.ProcessTypeRedeemERC {
		cd, err = ethereum.NewChainDriver(ethconfig, ethCtx.Logger, ethoptions.ERCContractAddress, ethoptions.ERCContractABI, ethereum.ERC)
		if err != nil {
			ethCtx.Logger.Error("err trying to get ChainDriver : ", j.GetJobID(), err, tracker.Type)
			return
		}
		reqParams, err := cd.ParseERC20Redeem(tracker.SignedETHTx, ethoptions.ERCContractABI)
		if err != nil {
			ethCtx.Logger.Error("Error in Parsing amount from rawTx (ERC20 Redeem)", j.GetJobID(), err)
			return
		}
		redeemAmount = reqParams.Amount
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
	txReceipt, err := cd.VerifyReceipt(tx.Hash())
	if err != nil {
		ethCtx.Logger.Error("Trying to confirm RedeemTX sent by User Receipt :", err)
		return
	}
	//Failed to delete old version of chainstate err version does not exist version: -900
	// Get receipt first ,then status [ other way around might cause ambiguity ]
	// If expired fail tracker

	// If Success is true and validator has send signature implies , Validator has signed , but his sign got reverted
	//     as this was the fourth sign . Retrycount > 0 means that the validator did sign , status was 0 before ,implies
	//     this is not an old redeem

	// Status is not 0 , and txreceipt is true , redeem tx present in Ethereum ,but redeem status is not ongoing .Fail tracker

	// HasValidatorsigned returns success , means signature confirmed

	success, err := cd.HasValidatorSigned(addr, msg.From())
	if err != nil {
		ethCtx.Logger.Error("Error connecting to HasValidatorSigned function in Smart Contract  :", j.GetJobID(), err)
	}
	status, err := cd.VerifyRedeem(addr, msg.From())
	if err != nil {
		ethCtx.Logger.Error("Error in verifying redeem :", j.GetJobID(), err)
		return
	}
	if success {
		ethCtx.Logger.Info("Validator Sign Confirmed | Validator Address :", ethCtx.ValidatorAddress.Humanize(), "| User Eth Address :", msg.From().Hex())
		j.Status = jobs.Completed
		return
	}
	if j.RetryCount >= 0 && !success {
		ethCtx.Logger.Info("Waiting for Validator SignTX to be mined")
	}
	if err == ethereum.ErrRedeemExpired {
		ethCtx.Logger.Info("Failing from sign : Redeem Expired")
		j.Status = jobs.Failed
		BroadcastReportFinalityETHTx(ctx.(*JobsContext), j.TrackerName, j.JobID, false)
	}
	if status != 0 {
		if status == 1 && j.RetryCount >= 1 {
			ethCtx.Logger.Info("Redeem TX successful , 67 % Votes have already been confirmed")
			j.Status = jobs.Completed
			return
		}
		if txReceipt == true {
			ethCtx.Logger.Info("Redeem Request not created by user | Current Status : ", status.String())
			j.Status = jobs.Failed
			BroadcastReportFinalityETHTx(ctx.(*JobsContext), j.TrackerName, j.JobID, false)
			return
		}
	}

	//Signing ony done once Request sent only once
	if j.RetryCount == 0 {

		redeemAddr := common.HexToAddress(tracker.To.String())
		tx, err = cd.SignRedeem(addr, redeemAmount, redeemAddr)
		if err != nil {
			ethCtx.Logger.Error("Error in creating signing transaction : ", j.GetJobID(), err)
			return
		}

		privkey := ethCtx.GetValidatorETHPrivKey()
		chainid, err := cd.ChainId()
		if err != nil {
			ethCtx.Logger.Error("Failed to get chain id ", err)
			return
		}

		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainid), privkey)
		privkey = nil
		_, err = cd.BroadcastTx(signedTx)
		if err != nil {
			ethCtx.Logger.Error("Unable to broadcast transaction :", j.GetJobID(), err, " | RetryCount : ", j.RetryCount)
			return
		}
		j.RetryCount += 1
		ethCtx.Logger.Info("Validator Sign Broadcasted | Validator Address :", ethCtx.ValidatorAddress.Humanize())
	}
	//j.Status = jobs.Completed
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

func (j *JobETHSignRedeem) IsFailed() bool {
	return j.Status == jobs.Failed
}
