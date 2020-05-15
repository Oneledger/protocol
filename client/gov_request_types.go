package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

type CreateProposalRequest struct {
	ProposalType   governance.ProposalType `json:"proposal_type"`
	Description    string                  `json:"description"`
	Proposer       keys.Address            `json:"proposer"`
	InitialFunding action.Amount           `json:"initial_funding"`
	GasPrice       action.Amount           `json:"gasPrice"`
	Gas            int64                   `json:"gas"`
}
