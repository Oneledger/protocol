/*

 */

package identity

import (
	"bytes"
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/client"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
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

func (j *JobAddSignature) GetType() string {
	return JobTypeAddSignature
}

type doJobData struct {
	BTCPrivKey       *btcec.PrivateKey
	Params           *chaincfg.Params
	ValidatorAddress action.Address
}

func (j *JobAddSignature) DoMyJob(ctx *JobsContext, data interface{}) {

	inp := data.(doJobData)

	tracker, err := ctx.trackers.Get(j.TrackerName)
	if err != nil {
		return
	}

	lockTx := wire.NewMsgTx(wire.TxVersion)
	lockTx.Deserialize(bytes.NewReader(tracker.ProcessTx))

	sig, err := txscript.RawTxInSignature(lockTx, 0, tracker.CurrentLockScript, txscript.SigHashAll,
		inp.BTCPrivKey)
	if err != nil {
		fmt.Println(err, "RawTxInSignature")
	}

	addrPubKey, err := btcutil.NewAddressPubKey(inp.BTCPrivKey.PubKey().SerializeCompressed(), inp.Params)

	addSigData := btc.AddSignature{
		TrackerName:      j.TrackerName,
		ValidatorPubKey:  addrPubKey,
		BTCSignature:     sig,
		ValidatorAddress: inp.ValidatorAddress,
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

func (j *JobAddSignature) IsMyJobDone(addr btcutil.AddressPubKey, ctx *JobsContext) bool {

	tracker, err := ctx.trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	return tracker.Multisig.HasAddressSigned(addr)
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
