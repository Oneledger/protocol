/*

 */

package btc

import (
	"bytes"
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type ExtMintOBTC struct {
	TrackerName  string
	OwnerAddress action.Address
	RandomBytes  []byte
}

var _ action.Msg = &ExtMintOBTC{}

func (em *ExtMintOBTC) Signers() []action.Address {
	return []action.Address{
		em.OwnerAddress,
	}
}

func (em *ExtMintOBTC) Type() action.Type {
	return action.BTC_EXT_MINT
}

func (em *ExtMintOBTC) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.BTC_REPORT_FINALITY_MINT.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: em.OwnerAddress.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(em.TrackerName),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (em *ExtMintOBTC) Marshal() ([]byte, error) {
	return json.Marshal(em)
}

func (em *ExtMintOBTC) Unmarshal(data []byte) error {
	return json.Unmarshal(data, em)
}

type extMintOBTCTx struct {
}

func (extMintOBTCTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	f := ReportFinalityMint{}
	err := f.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), f.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeeOpt, signedTx.Fee)
	if err != nil {
		return false, err
	}

	tracker, err := ctx.Trackers.Get(f.TrackerName)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(tracker.ProcessOwner, f.OwnerAddress) {
		return false, errors.New("tracker process not owned by user")
	}

	if tracker.State != bitcoin.BusyFinalizing {
		return false, errors.New("tracker not available for broadcast")
	}

	return true, nil
}

func (extMintOBTCTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	f := ExtMintOBTC{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.Trackers.Get(f.TrackerName)
	if err != nil {
		return false, action.Response{Log: "tracker not found" + f.TrackerName}
	}

	if !bytes.Equal(tracker.ProcessOwner, f.OwnerAddress) {
		return false, action.Response{Log: "tracker process not owned by user"}
	}

	if tracker.State != bitcoin.BusyFinalizing {
		return false, action.Response{Log: "tracker not ready for finalizing"}
	}

	// mint oBTC
	curr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	oBTCCoin := curr.NewCoinFromUnit(tracker.ProcessBalance - tracker.CurrentBalance)
	err = ctx.Balances.AddToAddress(f.OwnerAddress, oBTCCoin)
	if err != nil {
		return false, action.Response{Log: "error adding oBTC to address"}
	}

	validatorPubKeys, err := ctx.Validators.GetBitcoinKeys(&chaincfg.TestNet3Params)
	m := (len(validatorPubKeys) * 2 / 3) + 1

	_, lockScriptAddress, addressList, err := bitcoin2.CreateMultiSigAddress(m, validatorPubKeys, f.RandomBytes)

	// do final reset changes
	signers := make([]keys.Address, len(addressList))
	for i := range addressList {
		signers[i] = keys.Address(addressList[i])
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

	err = ctx.Trackers.SetTracker(f.TrackerName, tracker)
	if err != nil || !ok {
		return false, action.Response{Log: "error resetting tracker, try again"}
	}

	return true, action.Response{
		Tags: f.Tags(),
	}
}

func (extMintOBTCTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	f := ExtMintOBTC{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	tracker, err := ctx.Trackers.Get(f.TrackerName)
	if err != nil {
		return false, action.Response{Log: "tracker not found" + f.TrackerName}
	}

	if !bytes.Equal(tracker.ProcessOwner, f.OwnerAddress) {
		return false, action.Response{Log: "tracker process not owned by user"}
	}

	if tracker.State != bitcoin.BusyFinalizing {
		return false, action.Response{Log: "tracker not ready for finalizing"}
	}

	// mint oBTC
	curr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	oBTCCoin := curr.NewCoinFromUnit(tracker.ProcessBalance - tracker.CurrentBalance)
	err = ctx.Balances.AddToAddress(f.OwnerAddress, oBTCCoin)
	if err != nil {
		return false, action.Response{Log: "error adding oBTC to address"}
	}

	validatorPubKeys, err := ctx.Validators.GetBitcoinKeys(&chaincfg.TestNet3Params)
	m := (len(validatorPubKeys) * 2 / 3) + 1

	lockScript, lockScriptAddress, addressList, err := bitcoin2.CreateMultiSigAddress(m, validatorPubKeys, f.RandomBytes)

	// do final reset changes
	signers := make([]keys.Address, len(addressList))
	for i := range addressList {
		signers[i] = keys.Address(addressList[i])
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

	if ctx.LockScriptStore != nil {
		err := ctx.LockScriptStore.SaveLockScript(lockScriptAddress, lockScript)
		if err != nil {
			return false, action.Response{Log: "error setting lockscript to store"}
		}
	}

	err = ctx.Trackers.SetTracker(f.TrackerName, tracker)
	if err != nil || !ok {
		return false, action.Response{Log: "error resetting tracker, try again"}
	}

	return true, action.Response{
		Tags: f.Tags(),
	}
}

func (extMintOBTCTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
