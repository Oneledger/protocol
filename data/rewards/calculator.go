package rewards

import (
	"math/big"
	"time"

	"github.com/Oneledger/protocol/data/balance"
	tmstore "github.com/tendermint/tendermint/store"
)

type RewardYear struct {
	StartTime   time.Time
	CloseTime   time.Time
	Distributed *balance.Amount
}

type RewardYears struct {
	Years []RewardYear
}

type RewardCached struct {
	year      int
	cycleNo   int64
	burnedout bool
	amount    *balance.Amount
}

// cached data updates at every cycle
func NewRewardCached() RewardCached {
	return RewardCached{
		year:      -1,
		cycleNo:   0,
		burnedout: false,
		amount:    balance.NewAmount(0),
	}
}

type RewardCalculator struct {
	height      int64
	cached      RewardCached
	rewardYears RewardYears
	blockStore  *tmstore.BlockStore
	options     *Options
}

func NewRewardCalculator() *RewardCalculator {
	return &RewardCalculator{
		height: 0,
		cached: NewRewardCached(),
	}
}

func (calc *RewardCalculator) Reset(height int64, rewardYears RewardYears) {
	calc.height = height
	calc.rewardYears = rewardYears
}

func (calc *RewardCalculator) Burnedout() bool {
	return calc.cached.burnedout
}

func (calc *RewardCalculator) Calculate() (amt *balance.Amount, err error) {
	// return if all reward years already passed
	amt = balance.NewAmount(0)
	*amt = *calc.cached.amount
	if calc.cached.burnedout {
		return
	}

	// calculate again only at the beginning of a new cycle
	cycleNo, isNew := calc.enterNewCycle()
	if !isNew {
		return
	}

	// forcast how many more blocks can be generated before year close
	options := calc.options
	numofMoreBlocks, year := calc.numofMoreBlocksBeforeYearClose()
	if numofMoreBlocks == 0 {
		*amt = options.BurnoutRate
		calc.cacheResult(year, cycleNo, amt, true)
		return
	}

	// how much rewards left before year close
	yearSupply := options.YearBlockRewardShares[year]
	yearDistributed := calc.rewardYears.Years[year].Distributed
	yearLeft, err := yearSupply.Minus(*yearDistributed)
	if err != nil {
		// shouldn't happen by design
		logger.Errorf("Year rewards burned out unexpectedly, year= %v", year+1)
		return
	}

	// calculate rewards per block
	amt = balance.NewAmountFromBigInt(big.NewInt(0).Div(yearLeft.BigInt(), big.NewInt(numofMoreBlocks)))
	calc.cacheResult(year, cycleNo, amt, false)
	return
}

func (calc *RewardCalculator) secondsPerCycleLatest() (int64, time.Time) {
	tEnd := calc.blockStore.LoadBlockMeta(1).Header.Time.UTC()
	secsPerCycle := calc.options.EstimatedSecondsPerCycle
	if calc.height > calc.options.BlockSpeedCalculateCycle {
		// get speed calculation [begin, end] height
		cycle := calc.options.BlockSpeedCalculateCycle
		cycleEndHeight := (calc.height-1)/cycle*cycle + 1
		cycleBeginHeight := cycleEndHeight - cycle

		// duration of the cycle, in secs
		tBegin := calc.blockStore.LoadBlockMeta(cycleBeginHeight).Header.Time.UTC()
		tEnd = calc.blockStore.LoadBlockMeta(cycleEndHeight).Header.Time.UTC()
		secsPerCycle = int64(tEnd.Sub(tBegin).Seconds())
	}
	return secsPerCycle, tEnd
}

// return 0 blocks if all rewards years burned out
func (calc *RewardCalculator) numofMoreBlocksBeforeYearClose() (int64, int) {
	secsPerCycle, tCycleEnd := calc.secondsPerCycleLatest()

	numofMoreBlocks := int64(0)
	yearIndex := -1
	for i, rewardYear := range calc.rewardYears.Years {
		secsToClose := int64(rewardYear.CloseTime.Sub(tCycleEnd).Seconds())
		if secsToClose >= calc.options.YearCloseWindow {
			// calculate how many more blocks proportionally
			cycle := calc.options.BlockSpeedCalculateCycle
			numofMoreBlocks = int64(float64(secsToClose*cycle) / float64(secsPerCycle))
			if numofMoreBlocks == 0 {
				// this shouldn't happen if YearCloseWindow is set propoerly
				continue
			}
			yearIndex = i
			break
		}
	}
	return numofMoreBlocks, yearIndex
}

// see if we entered a new cycle
func (calc *RewardCalculator) enterNewCycle() (cycleNo int64, isNew bool) {
	isNew = false
	cycleNo = (calc.height-1)/calc.options.BlockSpeedCalculateCycle + 1
	if cycleNo > calc.cached.cycleNo {
		isNew = true
	}
	return
}

func (calc *RewardCalculator) cacheResult(year int, cycleNo int64, amt *balance.Amount, burnedout bool) {
	calc.cached.year = year
	calc.cached.cycleNo = cycleNo
	*calc.cached.amount = *amt
	calc.cached.burnedout = burnedout
}

func (rws *RewardCalculator) SetOptions(options *Options) {
	rws.options = options
}

func (calc *RewardCalculator) SetBlockStore(blockStore *tmstore.BlockStore) {
	calc.blockStore = blockStore
}
