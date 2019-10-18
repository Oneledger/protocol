/*

 */

package identity

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/client"
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

func NewAddSignatureJob(trackerName string) Job {

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

func (j *JobAddSignature) DoMyJob(ctx *JobsContext) {

	tracker, err := ctx.trackers.Get(j.TrackerName)
	if err != nil {
		return
	}

	lockTx := wire.NewMsgTx(wire.TxVersion)
	err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessTx))
	if err != nil {
		//
	}

	sig, err := txscript.RawTxInSignature(lockTx, 0, tracker.CurrentLockScript, txscript.SigHashAll,
		ctx.BTCPrivKey)
	if err != nil {
		fmt.Println(err, "RawTxInSignature")
	}

	addrPubKey, err := btcutil.NewAddressPubKey(ctx.BTCPrivKey.PubKey().SerializeCompressed(), ctx.Params)

	addSigData := btc.AddSignature{
		TrackerName:      j.TrackerName,
		ValidatorPubKey:  addrPubKey,
		BTCSignature:     sig,
		ValidatorAddress: ctx.ValidatorAddress,
		Memo:             j.JobID,
	}

	txData, err := addSigData.Marshal()
	if err != nil {

	}

	tx := action.RawTx{
		Type: action.BTC_ADD_SIGNATURE,
		Data: txData,
		Fee:  nil,
		Memo: j.JobID,
	}

	req := client.InternalBroadcastRequest{
		RawTx: tx,
	}
	rep := client.BroadcastReply{}

	err = ctx.service.InternalBroadcast(req, &rep)
	if err != nil {
		// TODO
	}
}

func (j *JobAddSignature) IsMyJobDone(ctx *JobsContext) bool {

	tracker, err := ctx.trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	addrPubKey, err := btcutil.NewAddressPubKey(ctx.BTCPrivKey.PubKey().SerializeCompressed(), ctx.Params)
	return tracker.Multisig.HasAddressSigned(*addrPubKey)
}

func (j *JobAddSignature) IsSufficient(ctx *JobsContext) bool {
	tracker, err := ctx.trackers.Get(j.TrackerName)
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
