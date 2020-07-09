package rewards

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/storage"
)

type Options struct {
	RewardInterval    int64          `json:"rewardInterval"`
	RewardPoolAddress string         `json:"rewardPoolAddress"`
	RewardCurrency    string         `json:"rewardCurrency"`
	CalculateInterval int64          `json:"calculateInterval"`
	AnnualSupply      balance.Amount `json:"annualSupply"`
	YearsOfSupply     int64          `json:"yearsOfSupply"`
}

type Interval struct {
	LastIndex  int64 `json:"lastIndex"`
	LastHeight int64 `json:"lastHeight"`
}

type RewardMasterStore struct {
	Reward   *RewardStore
	RewardCm *RewardCumulativeStore
}

func (rwz *RewardMasterStore) WithState(state *storage.State) *RewardMasterStore {
	rwz.Reward.WithState(state)
	rwz.RewardCm.WithState(state)
	return rwz
}

func (rwz *RewardMasterStore) SetOptions(options *Options) {
	rwz.Reward.SetOptions(options)
	rwz.RewardCm.SetOptions(options)
}

func (rwz *RewardMasterStore) GetOptions() *Options {
	return rwz.Reward.GetOptions()
}

func NewRewardMasterStore(rwz *RewardStore, rwzc *RewardCumulativeStore) *RewardMasterStore {
	return &RewardMasterStore{
		Reward:   rwz,
		RewardCm: rwzc,
	}
}
