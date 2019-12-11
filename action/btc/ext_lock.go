/*

 */

package btc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/data/bitcoin"
)

type Lock struct {
	// OLT address of the person locking the BTC
	Locker action.Address

	// Name of the tracker used to register this txn
	TrackerName string

	// BTC Txn as a byte array
	BTCTxn []byte

	// The amount in satoshi to lock
	LockAmount int64
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

	tracker, err := ctx.BTCTrackers.Get(lock.TrackerName)
	if err != nil {
		return false, err
	}

	if !tracker.IsAvailable() {
		return false, errors.New("tracker not available")
	}

	tx := wire.NewMsgTx(wire.TxVersion)

	buf := bytes.NewBuffer(lock.BTCTxn)
	err = tx.Deserialize(buf)
	if err != nil {
		return false, errors.New("err in deserializing btc txn")
	}

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

	if !bitcoin2.ValidateLock(tx, ctx.BlockCypherToken, ctx.BlockCypherChainType, tracker.CurrentTxId,
		tracker.ProcessLockScriptAddress, tracker.CurrentBalance, lock.LockAmount) {

		return false, errors.New("txn doesn't match tracker")
	}

	return true, nil
}

func (btcLockTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runBTCLock(ctx, tx)
}

func (btcLockTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runBTCLock(ctx, tx)
}

func (btcLockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.State.ConsumeUpfront(701660)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runBTCLock(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	lock := Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.BTCTrackers.Get(lock.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", lock.TrackerName)}
	}

	if !tracker.IsAvailable() {
		return false, action.Response{Log: fmt.Sprintf("tracker not available for lock: ", lock.TrackerName)}
	}

	btcTx := wire.NewMsgTx(wire.TxVersion)

	buf := bytes.NewBuffer(lock.BTCTxn)
	err = btcTx.Deserialize(buf)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("err desrializing btc txn ", lock.TrackerName)}
	}

	isFirstTxn := len(btcTx.TxIn) == 1
	op := btcTx.TxIn[0].PreviousOutPoint

	if !isFirstTxn && op.Hash != *tracker.CurrentTxId {
		return false, action.Response{Log: fmt.Sprintf("err incorrect btc txn", lock.TrackerName)}
	}
	if !isFirstTxn && op.Index != 0 {
		return false, action.Response{Log: fmt.Sprintf("err incorrect btc txn ", lock.TrackerName)}
	}
	if isFirstTxn && btcTx.TxOut[0].Value != lock.LockAmount+tracker.CurrentBalance {
		return false, action.Response{Log: fmt.Sprintf("err incorrect btc txn ", lock.TrackerName)}
	}

	if !bytes.Equal(btcTx.TxOut[0].PkScript, tracker.ProcessLockScriptAddress) {
		return false, action.Response{Log: fmt.Sprintf("err incorrect btc lock address ", lock.TrackerName)}
	}

	curr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	if !ok {
		return false, action.Response{Log: fmt.Sprintf("BTC currency not available", lock.TrackerName)}
	}

	lockCoin := curr.NewCoinFromUnit(lock.LockAmount)
	tally := action.Address(lockBalanceAddress)
	balCoin, err := ctx.Balances.GetBalanceForCurr(tally, &curr)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("unable to get btc lock total balance", lock.TrackerName)}
	}

	totalSupplyCoin := curr.NewCoinFromInt(totalBTCSupply)

	if !balCoin.Plus(lockCoin).LessThanEqualCoin(totalSupplyCoin) {
		return false, action.Response{Log: fmt.Sprintf("btc lock exceeded limit", lock.TrackerName)}
	}

	tracker.ProcessType = bitcoin.ProcessTypeLock
	tracker.ProcessOwner = lock.Locker
	tracker.Multisig.Msg = lock.BTCTxn

	tracker.ProcessBalance = tracker.CurrentBalance + lock.LockAmount
	tracker.ProcessUnsignedTx = lock.BTCTxn // with user signature
	tracker.State = bitcoin.Requested

	err = ctx.BTCTrackers.SetTracker(lock.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "failed to update tracker"}
	}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}
