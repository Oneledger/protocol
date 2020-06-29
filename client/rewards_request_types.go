package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type ListRewardsRequest struct {
	Validator string `json:"validator"`
}

type ListRewardsReply struct {
	Validator keys.Address     `json:"validator"`
	Rewards   []balance.Amount `json:"rewards"`
}

type WithdrawRewardsRequest struct {
	Validator keys.Address  `json:"validator"`
	GasPrice  action.Amount `json:"gasPrice"`
	Gas       int64         `json:"gas"`
}

type WithdrawRewardsRequest struct {
	GasPrice action.Amount `json:"gasPrice"`
	Gas      int64         `json:"gas"`
}
