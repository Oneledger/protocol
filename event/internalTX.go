package event

import (
	"errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/chains/ethereum"
	ethereum2 "github.com/Oneledger/protocol/data/ethereum"
)

func BroadcastReportFinalityETHTx(ethCtx *JobsContext, trackerName ethereum.TrackerName, jobID string, success bool) error {

	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.WithPrefixType(ethereum2.PrefixOngoing).Get(trackerName)
	index, _ := tracker.CheckIfVoted(ethCtx.ValidatorAddress)
	if index < 0 {
		return errors.New("Validator already Voted")
	}
	reportFailed := &eth.ReportFinality{
		TrackerName:      trackerName,
		Locker:           tracker.ProcessOwner,
		ValidatorAddress: ethCtx.ValidatorAddress,
		VoteIndex:        index,
		Success:          success,
	}

	txData, err := reportFailed.Marshal()
	if err != nil {
		ethCtx.Logger.Error("Error while preparing mint txn ", jobID, err)
		return err
	}

	internalFailedTx := action.RawTx{
		Type: action.ETH_REPORT_FINALITY_MINT,
		Data: txData,
		Fee:  action.Fee{},
		Memo: jobID,
	}

	req := InternalBroadcastRequest{
		RawTx: internalFailedTx,
	}
	rep := BroadcastReply{}
	err = ethCtx.Service.InternalBroadcast(req, &rep)

	if err != nil || !rep.OK {
		ethCtx.Logger.Error("Error while broadcasting vote to Fail transaction ", jobID, err, rep.Log)
		return err
	}
	return nil
}
