package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/network_delegation"
)

type NetworkDelegateRequest struct {
	DelegationAddress keys.Address  `json:"delegationAddress"`
	Amount            action.Amount `json:"amount"`
	GasPrice          action.Amount `json:"gasPrice"`
	Gas               int64         `json:"gas"`
}

type GetDelegRewardsRequest struct {
	Delegator   keys.Address `json:"delegator"`
	InclPending bool         `json:"inclPending"`
}

type GetDelegRewardsReply struct {
	Balance balance.Amount                       `json:"balance"`
	Pending []*network_delegation.PendingRewards `json:"pending"`
	Matured balance.Amount                       `json:"matured"`
	Height  int64                                `json:"height"`
}

type WithdrawDelegRewardsRequest struct {
	Delegator keys.Address   `json:"delegator"`
	Amount    balance.Amount `json:"amount"`
}
