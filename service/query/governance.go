package query

import (
	"errors"

	"github.com/Oneledger/protocol/client"
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

	*reply = client.ListProposalsReply{
		Proposals: proposals,
		State:     proposalState,
		Height:    svc.proposalMaster.Proposal.GetState().Version(),
	}

	return nil
}

// list single proposal by id
func (svc *Service) ListProposal(req client.ListProposalsRequest, reply *client.ListProposalsReply) error {
	proposalID := governance.ProposalID(req.ProposalId)
	proposal, state, err := svc.proposalMaster.Proposal.QueryAllStores(proposalID)
	if err != nil {
		svc.logger.Error("error getting proposal", err)
		return codes.ErrGetProposal
	}

	*reply = client.ListProposalsReply{
		Proposals: []governance.Proposal{*proposal},
		State:     state,
		Height:    svc.proposalMaster.Proposal.GetState().Version(),
	}

	return nil
}
