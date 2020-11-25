package client

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type ListRewardsRequest struct{}

type RewardsRequest struct {
	Validator string `json:"validator"`
}

type RewardRecord struct {
	Index  int64
	Amount balance.Amount
}

type ListRewardsReply struct {
	Validator keys.Address   `json:"validator"`
	Rewards   []RewardRecord `json:"rewards"`
	Height    int64          `json:"height"`
}

type ValidatorRewardStats struct {
	Address         keys.Address   `json:"address"`
	PendingAmount   balance.Amount `json:"pendingAmount"`
	WithdrawnAmount balance.Amount `json:"withdrawnAmount"`
	MatureBalance   balance.Amount `json:"matureBalance"`
	TotalAmount     balance.Amount `json:"totalAmount"`
}

type RewardStat struct {
	Validators   []ValidatorRewardStats `json:"validators"`
	TotalRewards balance.Amount         `json:"totalRewards"`
	Height       int64                  `json:"height"`
}

type WithdrawRewardsRequest struct {
	ValidatorAddress keys.Address   `json:"validatorSigningAddress"`
	WithdrawAmount   balance.Amount `json:"withdrawAmount"`
	//GasPrice         action.Amount  `json:"gasPrice"`
	//Gas              int64          `json:"gas"`
}

type WithdrawRewardsReply struct {
	RawTx []byte `json:"rawTx"`
	//Signature action.Signature `json:"signature"`
}
