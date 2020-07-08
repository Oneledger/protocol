package client

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type RewardsRequest struct {
	Validator string `json:"validator"`
}

type ListRewardsReply struct {
	Validator keys.Address     `json:"validator"`
	Rewards   []balance.Amount `json:"rewards"`
}

type ValidatorRewardStat struct {
	PendingAmount   balance.Amount `json:"pendingAmount"`
	WithdrawnAmount balance.Amount `json:"withdrawnAmount"`
	MatureBalance   balance.Amount `json:"matureBalance"`
	TotalAmount     balance.Amount `json:"totalAmount"`
}

type RewardStat struct {
	Validators   []ValidatorRewardStat `json:"validators"`
	TotalRewards balance.Amount        `json:"totalRewards"`
}
