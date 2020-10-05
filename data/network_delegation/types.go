package network_delegation

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type DelegationPrefixType int

type PendingRewards struct {
	Height int64          `json:"height"`
	Amount balance.Amount `json:"amount"`
}

type DelegPendingRewards struct {
	Address keys.Address      `json:"address"`
	Rewards []*PendingRewards `json:"rewards"`
}

type Delegator struct {
	Address keys.Address
	Amount  balance.Coin
}

type PendingDelegator struct {
	Address keys.Address
	Amount  balance.Coin
	Height  int64
}

type State struct {
	ActiveList  []Delegator
	MatureList  []Delegator
	PendingList []PendingDelegator
}

func dumpDelegators() {
}
