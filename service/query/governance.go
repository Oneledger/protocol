package query

import (
	"errors"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/governance"
	codes "github.com/Oneledger/protocol/status_codes"
)

// list single proposal by id
func (svc *Service) ListProposal(req client.ListProposalRequest, reply *client.ListProposalsReply) error {
	proposalID := governance.ProposalID(req.ProposalId)
	proposal, _, err := svc.proposalMaster.Proposal.QueryAllStores(proposalID)
	if err != nil {
		svc.logger.Error("error getting proposal", err)
		return codes.ErrGetProposal
	}

	options := svc.proposalMaster.Proposal.GetOptionsByType(proposal.Type)
	funds := svc.proposalMaster.ProposalFund.GetCurrentFundsForProposal(proposalID)
	stat, _ := svc.proposalMaster.ProposalVote.ResultSoFar(proposalID, options.PassPercentage)

	ps := client.ProposalStat{
		Proposal: *proposal,
		Funds:    *funds,
		Votes:    *stat,
	}
	*reply = client.ListProposalsReply{
		ProposalStats: []client.ProposalStat{ps},
		Height:        svc.proposalMaster.Proposal.GetState().Version(),
	}

	return nil
}

// list single proposal by id or list proposals
func (svc *Service) ListProposals(req client.ListProposalsRequest, reply *client.ListProposalsReply) error {
	// Validate parameters
	if len(req.Proposer) != 0 {
		err := req.Proposer.Err()
		if err != nil {
			return errors.New("invalid proposer address")
		}
	}

	// Query in single store if specified
	pms := svc.proposalMaster
	var proposals []governance.Proposal
	if req.State != governance.ProposalStateInvalid {
		proposals = pms.Proposal.FilterProposals(req.State, req.Proposer, req.ProposalType)
	} else { // Query in all stores otherwise
		active := pms.Proposal.FilterProposals(governance.ProposalStateActive, req.Proposer, req.ProposalType)
		passed := pms.Proposal.FilterProposals(governance.ProposalStatePassed, req.Proposer, req.ProposalType)
		failed := pms.Proposal.FilterProposals(governance.ProposalStateFailed, req.Proposer, req.ProposalType)
		finalized := pms.Proposal.FilterProposals(governance.ProposalStateFinalized, req.Proposer, req.ProposalType)
		finalizeFailed := pms.Proposal.FilterProposals(governance.ProposalStateFinalizeFailed, req.Proposer, req.ProposalType)
		proposals = append(proposals, active...)
		proposals = append(proposals, passed...)
		proposals = append(proposals, failed...)
		proposals = append(proposals, finalized...)
		proposals = append(proposals, finalizeFailed...)
	}

	// Organize reply packet:
	// Proposals and its current funds, votes(if available)
	proposalStats := make([]client.ProposalStat, len(proposals))
	for i, prop := range proposals {
		options := pms.Proposal.GetOptionsByType(prop.Type)
		funds := pms.ProposalFund.GetCurrentFundsForProposal(prop.ProposalID)
		stat, _ := pms.ProposalVote.ResultSoFar(prop.ProposalID, options.PassPercentage)
		ps := client.ProposalStat{
			Proposal: prop,
			Funds:    *funds,
			Votes:    *stat,
		}
		proposalStats[i] = ps
	}

	*reply = client.ListProposalsReply{
		ProposalStats: proposalStats,
		Height:        pms.Proposal.GetState().Version(),
	}

	return nil
}
