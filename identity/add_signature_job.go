/*

 */

package identity

import (
	"bytes"
	"fmt"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type JobAddSignature struct {
	Type string

	TrackerName string
}

func (j *JobAddSignature) GetType() string {
	return JobTypeAddSignature
}


type doJobData struct {
	BTCPrivKey *btcec.PrivateKey
	ValidatorPubKey keys.PublicKey
}
func (j *JobAddSignature) DoValidatorJob(ctx *JobsContext, data interface{}) {

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



	addSigData := btc.AddSignature{
		TrackerName: j.TrackerName,
		ValidatorPubKey: ,
		PubKey: inp.ValidatorPubKey,
		BTCSignature:sig,
		ValidatorAddress:
	}
}



func (j *JobAddSignature) IsMyJobDone(key keys.PrivateKey, ctx *JobsContext) bool {
	handler, err := key.GetHandler()
	if err != nil {
		return false
	}

	pubKeyHandler, err := handler.PubKey().GetHandler()
	if err != nil {
		return false
	}

	tracker, err := ctx.trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	return tracker.Multisig.HasAddressSigned(pubKeyHandler.Address())
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
}