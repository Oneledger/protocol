package query

import (
	"errors"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	codes "github.com/Oneledger/protocol/status_codes"
)

func translatePrefix(prefix string) governance.ProposalState {
	switch prefix {
	case "active":
		return governance.ProposalStateActive
	case "passed":
		return governance.ProposalStatePassed
	case "failed":
		return governance.ProposalStateFailed
	case "finalized":
		return governance.ProposalStateFinalized
	default:
		return governance.ProposalStateError
	}
}

// list single proposal by id or list proposals by state
func (svc *Service) ListProposals(req client.ListProposalsRequest, reply *client.ListProposalsReply) error {
	// List single proposal if ID is given
	if req.ProposalId != "" {
		return svc.ListProposal(req, reply)
	}

	// List proposals by given state
	proposalState := translatePrefix(req.State)
	if proposalState == governance.ProposalStateError {
		return errors.New("invalid proposal state")
	}
	proposalStore := svc.proposalMaster.Proposal.WithPrefixType(proposalState)
	proposals := make([]governance.Proposal, 0)
	proposalStore.Iterate(func(id governance.ProposalID, proposal *governance.Proposal) bool {
		proposals = append(proposals, *proposal)
		return false
	})

	// List current funds of each proposal
	proposalFunds := make([]balance.Amount, len(proposals))

	// List current vote status of each proposal if available
	proposalVotes := make([]governance.VoteStatus, len(proposals))
	for i, prop := range proposals {
		options := svc.proposalMaster.Proposal.GetOptionsByType(prop.Type)
		funds := governance.GetCurrentFunds(prop.ProposalID, svc.proposalMaster.ProposalFund)
		stat, _ := svc.proposalMaster.ProposalVote.ResultSoFar(prop.ProposalID, options.PassPercentage)
		proposalFunds[i] = *funds
		proposalVotes[i] = *stat
	}

	*reply = client.ListProposalsReply{
		Proposals:     proposals,
		ProposalFunds: proposalFunds,
		ProposalVotes: proposalVotes,
		State:         proposalState,
		Height:        svc.proposalMaster.Proposal.GetState().Version(),
	}

	return nil
}

func (svc *Service) ListProposal(req client.ListProposalsRequest, reply *client.ListProposalsReply) error {
	proposalID := governance.ProposalID(req.ProposalId)
	proposal, state, err := svc.proposalMaster.Proposal.QueryAllStores(proposalID)
	if err != nil {
		svc.logger.Error("error getting proposal", err)
		return codes.ErrGetProposal
	}

	options := svc.proposalMaster.Proposal.GetOptionsByType(proposal.Type)
	funds := governance.GetCurrentFunds(proposalID, svc.proposalMaster.ProposalFund)
	stat, _ := svc.proposalMaster.ProposalVote.ResultSoFar(proposalID, options.PassPercentage)

	*reply = client.ListProposalsReply{
		Proposals:     []governance.Proposal{*proposal},
		ProposalFunds: []balance.Amount{*funds},
		ProposalVotes: []governance.VoteStatus{*stat},
		State:         state,
		Height:        svc.proposalMaster.Proposal.GetState().Version(),
	}

	return nil
}
