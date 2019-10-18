/*

 */

package identity

import (
	"bytes"
	"fmt"

	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type JobBTCBroadcast struct {
	Type string

	TrackerName string

	JobID string

	Done bool
}

func (j *JobBTCBroadcast) DoMyJob(ctx *JobsContext) {
	tracker, err := ctx.trackers.Get(j.TrackerName)
	if err != nil {
		return
	}

	lockTx := wire.NewMsgTx(wire.TxVersion)
	err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessTx))
	if err != nil {
		//
	}

	type sign []byte
	btcSignatures := tracker.Multisig.GetSignatures()
	signatures := make([]sign, len(btcSignatures))
	for i := range btcSignatures {
		index := btcSignatures[i].Index
		signatures[index] = btcSignatures[i].Sign
	}

	builder := txscript.NewScriptBuilder().AddOp(txscript.OP_FALSE)
	for i := range signatures {
		builder.AddData(signatures[i])
		if i == tracker.Multisig.M {
			break
		}
	}
	builder.AddFullData(tracker.CurrentLockScript)
	sigScript, err := builder.Script()

	cd := bitcoin.NewChainDriver(ctx.BlockCypherToken)
	lockTx = cd.AddLockSignature(tracker.ProcessTx, sigScript)

	buf := bytes.NewBuffer([]byte{})
	lockTx.Serialize(buf)

	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18831",
		User:         "oltest01",
		Pass:         "olpass01",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	clt, err := rpcclient.New(connCfg, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	txHash, err := cd.BroadcastTx(lockTx, clt)

}

func (j *JobBTCBroadcast) IsMyJobDone(ctx *JobsContext) bool {
	panic("implement me")
}

func (j *JobBTCBroadcast) IsSufficient() bool {
	panic("implement me")
}

func (j *JobBTCBroadcast) DoFinalize() {
	panic("implement me")
}

/*
	simple getters
*/
func (j *JobBTCBroadcast) GetType() string {
	return JobTypeBTCBroadcast
}

func (j *JobBTCBroadcast) GetJobID() string {
	return j.JobID
}

func (j *JobBTCBroadcast) IsDone() bool {
	return j.Done
}
