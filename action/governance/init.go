package governance

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {
	serialize.RegisterConcrete(new(CreateProposal), "action_cp")
	serialize.RegisterConcrete(new(CancelProposal), "action_ccp")
	serialize.RegisterConcrete(new(VoteProposal), "action_vp")
	//todo maybe put this part into action init
	action.RegisterTxType(0x30, "PROPOSAL_CREATE")
	action.RegisterTxType(0x31, "PROPOSAL_CANCEL")
	action.RegisterTxType(0x32, "PROPOSAL_FUND")
	action.RegisterTxType(0x33, "PROPOSAL_VOTE")
	action.RegisterTxType(0x34, "PROPOSAL_FINALIZE")
	action.RegisterTxType(0x35, "EXPIRE_VOTES")
	action.RegisterTxType(0x36, "PROPOSAL_WITHDRAW_FUNDS")
}

func EnableGovernance(r action.Router) error {
	err := r.AddHandler(action.PROPOSAL_CREATE, CreateProposal{})
	if err != nil {
		return errors.Wrap(err, "CreateProposalTx")
	}
	err = r.AddHandler(action.PROPOSAL_CANCEL, cancelProposalTx{})
	if err != nil {
		return errors.Wrap(err, "cancelProposalTx")
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
	err = r.AddHandler(action.PROPOSAL_FINALIZE, FinalizeProposal{})
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
	err = r.AddHandler(action.PROPOSAL_FINALIZE, FinalizeProposal{})
	if err != nil {
		return errors.Wrap(err, "finalizeProposalTx")
	}
	return nil
}
