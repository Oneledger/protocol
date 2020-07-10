package query

import (
	"errors"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/governance"
	codes "github.com/Oneledger/protocol/status_codes"
)

// list single proposal by id
func (svc *Service) ListProposal(req client.ListProposalRequest, reply *client.ListProposalsReply) error {
	proposal, _, err := svc.proposalMaster.Proposal.QueryAllStores(req.ProposalId)
	if err != nil {
		svc.logger.Error("error getting proposal", err)
		return codes.ErrGetProposal
	}

	options := svc.proposalMaster.Proposal.GetOptionsByType(proposal.Type)
	funds := svc.proposalMaster.ProposalFund.GetCurrentFundsForProposal(req.ProposalId)
	stat, _ := svc.proposalMaster.ProposalVote.ResultSoFar(req.ProposalId, options.PassPercentage)

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

// list funds by funder for a proposal
func (svc *Service) GetFundsForProposalByFunder(req client.GetFundsForProposalByFunderRequest, reply *client.GetFundsForProposalByFunderReply) error {
	// Validate parameters
	err := req.Funder.Err()
	if err != nil {
		return errors.New("invalid funder address")
	}

	amount := svc.proposalMaster.ProposalFund.GetFundsForProposalByFunder(req.ProposalId, req.Funder)
	*reply = client.GetFundsForProposalByFunderReply{
		Amount: *amount,
	}

	return nil
}

func (svc *Service) GetGovernanceOptionsForHeight(req client.GovernanceOptionsRequest, reply *client.GovernanceOptionsReply) error {

	feeOpt, err := svc.governance.GetFeeOption()
	if err != nil {
		return err
	}
	propOpt, err := svc.governance.GetProposalOptions()
	if err != nil {
		return err
	}
	rewardOpt, err := svc.governance.GetRewardOptions()
	if err != nil {
		return err
	}
	ethOpt, err := svc.governance.GetETHChainDriverOption()
	if err != nil {
		return err
	}
	btcOpt, err := svc.governance.GetBTCChainDriverOption()
	if err != nil {
		return err
	}
	onsOpt, err := svc.governance.GetONSOptions()
	if err != nil {
		return err
	}
	*reply = client.GovernanceOptionsReply{
		GovOptions: governance.GovernanceState{
			FeeOption:     *feeOpt,
			ETHCDOption:   *ethOpt,
			BTCCDOption:   *btcOpt,
			ONSOptions:    *onsOpt,
			PropOptions:   *propOpt,
			RewardOptions: *rewardOpt,
		}}
	return nil
}

func (svc *Service) GetProposalOptions(_ client.ListTxTypesRequest, reply *client.GetProposalOptionsReply) error {

	options := *svc.proposalMaster.Proposal.GetOptions()
	height := svc.proposalMaster.Proposal.GetState().Version()

	*reply = client.GetProposalOptionsReply{
		ProposalOptions: options,
		Height:          height,
	}
	return nil
}
