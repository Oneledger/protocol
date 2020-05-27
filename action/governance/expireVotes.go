package governance

import (
	"encoding/json"
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

var _ action.Msg = &ExpireVotes{}

type ExpireVotes struct {
	ProposalID       governance.ProposalID
	ValidatorAddress action.Address
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

	err = action.ValidateFee(ctx.FeePool.GetOpt(), signedTx.Fee)
	if err != nil {
		return false, err
	}

	//Check if proposal id is valid
	if len(expireVotes.ProposalID) <= 0 {
		return false, errors.New("invalid proposal id")
	}

	//Check if validator address is valid
	if !ctx.Validators.Exists(e.ValidatorAddress) {
		return false, errors.New("signer is not a validator")
	}

	return true, nil
}

func (e ExpireVotes) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runExpireVotes(ctx, tx)
}

func (e ExpireVotes) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	return runExpireVotes(ctx, tx)
}

func (e ExpireVotes) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 1)
}

func runExpireVotes(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	expireVotes := ExpireVotes{}
	err := expireVotes.Unmarshal(tx.Data)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed_deserialize"),
		}
		return false, result
	}

	//Get proposal from active prefix
	activeProposals := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateActive)
	proposal, err := activeProposals.Get(expireVotes.ProposalID)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed"),
		}
		return false, result
	}

	//Update outcome and status of proposal
	proposal.Status = governance.ProposalStatusCompleted
	proposal.Outcome = governance.ProposalOutcomeInsufficientVotes

	//Set proposal in Failed prefix
	failedProposals := ctx.ProposalMasterStore.Proposal.WithPrefixType(governance.ProposalStateFailed)
	err = failedProposals.Set(proposal)
	if err != nil {
		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed"),
		}
		return false, result
	}

	//Delete proposal from active prefix
	deleted, err := activeProposals.Delete(proposal.ProposalID)
	if !deleted {
		//Need to delete proposal from failed prefix
		deleted, err = failedProposals.Delete(proposal.ProposalID)
		if !deleted {
			panic(errors.Wrap(err, "error deleting proposal from failed prefix."))
		}

		result := action.Response{
			Events: action.GetEvent(expireVotes.Tags(), "expire_votes_failed"),
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
