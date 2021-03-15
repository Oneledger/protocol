/*

 */

package btc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
)

type FailedBroadcastReset struct {
	TrackerName      string
	ValidatorAddress action.Address
}

func (fbr *FailedBroadcastReset) Signers() []action.Address {
	return []action.Address{
		fbr.ValidatorAddress,
	}
}

func (fbr *FailedBroadcastReset) Type() action.Type {
	return action.BTC_FAILED_BROADCAST_RESET
}

func (fbr *FailedBroadcastReset) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(fbr.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: []byte(fbr.ValidatorAddress.String()),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(fbr.TrackerName),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (fbr *FailedBroadcastReset) TagsFailed() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(fbr.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: []byte(fbr.ValidatorAddress.String()),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(fbr.TrackerName),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.lock_redeem_status"),
		Value: []byte("failure"),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

func (fbr *FailedBroadcastReset) Marshal() ([]byte, error) {
	return json.Marshal(fbr)
}

func (fbr *FailedBroadcastReset) Unmarshal(data []byte) error {
	return json.Unmarshal(data, fbr)
}

type btcBroadcastFailureReset struct {
}

func (btcBroadcastFailureReset) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	failedBroadcastReset := FailedBroadcastReset{}

	err := failedBroadcastReset.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), failedBroadcastReset.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	if !ctx.Validators.IsValidatorAddress(failedBroadcastReset.ValidatorAddress) {
		return false, errors.New("only validator can report a broadcast")
	}

	tracker, err := ctx.BTCTrackers.Get(failedBroadcastReset.TrackerName)
	if err != nil {
		return false, err
	}

	if tracker.State != bitcoin.BusyBroadcasting {
		return false, errors.New("tracker not broadcasting")
	}

	return true, nil
}

func (btcBroadcastFailureReset) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runBroadcastFailureReset(ctx, tx)
}

func (btcBroadcastFailureReset) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runBroadcastFailureReset(ctx, tx)
}

func (btcBroadcastFailureReset) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

func runBroadcastFailureReset(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	fbr := FailedBroadcastReset{}

	err := fbr.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	if !ctx.Validators.IsValidatorAddress(fbr.ValidatorAddress) {
		return false, action.Response{Log: "broadcast reporter not found in validator list"}
	}

	tracker, err := ctx.BTCTrackers.Get(fbr.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", fbr.TrackerName)}
	}

	if tracker.State != bitcoin.BusyBroadcasting {
		return false, action.Response{Log: fmt.Sprintf("tracker not broadcasting: ", fbr.TrackerName)}
	}

	validatorSignedFlag := false
	for _, fv := range tracker.ResetVotes {
		if bytes.Equal(fv, fbr.ValidatorAddress) {
			validatorSignedFlag = true
		}
	}

	if !validatorSignedFlag {
		tracker.ResetVotes = append(tracker.ResetVotes, fbr.ValidatorAddress)
	}

	valSet, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return false, action.Response{Log: "cannot get validator set"}
	}

	nValidators := len(valSet)
	votesThresholdForReset := (2 * nValidators) / 3

	// are there enough reset votes?
	if len(tracker.ResetVotes) < votesThresholdForReset {
		err = ctx.BTCTrackers.SetTracker(fbr.TrackerName, tracker)
		if err != nil {
			return false, action.Response{Log: "failed to save tracker"}
		}

		return true, action.Response{
			Events: action.GetEvent(fbr.Tags(), "btc_broadcast_reset_pending"),
		}
	}

	// if the process is redeem return the user oBTC
	if tracker.ProcessType == bitcoin.ProcessTypeRedeem {
		amount := tracker.CurrentBalance - tracker.ProcessBalance

		btcCurr, ok := ctx.Currencies.GetCurrencyByName("BTC")
		if !ok {
			return false, action.Response{Log: "failed to find currency BTC"}
		}
		coin := btcCurr.NewCoinFromUnit(amount)

		err = ctx.Balances.AddToAddress(tracker.ProcessOwner, coin)
		if err != nil {
			return false, action.Response{Log: "failed to add currency err:" + err.Error()}
		}

		tally := action.Address(ctx.BTCTrackers.GetOption().TotalSupplyAddr)
		err = ctx.Balances.AddToAddress(tally, coin)
		if err != nil {
			return false, action.Response{Log: "failed to add currency err:" + err.Error()}
		}
	}

	tracker.Multisig.Msg = nil
	tracker.Multisig.Signatures = []keys.BTCSignature{}

	tracker.State = bitcoin.Available
	tracker.ProcessTxId = nil
	tracker.ProcessBalance = 0
	tracker.ProcessUnsignedTx = nil

	tracker.ProcessOwner = nil
	tracker.ProcessType = bitcoin.ProcessTypeNone

	tracker.FinalityVotes = []keys.Address{}
	tracker.ResetVotes = []keys.Address{}

	err = ctx.BTCTrackers.SetTracker(fbr.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to save tracker"}
	}

	return true, action.Response{Events: action.GetEvent(fbr.Tags(), "btc_broadcast_reset_complete")}
}
