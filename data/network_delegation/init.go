package network_delegation

import (
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("network_delegation")
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
