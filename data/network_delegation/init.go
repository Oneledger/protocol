package network_delegation

import (
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

var logger *log.Logger
var prefixMap map[DelegationPrefixType]string

const (
	MatureType  DelegationPrefixType = 0x103
	PendingType DelegationPrefixType = 0x102
	ActiveType  DelegationPrefixType = 0x101

	DELEGATION_POOL_KEY       = "00000000000000000001"
	COMMISSION_PERCENTAGE     = 15
	BLOCK_PROPOSER_COMMISSION = 30

	//TODO this is hardcoded for now, will be changed in the future
	RewardsMaturityTime = 3
)

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("network_delegation")

	prefixMap = make(map[DelegationPrefixType]string)
	prefixMap[ActiveType] = "active_list"
	prefixMap[MatureType] = "mature_list"
	prefixMap[PendingType] = "pending_list"
}

type Options struct {
	RewardsMaturityTime int64 `json:"rewardsMaturityTime"`
}

type MasterStore struct {
	Deleg   *Store
	Rewards *DelegRewardStore
}

func NewMasterStore(pfxDeleg, pfxRewards string, state *storage.State) *MasterStore {
	return &MasterStore{
		Deleg:   NewStore(pfxDeleg, state),
		Rewards: NewDelegRewardStore(pfxRewards, state),
	}
}

func (master *MasterStore) WithState(state *storage.State) *MasterStore {
	master.Deleg.WithState(state)
	master.Rewards.WithState(state)
	return master
}
