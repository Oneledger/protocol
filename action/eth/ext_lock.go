package eth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bitcoin"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type Lock struct {
	Locker      action.Address
	TrackerName ethcommon.Hash
	ETHTxn      []byte
	LockAmount  int64
}

var _ action.Msg = &Lock{}

func (et Lock) Signers() []action.Address {
	return []action.Address{et.Locker}
}

func (et Lock) Type() action.Type {
	return action.ETH_LOCK
}

func (et Lock) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(et.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.locker"),
		Value: et.Locker.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (et Lock) Marshal() ([]byte, error) {
	return json.Marshal(et)
}

func (et *Lock) Unmarshal(data []byte) error {
	return json.Unmarshal(data, et)
}

type ethLockTx struct {
}

var _ action.Tx = ethLockTx{}

func (ethLockTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	lock := &Lock{}
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
    // Check lock fields for incoming trasaction
    // Verify if lockAmount == eth.tx.value
    //TODO : Verify beninfiaciary address in ETHTX == locker (Phase 2)
	return true, nil
}

func (ethLockTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	lock := &Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	//tracker, err := ctx.Trackers.Get(lock.TrackerName)

    // Create ethereum tracker
	tracker := NewTracker()
	tracker.State = bitcoin.BusySigningTrackerState
	tracker.ProcessOwner = lock.Locker
	tracker.ProcessUnsignedTx = lock.ETHTxn // with user signature

    // Save eth Tracker

	return true, action.Response{
		Tags: lock.Tags(),
	}
}

func (ethLockTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
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
	tracker.ProcessUnsignedTx = lock.ETHTxn

	currentUTXO := bitcoin.NewUTXO(tracker.CurrentTxId, 0, tracker.CurrentBalance)
	processUTXO := bitcoin.NewUTXO(nil, 0, tracker.CurrentBalance+lock.LockAmount)

	cd := bitcoin2.NewChainDriver("")
	tracker.ProcessUnsignedTx = cd.PrepareLock(currentUTXO, processUTXO, tracker.ProcessLockScriptAddress)

	err = ctx.Trackers.SetTracker(lock.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to update tracker"}
	}

	if ctx.JobStore != nil {
		//TODO write function ETHNewAddJOBSignature
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

func (ethLockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	panic("implement me")
}
