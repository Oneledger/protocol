package rewards

import (
	"errors"
	"time"

	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type RewardCumulativeStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte

	node          *consensus.Node
	currencies    *balance.CurrencySet
	rewardOptions *Options
}

func NewRewardCumulativeStore(prefix string, state *storage.State) *RewardCumulativeStore {
	return &RewardCumulativeStore{
		state:  state,
		prefix: storage.Prefix(prefix),
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (rws *RewardCumulativeStore) Init(node *consensus.Node, currencies *balance.CurrencySet) {
	rws.node = node
	rws.currencies = currencies
}

func (rws *RewardCumulativeStore) WithState(state *storage.State) *RewardCumulativeStore {
	rws.state = state
	return rws
}

// Pull a combined block rewards from total supply for all voting validators for given block height
// Pulled rewards MUST be 100% distributed since it has already been deducted from total supply
func (rws *RewardCumulativeStore) PullRewards(blockHeight int64, bftTime time.Time, poolAmt balance.Amount) (amount *balance.Amount, burnedout bool, err error) {
	// get year distributed amount till now
	yearRewards, err := rws.GetYearDistributedRewards()
	if err != nil {
		return
	}

	// calculate reward for each block
	calc := NewRewardCalculator(blockHeight, yearRewards, rws.node, rws.rewardOptions)
	amount, burnedout, year, err := calc.Calculate()
	if err != nil {
		return
	}
	if burnedout && poolAmt.LessThan(*amount) {
		*amount = poolAmt
	}

	// accumulates total distributed rewards
	err = rws.addTotalDistributedRewards(amount)
	if err != nil {
		return
	}

	// accumulates year distributed rewards
	rws.addYearDistributedRewards(yearRewards, year, amount)
	if err != nil {
		return
	}

	return
}

// Get total distributed rewards till now.
func (rws *RewardCumulativeStore) GetTotalDistributedRewards() (amt *balance.Amount, err error) {
	key := rws.getTotalDistributedKey()
	amt, err = rws.get(key)
	return
}

// Get a list of each year's distributed rewards till now
func (rws *RewardCumulativeStore) GetYearDistributedRewards() (yearRewards []YearReward, err error) {
	key := rws.getYearDistributedKey()
	if !rws.state.Exists(key) {
		yearRewards, err = rws.initYearRewards(key)
		return
	}
	yearRewards, err = rws.getYearRewards(key)
	return
}

// Get matured rewards balance, the widrawable rewards, till now.
func (rws *RewardCumulativeStore) GetMaturedBalance(validator keys.Address) (amt *balance.Amount, err error) {
	key := rws.getBalanceKey(validator)
	amt, err = rws.get(key)
	return
}

// Add an 'amount' of matured rewards to rewards balance
func (rws *RewardCumulativeStore) AddMaturedBalance(validator keys.Address, amount *balance.Amount) error {
	key := rws.getBalanceKey(validator)
	amt, err := rws.get(key)
	if err != nil {
		return err
	}

	err = rws.set(key, amt.Plus(*amount))
	return err
}

// Get total matured rewards till now, including withdrawn amount. This number is calculated on the fly
func (rws *RewardCumulativeStore) GetMaturedRewards(validator keys.Address) (amt *balance.Amount, err error) {
	keyBalance := rws.getBalanceKey(validator)
	amtBalance, err := rws.get(keyBalance)
	if err != nil {
		return
	}

	keyWithdrawn := rws.getWithdrawnKey(validator)
	amtWithdrawn, err := rws.get(keyWithdrawn)
	if err != nil {
		return
	}

	amt = amtBalance.Plus(*amtWithdrawn)
	return
}

// Get total rewards withdrawn till now
func (rws *RewardCumulativeStore) GetWithdrawnRewards(validator keys.Address) (amt *balance.Amount, err error) {
	key := rws.getWithdrawnKey(validator)
	amt, err = rws.get(key)
	return
}

// Withdraw an 'amount' of rewards from rewards balance
func (rws *RewardCumulativeStore) WithdrawRewards(validator keys.Address, amount *balance.Amount) error {

	err := rws.minusRewardsBalance(validator, amount)
	if err != nil {
		return errors.Wrap(err, "Minus from Matured Balance")
	}
	err = rws.addWithdrawnRewards(validator, amount)
	if err != nil {
		return errors.Wrap(err, "Add to Withdraw Balance")
	}

	return nil
}

func (rws *RewardCumulativeStore) SetOptions(options *Options) {
	rws.rewardOptions = options
}

func (rws *RewardCumulativeStore) GetOptions() *Options {
	return rws.rewardOptions
}

//-----------------------------helper functions defined below
//
// Set cumulative amount(s) by key
func (rws *RewardCumulativeStore) set(key storage.StoreKey, amts interface{}) error {
	dat, err := rws.szlr.Serialize(amts)
	if err != nil {
		return err
	}
	err = rws.state.Set(key, dat)
	return err
}

// Get cumulative amount by key
func (rws *RewardCumulativeStore) get(key storage.StoreKey) (amt *balance.Amount, err error) {
	dat, err := rws.state.Get(key)
	if err != nil {
		return
	}
	amt = balance.NewAmount(0)
	if len(dat) == 0 {
		return
	}
	err = rws.szlr.Deserialize(dat, amt)
	return
}

// Get year rewards
func (rws *RewardCumulativeStore) getYearRewards(key storage.StoreKey) (rewards []YearReward, err error) {
	dat, err := rws.state.Get(key)
	if err != nil {
		return
	}
	rewards = make([]YearReward, 0)
	if len(dat) == 0 {
		err = YearRewardsMissing
		return
	}
	err = rws.szlr.Deserialize(dat, rewards)
	return
}

// Calculate year from time
func (rws *RewardCumulativeStore) getRewardYear(bftTime time.Time) int {
	tBlock := bftTime.UTC()
	tStart := rws.node.BlockStore().LoadBlock(1).Header.Time.UTC()

	year := tStart.Year()
	tStart = tStart.AddDate(1, 0, 0)
	for tBlock.After(tStart) || tBlock.Equal(tStart) {
		year++
		tStart = tStart.AddDate(1, 0, 0)
	}
	return year
}

// Initialize each reward year's information
func (rws *RewardCumulativeStore) initYearRewards(key storage.StoreKey) (rewards []YearReward, err error) {
	tStart := rws.node.BlockStore().LoadBlock(1).Header.Time.UTC()
	numofYears := len(rws.rewardOptions.YearBlockRewardShares)

	// calculate each year's start/close time
	rewards = make([]YearReward, 0)
	for i := 0; i < numofYears; i++ {
		tClose := tStart.AddDate(1, 0, 0).UTC()
		reward := YearReward{
			StartTime:   tStart,
			CloseTime:   tClose,
			Distributed: balance.NewAmount(0),
		}
		tStart = tClose
		rewards = append(rewards, reward)
		logger.Infof("Initial year-%v start: %s, close: %s", i+1, tStart, tClose)
	}

	// save to DB
	err = rws.set(key, rewards)
	return
}

// Key for total distributed rewards
func (rws *RewardCumulativeStore) getTotalDistributedKey() []byte {
	key := string(rws.prefix) + "tdist"
	return storage.StoreKey(key)
}

// Key for each year's distributed rewards
func (rws *RewardCumulativeStore) getYearDistributedKey() []byte {
	key := string(rws.prefix) + "ydist"
	return storage.StoreKey(key)
}

// Key for balance
func (rws *RewardCumulativeStore) getBalanceKey(validator keys.Address) []byte {
	key := string(rws.prefix) + validator.String() + storage.DB_PREFIX + "balance"
	return storage.StoreKey(key)
}

// Key for withdrawn
func (rws *RewardCumulativeStore) getWithdrawnKey(validator keys.Address) []byte {
	key := string(rws.prefix) + validator.String() + storage.DB_PREFIX + "withdrawn"
	return storage.StoreKey(key)
}

// Key for block rewards calculator
func (rws *RewardCumulativeStore) getCalculatorKey(blockHeight int64) []byte {
	key := string(rws.prefix) + "calc"
	return storage.StoreKey(key)
}

// Add to total distributed rewards
func (rws *RewardCumulativeStore) addTotalDistributedRewards(amount *balance.Amount) error {
	key := rws.getTotalDistributedKey()
	amt, err := rws.get(key)
	if err != nil {
		return err
	}

	err = rws.set(key, amt.Plus(*amount))
	return err
}

// Add to total year distributed rewards
func (rws *RewardCumulativeStore) addYearDistributedRewards(yearRewards []YearReward, year int, amount *balance.Amount) error {
	key := rws.getYearDistributedKey()
	yearDist := yearRewards[year].Distributed
	yearRewards[year].Distributed = yearDist.Plus(*amount)

	err := rws.set(key, yearRewards)
	return err
}

// Deducts an 'amount' of rewards from rewards balance
func (rws *RewardCumulativeStore) minusRewardsBalance(validator keys.Address, amount *balance.Amount) error {
	key := rws.getBalanceKey(validator)
	amt, err := rws.get(key)
	if err != nil {
		return err
	}

	result, err := amt.Minus(*amount)
	if err != nil {
		return err
	}

	err = rws.set(key, result)
	return err
}

// Add to total rewards withdrawn till now
func (rws *RewardCumulativeStore) addWithdrawnRewards(validator keys.Address, amount *balance.Amount) error {
	key := rws.getWithdrawnKey(validator)
	amt, err := rws.get(key)
	if err != nil {
		return err
	}

	err = rws.set(key, amt.Plus(*amount))
	return err
}
