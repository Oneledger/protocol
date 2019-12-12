/*

 */

package btc

import (
	"bytes"
	"encoding/json"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
)

type ReportFinalityMint struct {
	TrackerName      string
	OwnerAddress     action.Address
	ValidatorAddress action.Address
	RandomBytes      []byte
}

var _ action.Msg = &ReportFinalityMint{}

func (m *ReportFinalityMint) Signers() []action.Address {
	return []action.Address{
		m.ValidatorAddress,
	}
}

func (m *ReportFinalityMint) Type() action.Type {
	return action.BTC_REPORT_FINALITY_MINT
}

func (m *ReportFinalityMint) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.BTC_REPORT_FINALITY_MINT.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: m.OwnerAddress.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(m.TrackerName),
	}
	tag4 := common.KVPair{
		Key:   []byte("tx.validator"),
		Value: m.ValidatorAddress.Bytes(),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
	return tags
}

func (m *ReportFinalityMint) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *ReportFinalityMint) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

type reportFinalityMintTx struct {
}

var _ action.Tx = reportFinalityMintTx{}

func (reportFinalityMintTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {

	f := ReportFinalityMint{}
	err := f.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), f.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	tracker, err := ctx.BTCTrackers.Get(f.TrackerName)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(tracker.ProcessOwner, f.OwnerAddress) {
		return false, errors.New("tracker process not owned by user")
	}

	if tracker.State != bitcoin.BusyFinalizing {
		return false, errors.New("tracker not available for finalizing")
	}

	return true, nil
}

func (reportFinalityMintTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runReportFinalityMint(ctx, tx)
}

func (reportFinalityMintTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runReportFinalityMint(ctx, tx)
}

func (reportFinalityMintTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	// return action.BasicFeeHandling(ctx, signedTx, start, size, 1)

	return true, action.Response{}
}

func runReportFinalityMint(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	f := ReportFinalityMint{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.BTCTrackers.Get(f.TrackerName)
	if err != nil {
		return false, action.Response{Log: "tracker not found" + f.TrackerName}
	}

	if !bytes.Equal(tracker.ProcessOwner, f.OwnerAddress) {
		return false, action.Response{Log: "tracker process not owned by user"}
	}

	if tracker.State != bitcoin.BusyFinalizing {
		return false, action.Response{Log: "tracker not ready for finalizing"}
	}

	valSet, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return false, action.Response{Log: "cannot get validator set"}
	}

	nValidators := len(valSet)
	votesThresholdForMint := (2 * nValidators) / 3

	isSenderAValidator := false
	for i := range valSet {
		if bytes.Equal(valSet[i].Address, f.ValidatorAddress) {
			isSenderAValidator = true
		}
	}

	if !isSenderAValidator {
		return false, action.Response{Log: "transaction sender not a validator"}
	}

	validatorSignedFlag := false
	for _, fv := range tracker.FinalityVotes {
		if bytes.Equal(fv, f.ValidatorAddress) {
			validatorSignedFlag = true
		}
	}

	if !validatorSignedFlag {
		tracker.FinalityVotes = append(tracker.FinalityVotes, f.ValidatorAddress)
	}

	// are there enough finality votes?
	if len(tracker.FinalityVotes) < votesThresholdForMint {

		// if not enough votes to mint end transaction processing

		err = ctx.BTCTrackers.SetTracker(f.TrackerName, tracker)
		if err != nil {
			return false, action.Response{Log: "tracker not ready for finalizing"}
		}

		return true, action.Response{
			Tags: f.Tags(),
		}
	}

	// if type is lock, then mint the oBTC
	if tracker.ProcessType == bitcoin.ProcessTypeLock {

		// mint oBTC
		curr, ok := ctx.Currencies.GetCurrencyByName("BTC")
		if !ok {

		}

		oBTCCoin := curr.NewCoinFromUnit(tracker.ProcessBalance - tracker.CurrentBalance)
		err = ctx.Balances.AddToAddress(f.OwnerAddress, oBTCCoin)
		if err != nil {
			ctx.Logger.Error(err)
			return false, action.Response{Log: "error adding oBTC to address"}
		}

		tally := keys.Address(lockBalanceAddress)
		err = ctx.Balances.AddToAddress(tally, oBTCCoin)
		if err != nil {
			ctx.Logger.Error(err)
			return false, action.Response{Log: "error adding oBTC to address"}
		}

		ctx.Logger.Info("btc coin minted to ", f.OwnerAddress)
	}

	// set the tracker to the new state

	validatorPubKeys, err := ctx.Validators.GetBitcoinKeys(&chaincfg.TestNet3Params)
	m := (len(validatorPubKeys) * 2 / 3) + 1

	lockScript, lockScriptAddress, addressList, err := bitcoin2.CreateMultiSigAddress(m, validatorPubKeys,
		f.RandomBytes, ctx.BTCChainType)

	// do final reset changes
	signers := make([]keys.Address, len(addressList))
	for i := range addressList {
		addr := base58.Decode(addressList[i])
		signers[i] = keys.Address(addr)
	}
	tracker.Multisig, err = keys.NewBTCMultiSig(nil, m, signers)

	tracker.State = bitcoin.Available

	tracker.CurrentTxId = tracker.ProcessTxId
	tracker.CurrentBalance = tracker.ProcessBalance
	tracker.CurrentLockScriptAddress = tracker.ProcessLockScriptAddress

	tracker.ProcessTxId = nil
	tracker.ProcessBalance = 0
	tracker.ProcessLockScriptAddress = lockScriptAddress
	tracker.ProcessUnsignedTx = nil
	tracker.ProcessOwner = nil
	tracker.FinalityVotes = nil
	tracker.ResetVotes = nil
	tracker.ProcessType = bitcoin.ProcessTypeNone

	// TODO check if node is validator
	if ctx.LockScriptStore != nil {
		err := ctx.LockScriptStore.SaveLockScript(lockScriptAddress, lockScript)
		if err != nil {
			return false, action.Response{Log: "error setting lockscript to store"}
		}
	}

	err = ctx.BTCTrackers.SetTracker(f.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: "error resetting tracker, try again"}
	}

	return true, action.Response{
		Tags: f.Tags(),
	}
}
