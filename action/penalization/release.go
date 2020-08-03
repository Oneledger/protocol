package penalization

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &Release{}

type Release struct {
	ValidatorAddress keys.Address
}

func (r Release) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Release) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r Release) Signers() []action.Address {
	return []action.Address{r.ValidatorAddress.Bytes()}
}

func (r Release) Type() action.Type {
	return action.RELEASE
}

func (r Release) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(r.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: r.ValidatorAddress.Bytes(),
	}

	tags = append(tags, tag, tag2)
	return tags
}

var _ action.Tx = releaseTx{}

type releaseTx struct{}

func (rtx releaseTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	r := &Release{}
	err := r.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}
	err = action.ValidateBasic(tx.RawBytes(), r.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	feeOpt, err := ctx.GovernanceStore.GetFeeOption()
	if err != nil {
		return false, gov.ErrGetFeeOptions
	}
	err = action.ValidateFee(feeOpt, tx.Fee)
	if err != nil {
		return false, err
	}

	if err := r.ValidatorAddress.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (rtx releaseTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing 'release' transaction for ProcessCheck", tx)
	ok, result = runReleaseTransaction(ctx, tx)
	ctx.Logger.Debug("Result 'release' transaction for ProcessCheck", ok, result)
	return
}

func (rtx releaseTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing 'release' transaction for ProcessDeliver", tx)
	ok, result = runReleaseTransaction(ctx, tx)
	ctx.Logger.Debug("Result 'release' transaction for ProcessDeliver", ok, result)
	return
}

func (rtx releaseTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.Logger.Debug("Processing 'release' Transaction for ProcessFee", signedTx)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runReleaseTransaction(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	r := &Release{}
	err := r.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	blockHeight := ctx.Header.GetHeight()
	blockCreatedAt := ctx.Header.GetTime()

	options, err := ctx.GovernanceStore.GetEvidenceOptions()
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.EvidenceStore.HandleRelease(options, r.ValidatorAddress, blockHeight, blockCreatedAt)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	return true, action.Response{Events: action.GetEvent(r.Tags(), "release")}
}
