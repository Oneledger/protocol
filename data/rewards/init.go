package rewards

import (
	"os"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("rewards")
}

type Options struct {
	RewardInterval           int64            `json:"rewardInterval"`
	RewardPoolAddress        string           `json:"rewardPoolAddress"`
	RewardCurrency           string           `json:"rewardCurrency"`
	EstimatedSecondsPerCycle int64            `json:"estimatedSecondsPerCycle"`
	BlockSpeedCalculateCycle int64            `json:"blockSpeedCalculateCycle"`
	YearCloseWindow          int64            `json:"yearCloseWindow"`
	YearBlockRewardShares    []balance.Amount `json:"yearBlockRewardShares"`
	BurnoutRate              balance.Amount   `json:"burnoutRate"`
}

type Interval struct {
	LastIndex  int64 `json:"lastIndex"`
	LastHeight int64 `json:"lastHeight"`
}

type RewardMasterState struct {
	RewardState *RewardState     `json:"rewardState"`
	CumuState   *RewardCumuState `json:"cumuState"`
}

func NewRewardMasterState() *RewardMasterState {
	return &RewardMasterState{
		RewardState: NewRewardState(),
		CumuState:   NewRewardCumuState(),
	}
}

type RewardMasterStore struct {
	Reward   *RewardStore
	RewardCm *RewardCumulativeStore
}

func NewRewardMasterStore(rwz *RewardStore, rwzc *RewardCumulativeStore) *RewardMasterStore {
	return &RewardMasterStore{
		Reward:   rwz,
		RewardCm: rwzc,
	}
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

func (rwz *RewardMasterStore) DumpState() (masterState *RewardMasterState, succeed bool) {
	masterState = NewRewardMasterState()
	rwdState, err := rwz.Reward.dumpState()
	if err != nil {
		return
	}
	cumuState, err := rwz.RewardCm.dumpState()
	if err != nil {
		return
	}
	masterState.RewardState = rwdState
	masterState.CumuState = cumuState
	succeed = true
	return
}

func (rwz *RewardMasterStore) LoadState(masterState RewardMasterState) (succeed bool) {
	err := rwz.Reward.loadState(masterState.RewardState)
	if err != nil {
		return
	}
	err = rwz.RewardCm.loadState(masterState.CumuState)
	if err != nil {
		return
	}
	succeed = true
	return
}
