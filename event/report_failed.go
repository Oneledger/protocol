package event

import (
	"errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/eth"
	"github.com/Oneledger/protocol/chains/ethereum"
)

func BroadcastReportFailedETHTx(ctx interface{},trackerName ethereum.TrackerName ,jobID string) (error){
	ethCtx, _ := ctx.(*JobsContext)
	trackerStore := ethCtx.EthereumTrackers
	tracker, err := trackerStore.Get(trackerName)
	index, _ := tracker.CheckIfVoted(ethCtx.ValidatorAddress)
	if index < 0 {
		return errors.New("Validator already Voted")
	}
	reportFailed := &eth.ReportFinality{
		TrackerName:      trackerName,
		Locker:           tracker.ProcessOwner,
		ValidatorAddress: ethCtx.ValidatorAddress,
		VoteIndex:        index,
		Refund:           false,
	}

	txData, err := reportFailed.Marshal()
	if err != nil {
		ethCtx.Logger.Error("Error while preparing mint txn ",jobID, err)
		return err
	}

	internalFailedTx := action.RawTx{
		Type: action.ETH_REPORT_FAILED,
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
		ethCtx.Logger.Error("error while broadcasting finality vote and mint txn ", jobID, err, rep.Log)
		return err
	}
	return nil
}
