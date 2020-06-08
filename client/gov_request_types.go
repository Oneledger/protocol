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

type ListProposalRequest struct {
	ProposalId governance.ProposalID `json:"proposal_id"`
}

type ListProposalReply struct {
	Proposal governance.Proposal   `json:"proposal"`
	Fund     balance.Amount        `json:"fund"`
	Vote     governance.VoteStatus `json:"vote"`
	Height   int64                 `json:"height"`
}

type ListProposalsRequest struct {
	State        string       `json:"state"`
	Proposer     keys.Address `json:"proposer"`
	ProposalType string       `json:"proposal_type"`
}

type ListProposalsReply struct {
	Proposals     []governance.Proposal   `json:"proposals"`
	ProposalFunds []balance.Amount        `json:"proposal_funds"`
	ProposalVotes []governance.VoteStatus `json:"proposal_votes"`
	Height        int64                   `json:"height"`
}

type VoteProposalRequest struct {
	ProposalId string        `json:"proposal_id"`
	Opinion    string        `json:"opinion"`
	Address    keys.Address  `json:"address"`
	GasPrice   action.Amount `json:"gasPrice"`
	Gas        int64         `json:"gas"`
}

type VoteProposalReply struct {
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
	Contributor   keys.Address          `json:"contributor_address"`
	WithdrawValue action.Amount         `json:"withdraw_value"`
	Beneficiary   keys.Address          `json:"beneficiary_address"`
	GasPrice      action.Amount         `json:"gasPrice"`
	Gas           int64                 `json:"gas"`
}
