/*

 */

package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/wire"

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

	isFirstTxn := len(tx.TxIn) == 1
	op := tx.TxIn[0].PreviousOutPoint

	if !isFirstTxn && op.Hash != *tracker.CurrentTxId {
		return false, errors.New("txn doesn't match tracker")
	}
	if !isFirstTxn && op.Index != 0 {
		return false, errors.New("txn doesn't match tracker")
	}
	if isFirstTxn && tx.TxOut[0].Value != lock.LockAmount+tracker.CurrentBalance {
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

	vs, err := ctx.Validators.GetValidatorSet()
	threshold := (len(vs) * 2 / 3) + 1
	list := make([]keys.Address, 0, len(vs))
	for i := range vs {
		list = append(list, vs[i].BTCAddressPubKey.ScriptAddress())
	}

	tracker.State = bitcoin.BusySigningTrackerState
	tracker.ProcessOwner = lock.Locker
	tracker.ProcessUnsignedTx = lock.BTCTxn // with user signature
	tracker.Multisig, err = keys.NewBTCMultiSig(lock.BTCTxn, threshold, list)
	tracker.ProcessBalance = tracker.CurrentBalance + lock.LockAmount

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

	ctx.Logger.Debug(hex.EncodeToString(lock.BTCTxn))

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
	tracker.Multisig.Msg = lock.BTCTxn
	tracker.ProcessBalance = tracker.CurrentBalance + lock.LockAmount

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
		Info: fmt.Sprintf("tracker: %s", lock.TrackerName),
	}
}

func (btcLockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
	// return true, action.Response{}
}
