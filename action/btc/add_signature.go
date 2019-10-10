/*

 */

package btc

import (
	"encoding/json"
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type AddSignature struct {
	TrackerName      string
	ValidatorPubKey  *btcutil.AddressPubKey
	PubKey           keys.PublicKey
	BTCSignature     []byte
	ValidatorAddress action.Address
}

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

	_, err = action.ValidateBasic(signedTx.RawBytes(), addSignature.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	tracker, err := ctx.Trackers.Get(addSignature.TrackerName)
	if err != nil {
		return false, err
	}

	if tracker.State != bitcoin.BusySigningTrackerState {
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

	tracker, err := ctx.Trackers.Get(addSignature.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", addSignature.TrackerName)}
	}

	if tracker.State != bitcoin.BusySigningTrackerState {
		return false, action.Response{Log: fmt.Sprintf("tracker not accepting signatures: ", addSignature.TrackerName)}
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

	tracker, err := ctx.Trackers.Get(addSignature.TrackerName)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("tracker not found: %s", addSignature.TrackerName)}
	}

	if tracker.State != bitcoin.BusySigningTrackerState {
		return false, action.Response{Log: fmt.Sprintf("tracker not accepting signatures: %s", addSignature.TrackerName)}
	}

	err = tracker.AddSignature(addSignature.PubKey, addSignature.BTCSignature, addSignature.ValidatorPubKey)
	if err != nil {
		return false, action.Response{Log: fmt.Sprintf("error adding signature: %s, error: ", addSignature.TrackerName, err)}
	}

	return true, action.Response{
		Tags: addSignature.Tags(),
	}
}

func (ast btcAddSignatureTx) ProcessFee(ctx *action.Context, fee action.Fee) (bool, action.Response) {
	panic("implement me")
}
