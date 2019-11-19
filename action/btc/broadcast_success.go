/*

 */

package btc

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/tendermint/tendermint/libs/common"
)

type BroadcastSuccess struct {
	TrackerName      string
	ValidatorAddress action.Address
	BTCTxID          chainhash.Hash
}

func (b BroadcastSuccess) Signers() []action.Address {
	return []action.Address{
		b.ValidatorAddress,
	}
}

func (b BroadcastSuccess) Type() action.Type {
	return action.BTC_BROADCAST_SUCCESS
}

func (b BroadcastSuccess) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(b.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.locker"),
		Value: []byte(b.ValidatorAddress.String()),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (b BroadcastSuccess) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b BroadcastSuccess) Unmarshal(data []byte) error {
	return json.Unmarshal(data, b)
}

type btcBroadcastSuccessTx struct {
}

func (b *btcBroadcastSuccessTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {

	broadcastSuccess := BroadcastSuccess{}

	err := broadcastSuccess.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), broadcastSuccess.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	if !ctx.Validators.IsValidatorAddress(broadcastSuccess.ValidatorAddress) {
		return false, errors.New("only validator can report a broadcast")
	}

	tracker, err := ctx.BTCTrackers.Get(broadcastSuccess.TrackerName)
	if err != nil {
		return false, err
	}

	if tracker.State != bitcoin.BusyBroadcasting {
		return false, errors.New("tracker not broadcasting")
	}

	return true, nil
}

func (b *btcBroadcastSuccessTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return b.process(ctx, tx)
}

func (b *btcBroadcastSuccessTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return b.process(ctx, tx)
}

func (b *btcBroadcastSuccessTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

func (b *btcBroadcastSuccessTx) process(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	broadcastSuccess := BroadcastSuccess{}
	err := broadcastSuccess.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	if !ctx.Validators.IsValidatorAddress(broadcastSuccess.ValidatorAddress) {
		return false, action.Response{Log: "broadcast reporter not found in validator list"}
	}

	tracker, err := ctx.BTCTrackers.Get(broadcastSuccess.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", broadcastSuccess.TrackerName)}
	}

	if tracker.State != bitcoin.BusyBroadcasting {
		return false, action.Response{Log: fmt.Sprintf("tracker not broadcasting: ", broadcastSuccess.TrackerName)}
	}

	tracker.ProcessTxId = &broadcastSuccess.BTCTxID

	dat := bitcoin.BTCTransitionContext{Tracker: tracker}
	_, err = bitcoin.Engine.Process("reportBroadcastSuccess", dat, tracker.State)
	if err != nil {

	}

	err = ctx.BTCTrackers.SetTracker(broadcastSuccess.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error updating tracker store: %s, error: ", broadcastSuccess.TrackerName, err)}
	}

	return true, action.Response{
		Tags: broadcastSuccess.Tags(),
	}
}
