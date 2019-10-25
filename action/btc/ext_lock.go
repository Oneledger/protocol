/*

 */

package btc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/wire"

	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/data/bitcoin"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/tendermint/tendermint/libs/common"
)

type Lock struct {
	Locker      action.Address
	TrackerName string
	BTCTxn      []byte
	LockAmount  int64
}

var _ action.Msg = &Lock{}

func (bl Lock) Signers() []action.Address {
	return []action.Address{bl.Locker}
}

func (bl Lock) Type() action.Type {
	return action.BTC_LOCK
}

func (bl Lock) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(bl.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.locker"),
		Value: bl.Locker.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (bl Lock) Marshal() ([]byte, error) {
	return json.Marshal(bl)
}

func (bl *Lock) Unmarshal(data []byte) error {
	return json.Unmarshal(data, bl)
}

type btcLockTx struct {
}

var _ action.Tx = btcLockTx{}

func (btcLockTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	lock := Lock{}
	err := lock.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), lock.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeeOpt, signedTx.Fee)
	if err != nil {
		return false, err
	}

	tracker, err := ctx.Trackers.Get(lock.TrackerName)
	if err != nil {
		return false, err
	}

	if !tracker.IsAvailable() {
		return false, errors.New("tracker not available")
	}

	tx := wire.NewMsgTx(wire.TxVersion)

	buf := bytes.NewBuffer(lock.BTCTxn)
	tx.Deserialize(buf)

	op := tx.TxIn[0].PreviousOutPoint
	if op.Hash != *tracker.CurrentTxId {
		return false, errors.New("txn doesn't match tracker")
	}
	if op.Index != 0 {
		return false, errors.New("txn doesn't match tracker")
	}
	if tx.TxOut[0].Value != lock.LockAmount+tracker.CurrentBalance {
		return false, errors.New("txn doesn't match tracker")
	}

	return true, nil
}

func (btcLockTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	lock := Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.Trackers.Get(lock.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", lock.TrackerName)}
	}

	if !tracker.IsAvailable() {
		return false, action.Response{Log: fmt.Sprintf("tracker not available for lock: ", lock.TrackerName)}
	}

	tracker.State = bitcoin.BusySigningTrackerState
	tracker.ProcessOwner = lock.Locker
	tracker.ProcessUnsignedTx = lock.BTCTxn // with user signature

	err = ctx.Trackers.SetTracker(lock.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to update tracker"}
	}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}

func (btcLockTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	lock := Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.Trackers.Get(lock.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", lock.TrackerName)}
	}

	if !tracker.IsAvailable() {
		return false, action.Response{Log: fmt.Sprintf("tracker not available for lock: ", lock.TrackerName)}
	}

	tracker.State = bitcoin.BusySigningTrackerState
	tracker.ProcessOwner = lock.Locker
	tracker.ProcessUnsignedTx = lock.BTCTxn

	currentUTXO := bitcoin.NewUTXO(tracker.CurrentTxId, 0, tracker.CurrentBalance)
	processUTXO := bitcoin.NewUTXO(nil, 0, tracker.CurrentBalance+lock.LockAmount)

	cd := bitcoin2.NewChainDriver("")
	tracker.ProcessUnsignedTx = cd.PrepareLock(currentUTXO, processUTXO, tracker.ProcessLockScriptAddress)

	err = ctx.Trackers.SetTracker(lock.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to update tracker"}
	}

	if ctx.JobStore != nil {
		job := NewAddSignatureJob(lock.TrackerName)
		err = ctx.JobStore.SaveJob(job)
		if err != nil {
			return false, action.Response{Log: "job serialization failed"}
		}
	}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}

func (btcLockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
