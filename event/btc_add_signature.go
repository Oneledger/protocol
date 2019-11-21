/*

 */

package event

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/data/jobs"
)

type JobAddSignature struct {
	Type string

	TrackerName string

	JobID string

	Status jobs.Status
}

func NewAddSignatureJob(trackerName string) jobs.Job {

	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	return &JobAddSignature{
		Type:        JobTypeAddSignature,
		TrackerName: trackerName,
		JobID:       id,
		Status:      jobs.New,
	}
}

func (j *JobAddSignature) GetType() string {
	return JobTypeAddSignature
}

func (j *JobAddSignature) DoMyJob(ctxI interface{}) {
	ctx, _ := ctxI.(*JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		ctx.Logger.Error("error while getting tracker ", err, j.TrackerName)
		return
	}

	ctx.Logger.Info(hex.EncodeToString(tracker.ProcessUnsignedTx))

	lockTx := wire.NewMsgTx(wire.TxVersion)
	err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessUnsignedTx))
	if err != nil {
		ctx.Logger.Error("error while deserializing lock", err)
		return
	}

	lockScript, err := ctx.LockScripts.GetLockScript(tracker.CurrentLockScriptAddress)
	if err != nil {
		ctx.Logger.Error("erroring in reading lockscript", err)
		return
	}

	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), ctx.BTCPrivKey.Data)

	sig, err := txscript.RawTxInSignature(lockTx, 0, lockScript, txscript.SigHashAll, pk)
	if err != nil {
		ctx.Logger.Error(err, "RawTxInSignature")
		return
	}

	addSigData := btc.AddSignature{
		TrackerName:      j.TrackerName,
		ValidatorPubKey:  pk.PubKey().SerializeCompressed(),
		BTCSignature:     sig,
		ValidatorAddress: ctx.ValidatorAddress,
		Memo:             j.JobID,
	}

	txData, err := addSigData.Marshal()
	if err != nil {
		ctx.Logger.Error("error in marshalling txn", err)
		return
	}

	tx := action.RawTx{
		Type: action.BTC_ADD_SIGNATURE,
		Data: txData,
		Fee:  action.Fee{},
		Memo: j.JobID,
	}

	req := InternalBroadcastRequest{
		RawTx: tx,
	}
	rep := BroadcastReply{}

	err = ctx.Service.InternalBroadcast(req, &rep)
	if err != nil {
		ctx.Logger.Error("error in broadcasting internal addsignature", err)
		return
	}

}

func (j *JobAddSignature) GetJobID() string {
	return j.JobID
}

func (j JobAddSignature) IsDone() bool {
	return j.Status == jobs.Completed
}
