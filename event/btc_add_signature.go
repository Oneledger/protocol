/*

 */

package event

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

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

func NewAddSignatureJob(trackerName, id string) jobs.Job {

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

	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), ctx.BTCPrivKey.Data)

	opt := ctx.Trackers.GetOptions()
	addressPubKey, err := btcutil.NewAddressPubKey(pk.PubKey().SerializeCompressed(), opt.BTCParams)
	if err != nil {
		ctx.Logger.Error("error while generating btc address", err)
		return
	}

	if tracker.Multisig.HasAddressSigned(addressPubKey.ScriptAddress()) {

		j.Status = jobs.Completed
		return
	}

	lockTx := wire.NewMsgTx(wire.TxVersion)
	err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessUnsignedTx))
	if err != nil {
		ctx.Logger.Error("error while deserializing lock", err)
		return
	}

	lockScript, err := ctx.LockScripts.GetLockScript(tracker.CurrentLockScriptAddress)
	if err != nil {
		ctx.Logger.Error("error in reading lockscript", err)
		return
	}

	if len(lockScript) == 0 && tracker.CurrentTxId != nil {
		ctx.Logger.Error("error in reading lockscript", err)
		return
	}

	sig, err := txscript.RawTxInSignature(lockTx, 0, lockScript, txscript.SigHashAll, pk)
	if err != nil {
		ctx.Logger.Error(err, "RawTxInSignature")
		ctx.Logger.Error(hex.EncodeToString(lockScript), hex.EncodeToString(tracker.CurrentLockScriptAddress))
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
	if err != nil || !rep.OK {
		ctx.Logger.Error("error in broadcasting internal addsignature", err, rep.Log)
		return
	}
}

func (j *JobAddSignature) GetJobID() string {
	return j.JobID
}

func (j JobAddSignature) IsDone() bool {
	return j.Status == jobs.Completed
}
