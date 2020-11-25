package network_delegation

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"math/big"
)

type DelegationPrefixType int

//------------------------ State Types ------------------------------
type Delegator struct {
	Address *keys.Address `json:"address"`
	Amount  *balance.Coin `json:"amount"`
}

type PendingDelegator struct {
	Address *keys.Address `json:"address"`
	Amount  *balance.Coin `json:"amount"`
	Height  int64         `json:"height"`
}

type Reward struct {
	Amount  *balance.Amount
	Address keys.Address
}
type PendingReward struct {
	Amount  *balance.Amount
	Address keys.Address
	Height  int64
}

type State struct {
	ActiveList  []Delegator        `json:"active_list"`
	MatureList  []Delegator        `json:"mature_list"`
	PendingList []PendingDelegator `json:"pending_list"`
}

type RewardState struct {
	BalanceList []Reward        `json:"balance_list"`
	MatureList  []Reward        `json:"mature_list"`
	PendingList []PendingReward `json:"pending_list"`
}

//------------------------ State Types End ------------------------------

type PendingRewards struct {
	Height int64          `json:"height"`
	Amount balance.Amount `json:"amount"`
}

type DelegPendingRewards struct {
	Address keys.Address      `json:"address"`
	Rewards []*PendingRewards `json:"rewards"`
}

type DelegationRewardCtx struct {
	TotalRewards    *balance.Amount
	DelegationPower *big.Int
	TotalPower      *big.Int
	Height          int64
	ProposerAddress keys.Address
}

type DelegationRewardResponse struct {
	DelegationRewards *balance.Amount
	ProposerReward    *balance.Amount
	Commission        *balance.Amount
}

func (prefix DelegationPrefixType) GetJSONPrefix() string {
	return prefixMap[prefix]
}
