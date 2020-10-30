package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/network_delegation"
)

//-------------query
type GetUndelegatedRequest struct {
	Delegator keys.Address `json:"delegator"`
}

type GetTotalNetwkDelegation struct {
	OnlyActive int `json:"onlyActive"`
}

type SinglePendingAmount struct {
	Amount       balance.Amount `json:"amount"`
	MatureHeight int64          `json:"matureHeight"`
}

type GetUndelegatedReply struct {
	PendingAmounts []SinglePendingAmount `json:"pendingAmount"`
	MaturedAmount  balance.Amount        `json:"maturedAmount"`
	TotalAmount    balance.Amount        `json:"totalAmount"`
	Height         int64                 `json:"height"`
}

type GetTotalNetwkDelgReply struct {
	ActiveAmount  balance.Amount `json:"activeAmount"`
	PendingAmount balance.Amount `json:"pendingAmount"`
	MaturedAmount balance.Amount `json:"maturedAmount"`
	TotalAmount   balance.Amount `json:"totalAmount"`
	Height        int64          `json:"height"`
}

type GetDelegRewardsRequest struct {
	Delegator   keys.Address `json:"delegator"`
	InclPending bool         `json:"inclPending"`
}

type GetTotalDelegRewardsRequest struct{}

type GetTotalDelegRewardsReply struct {
	TotalRewards balance.Amount `json:"totalRewards"`
	Height       int64          `json:"height"`
}

type GetDelegRewardsReply struct {
	Balance balance.Amount                       `json:"balance"`
	Pending []*network_delegation.PendingRewards `json:"pending"`
	Matured balance.Amount                       `json:"matured"`
	Height  int64                                `json:"height"`
}

//-------------TX
type NetworkDelegateRequest struct {
	DelegationAddress keys.Address  `json:"delegationAddress"`
	Amount            action.Amount `json:"amount"`
	GasPrice          action.Amount `json:"gasPrice"`
	Gas               int64         `json:"gas"`
}

type NetUndelegateRequest struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
	GasPrice  action.Amount `json:"gasPrice"`
	Gas       int64         `json:"gas"`
}

type WithdrawDelegRewardsRequest struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
}

type ReinvestDelegRewardsRequest struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
}

type ListDelegationRequest struct {
	DelegationAddress keys.Address `json:"delegationAddress"`
}

type ListDelegationReply struct {
	DelegationStats DelegationStats `json:"delegationStats"`
	Height          int64           `json:"height"`
}

type DelegationStats struct {
	Active  string `json:"active"`
	Pending string `json:"pending"`
	Matured string `json:"matured"`
}

type FinalizeRewardsRequest struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
	GasPrice  action.Amount `json:"gasPrice"`
	Gas       int64         `json:"gas"`
}
