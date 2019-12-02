package eth

import (
	"encoding/json"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	ethchaindriver "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/ethereum"
)

type Lock struct {
	Locker action.Address
	ETHTxn []byte
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

	cd, err := ethchaindriver.NewChainDriver(config.DefaultEthConfig(), ctx.Logger,ctx.ETHOptions)
	if err != nil {
		ctx.Logger.Error("err trying to get ChainDriver : ", err)
		return false,action.Response{
			Log:       "Unable to get Chain Driver",
		}
	}// parse eth transaction
	ethTx,err := cd.DecodeTransaction(lock.ETHTxn)
	if ethTx.To() != &ctx.ETHOptions.ContractAddress {
		return false,action.Response{
			Log:       "Invalid transaction ,To field of Transaction does not match Contract address",
		}
	}
	if err != nil {
		return false, action.Response{Log: "err decoding txn: " + err.Error()}
	}

	val, err := ctx.Validators.GetValidatorsAddress()
	if err != nil {
		return false, action.Response{Log: "error in getting validator addresses" + err.Error()}
	}

	// Create ethereum tracker
	tracker := ethereum.NewTracker(
		ethereum.ProcessTypeLock,
		lock.Locker,
		lock.ETHTxn,
		ethcommon.BytesToHash(lock.ETHTxn),
		val,
	)

	tracker.State = ethereum.New
	tracker.ProcessOwner = lock.Locker
	tracker.SignedETHTx = lock.ETHTxn

	// Save eth Tracker
	err = ctx.ETHTrackers.Set(tracker)
	fmt.Println("Setting the tracker")
	if err != nil {
		return false, action.Response{Log: "error saving eth tracker: " + err.Error()}
	}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}
