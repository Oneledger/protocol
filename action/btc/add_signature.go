/*

 */

package btc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bitcoin"
)

// AddSignature is an internal transaction on the OneLedger Network. This transaction is used to add validator/witness
// signatures to the bitcoin lock or redeem transaction.
type AddSignature struct {
	TrackerName string

	ValidatorPubKey []byte
	BTCSignature    []byte

	ValidatorAddress action.Address
	Memo             string
}

var _ action.Msg = &AddSignature{}

func (as *AddSignature) Signers() []action.Address {
	return []action.Address{
		as.ValidatorAddress,
	}
}

func (AddSignature) Type() action.Type {
	return action.BTC_ADD_SIGNATURE
}

func (as *AddSignature) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(as.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.locker"),
		Value: []byte(as.ValidatorAddress.String()),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(as.TrackerName),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (as *AddSignature) Marshal() ([]byte, error) {
	return json.Marshal(as)
}

func (as *AddSignature) Unmarshal(data []byte) error {
	return json.Unmarshal(data, as)
}

type btcAddSignatureTx struct {
}

func (ast btcAddSignatureTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {

	addSignature := AddSignature{}

	err := addSignature.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), addSignature.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	if !ctx.Validators.IsValidatorAddress(addSignature.ValidatorAddress) {
		return false, errors.New("only validator can add a signature")
	}

	tracker, err := ctx.BTCTrackers.Get(addSignature.TrackerName)
	if err != nil {
		return false, err
	}

	if tracker.State != bitcoin.BusySigning {
		return false, errors.New("tracker not accepting signatures")
	}

	if tracker.HasEnoughSignatures() {
		return false, errors.New("tracker has sufficient signatures")
	}

	return true, nil
}

func (ast btcAddSignatureTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	return runAddSignature(ctx, tx)
}

func (ast btcAddSignatureTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runAddSignature(ctx, tx)
}

func (ast btcAddSignatureTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}

func runAddSignature(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

	addSignature := AddSignature{}
	err := addSignature.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}

	if !ctx.Validators.IsValidatorAddress(addSignature.ValidatorAddress) {
		return false, action.Response{Log: "signer not found in validator list"}
	}

	tracker, err := ctx.BTCTrackers.Get(addSignature.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", addSignature.TrackerName)}
	}

	if tracker.HasEnoughSignatures() {
		return false, action.Response{Log: "tracker has sufficient signatures"}
	}

	if tracker.State != bitcoin.BusySigning {
		return false, action.Response{Log: fmt.Sprintf("tracker not accepting signatures: %s", addSignature.TrackerName)}
	}

	opt := ctx.BTCTrackers.GetConfig()
	addressPubKey, err := btcutil.NewAddressPubKey(addSignature.ValidatorPubKey, opt.BTCParams)
	if err != nil {
		return false, action.Response{Log: "error creating validator btc pubkey " + err.Error()}
	}

	btcTx := wire.NewMsgTx(wire.TxVersion)
	buf := bytes.NewBuffer(tracker.ProcessUnsignedTx)
	err = btcTx.Deserialize(buf)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "error parsing btc txn").Error()}
	}

	// individual signature verification
	isFirstLock := tracker.CurrentTxId == nil
	if !isFirstLock {

		pk, err := btcec.ParsePubKey(addSignature.ValidatorPubKey, btcec.S256())
		sign, err := btcec.ParseSignature(addSignature.BTCSignature, btcec.S256())
		if err != nil {
			return false, action.Response{Log: errors.Wrap(err, "failed parse signature").Error()}
		}

		sc, err := ctx.LockScriptStore.Get(tracker.CurrentLockScriptAddress)
		if err != nil {
			return false, action.Response{Log: errors.Wrap(err, "cannot find lockscript").Error()}
		}

		hash, err := txscript.CalcSignatureHash(sc, txscript.SigHashAll, btcTx, 0)
		if err != nil {
			return false, action.Response{Log: errors.Wrap(err, "failed to calc signature hash").Error()}
		}

		ok := sign.Verify(hash, pk)
		if !ok {
			return false, action.Response{Log: "invalid validator signature"}
		}
	}

	err = tracker.AddSignature(addSignature.BTCSignature, addressPubKey.ScriptAddress())
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error adding signature: %s, error: %s", addSignature.TrackerName, err.Error())}
	}

	if tracker.HasEnoughSignatures() {
		tracker.State = bitcoin.BusyScheduleBroadcasting
	}

	err = ctx.BTCTrackers.SetTracker(addSignature.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error updating tracker store: %s, error: %s", addSignature.TrackerName, err.Error())}
	}

	return true, action.Response{
		Events: action.GetEvent(addSignature.Tags(), "btc_add_signature"),
	}

}
