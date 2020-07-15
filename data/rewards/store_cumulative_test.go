package rewards

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	tmstore "github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	memDb      db.DB
	store      *RewardCumulativeStore
	cs         *storage.State
	validator1 keys.Address
	validator2 keys.Address
	zero       *balance.Amount
	amt1       *balance.Amount
	amt2       *balance.Amount
	amt3       *balance.Amount
	amt4       *balance.Amount
	amt5       *balance.Amount
	withdraw1  *balance.Amount
	withdraw2  *balance.Amount

	estimatedSecondsPerBlock  = int64(17280)          // 0.2 day per block
	estimatedSecondsPerCycle  = int64(3600 * 24 * 10) // 1 cycle == 10 days
	blockSpeedCalculateCycle  = int64(50)             // calculate speed every 50 blocks
	burnoutRate, _            = balance.NewAmountFromString("5", 10)
	yearCloseWindow           = int64(3600 * 24) // 1 day
	yearBlockRewardShare_1, _ = balance.NewAmountFromString("70000000", 10)
	yearBlockRewardShare_2, _ = balance.NewAmountFromString("70000000", 10)
	yearBlockRewardShare_3, _ = balance.NewAmountFromString("40000000", 10)
	yearBlockRewardShare_4, _ = balance.NewAmountFromString("40000000", 10)
	yearBlockRewardShare_5, _ = balance.NewAmountFromString("10000000", 10)
	yearBlockRewardShares     = []balance.Amount{
		*yearBlockRewardShare_1,
		*yearBlockRewardShare_2,
		*yearBlockRewardShare_3,
		*yearBlockRewardShare_4,
		*yearBlockRewardShare_5,
	}

	rewzOpt = &Options{
		RewardCurrency:           "OLT",
		EstimatedSecondsPerCycle: estimatedSecondsPerCycle,
		BlockSpeedCalculateCycle: blockSpeedCalculateCycle,
		YearCloseWindow:          yearCloseWindow,
		YearBlockRewardShares:    yearBlockRewardShares,
		BurnoutRate:              *burnoutRate,
	}
)

func makeFakeBlock(blockStore *tmstore.BlockStore, height int64, bftTime time.Time) *types.Block {
	header := types.Header{Height: height, Time: bftTime}
	block := &types.Block{Header: header, LastCommit: &types.Commit{}}
	blockStore.SaveBlock(block, &types.PartSet{}, &types.Commit{})
	return block
}

func setupBlockStore(years int) time.Time {
	blockStore := tmstore.NewBlockStore(memDb)

	// seed
	tNow := time.Now()
	tStart := tNow
	rand.Seed(tNow.UnixNano())

	// simulates randomly generating at approximately 0.2day(4.3~5.3hours) per block
	secsPerBlock := int64(estimatedSecondsPerBlock)
	numofBlocks := 2000 * years
	for i := 0; i < numofBlocks; i++ { // generate enough blocks for 5 years
		secs := int64(0)
		if i > 0 {
			secs = secsPerBlock + rand.Int63n(3600) - 1800
			tNow = tNow.Add(time.Second * time.Duration(secs))
		}
		makeFakeBlock(blockStore, int64(i+1), tNow)
	}
	store.SetBlockStore(blockStore)
	return tStart
}

func setupRewardYears(tStart time.Time) RewardYears {
	numofYears := len(store.rewardOptions.YearBlockRewardShares)
	rewards := RewardYears{
		Years: make([]RewardYear, 0),
	}
	for i := 0; i < numofYears; i++ {
		tClose := tStart.AddDate(1, 0, 0).UTC()
		reward := RewardYear{
			StartTime:   tStart,
			CloseTime:   tClose,
			Distributed: balance.NewAmount(0),
		}
		tStart = tClose
		rewards.Years = append(rewards.Years, reward)
	}
	return rewards
}

func setup() {
	fmt.Println("####### Testing  rewards cumulative store #######")
	memDb = db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("rewards", memDb))

	store = NewRewardCumulativeStore("rwcum", cs)
	store.SetOptions(rewzOpt)
	generateAddresses()
}

func generateAddresses() {
	pub1, _, _ := keys.NewKeyPairFromTendermint()
	h1, _ := pub1.GetHandler()
	validator1 = h1.Address()

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()
	validator2 = h2.Address()

	zero = balance.NewAmount(0)
	amt1 = balance.NewAmount(100)
	amt2 = balance.NewAmount(200)
	amt3 = balance.NewAmount(377)
	withdraw1 = balance.NewAmount(163)
	withdraw2 = balance.NewAmount(499)
}

func TestNewRewardsCumulativeStore(t *testing.T) {
	setup()
	mutured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	withdrawn, err := store.GetWithdrawnRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, zero, mutured)
	assert.Equal(t, zero, balance)
	assert.Equal(t, zero, withdrawn)
}

func TestRewardsCumulativeStore_RewardYears(t *testing.T) {
	setup()
	tStart := setupBlockStore(1).UTC()
	years := setupRewardYears(tStart)

	// check pulled reward amount
	amount, burnedout, year, err := store.PullRewards(1, balance.NewAmount(0))
	assert.Nil(t, err)
	assert.False(t, burnedout)
	assert.Equal(t, 0, year)
	assert.True(t, zero.LessThan(*amount))

	// set consumed
	err = store.ConsumeRewards(amount, burnedout, year)
	assert.Nil(t, err)

	// check consumed
	rewardYears, err := store.GetYearDistributedRewards()
	years.Years[0].Distributed = years.Years[0].Distributed.Plus(*amount)
	assert.Nil(t, err)
	assert.Equal(t, years, rewardYears)
}

func TestRewardsCumulativeStore_AddGetMaturedBalance(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected := amt1.Plus(*amt2)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, matured, expected)
}

func TestRewardsCumulativeStore_WithdrawRewards(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	store.WithdrawRewards(validator1, withdraw1)
	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected, _ := amt1.Plus(*amt2).Minus(*withdraw1)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	expected = amt1.Plus(*amt2)
	assert.Equal(t, matured, expected)
}

func TestRewardsCumulativeStore_GetWithdrawnRewards(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	store.AddMaturedBalance(validator2, amt1)
	store.WithdrawRewards(validator1, withdraw1)
	store.AddMaturedBalance(validator1, amt3)
	store.AddMaturedBalance(validator2, amt2)
	store.WithdrawRewards(validator1, withdraw2)

	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected, _ := amt1.Plus(*amt2).Plus(*amt3).Minus(*withdraw1.Plus(*withdraw2))
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	expected = amt1.Plus(*amt2).Plus(*amt3)
	assert.Equal(t, matured, expected)

	withdrawn, err := store.GetWithdrawnRewards(validator1)
	assert.Nil(t, err)
	expected = withdraw1.Plus(*withdraw2)
	assert.Equal(t, withdrawn, expected)
}

func TestRewardsCumulativeStore_WithdrawOthers(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	err := store.WithdrawRewards(validator2, withdraw1)
	assert.NotNil(t, err)

	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected := amt1.Plus(*amt2)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, matured, expected)
}

func TestRewardsCumulativeStore_OverWithdraw(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	err := store.WithdrawRewards(validator1, withdraw2)
	assert.NotNil(t, err)

	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected := amt1.Plus(*amt2)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, matured, expected)
}
