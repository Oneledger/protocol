package governance

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {
	serialize.RegisterConcrete(new(CreateProposal), "action_cp")
	serialize.RegisterConcrete(new(VoteProposal), "action_vp")
}

func EnableGovernance(r action.Router) error {
	err := r.AddHandler(action.PROPOSAL_CREATE, CreateProposal{})
	if err != nil {
		return errors.Wrap(err, "CreateProposalTx")
	}
	err = r.AddHandler(action.PROPOSAL_VOTE, voteProposalTx{})
	if err != nil {
		return errors.Wrap(err, "voteProposalTx")
	}
	err = r.AddHandler(action.PROPOSAL_FUND, fundProposalTx{})
	if err != nil {
		return errors.Wrap(err, "fundProposalTx")
	}
	err = r.AddHandler(action.PROPOSAL_WITHDRAW_FUNDS, WithdrawFunds{})
	if err != nil {
		return errors.Wrap(err, "WithdrawProposalTx")
	}
	err = r.AddHandler(action.EXPIRE_VOTES, ExpireVotes{})
	if err != nil {
		return errors.Wrap(err, "ExpireVotesTx")
	}
	err = r.AddHandler(action.PROPOSAL_FINALIZE, finalizeProposalTx{})
	if err != nil {
		return errors.Wrap(err, "finalizeProposalTx")
	}
	return nil
}

func EnableInternalGovernance(r action.Router) error {
	err := r.AddHandler(action.EXPIRE_VOTES, ExpireVotes{})
	if err != nil {
		return errors.Wrap(err, "ExpireVotesTx")
	}
	return nil
}
