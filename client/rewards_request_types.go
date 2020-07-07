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
	ValidatorSigningAddress keys.Address  `json:"validatorSigningAddress"`
	WithdrawAmount          action.Amount `json:"withdrawAmount"`
	GasPrice                action.Amount `json:"gasPrice"`
	Gas                     int64         `json:"gas"`
}
