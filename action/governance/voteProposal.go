package governance

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	gov "github.com/Oneledger/protocol/data/governance"
)

var _ action.Msg = &VoteProposal{}

type VoteProposal struct {
	ProposalID       gov.ProposalID
	Address          action.Address
	ValidatorAddress action.Address
	Opinion          gov.VoteOpinion
}

var _ action.Tx = voteProposalTx{}

type voteProposalTx struct {
}

func (a voteProposalTx) Validate(ctx *action.Context, tx action.SignedTx) (bool, error) {
	ctx.Logger.Debug("Validate voteProposalTx transaction for CheckTx", tx)

	vote := &VoteProposal{}
	err := vote.Unmarshal(tx.Data)
	if err != nil {
		return false, errors.Wrap(action.ErrWrongTxType, err.Error())
	}

	// validate basic signature
	err = action.ValidateBasic(tx.RawBytes(), vote.Signers(), tx.Signatures)
	if err != nil {
		return false, err
	}

	err = action.ValidateFee(ctx.FeePool.GetOpt(), tx.Fee)
	if err != nil {
		return false, err
	}

	// validate params
	if len(vote.ProposalID) == 0 {
		return false, errors.New("empty proposalID")
	}
	if err = vote.Address.Err(); err != nil {
		return false, errors.Wrap(err, "invalid voter address")
	}
	if !ctx.Validators.IsValidatorAddress(vote.ValidatorAddress) {
		return false, errors.Wrap(err, "not a validator address")
	}
	if err = vote.Opinion.Err(); err != nil {
		return false, errors.Wrap(err, "invalid vote opinion")
	}

	return true, nil
}

func (a voteProposalTx) ProcessCheck(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessCheck voteProposalTx transaction for CheckTx", tx)
	return runVote(ctx, tx)
}

func (a voteProposalTx) ProcessFee(ctx *action.Context, signedTx action.SignedTx, start action.Gas, size action.Gas) (bool, action.Response) {
	return action.BasicFeeHandling(ctx, signedTx, start, size, 2)
}

func (a voteProposalTx) ProcessDeliver(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	ctx.Logger.Debug("ProcessDeliver voteProposalTx transaction for DeliverTx", tx)
	return runVote(ctx, tx)
}

