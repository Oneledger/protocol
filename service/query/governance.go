package query

import (
	"errors"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	codes "github.com/Oneledger/protocol/status_codes"
)

// list single proposal by id
func (svc *Service) ListProposal(req client.ListProposalRequest, reply *client.ListProposalReply) error {
	proposalID := governance.ProposalID(req.ProposalId)
	proposal, _, err := svc.proposalMaster.Proposal.QueryAllStores(proposalID)
	if err != nil {
		svc.logger.Error("error getting proposal", err)
		return codes.ErrGetProposal
	}

	options := svc.proposalMaster.Proposal.GetOptionsByType(proposal.Type)
	funds := governance.GetCurrentFunds(proposalID, svc.proposalMaster.ProposalFund)
	stat, _ := svc.proposalMaster.ProposalVote.ResultSoFar(proposalID, options.PassPercentage)

	*reply = client.ListProposalReply{
		Proposal: *proposal,
		Fund:     *funds,
		Vote:     *stat,
		Height:   svc.proposalMaster.Proposal.GetState().Version(),
	}

	return nil
}

// list single proposal by id or list proposals
func (svc *Service) ListProposals(req client.ListProposalsRequest, reply *client.ListProposalsReply) error {
	// Validate parameters
	pState := governance.NewProposalState(req.State)
	pType := governance.NewProposalType(req.ProposalType)
	if req.State != "" {
		if pState == governance.ProposalStateInvalid {
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
	if pState == governance.ProposalStateInvalid &&
		pType == governance.ProposalTypeInvalid && len(req.Proposer) == 0 {
		return errors.New("invalid request parameters")
	}

	// Query in single store if specified
	pms := svc.proposalMaster
	var proposals []governance.Proposal
	if pState != governance.ProposalStateInvalid {
		proposals = pms.Proposal.FilterProposals(pState, req.Proposer, pType)
	} else { // Query in all stores otherwise
		active := pms.Proposal.FilterProposals(governance.ProposalStateActive, req.Proposer, pType)
		passed := pms.Proposal.FilterProposals(governance.ProposalStatePassed, req.Proposer, pType)
		failed := pms.Proposal.FilterProposals(governance.ProposalStateFailed, req.Proposer, pType)
		proposals = append(proposals, active...)
		proposals = append(proposals, passed...)
		proposals = append(proposals, failed...)
	}

	// Organize reply packet:
	// Proposals and its current funds, votes(if available) and state of each proposal
	proposalFunds := make([]balance.Amount, len(proposals))
	proposalVotes := make([]governance.VoteStatus, len(proposals))
	for i, prop := range proposals {
		options := pms.Proposal.GetOptionsByType(prop.Type)
		funds := governance.GetCurrentFunds(prop.ProposalID, pms.ProposalFund)
		stat, _ := pms.ProposalVote.ResultSoFar(prop.ProposalID, options.PassPercentage)
		proposalFunds[i] = *funds
		proposalVotes[i] = *stat
	}

	*reply = client.ListProposalsReply{
		Proposals:     proposals,
		ProposalFunds: proposalFunds,
		ProposalVotes: proposalVotes,
		Height:        pms.Proposal.GetState().Version(),
	}

	return nil
}
