package penalization

import (
	"encoding/binary"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &Allegation{}

type Allegation struct {
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

	tags = append(tags, tag, tag2, tag3, tag4, tag5)
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

	if err := r.MaliciousAddress.Err(); err != nil {
		return false, err
	}

	// TODO: Move to runAllegationTransaction
	if ctx.EvidenceStore.IsFrozenValidator(r.MaliciousAddress) {
		return false, action.ErrFrozenValidator
	}

	if !ctx.EvidenceStore.IsActiveValidator(r.ValidatorAddress) {
		return false, action.ErrNonActiveValidator
	}

	if r.ValidatorAddress.Equal(r.MaliciousAddress) {
		return false, action.ErrInvalidAddress
	}
	return true, nil
}

func (atx allegationTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing 'allegation' transaction for ProcessCheck", tx)
	ok, result = runAllegationTransaction(ctx, tx)
	ctx.Logger.Debug("Result 'allegation' transaction for ProcessCheck", ok, result)
	return
}

func (atx allegationTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing 'allegation' transaction for ProcessDeliver", tx)
	ok, result = runAllegationTransaction(ctx, tx)
	ctx.Logger.Debug("Result 'allegation' transaction for ProcessDeliver", ok, result)
	return
}

func (atx allegationTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.Logger.Debug("Processing 'allegation' Transaction for ProcessFee", signedTx)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runAllegationTransaction(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	// TODO: Check if validator staking address is matched with the requested
	al := &Allegation{}
	err := al.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.EvidenceStore.PerformAllegation(al.ValidatorAddress, al.MaliciousAddress, al.BlockHeight, al.ProofMsg)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	return true, action.Response{Events: action.GetEvent(al.Tags(), "allegation")}
}
