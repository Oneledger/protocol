/*

 */

package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcec"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type JobAddSignature struct {
	Type string

	TrackerName string

	JobID string

	Done bool

	RetryCount int8
}

func NewAddSignatureJob(trackerName string) jobs.Job {

	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	return &JobAddSignature{
		JobTypeAddSignature,
		trackerName,
		id,
		false,
		0,
	}
}

func (j *JobAddSignature) GetType() string {
	return JobTypeAddSignature
}

type doJobData struct {
}

func (j *JobAddSignature) DoMyJob(ctxI interface{}) {
	ctx, _ := ctxI.(*action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		ctx.Logger.Error("error while getting tracker ", err, j.TrackerName)
		j.RetryCount += 1
		return
	}

	ctx.Logger.Info(hex.EncodeToString(tracker.ProcessUnsignedTx))

	lockTx := wire.NewMsgTx(wire.TxVersion)
	err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessUnsignedTx))
	if err != nil {
		j.RetryCount += 1
		return
	}

	lockScript, err := ctx.LockScripts.GetLockScript(tracker.CurrentLockScriptAddress)
	if err != nil {
		j.RetryCount += 1
		return
	}

	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), ctx.BTCPrivKey.Data)

	sig, err := txscript.RawTxInSignature(lockTx, 0, lockScript, txscript.SigHashAll, pk)
	if err != nil {
		fmt.Println(err, "RawTxInSignature")
		j.RetryCount += 1
		return
	}

	addSigData := AddSignature{
		TrackerName:      j.TrackerName,
		ValidatorPubKey:  pk.PubKey().SerializeCompressed(),
		BTCSignature:     sig,
		ValidatorAddress: ctx.ValidatorAddress,
		Memo:             j.JobID,
	}

	txData, err := addSigData.Marshal()
	if err != nil {
		// retry later
		j.RetryCount += 1
		return
	}

	tx := action.RawTx{
		Type: action.BTC_ADD_SIGNATURE,
		Data: txData,
		Fee:  action.Fee{},
		Memo: j.JobID,
	}

	req := action.InternalBroadcastRequest{
		RawTx: tx,
	}
	rep := action.BroadcastReply{}

	err = ctx.Service.InternalBroadcast(req, &rep)
	if err != nil {
		j.RetryCount += 1
		return
	}

}

func (j *JobAddSignature) IsMyJobDone(ctxI interface{}) bool {
	ctx, _ := ctxI.(*action.JobsContext)

	if j.RetryCount > MaxJobRetries {
		return true
	}

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	return tracker.Multisig.HasAddressSigned(ctx.ValidatorAddress)
}

func (j *JobAddSignature) IsSufficient(ctxI interface{}) bool {

	ctx, _ := ctxI.(*action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	return tracker.Multisig.IsValid()
}

func (j *JobAddSignature) DoFinalize() {
	// TODO:

	j.Done = true
}

func (j *JobAddSignature) GetJobID() string {
	return j.JobID
}

func (j *JobAddSignature) IsDone() bool {
	return j.Done
}
