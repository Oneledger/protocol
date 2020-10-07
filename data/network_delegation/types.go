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
	Address *keys.Address `json:"address"`
	Amount  *balance.Coin `json:"amount"`
}

type PendingDelegator struct {
	Address *keys.Address `json:"address"`
	Amount  *balance.Coin `json:"amount"`
	Height  int64         `json:"height"`
}

type State struct {
	ActiveList  []Delegator        `json:"active_list"`
	MatureList  []Delegator        `json:"mature_list"`
	PendingList []PendingDelegator `json:"pending_list"`
}

func (prefix DelegationPrefixType) GetJSONPrefix() string {
	return prefixMap[prefix]
}
