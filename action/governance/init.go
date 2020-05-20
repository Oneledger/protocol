package governance

import (
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
)

func init() {
	serialize.RegisterConcrete(new(CreateProposal), "action_cp")
}

func EnableGovernance(r action.Router) error {
	err := r.AddHandler(action.PROPOSAL_CREATE, CreateProposal{})
	if err != nil {
		return errors.Wrap(err, "CreateProposalTx")
	}
	err = r.AddHandler(action.PROPOSAL_FUND, fundProposalTx{})
	if err != nil {
		return errors.Wrap(err, "fundProposalTx")
	}
	err = r.AddHandler(action.PROPOSAL_FUND, finalizeProposalTx{})
	if err != nil {
		return errors.Wrap(err, "finalizeProposalTx")
	}

	return nil
}
