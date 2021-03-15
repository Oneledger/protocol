/*

 */

package btc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

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

func (bl Lock) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(bl.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.locker"),
		Value: bl.Locker.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(bl.TrackerName),
	}
	la := strconv.FormatInt(bl.LockAmount, 10)
	tag4 := kv.Pair{
		Key:   []byte("tx.lock_amount"),
		Value: []byte(la),
	}
	tag5 := kv.Pair{
		Key:   []byte("tx.lock_currency"),
		Value: []byte("BTC"),
	}

	tags = append(tags, tag, tag2, tag3, tag4, tag5)
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

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
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

	opt := ctx.BTCTrackers.GetConfig()

	if !ValidateExtLockStructure(tracker, tx, opt.BTCParams) {
		return false, errors.New("err in ext lock txn")
	}

	isFirstLock := tracker.CurrentTxId == nil
	if !bitcoin2.ValidateLock(tx, opt.BlockCypherToken, opt.BlockCypherChainType, tracker.ProcessLockScriptAddress,
		tracker.CurrentBalance, lock.LockAmount, isFirstLock) {

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

func (btcLockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
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

	opt := ctx.BTCTrackers.GetConfig()
	cdOption := ctx.BTCTrackers.GetOption()

	if !ValidateExtLockStructure(tracker, btcTx, opt.BTCParams) {
		return false, action.Response{Log: "err in ext lock txn structure"}
	}

	newTrackerBalance := btcTx.TxOut[0].Value
	lockAmount := newTrackerBalance - tracker.CurrentBalance
	if lockAmount != lock.LockAmount {
		return false, action.Response{Log: "err in lock amount"}
	}

	curr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	if !ok {
		return false, action.Response{Log: fmt.Sprintf("BTC currency not available", lock.TrackerName)}
	}

	lockCoin := curr.NewCoinFromUnit(lock.LockAmount)
	tally := action.Address(cdOption.TotalSupplyAddr)
	balCoin, err := ctx.Balances.GetBalanceForCurr(tally, &curr)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("unable to get btc lock total balance", lock.TrackerName)}
	}

	totalSupplyCoin := curr.NewCoinFromString(cdOption.TotalSupply)

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
		Events: action.GetEvent(lock.Tags(), "btc_lock"),
	}
}
