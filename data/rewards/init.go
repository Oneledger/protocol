package rewards

type Options struct {
	RewardInterval    int64  `json:"rewardInterval"`
	RewardPoolAddress string `json:"rewardPoolAddress"`
}

type Interval struct {
	LastIndex  int64 `json:"lastIndex"`
	LastHeight int64 `json:"lastHeight"`
}
