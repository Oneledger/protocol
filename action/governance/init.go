package governance

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
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

	return nil
}
