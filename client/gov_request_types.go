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

type ListProposalsRequest struct {
	ProposalId governance.ProposalID `json:"proposal_id"`
	State      string                `json:"state"`
}

type ListProposalsReply struct {
	Proposals     []governance.Proposal    `json:"proposals"`
	ProposalFunds []balance.Amount         `json:"proposal_funds"`
	ProposalVotes []governance.VoteStatus  `json:"proposal_votes"`
	State         governance.ProposalState `json:"state"`
	Height        int64                    `json:"height"`
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

type CancelProposalRequest struct {
	ProposalId governance.ProposalID `json:"proposal_id"`
	Proposer   keys.Address          `json:"proposer"`
	Reason     string                `json:"reason"`
	GasPrice   action.Amount         `json:"gasPrice"`
	Gas        int64                 `json:"gas"`
}

type WithdrawFundsRequest struct {
	ProposalID    governance.ProposalID `json:"proposal_id"`
	Contributor   keys.Address          `json:"contributor_address"`
	WithdrawValue action.Amount         `json:"withdraw_value"`
	Beneficiary   keys.Address          `json:"beneficiary_address"`
	GasPrice      action.Amount         `json:"gasPrice"`
	Gas           int64                 `json:"gas"`
}

type FinalizeProposalRequest struct {
	ProposalId governance.ProposalID `json:"proposal_id"`
	Proposer   action.Address        `json:"proposer"`
	GasPrice   action.Amount         `json:"gasPrice"`
	Gas        int64                 `json:"gas"`
}
