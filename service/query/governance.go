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
	funds := governance.GetCurrentFunds(proposalID, svc.proposalMaster.ProposalFund)
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
	pState := governance.NewProposalState(req.State)
	pType := governance.NewProposalType(req.ProposalType)
	if req.State != "" {
		if pState == governance.ProposalStateError {
			return errors.New("invalid proposal state")
		}
	}
	if req.ProposalType != "" {
		if pType == governance.ProposalTypeInvalid {
			return errors.New("invalid proposal type")
		}
	}
	if len(req.Proposer) != 0 {
		err := req.Proposer.Err()
		if err != nil {
			return errors.New("invalid proposer address")
		}
	}
	if pState == governance.ProposalStateError &&
		pType == governance.ProposalTypeInvalid && len(req.Proposer) == 0 {
		return errors.New("invalid request parameters")
	}

	// Query in single store if specified
	pms := svc.proposalMaster
	var proposals []governance.Proposal
	if pState != governance.ProposalStateError {
		proposals = pms.Proposal.FilterProposals(pState, req.Proposer, pType)
	} else { // Query in all stores otherwise
		active := pms.Proposal.FilterProposals(governance.ProposalStateActive, req.Proposer, pType)
		passed := pms.Proposal.FilterProposals(governance.ProposalStatePassed, req.Proposer, pType)
		failed := pms.Proposal.FilterProposals(governance.ProposalStateFailed, req.Proposer, pType)
		finalized := pms.Proposal.FilterProposals(governance.ProposalStateFinalized, req.Proposer, pType)
		finalizeFailed := pms.Proposal.FilterProposals(governance.ProposalStateFinalizeFailed, req.Proposer, pType)
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
		funds := governance.GetCurrentFunds(prop.ProposalID, pms.ProposalFund)
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
