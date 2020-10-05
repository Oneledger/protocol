package network_delegation

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type PendingRewards struct {
	Height int64          `json:"height"`
	Amount balance.Amount `json:"amount"`
}

type DelegPendingRewards struct {
	Address keys.Address      `json:"address"`
	Rewards []*PendingRewards `json:"rewards"`
}
