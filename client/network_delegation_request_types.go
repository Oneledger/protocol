package client

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type WithdrawDelegRewardsRequest struct {
	Delegator keys.Address   `json:"delegator"`
	Amount    balance.Amount `json:"amount"`
}
