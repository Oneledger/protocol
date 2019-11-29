/*

 */

package event

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/chains/bitcoin"
	bitcoin2 "github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/jobs"
)

type JobBTCBroadcast struct {
	Type string

	TrackerName string

	JobID string

	Status jobs.Status
}

func NewBTCBroadcastJob(trackerName, id string) jobs.Job {

	return &JobBTCBroadcast{
		Type:        JobTypeBTCBroadcast,
		TrackerName: trackerName,
		JobID:       id,
		Status:      jobs.New,
	}
}

func (j *JobBTCBroadcast) DoMyJob(ctxI interface{}) {

	ctx, _ := ctxI.(*JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		ctx.Logger.Error("err trying to deserialize tracker: ", j.TrackerName, err)
		return
	}

	if tracker.State != bitcoin2.BusyBroadcasting {
		j.Status = jobs.Completed
		return
	}

	lockTx := wire.NewMsgTx(wire.TxVersion)
	err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessUnsignedTx))
	if err != nil {
		ctx.Logger.Error("err trying to deserialize btc txn: ", err, j.TrackerName)
		return
	}

	signatures := tracker.Multisig.GetSignaturesInOrder()

	builder := txscript.NewScriptBuilder().AddOp(txscript.OP_FALSE)
	for i := range signatures {
		if i == tracker.Multisig.M {
			// break
		}

		builder.AddData(signatures[i])
	}

	lockScript, err := ctx.LockScripts.GetLockScript(tracker.CurrentLockScriptAddress)
	if err != nil {
		ctx.Logger.Error("err trying to get lockscript ", err, j.TrackerName)
		return
	}

	builder.AddData(lockScript)
	sigScript, err := builder.Script()
	if err != nil {
		ctx.Logger.Error("error in building sig script", err)
	}

	cd := bitcoin.NewChainDriver(ctx.BlockCypherToken)
	lockTx = cd.AddLockSignature(tracker.ProcessUnsignedTx, sigScript)

	buf := bytes.NewBuffer([]byte{})
	err = lockTx.Serialize(buf)
	if err != nil {
		ctx.Logger.Error("err trying to serialize btc final txn ", err, j.TrackerName)
		return
	}

	var txBytes []byte
	buf = bytes.NewBuffer(txBytes)
	lockTx.Serialize(buf)
	txBytes = buf.Bytes()

	ctx.Logger.Debug(hex.EncodeToString(txBytes))

	if !(len(lockTx.TxIn) == 1 && tracker.ProcessType == bitcoin2.ProcessTypeLock) {

		vm, err := txscript.NewEngine(tracker.CurrentLockScriptAddress, lockTx, 0, txscript.StandardVerifyFlags, nil, nil, tracker.CurrentBalance)
		if err != nil {
			fmt.Println("new engine", err)
			ctx.Logger.Error("error in test engine")
			return
		}
		if err := vm.Execute(); err != nil {
			fmt.Println("vm Execute", err)
			ctx.Logger.Error("error in vm execute")
			return
		}
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         ctx.BTCNodeAddress + ":" + ctx.BTCRPCPort,
		User:         ctx.BTCRPCUsername,
		Pass:         ctx.BTCRPCPassword,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	clt, err := rpcclient.New(connCfg, nil)
	if err != nil {
		ctx.Logger.Error("err trying to connect to bitcoin node", j.TrackerName)
		return
	}

	hash, err := cd.BroadcastTx(lockTx, clt)
	// use dummy hash for testing without broadcasting
	// fmt.Println(clt)
	// hash, err := chainhash.NewHashFromStr("cb0eee8e68b474cd1e845702052847dcbf248eb5a08ec498e887108842019d06")
	if err == nil {

		ctx.Logger.Info("bitcoin tx successful", hash)

		bs := btc.BroadcastSuccess{
			tracker.Name,
			ctx.ValidatorAddress,
			*hash,
		}

		txData, err := bs.Marshal()
		if err != nil {
			ctx.Logger.Error("error while preparing mint txn ", err, j.TrackerName)
			return
		}
		tx := action.RawTx{
			Type: action.BTC_BROADCAST_SUCCESS,
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
			ctx.Logger.Error("error while broadcasting finality vote and mint txn ", err, j.TrackerName)
			return
		}

		ctx.Logger.Info("btc success internal tx broadcast")
		ctx.Logger.Infof("%#v \n", rep)
		j.Status = jobs.Completed

	} else {
		ctx.Logger.Error("broadcast failed err: ", err, " tracker: ", j.TrackerName)
	}

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

func (j JobBTCBroadcast) IsDone() bool {
	return j.Status == jobs.Completed
}
