package client

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/network_delegation"
)

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
