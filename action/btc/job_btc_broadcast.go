/*

 */

package btc

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type JobBTCBroadcast struct {
	Type string

	TrackerName string

	JobID string

	BroadcastSuccessful bool

	Done bool
}

func (j *JobBTCBroadcast) DoMyJob(ctxI interface{}) {

	ctx, _ := ctxI.(*action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return
	}

	if !j.BroadcastSuccessful {

		lockTx := wire.NewMsgTx(wire.TxVersion)
		err = lockTx.Deserialize(bytes.NewReader(tracker.ProcessUnsignedTx))
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

		lockScript, err := ctx.LockScripts.GetLockScript(tracker.CurrentLockScriptAddress)
		if err != nil {

		}

		builder.AddFullData(lockScript)
		sigScript, err := builder.Script()

		cd := bitcoin.NewChainDriver(ctx.BlockCypherToken)
		lockTx = cd.AddLockSignature(tracker.ProcessUnsignedTx, sigScript)

		buf := bytes.NewBuffer([]byte{})
		lockTx.Serialize(buf)

		// TODO load from config
		connCfg := &rpcclient.ConnConfig{
			Host:         ctx.BTCNodeAddress + ":" + ctx.BTCRPCPort,
			User:         ctx.BTCRPCUsername,
			Pass:         ctx.BTCRPCPassword,
			HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
			DisableTLS:   true, // Bitcoin core does not provide TLS by default
		}

		clt, err := rpcclient.New(connCfg, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = cd.BroadcastTx(lockTx, clt)
		if err != nil {
			j.BroadcastSuccessful = true
		}

	} else {

		cd := bitcoin.NewChainDriver(ctx.BlockCypherToken)

		chain := "test3"
		switch ctx.BTCChainnet {
		case "testnet3":
			chain = "test3"
		case "testnet":
			chain = "test"
		case "mainnet":
			chain = "main"
		}

		ok, _ := cd.CheckFinality(tracker.ProcessTxId, ctx.BlockCypherToken, chain)
		if !ok {
			return
		}

		data := [4]byte{}
		_, err = io.ReadFull(rand.Reader, data[:])
		if err != nil {
			return
		}

		reportFinalityMint := ReportFinalityMint{
			TrackerName:      j.TrackerName,
			OwnerAddress:     tracker.ProcessOwner,
			ValidatorAddress: ctx.ValidatorAddress,
			RandomBytes:      data[:],
		}

		txData, err := reportFinalityMint.Marshal()
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
}

func (j *JobBTCBroadcast) IsMyJobDone(ctxI interface{}) bool {
	ctx, _ := ctxI.(*action.JobsContext)

	tracker, err := ctx.Trackers.Get(j.TrackerName)
	if err != nil {
		return false
	}

	if tracker.IsAvailable() {
		return true
	}

	for _, fv := range tracker.FinalityVotes {
		if bytes.Equal(fv, ctx.ValidatorAddress) {
			return true
		}
	}

	return false
}

func (j *JobBTCBroadcast) IsSufficient(ctxI interface{}) bool {
	return j.IsMyJobDone(ctxI)
}

func (j *JobBTCBroadcast) DoFinalize() {
	j.Done = true
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
