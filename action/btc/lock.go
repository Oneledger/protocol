/*

 */

package btc

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/data/bitcoin"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/tendermint/tendermint/libs/common"
)

type Lock struct {
	Locker      action.Address
	Amount      int64
	TrackerName string
	ProcessUTXO bitcoin.UTXO
}

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

func (btcLockTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	lock := Lock{}
	err := lock.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	_, err = action.ValidateBasic(signedTx.RawBytes(), lock.Signers(), signedTx.Signatures)
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

	tracker.State = bitcoin.BusyBroadcastingTrackerState

	tracker.ProcessUTXO = &lock.ProcessUTXO
	tracker.ProcessOwner = lock.Locker

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

	tracker.State = bitcoin.BusyBroadcastingTrackerState

	tracker.ProcessUTXO = &lock.ProcessUTXO
	tracker.ProcessOwner = lock.Locker

	err = ctx.Trackers.SetTracker(lock.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to update tracker"}
	}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}

func (btcLockTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
