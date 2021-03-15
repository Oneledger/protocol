package governance

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
)

var _ action.Msg = &ExpireVotes{}

type ExpireVotes struct {
	ProposalID       governance.ProposalID `json:"proposalId"`
	ValidatorAddress action.Address        `json:"validatorAddress"`
}

func (e ExpireVotes) Validate(ctx *action.Context, signedTx action.SignedTx) (bool, error) {
	expireVotes := ExpireVotes{}
	err := expireVotes.Unmarshal(signedTx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	//validate basic signature
	err = action.ValidateBasic(signedTx.RawBytes(), expireVotes.Signers(), signedTx.Signatures)
	if err != nil {
		return false, err
	}

	//Check if Proposal ID is valid
	if err = expireVotes.ProposalID.Err(); err != nil {
		return false, governance.ErrInvalidProposalId
	}

	return true, nil
}

func (e ExpireVotes) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runExpireVotes(ctx, tx)
}

func (e ExpireVotes) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Detail("Processing ExpireVotes Transaction for DeliverTx", tx)
	return runExpireVotes(ctx, tx)
}

func (e ExpireVotes) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas, gasUsed action.Gas) (bool, action.Response) {
	ctx.State.ConsumeVerifySigGas(1)
	ctx.State.ConsumeStorageGas(size)

	// check the used gas for the tx
	final := ctx.Balances.State.ConsumedGas()
	used := int64(final - start)
	ctx.Logger.Detail("Gas Used : ", used)

	return true, action.Response{}
}

func runExpireVotes(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	active := governance.ProposalStateActive
	failed := governance.ProposalStateFailed

	expireVotes := ExpireVotes{}
	err := expireVotes.Unmarshal(tx.Data)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed_deserialize"),
			Log:    action.ErrWrongTxType.Wrap(err).Marshal(),
		}
		return false, result
	}

	//Get proposal from active prefix
	proposal, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(active).Get(expireVotes.ProposalID)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed"),
			Log:    governance.ErrProposalNotExists.Wrap(err).Marshal(),
		}
		return false, result
	}

	//Update outcome and status of proposal
	proposal.Status = governance.ProposalStatusCompleted
	proposal.Outcome = governance.ProposalOutcomeInsufficientVotes

	//Set proposal in Failed prefix
	err = ctx.ProposalMasterStore.Proposal.WithPrefixType(failed).Set(proposal)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed"),
			Log:    governance.ErrAddingProposalToFailedStore.Wrap(err).Marshal(),
		}
		return false, result
	}

	//Delete proposal from active prefix
	deleted, err := ctx.ProposalMasterStore.Proposal.WithPrefixType(active).Delete(proposal.ProposalID)
	if !deleted {
		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed"),
			Log:    governance.ErrDeletingProposalFromFailedStore.Marshal(),
		}
		return false, result
	}

	result := action.Response{
		Events: action.GetEvent(expireVotes.Tags(), "expire_votes_success"),
	}
	return true, result
}

func (e ExpireVotes) Signers() []action.Address {
	return []action.Address{e.ValidatorAddress}
}

func (e ExpireVotes) Type() action.Type {
	return action.EXPIRE_VOTES
}

func (e ExpireVotes) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(action.EXPIRE_VOTES.String()),
	}

	tag1 := kv.Pair{
		Key:   []byte("tx.proposal_id"),
		Value: []byte(e.ProposalID),
	}

	tag2 := kv.Pair{
		Key:   []byte("tx.validator"),
		Value: []byte(e.ValidatorAddress.String()),
	}

	tags = append(tags, tag, tag1, tag2)
	return tags
}

func (e ExpireVotes) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ExpireVotes) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, e)
}
