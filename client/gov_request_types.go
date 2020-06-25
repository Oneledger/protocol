package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

type CreateProposalRequest struct {
	ProposalID     string                     `json:"proposalId"`
	ProposalType   string                     `json:"proposalType"`
	Headline       string                     `json:"headline"`
	Description    string                     `json:"description"`
	Proposer       keys.Address               `json:"proposer"`
	InitialFunding action.Amount              `json:"initialFunding"`
	ConfigUpdate   governance.GovernanceState `json:"configUpdate"`
	GasPrice       action.Amount              `json:"gasPrice"`
	Gas            int64                      `json:"gas"`
}

type ListProposalRequest struct {
	ProposalId governance.ProposalID `json:"proposalId"`
}

type ListProposalsRequest struct {
	State        governance.ProposalState `json:"state"`
	Proposer     keys.Address             `json:"proposer"`
	ProposalType governance.ProposalType  `json:"proposalType"`
}

type ProposalStat struct {
	Proposal governance.Proposal   `json:"proposal"`
	Funds    balance.Amount        `json:"funds"`
	Votes    governance.VoteStatus `json:"votes"`
}

type ListProposalsReply struct {
	ProposalStats []ProposalStat `json:"proposalStats"`
	Height        int64          `json:"height"`
}

type GovernanceOptionsRequest struct {
}
type GovernanceOptionsReply struct {
	GovOptions governance.GovernanceState `json:"govOptions"`
}

type VoteProposalRequest struct {
	ProposalId string                 `json:"proposalId"`
	Opinion    governance.VoteOpinion `json:"opinion"`
	Address    keys.Address           `json:"address"`
	GasPrice   action.Amount          `json:"gasPrice"`
	Gas        int64                  `json:"gas"`
}

type VoteProposalReply struct {
	RawTx     []byte           `json:"rawTx"`
	Signature action.Signature `json:"signature"`
}

type FundProposalRequest struct {
	ProposalId    governance.ProposalID `json:"proposalId"`
	FundValue     action.Amount         `json:"fundValue"`
	FunderAddress action.Address        `json:"funderAddress"`
	GasPrice      action.Amount         `json:"gasPrice"`
	Gas           int64                 `json:"gas"`
}

type CancelProposalRequest struct {
	ProposalId governance.ProposalID `json:"proposalId"`
	Proposer   keys.Address          `json:"proposer"`
	Reason     string                `json:"reason"`
	GasPrice   action.Amount         `json:"gasPrice"`
	Gas        int64                 `json:"gas"`
}

type WithdrawFundsRequest struct {
	ProposalID    governance.ProposalID `json:"proposalId"`
	Funder        keys.Address          `json:"funderAddress"`
	WithdrawValue action.Amount         `json:"withdrawValue"`
	Beneficiary   keys.Address          `json:"beneficiaryAddress"`
	GasPrice      action.Amount         `json:"gasPrice"`
	Gas           int64                 `json:"gas"`
}

type FinalizeProposalRequest struct {
	ProposalId governance.ProposalID `json:"proposalId"`
	Proposer   action.Address        `json:"proposer"`
	GasPrice   action.Amount         `json:"gasPrice"`
	Gas        int64                 `json:"gas"`
}
