/*

 */

package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/bitcoin"

	"github.com/Oneledger/protocol/data/keys"

	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
)

type Redeem struct {
	Redeemer     action.Address
	TrackerName  string
	BTCTxn       []byte
	RedeemAmount int64
}

var _ action.Msg = &Redeem{}

func (bl Redeem) Signers() []action.Address {
	return []action.Address{bl.Redeemer}
}

func (bl Redeem) Type() action.Type {
	return action.BTC_REDEEM
}

func (bl Redeem) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(bl.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.redeem"),
		Value: bl.Redeemer.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

func (bl Redeem) Marshal() ([]byte, error) {
	return json.Marshal(bl)
}

func (bl *Redeem) Unmarshal(data []byte) error {
	return json.Unmarshal(data, bl)
}

type btcRedeemTx struct {
}

var _ action.Tx = btcRedeemTx{}

func (btcRedeemTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
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

	tracker, err := ctx.BTCTrackers.Get(lock.TrackerName)
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

func (btcRedeemTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	redeem := Redeem{}
	err := redeem.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.BTCTrackers.Get(redeem.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", redeem.TrackerName)}
	}

	if !tracker.IsAvailable() {
		return false, action.Response{Log: fmt.Sprintf("tracker not available for redeem: ", redeem.TrackerName)}
	}

	vs, err := ctx.Validators.GetValidatorSet()
	threshold := (len(vs) * 2 / 3) + 1
	list := make([]keys.Address, 0, len(vs))

	for i := range vs {
		ctx.Logger.Debug(i, vs[i].ECDSAPubKey.KeyType)

		addr, err := vs[i].GetBTCScriptAddress(ctx.BTCChainType)
		if err != nil {

		}
		list = append(list, addr)
	}

	tracker.ProcessType = bitcoin.ProcessTypeLock
	tracker.ProcessOwner = redeem.Redeemer
	tracker.Multisig, err = keys.NewBTCMultiSig(redeem.BTCTxn, threshold, list)
	tracker.ProcessBalance = tracker.CurrentBalance - redeem.RedeemAmount
	tracker.ProcessUnsignedTx = redeem.BTCTxn // with user signature

	dat := bitcoin.BTCTransitionContext{Tracker: tracker}
	//_, err = bitcoin.Engine.Process("reserveTracker", dat, tracker.State)
	//if err != nil {
	//	return false, action.Response{Log: "failed transition " + err.Error()}
	//}
	_ = dat
	err = ctx.BTCTrackers.SetTracker(redeem.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to update tracker err:" + err.Error()}
	}

	btcCurr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	if !ok {
		return false, action.Response{Log: "failed to find currency BTC"}
	}
	coin := btcCurr.NewCoinFromInt(redeem.RedeemAmount)

	err = ctx.Balances.MinusFromAddress(redeem.Redeemer, coin)
	if err != nil {
		return false, action.Response{Log: "failed to subtract currency err:" + err.Error()}
	}

	return true, action.Response{
		Tags: redeem.Tags(),
	}
}

func (btcRedeemTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	redeem := Redeem{}
	err := redeem.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	ctx.Logger.Debug(hex.EncodeToString(redeem.BTCTxn))

	vs, err := ctx.Validators.GetValidatorSet()
	threshold := (len(vs) * 2 / 3) + 1
	list := make([]keys.Address, 0, len(vs))

	tracker, err := ctx.BTCTrackers.Get(redeem.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", redeem.TrackerName)}
	}

	if !tracker.IsAvailable() {
		return false, action.Response{Log: fmt.Sprintf("tracker not available for redeem: ", redeem.TrackerName)}
	}
	for i := range vs {
		addr, err := vs[i].GetBTCScriptAddress(ctx.BTCChainType)
		if err != nil {

		}
		list = append(list, addr)
	}

	tracker.ProcessType = bitcoin.ProcessTypeRedeem
	tracker.ProcessOwner = redeem.Redeemer
	tracker.Multisig, err = keys.NewBTCMultiSig(redeem.BTCTxn, threshold, list)
	tracker.ProcessBalance = tracker.CurrentBalance - redeem.RedeemAmount
	tracker.ProcessUnsignedTx = redeem.BTCTxn // with user signature

	dat := bitcoin.BTCTransitionContext{Tracker: tracker}
	//_, err = bitcoin.Engine.Process("reserveTracker", dat, tracker.State)
	//if err != nil {
	//	return false, action.Response{Log: "failed transition " + err.Error()}
	//}
	_ = dat

	err = ctx.BTCTrackers.SetTracker(redeem.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to update tracker"}
	}

	btcCurr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	if !ok {
		return false, action.Response{Log: "failed to find currency BTC"}
	}
	coin := btcCurr.NewCoinFromInt(redeem.RedeemAmount)

	err = ctx.Balances.MinusFromAddress(redeem.Redeemer, coin)
	if err != nil {
		return false, action.Response{Log: "failed to subtract currency err:" + err.Error()}
	}

	return true, action.Response{
		Tags: redeem.Tags(),
		Info: fmt.Sprintf("tracker: %s", redeem.TrackerName),
	}
}

func (btcRedeemTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
	// return true, action.Response{}
}
