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

type MatureRewardsReply struct {
	Validator     keys.Address   `json:"validator"`
	MatureRewards balance.Amount `json:"matureRewards"`
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
