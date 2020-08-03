package penalization

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	gov "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &AllegationVote{}

type AllegationVote struct {
	RequestID int64
	Address   keys.Address
	Choice    int8
}

func (av AllegationVote) Marshal() ([]byte, error) {
	return json.Marshal(av)
}

func (av *AllegationVote) Unmarshal(data []byte) error {
	return json.Unmarshal(data, av)
}

func (av AllegationVote) Signers() []action.Address {
	return []action.Address{av.Address.Bytes()}
}

func (av AllegationVote) Type() action.Type {
	return action.ALLEGATION_VOTE
}

func (av AllegationVote) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(av.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.voter"),
		Value: av.Address.Bytes(),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.choice"),
		Value: []byte(strconv.Itoa(int(av.Choice))),
	}

	tags = append(tags, tag, tag2, tag3)
	return tags
}

var _ action.Tx = allegationVoteTx{}

type allegationVoteTx struct{}

func (atx allegationVoteTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	r := &AllegationVote{}
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

	if err := r.Address.Err(); err != nil {
		return false, err
	}

	if ctx.EvidenceStore.IsFrozenValidator(r.Address) {
		return false, action.ErrFrozenValidator
	}

	if !ctx.EvidenceStore.IsActiveValidator(r.Address) {
		return false, action.ErrNonActiveValidator
	}
	return true, nil
}

func (atx allegationVoteTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing 'allegation_vote' transaction for ProcessCheck", tx)
	ok, result = runAllegationVoteTransaction(ctx, tx)
	ctx.Logger.Debug("Result 'allegation_vote' transaction for ProcessCheck", ok, result)
	return
}

func (atx allegationVoteTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Debug("Processing 'allegation_vote' transaction for ProcessDeliver", tx)
	ok, result = runAllegationVoteTransaction(ctx, tx)
	ctx.Logger.Debug("Result 'allegation_vote' transaction for ProcessDeliver", ok, result)
	return
}

func (atx allegationVoteTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	ctx.Logger.Debug("Processing 'allegation_vote' Transaction for ProcessFee", signedTx)
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runAllegationVoteTransaction(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	al := &AllegationVote{}
	err := al.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	err = ctx.EvidenceStore.Vote(al.RequestID, al.Address, al.Choice)
	if err != nil {
		return false, action.Response{Log: err.Error()}
	}

	return true, action.Response{Events: action.GetEvent(al.Tags(), "allegation_vote")}
}
