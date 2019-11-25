package eth

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
)

type ExtMintOETH struct {
	TrackerName ethereum.TrackerName
	Locker      action.Address
	// Set locked amount from Tracker.signed TX ?
}

var _ action.Msg = &ExtMintOETH{}

func (eem ExtMintOETH) Signers() []action.Address {
	return []action.Address{
		eem.Locker,
	}
}

func (eem ExtMintOETH) Type() action.Type {
	return action.ETH_MINT
}

func (eem ExtMintOETH) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.ETH_REPORT_FINALITY_MINT.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: eem.Locker.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(eem.TrackerName.Hex()),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

func (eem ExtMintOETH) Marshal() ([]byte, error) {
	return json.Marshal(eem)
}

func (eem ExtMintOETH) Unmarshal(data []byte) error {
	return json.Unmarshal(data, eem)
}

type ethExtMintTx struct {
}

func (ethExtMintTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	//Implement check Finality first
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

	tracker, err := ctx.ETHTrackers.Get(f.TrackerName)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(tracker.ProcessOwner, f.Locker) {
		return false, errors.New("tracker process not owned by user")
	}

	if tracker.State != trackerlib.Finalized {
		return false, errors.New("Tracker for transaction not yet Finalized ")
	}
	return true, nil

}
func (ethExtMintTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	f := ExtMintOETH{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}
	tracker, err := ctx.ETHTrackers.Get(f.TrackerName)
	if err != nil {
		return false, action.Response{Log: "tracker not found" + f.TrackerName.Hex()}
	}
	if !bytes.Equal(tracker.ProcessOwner, f.Locker) {
		return false, action.Response{Log: "tracker process not owned by user"}
	}
	if !(tracker.State == trackerlib.Finalized) {
		return false, action.Response{Log: "Not enough votes collected for this tracker"}
	}
	// Mint OETH
	return true, action.Response{Log: "External Mint Successful"}
}

func (ethExtMintTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	f := ExtMintOETH{}
	err := f.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: "wrong tx type"}
	}
	tracker, err := ctx.ETHTrackers.Get(f.TrackerName)
	if err != nil {
		return false, action.Response{Log: "tracker not found" + f.TrackerName.Hex()}
	}
	if !bytes.Equal(tracker.ProcessOwner, f.Locker) {
		return false, action.Response{Log: "tracker process not owned by user"}
	}
	if !(tracker.State == trackerlib.Finalized) {
		return false, action.Response{Log: "Not enough votes collected for this tracker"}
	}
	// Mint OETH
	return true, action.Response{Log: "External Mint Successful"}
}

func (ethExtMintTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}