func runVote(ctx *action.Context, tx action.RawTx) (bool, action.Response) {
	vote := &VoteProposal{}
	err := vote.Unmarshal(tx.Data)
	if err != nil {
		return false, action.Response{
			Log: action.ErrorMarshal(action.ErrProposalUnmarshal.Code, errors.Wrap(action.ErrProposalUnmarshal, err.Error()).Error()),
		}
	}

	// Get Proposal from proposal ACTIVE store
	pms := ctx.ProposalMasterStore
	proposal, err := pms.Proposal.WithPrefixType(gov.ProposalStateActive).Get(vote.ProposalID)
	if err != nil {
		return false, action.Response{
			Log: action.ErrorMarshal(action.ErrProposalExists.Code, errors.Wrap(action.ErrProposalExists, err.Error()).Error()),
		}
	}

	// Check if proposal is in VOTING status
	if proposal.Status != gov.ProposalStatusVoting {
		return false, action.Response{
			Log: action.ErrorMarshal(action.ErrNotInVoting.Code, action.ErrNotInVoting.Msg),
		}
	}

	// Check if proposal voting height is passed
	if ctx.Header.Height > proposal.VotingDeadline {
		return false, action.Response{
			Log: action.ErrorMarshal(action.ErrVotingHeightReached.Code, action.ErrVotingHeightReached.Msg),
		}
	}

	// Get validator's voting power
	validator, err := ctx.Validators.Get(vote.ValidatorAddress)
	if err != nil {
		return false, action.Response{
			Log: action.ErrorMarshal(action.ErrGettingValidatorList.Code, errors.Wrap(action.ErrGettingValidatorList, err.Error()).Error()),
		}
	}

	// Add this vote to proposal vote store
	pv := gov.NewProposalVote(vote.ValidatorAddress, vote.Opinion, validator.Power)
	err = ctx.ProposalMasterStore.ProposalVote.Update(vote.ProposalID, pv)
	if err != nil {
		return false, action.Response{
			Log: action.ErrorMarshal(action.ErrAddingVoteToVoteStore.Code, errors.Wrap(action.ErrAddingVoteToVoteStore, err.Error()).Error()),
		}
	}

	// Peek vote result based on collected votes so far
	options := pms.Proposal.GetOptionsByType(proposal.Type)
	result, err := pms.ProposalVote.ResultSoFar(vote.ProposalID, options.PassPercentage)
	if err != nil {
		return false, action.Response{
			Log: action.ErrorMarshal(action.ErrPeekingVoteResult.Code, errors.Wrap(action.ErrPeekingVoteResult, err.Error()).Error()),
		}
	}

	// Pass or fail this proposal if possible
	if result == gov.VOTE_RESULT_PASSED {
		proposal.Status = gov.ProposalStatusCompleted
		proposal.Outcome = gov.ProposalOutcomeCompleted
		err = pms.Proposal.WithPrefixType(gov.ProposalStatePassed).Set(proposal)
		if err != nil {
			return false, action.Response{
				Log: action.ErrorMarshal(action.ErrAddingProposalToDB.Code, errors.Wrap(action.ErrAddingProposalToDB, err.Error()).Error()),
			}
		}
	} else if result == gov.VOTE_RESULT_FAILED {
		proposal.Status = gov.ProposalStatusCompleted
		proposal.Outcome = gov.ProposalOutcomeInsufficientVotes
		err = pms.Proposal.WithPrefixType(gov.ProposalStateFailed).Set(proposal)
		if err != nil {
			return false, action.Response{
				Log: action.ErrorMarshal(action.ErrAddingProposalToDB.Code, errors.Wrap(action.ErrAddingProposalToDB, err.Error()).Error()),
			}
		}
	}

	// Delete proposal in ACTIVE store
	if result != gov.VOTE_RESULT_TBD {
		ok, err := pms.Proposal.WithPrefixType(gov.ProposalStateActive).Delete(vote.ProposalID)
		if err != nil || !ok {
			return false, action.Response{
				Log: action.ErrorMarshal(action.ErrDeletingProposalFromDB.Code, action.ErrDeletingProposalFromDB.Msg),
			}
		}
	}

	return true, action.Response{Events: action.GetEvent(vote.Tags(), "vote_proposal_success")}
}

func (vote VoteProposal) Signers() []action.Address {
	return []action.Address{vote.Address.Bytes(), vote.ValidatorAddress.Bytes()}
}

func (vote VoteProposal) Type() action.Type {
	return action.PROPOSAL_VOTE
}

func (vote VoteProposal) Tags() kv.Pairs {
	tags := make([]kv.Pair, 0)

	tag := kv.Pair{
		Key:   []byte("tx.type"),
		Value: []byte(vote.Type().String()),
	}
	tag2 := kv.Pair{
		Key:   []byte("tx.proposalID"),
		Value: []byte(vote.ProposalID),
	}
	tag3 := kv.Pair{
		Key:   []byte("tx.voter"),
		Value: vote.Address.Bytes(),
	}
	tag4 := kv.Pair{
		Key:   []byte("tx.address"),
		Value: vote.ValidatorAddress.Bytes(),
	}
	tag5 := kv.Pair{
		Key:   []byte("tx.opinion"),
		Value: []byte(string(vote.Opinion)),
	}

	tags = append(tags, tag, tag2, tag3, tag4, tag5)
	return tags
}

func (vote *VoteProposal) Marshal() ([]byte, error) {
	return json.Marshal(vote)
}

func (vote *VoteProposal) Unmarshal(data []byte) error {
	return json.Unmarshal(data, vote)
}
