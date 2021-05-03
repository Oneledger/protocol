package evidence

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/helpers"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/keys"
)

var _ action.Msg = &AllegationVote{}

type AllegationVote struct {
	RequestID string
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
	tag4 := kv.Pair{
		Key:   []byte("tx.requestID"),
		Value: []byte(av.RequestID),
	}

	tags = append(tags, tag, tag2, tag3, tag4)
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

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	if err := r.Address.Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (atx allegationVoteTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing 'allegation_vote' transaction for ProcessCheck", tx)
	ok, result = runAllegationVoteTransaction(ctx, tx)
	ctx.Logger.Detail("Result 'allegation_vote' transaction for ProcessCheck", ok, result)
	return
}

func (atx allegationVoteTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (ok bool, result action.Response) {
	ctx.Logger.Detail("Processing 'allegation_vote' transaction for ProcessDeliver", tx)
	ok, result = runAllegationVoteTransaction(ctx, tx)
	ctx.Logger.Detail("Result 'allegation_vote' transaction for ProcessDeliver", ok, result)
	return
}

func (atx allegationVoteTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.Logger.Detail("Processing 'allegation_vote' Transaction for ProcessFee", signedTx)
	r := &AllegationVote{}
	err := r.Unmarshal(signedTx.Data)
	if err != nil {
		return false, action.Response{Log: errors.Wrap(err, "failed to unmarshal").Error()}
	}
	return action.StakingPayerFeeHandling(ctx, r.Address, signedTx, start, size, 1)
}

func runAllegationVoteTransaction(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	al := &AllegationVote{}
	err := al.Unmarshal(tx.Data)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, al.Tags(), err)
	}

	if ctx.EvidenceStore.IsFrozenValidator(al.Address) {
		return helpers.LogAndReturnFalse(ctx.Logger, evidence.ErrFrozenValidator, al.Tags(), err)
	}

	if !ctx.EvidenceStore.IsActiveValidator(al.Address) {
		return helpers.LogAndReturnFalse(ctx.Logger, evidence.ErrNonActiveValidator, al.Tags(), err)
	}
	ctx.Logger.Detail("Vote for :", al.Choice, al.Address)
	err = ctx.EvidenceStore.Vote(al.RequestID, al.Address, al.Choice)
	if err != nil {
		return helpers.LogAndReturnFalse(ctx.Logger, action.ErrWrongTxType, al.Tags(), err)
	}

	return helpers.LogAndReturnTrue(ctx.Logger, al.Tags(), "allegation_vote")
}
