//Package for transactions related to Etheruem
package eth

import (
	"bytes"
	"encoding/json"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	ethchaindriver "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/ethereum"
)

// Lock is a struct for one-Ledger transaction for Ether Lock
type Lock struct {
	Locker action.Address
	ETHTxn []byte
}

var _ action.Msg = &Lock{}

//Signers return the Address of the owner who created the transaction
func (et Lock) Signers() []action.Address {
	return []action.Address{et.Locker}
}

// Type returns the type of current action
func (et Lock) Type() action.Type {
	return action.ETH_LOCK
}

// Tags creates the tags to associate with the transaction
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
		Key:   []byte("tx.tracker"),
		Value: ethcommon.BytesToHash(et.ETHTxn).Bytes(),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

//Marshal Lock to byte array
func (et Lock) Marshal() ([]byte, error) {
	return json.Marshal(et)
}

func (et *Lock) Unmarshal(data []byte) error {
	return json.Unmarshal(data, et)
}

type ethLockTx struct {
}

var _ action.Tx = ethLockTx{}

// Validate provides basic validation for transaction Type and Fee
func (ethLockTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {

	// unmarshal the tx message
	lock := &Lock{}
	err := lock.Unmarshal(signedTx.Data)
	if err != nil {
		ctx.Logger.Error("error wrong tx type")
		return false, errors.Wrap(err, action.ErrWrongTxType.Msg)
	}

	// validate basic
	err = action.ValidateBasic(signedTx.RawBytes(), lock.Signers(), signedTx.Signatures)
	if err != nil {
		ctx.Logger.Error("validate basic failed", err)
		return false, err
	}

	// validate fee
	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		ctx.Logger.Error("validate fee failed", err)
		return false, err
	}

	// Check lock fields for incoming transaction

	//TODO : Verify beninfiaciary address in ETHTX == locker (Phase 2)
	return true, nil
}

// ProcessCheck runs checks on the transaction without commiting it .
func (e ethLockTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	lock := &Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		ctx.Logger.Error("wrong tx type", err)
		return false, action.Response{Log: "wrong tx type"}
	}

	return runLock(ctx, lock)
}

// ProcessDeliver run checks on transaction and commits it to a new block
func (e ethLockTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	lock := &Lock{}
	err := lock.Unmarshal(tx.Data)
	if err != nil {
		ctx.Logger.Error("wrong tx type", err)
		return false, action.Response{Log: "wrong tx type"}
	}

	ok, resp := runLock(ctx, lock)
	if !ok {
		return ok, resp
	}

	return true, action.Response{
		Tags: lock.Tags(),
	}
}

// ProcessFee process the transaction Fee in OLT
func (ethLockTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.State.ConsumeUpfront(237600)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

// runLock has the common functionality for ProcessCheck and ProcessDeliver
// Provides security checks for transaction
func runLock(ctx *action.Context, lock *Lock) (bool, action.Response) {

	ethTx, err := ethchaindriver.DecodeTransaction(lock.ETHTxn)
	if err != nil {
		ctx.Logger.Error("decode eth txn err", err)
		return false, action.Response{
			Log: "decode eth txn error" + err.Error(),
		}
	}

	cdOptions := ctx.ETHTrackers.GetOption()

	ok, err := ethchaindriver.VerifyLock(ethTx, cdOptions.ContractABI)
	if err != nil {
		ctx.Logger.Error("Unable to Verify Data for Ethereum Lock")
		return false, action.Response{
			Log: "Unable to verify lock trasaction" + err.Error(),
		}
	}
	if !ok {
		return false, action.Response{
			Log: "Bytes data does not match (function name field is different)",
		}
	}

	if !bytes.Equal(ethTx.To().Bytes(), cdOptions.ContractAddress.Bytes()) {
		ctx.Logger.Error("Contract address from eth Options : ", cdOptions.ContractAddress.Hex())
		ctx.Logger.Error("Contract address from eth TX : ", ethTx.To().Hex())
		ctx.Logger.Error("to field does not match contract address")
		return false, action.Response{
			Log: "Contract address does not match",
		}
	}

	val, err := ctx.Validators.GetValidatorsAddress()
	if err != nil {

		ctx.Logger.Error("err in getting validator address", err)
		return false, action.Response{Log: "error in getting validator addresses" + err.Error()}
	}

	curr, ok := ctx.Currencies.GetCurrencyByName("ETH")
	if !ok {
		return false, action.Response{Log: fmt.Sprintf("ETH currency not available", lock.Locker)}
	}
	lockCoin := curr.NewCoinFromString(ethTx.Value().String())
	// Adding lock amount to common address to maintain count of total oEth minted
	ethSupply := action.Address(cdOptions.TotalSupplyAddr)
	balCoin, err := ctx.Balances.GetBalanceForCurr(ethSupply, &curr)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("Unable to get Eth lock total balance", lock.Locker)}
	}

	totalSupplyCoin := curr.NewCoinFromString(cdOptions.TotalSupply)

	if !balCoin.Plus(lockCoin).LessThanEqualCoin(totalSupplyCoin) {
		return false, action.Response{Log: fmt.Sprintf("Eth lock exceeded limit", lock.Locker)}
	}
	name := ethcommon.BytesToHash(lock.ETHTxn)
	if ctx.ETHTrackers.WithPrefixType(ethereum.PrefixOngoing).Exists(name) || ctx.ETHTrackers.WithPrefixType(ethereum.PrefixPassed).Exists(name) {
		return false, action.Response{
			Log: "Tracker already exists / Lock for this ETHTX in progress or has completed successfully",
		}
	}
	if ctx.ETHTrackers.WithPrefixType(ethereum.PrefixFailed).Exists(name) {
		ctx.ETHTrackers.WithPrefixType(ethereum.PrefixFailed).Delete(name)
	}
	// Create ethereum tracker
	tracker := ethereum.NewTracker(
		ethereum.ProcessTypeLock,
		lock.Locker,
		lock.ETHTxn,
		name,
		val,
	)

	tracker.State = ethereum.New
	tracker.ProcessOwner = lock.Locker
	tracker.SignedETHTx = lock.ETHTxn
	// Save eth Tracker
	err = ctx.ETHTrackers.WithPrefixType(ethereum.PrefixOngoing).Set(tracker)
	if err != nil {
		ctx.Logger.Error("error saving eth tracker", err)
		return false, action.Response{Log: "error saving eth tracker: " + err.Error()}
	}
	return true, action.Response{
		Tags: lock.Tags(),
	}
}
