package network_delegation

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type DelegationPrefixType int

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
