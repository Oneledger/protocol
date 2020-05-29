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

func (svc *Service) GetProposals(req client.GetProposalsRequest, reply *client.GetProposalsResponse) error {
	proposalState := translatePrefix(req.Prefix)
	if proposalState == governance.ProposalStateError {
		return errors.New("invalid proposal state")
	}
	proposalStore := svc.proposalMaster.Proposal.WithPrefixType(proposalState)
	proposals := make([]governance.Proposal, 0)

	proposalStore.Iterate(func(id governance.ProposalID, proposal *governance.Proposal) bool {
		proposals = append(proposals, *proposal)
		return false
	})

	*reply = client.GetProposalsResponse{
		Proposals: proposals,
		Height:    svc.proposalMaster.Proposal.GetState().Version(),
	}

	return nil
}

func (svc *Service) GetProposalByID(req client.GetProposalByIDRequest, reply *client.GetProposalByIDReply) error {
	proposalID := governance.ProposalID(req.ProposalID)
	proposal, state, err := svc.proposalMaster.Proposal.QueryAllStores(proposalID)
	if err != nil {
		svc.logger.Error("error getting proposal", err)
		return codes.ErrGetProposal
	}

	proposalFunds := governance.GetCurrentFunds(proposalID, svc.proposalMaster.ProposalFund)

	*reply = client.GetProposalByIDReply{
		Proposal: *proposal,
		ProposalFunds: *proposalFunds,
		State:    state,
	}

	return nil
}
