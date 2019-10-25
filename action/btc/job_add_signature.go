/*

 */

package btc

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/jobs"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type JobAddSignature struct {
	Type string

	TrackerName string

	JobID string

	Done bool
}

func NewAddSignatureJob(trackerName string) jobs.Job {

	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	return &JobAddSignature{
		JobTypeAddSignature,
		trackerName,
		id,
		false,
	}
}

func (j *JobAddSignature) GetType() string {
	return JobTypeAddSignature
}

type doJobData struct {
}

func (j *JobAddSignature) DoMyJob(ctxI interface{}) {
	ctx, _ := ctxI.(action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return
	}

	lockTx := wire.NewMsgTx(wire.TxVersion)
	err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessUnsignedTx))
	if err != nil {
		//
	}

	lockScript, err := ctx.LockScripts.GetLockScript(tracker.CurrentLockScriptAddress)
	if err != nil {

	}

	sig, err := txscript.RawTxInSignature(lockTx, 0, lockScript, txscript.SigHashAll,
		ctx.BTCPrivKey)
	if err != nil {
		fmt.Println(err, "RawTxInSignature")
	}

	addrPubKey, err := btcutil.NewAddressPubKey(ctx.BTCPrivKey.PubKey().SerializeCompressed(), ctx.Params)

	addSigData := AddSignature{
		TrackerName:      j.TrackerName,
		ValidatorPubKey:  addrPubKey,
		BTCSignature:     sig,
		ValidatorAddress: ctx.ValidatorAddress,
		Memo:             j.JobID,
	}

	txData, err := addSigData.Marshal()
	if err != nil {
		// retry later
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
		// retry later
		return
	}
}

func (j *JobAddSignature) IsMyJobDone(ctxI interface{}) bool {
	ctx, _ := ctxI.(action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	addrPubKey, err := btcutil.NewAddressPubKey(ctx.BTCPrivKey.PubKey().SerializeCompressed(), ctx.Params)
	return tracker.Multisig.HasAddressSigned(*addrPubKey)
}

func (j *JobAddSignature) IsSufficient(ctxI interface{}) bool {

	ctx, _ := ctxI.(action.JobsContext)

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
