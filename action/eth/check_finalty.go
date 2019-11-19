package eth

import (
	"bytes"
	"encoding/json"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/chains/ethereum"
	trackerlib "github.com/Oneledger/protocol/data/ethereum"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/common"
)

type ReportFinalityMint struct {
	TrackerName      ethereum.TrackerName
	Locker           action.Address
	ValidatorAddress action.Address
}

var _ action.Msg = &ReportFinalityMint{}

func (m *ReportFinalityMint) Signers() []action.Address {
	return []action.Address{
		m.ValidatorAddress,
	}
}

func (m *ReportFinalityMint) Type() action.Type {
	return action.ETH_REPORT_FINALITY_MINT
}

func (m *ReportFinalityMint) Tags() common.KVPairs {
	tags := make([]common.KVPair, 0)

	tag := common.KVPair{
		Key:   []byte("tx.type"),
		Value: []byte(action.ETH_REPORT_FINALITY_MINT.String()),
	}
	tag2 := common.KVPair{
		Key:   []byte("tx.owner"),
		Value: m.Locker.Bytes(),
	}
	tag3 := common.KVPair{
		Key:   []byte("tx.tracker_name"),
		Value: []byte(m.TrackerName.Hex()),
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

func (r reportFinalityMintTx) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	f := ReportFinalityMint{}
	err := f.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	err = action.ValidateBasic(signedTx.RawBytes(), f.Signers(), signedTx.Signatures)
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
	if tracker.State != trackerlib.BusyBroadcasting {
		return false, errors.New("tracker not available for finalizing")
	}

	return true, nil
}

func (r reportFinalityMintTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	f := ReportFinalityMint{}
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
	valSet, err := ctx.Validators.GetValidatorSet()
	if err != nil {
		return false, action.Response{Log: "cannot get validator set"}
	}
	isSenderAValidator := false
	for i := range valSet {
		if bytes.Equal(valSet[i].Address, f.ValidatorAddress) {
			isSenderAValidator = true
		}
	}

	if !isSenderAValidator {
		return false, action.Response{Log: "transaction sender not a validator"}
	}

	for index, fv := range tracker.Validators {
		if bytes.Equal(fv, f.ValidatorAddress) {
			tracker.AddVote(fv, int64(index))
		}
	}
	// Are there Enough votes ?
	// If not LOG it and return
	// IF yes Check if status of minting
	// CHeck if it has been minted already, if not
	//    MINTING -> Mint OTETh and update status
	// If it has been minted
	//   MINTED -> Skip
	if !tracker.Finalized() {
		return false, action.Response{Log: "Not Enough votes to finalize a transaction"}
	}
	if !tracker.IsTaskCompleted() {
		ctx.Logger.Info("ready to mint")
		//Mint Tokens OETH
		tracker.CompleteTask()
	}

	return true, action.Response{}
}

func (r reportFinalityMintTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	panic("implement me")
}

func (r reportFinalityMintTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	panic("implement me")
}

var _ action.Tx = reportFinalityMintTx{}
