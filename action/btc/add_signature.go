/*

 */

package btc

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/data/keys"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type AddSignature struct {
	TrackerName      string
	ValidatorPubKey  []byte
	BTCSignature     []byte
	ValidatorAddress action.Address
	Params           *chaincfg.Params
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

func (as *AddSignature) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(as.Type().String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.locker"),
		Value: []byte(as.ValidatorAddress.String()),
	}

	tags = append(tags, tag, tag2)
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

	return true, nil
}

func (ast btcAddSignatureTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

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

	if tracker.State != bitcoin.BusySigning {
		return false, action.Response{Log: fmt.Sprintf("tracker not accepting signatures: ", addSignature.TrackerName)}
	}

	addressPubkey, err := btcutil.NewAddressPubKey(addSignature.ValidatorPubKey, ctx.BTCChainType)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("%s, error generating btc public key", err.Error())}
	}

	err = tracker.AddSignature(addSignature.BTCSignature, keys.Address(addressPubkey.EncodeAddress()))
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error adding signature: %s, error: ", addSignature.TrackerName, err)}
	}

	ctx.Logger.Info("before has enough signatures", tracker.HasEnoughSignatures(), ctx.JobStore)
	//if tracker.HasEnoughSignatures() &&
	//	ctx.JobStore != nil {
	//
	//	ctx.Logger.Info("in has enough signatures")
	//
	//	tracker.State = bitcoin.BusyBroadcasting
	//
	//	id := strconv.Itoa(int(time.Now().UnixNano()))
	//	job := event.JobBTCBroadcast{
	//		event.JobTypeBTCBroadcast,
	//		tracker.Name,
	//		id,
	//		false,
	//		nil,
	//		false,
	//		0,
	//	}
	//
	////	err := ctx.JobStore.SaveJob(&job)
	//	ctx.Logger.Error("error while scheduling bitcoin broadcast job", err)
	//}

	err = ctx.BTCTrackers.SetTracker(addSignature.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error updating tracker store: %s, error: ", addSignature.TrackerName, err)}
	}

	return true, action.Response{
		Tags: addSignature.Tags(),
	}
}

func (ast btcAddSignatureTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {

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

	if tracker.State != bitcoin.BusySigning {
		return false, action.Response{Log: fmt.Sprintf("tracker not accepting signatures: %s", addSignature.TrackerName)}
	}

	addressPubKey, err := btcutil.NewAddressPubKey(addSignature.ValidatorPubKey, ctx.BTCChainType)
	if err != nil {
		return false, action.Response{Log: "error creating validator btc pubkey " + err.Error()}
	}

	err = tracker.AddSignature(addSignature.BTCSignature, addressPubKey.ScriptAddress())
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error adding signature: %s, error: ", addSignature.TrackerName, err)}
	}

	if tracker.HasEnoughSignatures() {
		tracker.State = bitcoin.BusyBroadcasting
	}

	err = ctx.BTCTrackers.SetTracker(addSignature.TrackerName, tracker)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error updating tracker store: %s, error: ", addSignature.TrackerName, err)}
	}

	return true, action.Response{
		Tags: addSignature.Tags(),
	}
}

func (ast btcAddSignatureTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return true, action.Response{}
}
