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
	TotalAmount    balance.Amount        `json:"totalAmount"`
	Height         int64                 `json:"height"`
}

type GetTotalNetwkDelgReply struct {
	ActiveAmount  balance.Amount `json:"activeAmount"`
	PendingAmount balance.Amount `json:"pendingAmount"`
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
	//Matured balance.Amount                       `json:"matured"`
	Height int64 `json:"height"`
}

type ListDelegationRequest struct {
	DelegationAddresses []keys.Address `json:"delegationAddresses"`
}

type ListDelegationReply struct {
	AllDelegStats []*FullDelegStats `json:"allDelegStats"`
	Height        int64             `json:"height"`
}

type DelegStats struct {
	Active  balance.Amount `json:"active"`
	Pending balance.Amount `json:"pending"`
}

type DelegRewardsStats struct {
	Active  balance.Amount `json:"active"`
	Pending balance.Amount `json:"pending"`
}

type FullDelegStats struct {
	DelegAddress      keys.Address      `json:"delegationAddress"`
	DelegStats        DelegStats        `json:"delegationStats"`
	DelegRewardsStats DelegRewardsStats `json:"delegationRewardsStats"`
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

//type FinalizeRewardsRequest struct {
//	Delegator keys.Address  `json:"delegator"`
//	Amount    action.Amount `json:"amount"`
//	GasPrice  action.Amount `json:"gasPrice"`
//	Gas       int64         `json:"gas"`
//}

type ReinvestDelegRewardsRequest struct {
	Delegator keys.Address  `json:"delegator"`
	Amount    action.Amount `json:"amount"`
}
