package evidence

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/evidence"
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

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if err := r.ValidatorAddress.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (rtx releaseTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing 'release' transaction for ProcessCheck", tx)
	ok, result = runReleaseTransaction(ctx, tx)
	ctx.Logger.Detail("Result 'release' transaction for ProcessCheck", ok, result)
	return
}

func (rtx releaseTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing 'release' transaction for ProcessDeliver", tx)
	ok, result = runReleaseTransaction(ctx, tx)
	ctx.Logger.Detail("Result 'release' transaction for ProcessDeliver", ok, result)
	return
}

func (rtx releaseTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.Logger.Debug("Processing 'release' Transaction for ProcessFee", signedTx)
	r := &Release{}
	err := r.Unmarshal(signedTx.Data)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to unmarshal").Error()}
	}
	return action.StakingPayerFeeHandling(ctx, r.ValidatorAddress, signedTx, start, size, 1)
}

func runReleaseTransaction(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	r := &Release{}
	err := r.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, r.Tags(), err)
	}

	blockHeight := ctx.Header.GetHeight()
	blockCreatedAt := ctx.Header.GetTime()

	options, err := ctx.GovernanceStore.GetEvidenceOptions()
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, r.Tags(), err)
	}

	err = ctx.EvidenceStore.HandleRelease(options, r.ValidatorAddress, blockHeight, blockCreatedAt)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, evidence.ErrHandleReleaseFailed, r.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, r.Tags(), "release")
}
