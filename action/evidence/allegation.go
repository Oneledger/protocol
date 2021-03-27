package evidence

import (
	"encoding/binary"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &Allegation{}

type Allegation struct {
	RequestID        string
	ValidatorAddress keys.Address
	MaliciousAddress keys.Address
	BlockHeight      int64
	ProofMsg         string
}

func (r Allegation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Allegation) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r Allegation) Signers() []action.Address {
	return []action.Address{r.ValidatorAddress.Bytes()}
}

func (r Allegation) Type() action.Type {
	return action.ALLEGATION
}

func (r Allegation) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(r.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: r.ValidatorAddress.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.malicious"),
		Value: r.MaliciousAddress.Bytes(),
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(r.BlockHeight))
	tag4 := kv.Pair{
		Key:   []byte("tx.height"),
		Value: b,
	}
	tag5 := kv.Pair{
		Key:   []byte("tx.proof"),
		Value: []byte(r.ProofMsg),
	}
	tag6 := kv.Pair{
		Key:   []byte("tx.requestID"),
		Value: []byte(r.RequestID),
	}

	tags = append(tags, tag, tag2, tag3, tag4, tag5, tag6)
	return tags
}

var _ action.Tx = allegationTx{}

type allegationTx struct{}

func (atx allegationTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	r := &Allegation{}
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

	if err := r.MaliciousAddress.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (atx allegationTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing 'allegation' transaction for ProcessCheck", tx)
	ok, result = runAllegationTransaction(ctx, tx)
	ctx.Logger.Detail("Result 'allegation' transaction for ProcessCheck", ok, result)
	return
}

func (atx allegationTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing 'allegation' transaction for ProcessDeliver", tx)
	ok, result = runAllegationTransaction(ctx, tx)
	ctx.Logger.Detail("Result 'allegation' transaction for ProcessDeliver", ok, result)
	return
}

func (atx allegationTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.Logger.Detail("Processing 'allegation' Transaction for ProcessFee", signedTx)
	r := &Allegation{}
	err := r.Unmarshal(signedTx.Data)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to unmarshal").Error()}
	}
	return action.StakingPayerFeeHandling(ctx, r.ValidatorAddress, signedTx, start, size, 1)
}

func runAllegationTransaction(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	al := &Allegation{}
	err := al.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, al.Tags(), err)
	}

	if al.BlockHeight > ctx.Header.Height {
		return helpers.LogAndReturnFalse(ctx.Logger, evidence.ErrInvalidHeight, al.Tags(), err)
	}

	if ctx.EvidenceStore.IsFrozenValidator(al.MaliciousAddress) {
		return helpers.LogAndReturnFalse(ctx.Logger, evidence.ErrFrozenValidator, al.Tags(), err)
	}

	if !ctx.EvidenceStore.IsActiveValidator(al.ValidatorAddress) {
		return helpers.LogAndReturnFalse(ctx.Logger, evidence.ErrNonActiveValidator, al.Tags(), err)
	}

	if al.ValidatorAddress.Equal(al.MaliciousAddress) {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrInvalidAddress, al.Tags(), err)
	}
	ctx.Logger.Detail("Performing allegation : ", al.ValidatorAddress, " | on :", al.MaliciousAddress)
	err = ctx.EvidenceStore.PerformAllegation(al.ValidatorAddress, al.MaliciousAddress, al.RequestID, al.BlockHeight, al.ProofMsg)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, evidence.ErrCreateAllegationFailed, al.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, al.Tags(), "allegation")
}
