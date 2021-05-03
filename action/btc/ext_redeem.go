/*

 */

package btc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/data/bitcoin"
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

func (bl Redeem) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(bl.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.redeem"),
		Value: bl.Redeemer.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.redeemer"),
		Value: bl.Redeemer,
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.redeem_currency"),
		Value: []byte("BTC"),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
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
	redeem := Redeem{}
	err := redeem.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), redeem.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	tracker, err := ctx.BTCTrackers.Get(redeem.TrackerName)
	if err != nil {
		return false, err
	}

	if !tracker.IsAvailable() {
		return false, errors.New("tracker not available")
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	buf := bytes.NewBuffer(redeem.BTCTxn)
	err = tx.Deserialize(buf)
	if err != nil {
		return false, errors.New("err in the deserialized btc txn")
	}

	if len(tx.TxIn) != 1 {
		return false, errors.New("input should be 1 in txn")
	}
	op := tx.TxIn[0].PreviousOutPoint

	// check if the source 0 in the txn is our tracker
	if op.Hash != *tracker.CurrentTxId {
		return false, errors.New("txn doesn't match tracker")
	}

	if op.Index != 0 {
		return false, errors.New("txn doesn't match tracker")
	}

	opt := ctx.BTCTrackers.GetConfig()
	if !bitcoin2.ValidateRedeem(tx, opt.BlockCypherToken, opt.BlockCypherChainType, tracker.CurrentTxId,
		tracker.ProcessLockScriptAddress, tracker.CurrentBalance, redeem.RedeemAmount) {

		return false, errors.New("txn doesn't match tracker")
	}

	return true, nil
}

func (btcRedeemTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runExtRedeem(ctx, tx)
}

func (btcRedeemTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runExtRedeem(ctx, tx)
}

func (btcRedeemTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.State.ConsumeUpfront(701660)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runExtRedeem(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

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

	btcTx := wire.NewMsgTx(wire.TxVersion)
	buf := bytes.NewBuffer(redeem.BTCTxn)
	err = btcTx.Deserialize(buf)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("err bad btc txn: %s", redeem.TrackerName)}
	}

	if len(btcTx.TxIn) != 1 {
		return false, action.Response{Log: fmt.Sprintf("err bad btc txn: %s", redeem.TrackerName)}
	}
	op := btcTx.TxIn[0].PreviousOutPoint

	// check if the source 0 in the txn is our tracker
	if op.Hash != *tracker.CurrentTxId {
		return false, action.Response{Log: fmt.Sprintf("err bad btc txn: %s", redeem.TrackerName)}
	}

	if op.Index != 0 {
		return false, action.Response{Log: fmt.Sprintf("err bad btc txn: %s", redeem.TrackerName)}
	}

	if !bytes.Equal(btcTx.TxOut[0].PkScript, tracker.ProcessLockScriptAddress) {
		return false, action.Response{Log: fmt.Sprintf("err incorrect btc lock address ", redeem.TrackerName)}
	}

	tracker.ProcessType = bitcoin.ProcessTypeRedeem
	tracker.ProcessOwner = redeem.Redeemer

	tracker.Multisig.Msg = redeem.BTCTxn
	tracker.ProcessBalance = tracker.CurrentBalance - redeem.RedeemAmount
	tracker.ProcessUnsignedTx = redeem.BTCTxn // with user signature
	tracker.State = bitcoin.Requested

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
	coin := btcCurr.NewCoinFromUnit(redeem.RedeemAmount)

	err = ctx.Balances.MinusFromAddress(redeem.Redeemer, coin)
	if err != nil {
		return false, action.Response{Log: "failed to subtract currency err:" + err.Error()}
	}

	tally := action.Address(ctx.BTCTrackers.GetOption().TotalSupplyAddr)
	err = ctx.Balances.MinusFromAddress(tally, coin)
	if err != nil {
		return false, action.Response{Log: "failed to subtract currency err:" + err.Error()}
	}

	return true, action.Response{
		Events: action.GetEvent(redeem.Tags(), "btc_redeem"),
	}
}
