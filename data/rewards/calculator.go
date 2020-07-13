package rewards

import (
	"math/big"
	"time"

	"github.com/Oneledger/protocol/data/balance"
	tmnode "github.com/tendermint/tendermint/node"
)

type YearReward struct {
	StartTime   time.Time
	CloseTime   time.Time
	Distributed *balance.Amount
}

type RewardCalculator struct {
	height      int64
	yearRewards []YearReward
	node        *tmnode.Node
	options     *Options
}

func NewRewardCalculator(height int64, yearRewards []YearReward, node *tmnode.Node, options *Options) *RewardCalculator {
	return &RewardCalculator{
		height:      height,
		yearRewards: yearRewards,
		node:        node,
		options:     options,
	}
}

func (calc *RewardCalculator) Calculate() (amt *balance.Amount, burnedout bool, year int, err error) {
	amt = balance.NewAmount(0)
	burnedout = false

	// forcast how many more blocks can be generated before year close
	options := calc.options
	numofMoreBlocks, year := calc.numofMoreBlocksBeforeYearClose()
	if numofMoreBlocks == 0 {
		*amt = options.BurnoutRate
		burnedout = true
		return
	}

	// how much rewards left before year close
	yearRewardSupply := options.YearBlockRewardShares[year]
	yearRewardDistributed := calc.yearRewards[year].Distributed
	yearRewardLeft, err := yearRewardSupply.Minus(*yearRewardDistributed)
	if err != nil {
		// shouldn't happen by design
		logger.Errorf("Year rewards burned out unexpectedly, year= %v", year+1)
		return
	}

	// calculate rewards per block
	amt = balance.NewAmountFromBigInt(big.NewInt(0).Div(yearRewardLeft.BigInt(), big.NewInt(numofMoreBlocks)))
	return
}

func (calc *RewardCalculator) secondsPerCycleLatest() (int64, time.Time) {
	tEnd := calc.node.BlockStore().LoadBlockMeta(1).Header.Time.UTC()
	secsPerCycle := calc.options.EstimatedSecondsPerCycle
	if calc.height > calc.options.BlockSpeedCalculateCycle {
		// get speed calculation [begin, end] height
		cycle := calc.options.BlockSpeedCalculateCycle
		cycleEnd := (calc.height-1)/cycle*cycle + 1
		cycleBegin := cycleEnd - cycle

		// duration of the cycle, in secs
		tBegin := calc.node.BlockStore().LoadBlockMeta(cycleBegin).Header.Time.UTC()
		tEnd = calc.node.BlockStore().LoadBlockMeta(cycleEnd).Header.Time.UTC()
		secsPerCycle = int64(tEnd.Sub(tBegin).Seconds())
	}
	return secsPerCycle, tEnd
}

// return 0 blocks if all rewards years burned out
func (calc *RewardCalculator) numofMoreBlocksBeforeYearClose() (int64, int) {
	secsPerCycle, tCycleEnd := calc.secondsPerCycleLatest()

	numofMoreBlocks := int64(0)
	yearIndex := -1
	for i, rewardYear := range calc.yearRewards {
		secsToClose := int64(rewardYear.CloseTime.Sub(tCycleEnd).Seconds())
		if secsToClose >= calc.options.YearCloseWindow {
			// calculate how many more blocks proportionally
			cycle := calc.options.BlockSpeedCalculateCycle
			numofMoreBlocks = int64(float64(secsToClose) / float64(secsPerCycle) * float64(cycle))
			if numofMoreBlocks == 0 {
				// this shouldn't happen if YearCloseWindow is set propoerly
				continue
			}
			yearIndex = i
		}
	}
	return numofMoreBlocks, yearIndex
}
