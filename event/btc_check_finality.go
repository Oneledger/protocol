/*

 */

package event

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/chains/bitcoin"
	bitcoin2 "github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/jobs"
)

const (
	TwoMinutes  = 60 * 2
	FiveMinutes = 60 * 5

	SixtyMinutes = 60 * 60
)

type JobBTCCheckFinality struct {
	Type string

	TrackerName string
	JobID       string
	CheckAfter  int64

	Status jobs.Status
}

func NewBTCCheckFinalityJob(trackerName, id string) jobs.Job {

	return &JobBTCCheckFinality{
		Type:        JobTypeBTCCheckFinality,
		TrackerName: trackerName,
		JobID:       id,
		CheckAfter:  time.Now().Unix() + SixtyMinutes,
		Status:      jobs.New,
	}
}

func (cf *JobBTCCheckFinality) DoMyJob(ctxI interface{}) {

	ctx, _ := ctxI.(*JobsContext)

	if time.Now().Unix() < cf.CheckAfter {
		return
	}

	tracker, err := ctx.Trackers.Get(cf.TrackerName)
	if err != nil {
		ctx.Logger.Error("err trying to deserialize tracker: ", cf.TrackerName, err)
		return
	}

	if tracker.State != bitcoin2.BusyFinalizing ||
		tracker.HasVotedFinality(ctx.ValidatorAddress) {

		cf.Status = jobs.Completed
		return
	}

	cd := bitcoin.NewChainDriver(ctx.Trackers.Config.BlockCypherToken)

	chain := bitcoin.GetBlockCypherChainType(ctx.Trackers.Config.BTCChainnet)

	ctx.Logger.Info("checking btc finality for ", tracker.ProcessTxId)
	ok, err := cd.CheckFinality(tracker.ProcessTxId, ctx.Trackers.Config.BlockCypherToken, chain)
	if err != nil {
		ctx.Logger.Error("error while checking finality", err, cf.TrackerName)
		return
	}

	if !ok {
		cf.CheckAfter = time.Now().Unix() + FiveMinutes
		ctx.Logger.Info("not finalized yet", cf.TrackerName)
		return
	}

	data := [4]byte{}
	_, err = io.ReadFull(rand.Reader, data[:])
	if err != nil {
		ctx.Logger.Error("error while reading random bytes for minting", err, cf.TrackerName)
		return
	}

	reportFinalityMint := btc.ReportFinalityMint{
		TrackerName:      cf.TrackerName,
		OwnerAddress:     tracker.ProcessOwner,
		ValidatorAddress: ctx.ValidatorAddress,
		RandomBytes:      data[:],
	}

	txData, err := reportFinalityMint.Marshal()
	if err != nil {
		ctx.Logger.Error("error while preparing mint txn ", err, cf.TrackerName)
		return
	}

	tx := action.RawTx{
		Type: action.BTC_REPORT_FINALITY_MINT,
		Data: txData,
		Fee:  action.Fee{},
		Memo: cf.JobID,
	}

	req := InternalBroadcastRequest{
		RawTx: tx,
	}
	rep := BroadcastReply{}

	err = ctx.Service.InternalBroadcast(req, &rep)
	if err != nil || !rep.OK {
		ctx.Logger.Error("error while broadcasting finality vote and mint txn ", err, cf.TrackerName)
		return
	}

	cf.Status = jobs.Completed
}

func (cf *JobBTCCheckFinality) GetType() string {
	return JobTypeBTCCheckFinality
}

func (cf *JobBTCCheckFinality) GetJobID() string {
	return cf.JobID
}

func (cf *JobBTCCheckFinality) IsDone() bool {
	return cf.Status == jobs.Completed
}
