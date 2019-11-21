package eth

import (
	"encoding/json"
	"fmt"
	"strconv"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/ethereum"
)

type Lock struct {
	Locker      action.Address
	TrackerName ethcommon.Hash
	ETHTxn      []byte
	LockAmount  int64
}

var _ action.Msg = &Lock{}

// Signers for the ethereum ext lock is the user who wishes to lock his ether
func (et Lock) Signers() []action.Address {
	return []action.Address{et.Locker}
}

// Type for ethlock
func (et Lock) Type() action.Type {
	return action.ETH_LOCK
}

// Tags for ethereum lock
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
	tag3 := common.KVPair{
		Key:   []byte("tx.amount"),
		Value: []byte(strconv.FormatInt(et.LockAmount, 10)),
	}

	tags = append(tags, tag, tag2, tag3)
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

// Validate
func (ethLockTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {

	// unmarshal the tx message
	lock := &Lock{}
	err := lock.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic
	err = action.ValidateBasic(signedTx.RawBytes(), lock.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	// validate fee
	err = action.ValidateFee(ctx.FeeOpt, signedTx.Fee)
	if err != nil {
		return false, err
	}

	// Check lock fields for incoming trasaction
	ethTx := &types.Transaction{}
	err = rlp.DecodeBytes(lock.ETHTxn, ethTx)
	if err != nil {
		return false, errors.Wrap(err, "eth txn decode failed")
	}

	// check if lock amount in ethereum txn
	if ethTx.Value().Int64() != lock.LockAmount {
		return false, errors.New("incorrect lock amount in eth txn")
	}

	//TODO : Verify beninfiaciary address in ETHTX == locker (Phase 2)
	return true, nil
}

// ProcessCheck
func (e ethLockTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	lock := &Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	return e.processCommon(ctx, tx, lock)
}

// ProcessDeliver
func (e ethLockTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	lock := &Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	ok, resp := e.processCommon(ctx, tx, lock)
	if !ok {
		return ok, resp
	}

	//ethereum.NewTracker(lock.Locker,lock.ETHTxn,lock.TrackerName,)
	//todo: don't do job related work in delivery, just create tracker
	//if ctx.JobStore != nil {
	//
	//	job := event.JobETHBroadcast{lock.TrackerName}
	//	err = ctx.JobStore.SaveJob(job)
	//	if err != nil {
	//		return false, action.Response{Log: "job serialization failed err: " + err.Error()}
	//	}
	//}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}

// ProcessFee
func (ethLockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

// processCommon
func (ethLockTx) processCommon(ctx *action.Context, tx action.RawTx, lock *Lock) (bool, action.Response) {

	// parse eth transaction
	ethTx := &types.Transaction{}
	err := rlp.DecodeBytes(lock.ETHTxn, ethTx)
	if err != nil {
		return false, action.Response{Log: "err decoding txn: " + err.Error()}
	}

	// verify lock amount in the txn
	if ethTx.Value().Int64() != lock.LockAmount {
		return false, action.Response{Log: "incorrect lock amount in txn"}
	}

	val, err := ctx.Validators.GetValidatorsAddress()
	if err != nil {
		return false, action.Response{Log: "error in getting validator addresses" + err.Error()}
	}

	// Create ethereum tracker
	tracker := ethereum.NewTracker(
		lock.Locker, lock.ETHTxn,
		lock.TrackerName, val,
	)

	tracker.State = ethereum.New
	tracker.ProcessOwner = lock.Locker
	tracker.SignedETHTx = lock.ETHTxn

	// Save eth Tracker
	err = ctx.ETHTrackers.Set(*tracker)
	fmt.Println("Setting the tracker")
	if err != nil {
		return false, action.Response{Log: "error saving eth tracker: " + err.Error()}
	}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}
