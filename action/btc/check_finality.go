/*

 */

package btc

import (
	"bytes"
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type CheckFinality struct {
	TrackerName  string
	OwnerAddress action.Address
}

var _ action.Msg = &CheckFinality{}

func (bcf *CheckFinality) Signers() []action.Address {
	return []action.Address{
		bcf.OwnerAddress,
	}
}

func (bcf *CheckFinality) Type() action.Type {
	return action.BTC_CHECK_FINALITY
}

func (bcf *CheckFinality) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.BTC_CHECK_FINALITY.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: bcf.OwnerAddress.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.domain_name"),
		Value: []byte(bcf.TrackerName),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (bcf *CheckFinality) Marshal() ([]byte, error) {
	return json.Marshal(bcf)
}

func (bcf *CheckFinality) Unmarshal(data []byte) error {
	return json.Unmarshal(data, bcf)
}

type btcCheckFinalityTx struct {
}

var _ action.Tx = btcCheckFinalityTx{}

func (btcCheckFinalityTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	f := CheckFinality{}
	err := f.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), f.Signers(), signedTx.Signatures)
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

	if tracker.State != bitcoin.BusyFinalizingTrackerState {
		return false, errors.New("tracker not available for broadcast")
	}

	return true, nil
}

func (btcCheckFinalityTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	f := CheckFinality{}
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

	if tracker.State != bitcoin.BusyFinalizingTrackerState {
		return false, action.Response{Log: "tracker not ready for finalizing"}
	}

	cd := bitcoin2.NewChainDriver("abcd")
	ok, err := cd.CheckFinality(*tracker.ProcessUTXO.TxID)
	if err != nil || !ok {
		return false, action.Response{Log: "tracker not finalized"}
	}

	// mint oBTC
	curr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	oBTCCoin := curr.NewCoinFromUnit(tracker.ProcessUTXO.Balance - tracker.CurrentUTXO.Balance)
	err = ctx.Balances.AddToAddress(f.OwnerAddress, oBTCCoin)
	if err != nil {
		return false, action.Response{Log: "error adding oBTC to address"}
	}

	// do final reset changes
	tracker.State = bitcoin.AvailableTrackerState
	tracker.CurrentUTXO = tracker.ProcessUTXO
	tracker.ProcessUTXO = nil
	tracker.ProcessOwner = nil
	tracker.Multisig = nil
	tracker.ProcessTx = nil

	validatorPubKeys, err := ctx.Validators.GetBitcoinKeys(&chaincfg.TestNet3Params)
	m := (len(validatorPubKeys) * 2 / 3) + 1

	lockScript, lockScriptAddress, err := bitcoin2.CreateMultiSigAddress(m, validatorPubKeys)

	tracker.NextLockScript = lockScript
	tracker.NextLockScriptAddress = lockScriptAddress

	err = ctx.Trackers.SetTracker(f.TrackerName, tracker)
	if err != nil || !ok {
		return false, action.Response{Log: "error resetting tracker, try again"}
	}

	return true, action.Response{
		Tags: f.Tags(),
	}
}

func (btcCheckFinalityTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	f := CheckFinality{}
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

	if tracker.State != bitcoin.BusyFinalizingTrackerState {
		return false, action.Response{Log: "tracker not ready for finalizing"}
	}

	// mint oBTC
	curr, ok := ctx.Currencies.GetCurrencyByName("BTC")
	oBTCCoin := curr.NewCoinFromUnit(tracker.ProcessUTXO.Balance - tracker.CurrentUTXO.Balance)
	err = ctx.Balances.AddToAddress(f.OwnerAddress, oBTCCoin)
	if err != nil {
		return false, action.Response{Log: "error adding oBTC to address"}
	}

	// do final reset changes
	tracker.State = bitcoin.AvailableTrackerState
	tracker.CurrentUTXO = tracker.ProcessUTXO
	tracker.ProcessUTXO = nil
	tracker.ProcessOwner = nil
	tracker.Multisig = nil
	tracker.ProcessTx = nil

	validatorPubKeys, err := ctx.Validators.GetBitcoinKeys(&chaincfg.TestNet3Params)
	m := (len(validatorPubKeys) * 2 / 3) + 1

	lockScript, lockScriptAddress, err := bitcoin2.CreateMultiSigAddress(m, validatorPubKeys)

	tracker.NextLockScript = lockScript
	tracker.NextLockScriptAddress = lockScriptAddress

	err = ctx.Trackers.SetTracker(f.TrackerName, tracker)
	if err != nil || !ok {
		return false, action.Response{Log: "error resetting tracker, try again"}
	}

	return true, action.Response{
		Tags: f.Tags(),
	}
}

func (btcCheckFinalityTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
