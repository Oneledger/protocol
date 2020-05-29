package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

type CreateProposalRequest struct {
	ProposalID     string        `json:"proposal_id"`
	ProposalType   string        `json:"proposal_type"`
	Description    string        `json:"description"`
	Proposer       keys.Address  `json:"proposer"`
	InitialFunding action.Amount `json:"initial_funding"`
	GasPrice       action.Amount `json:"gasPrice"`
	Gas            int64         `json:"gas"`
}

type GetProposalsRequest struct {
	Prefix   string        `json:"prefix"`
	GasPrice action.Amount `json:"gasPrice"`
	Gas      int64         `json:"gas"`
}

type GetProposalsResponse struct {
	Proposals []governance.Proposal `json:"proposals"`
	Height    int64                 `json:"height"`
}

type GetProposalByIDRequest struct {
	ProposalID string `json:"proposal_id"`
}

type GetProposalByIDReply struct {
	Proposal governance.Proposal      `json:"proposal"`
	ProposalFunds balance.Amount	  `json:"proposal_funds"`
	State    governance.ProposalState `json:"state"`
}

type CreateVoteRequest struct {
	ProposalID string        `json:"proposal_id"`
	Opinion    string        `json:"opinion"`
	Address    keys.Address  `json:"address"`
	GasPrice   action.Amount `json:"gasPrice"`
	Gas        int64         `json:"gas"`
}

type CreateVoteReply struct {
	RawTx     []byte           `json:"rawTx"`
	Signature action.Signature `json:"signature"`
}

type FundProposalRequest struct {
	ProposalId    governance.ProposalID `json:"proposal_id"`
	FundValue     action.Amount         `json:"fund_value"`
	FunderAddress action.Address        `json:"funder_address"`
	GasPrice      action.Amount         `json:"gasPrice"`
	Gas           int64                 `json:"gas"`
}

type WithdrawFundsRequest struct {
	ProposalID    governance.ProposalID `json:"proposal_id"`
	Contributor   keys.Address			`json:"contributor_address"`
	WithdrawValue action.Amount			`json:"withdraw_value"`
	Beneficiary   keys.Address			`json:"beneficiary_address"`
	GasPrice      action.Amount         `json:"gasPrice"`
	Gas           int64                 `json:"gas"`
}
