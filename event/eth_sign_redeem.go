package event

import (
	"fmt"
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
	//Get receipt first ,then status [ other way around might cause ambiguity ]
	//Confirm redeem tx has been sent by user
	//Verify Redeem should only be checked after Verify Receipt is confirmed , otherwise we might get stale verify redeem
	//Verify Receipt returns
	txReceipt, err := cd.VerifyReceipt(tx.Hash())
	if err != nil {
		ethCtx.Logger.Debug("Error in Getting TX receipt:", err)
		panic("Shutting down node ,unable to process redeem Transaction")
	}
	if txReceipt == ethereum.Failed {
		// TX included in uncle block , or TX reverted ( not enough redeem fee)
		ethCtx.Logger.Debug("Transaction receipt Failed  | Failing Tracker:", err)
		j.Status = jobs.Failed
		err := BroadcastReportFinalityETHTx(ctx.(*JobsContext), j.TrackerName, j.JobID, false)
		if err != nil {
			panic(fmt.Sprintf("Unable to broadcast failed TX for : %s ", j.JobID))
		}
		return
	}
	if txReceipt == ethereum.NotFound {
		ethCtx.Logger.Debug("Waiting for User Redeem TX to be mined", err)
		return
	}

	/*
		              Points to note

		              If expired fail tracker

		              Success ( Validator sign has been broadcasted and cofirmed on the ethereum blockchain)
					        If Success is true and validator has send signature implies , validator has signed , but his sign got reverted
							as this was the fourth sign .

			          Retrycount incremented only when validator broadacasts signature

				      Status is not ongoing , and txreceipt is true , redeem tx present in Ethereum ,but redeem status is not ongoing Fail tracker .

		              HasValidatorsigned returns success , means signature confirmed
	*/

	//Checking for confirmation of Vote
	success, err := cd.HasValidatorSigned(addr, msg.From())
	if err != nil {
		ethCtx.Logger.Error("Error connecting to HasValidatorSigned function in Smart Contract  :", j.GetJobID(), err)
		panic("Error connecting to HasValidatorSigned function in Smart Contract ")
	}
	//Signature confirmed
	if success {
		ethCtx.Logger.Debug("validator Sign Confirmed | validator Address (SIGNER):", ethCtx.GetValidatorETHAddress().Hex(), "| User Eth Address :", msg.From().Hex())
		j.Status = jobs.Completed
		return
	}
	//Log print debugger , Sign has been broadcast but not mined yet
	if j.RetryCount >= 0 {
		ethCtx.Logger.Debug("Waiting for Validator SignTX to be mined")
	}

	//Checking for Status of redeem request (From Ethereum smart contract)
	status := cd.VerifyRedeem(addr, msg.From())
	//Ethereum connectivity issue
	if status == ethereum.ErrorConnecting {
		ethCtx.Logger.Error("Error connecting to HasValidatorSigned function in Smart Contract  :", j.GetJobID(), err)
		panic(fmt.Sprintf("Error connecting to HasValidatorSigned function in Smart Contract %s, %s :", j.GetJobID(), err))
	}

	// Redeem request has expired
	if status == ethereum.Expired && txReceipt == ethereum.Found {
		ethCtx.Logger.Info("Failing from sign : Redeem Expired")
		j.Status = jobs.Failed
		err := BroadcastReportFinalityETHTx(ctx.(*JobsContext), j.TrackerName, j.JobID, false)
		if err != nil {
			panic(fmt.Sprintf("Unable to broadcast failed TX for : %s ", j.JobID))
		}
		return
	}
	// Status of redeem is Success but USER's Sign has not been confirmed (success is not true yet) = Redeem has been confirmed but this validators vote was reverted .
	if status == ethereum.Success && j.RetryCount >= 1 && txReceipt == ethereum.Found {
		ethCtx.Logger.Debug("Redeem TX successful , 67 % Votes have already been confirmed")
		j.Status = jobs.Completed
		return
	}
	// Status in not ongoing ,but User Redeem Request is verifiable on etherum = User is sending in an old redeem transaction
	if status != ethereum.Ongoing && txReceipt == ethereum.Found {
		ethCtx.Logger.Info("Redeem Request not created by user | Current Status : ", status.String())
		j.Status = jobs.Failed
		err := BroadcastReportFinalityETHTx(ctx.(*JobsContext), j.TrackerName, j.JobID, false)
		if err != nil {
			panic(fmt.Sprintf("Unable to broadcast failed TX for : %s ", j.JobID))
		}
		return
	}

	//Signing done only once after redeem receipt has been confirmed
	if j.RetryCount == 0 && txReceipt == ethereum.Found {

		redeemAddr := common.BytesToAddress(tracker.To)
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
			ethCtx.Logger.Error("Unable to broadcast transaction :", j.GetJobID(), err)
			return
		}
		j.RetryCount += 1
		ethCtx.Logger.Debug("Validator Sign Broadcasted | validator Address | (OL):", ethCtx.ValidatorAddress.Humanize(), "ETH ", ethCtx.GetValidatorETHAddress().Hex())
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
