package rewards

import (
	"time"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type RewardCumulativeStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte

	rewardOptions *Options
}

func NewRewardCumulativeStore(prefix string, state *storage.State) *RewardCumulativeStore {
	return &RewardCumulativeStore{
		state:  state,
		prefix: storage.Prefix(prefix),
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (rws *RewardCumulativeStore) WithState(state *storage.State) *RewardCumulativeStore {
	rws.state = state
	return rws
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

// Calculate combined rewards for all voting validators for given block
func (rws *RewardCumulativeStore) CalcBlockRewards(blockHeight int64, bftTime time.Time) {

}

func (rws *RewardCumulativeStore) SetOptions(options *Options) {
	rws.rewardOptions = options
}

func (rws *RewardCumulativeStore) GetOptions() *Options {
	return rws.rewardOptions
}

//-----------------------------helper functions defined below
//
// Set cumulative amount by key
func (rws *RewardCumulativeStore) set(key storage.StoreKey, amt *balance.Amount) error {
	dat, err := rws.szlr.Serialize(amt)
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
