package network_delegation

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type DelegRewardStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte
}

func NewDelegRewardStore(prefix string, state *storage.State) *DelegRewardStore {
	return &DelegRewardStore{
		state:  state,
		prefix: storage.Prefix(prefix),
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (drs *DelegRewardStore) WithState(state *storage.State) *DelegRewardStore {
	drs.state = state
	return drs
}

// Add rewards balance
func (drs *DelegRewardStore) AddRewardsBalance(delegator keys.Address, amount *balance.Amount) error {
	key := drs.getRewardsBalanceKey(delegator)
	amt, err := drs.get(key)
	if err != nil {
		return err
	}

	err = drs.set(key, amt.Plus(*amount))
	return err
}

// Get rewards balance
func (drs *DelegRewardStore) GetRewardsBalance(delegator keys.Address) (amt *balance.Amount, err error) {
	key := drs.getRewardsBalanceKey(delegator)
	amt, err = drs.get(key)
	return
}

// Initiate a withdrawal of an 'amount' of rewards
func (drs *DelegRewardStore) Withdraw(delegator keys.Address, amount *balance.Amount, matureHeight int64) error {
	err := drs.minusRewardsBalance(delegator, amount)
	if err != nil {
		return errors.Wrap(err, "Minus from rewards balance")
	}
	err = drs.addPendingRewards(delegator, amount, matureHeight)
	if err != nil {
		return errors.Wrap(err, "Add to pending rewards")
	}

	return nil
}

// Get pending withdrawn rewards
func (drs *DelegRewardStore) GetPendingRewards(delegator keys.Address, height, blocks int64) (pdRewards *DelegPendingRewards, err error) {
	pdRewards = &DelegPendingRewards{Address: delegator}
	for h := height; h < height+blocks; h++ {
		key := drs.getPendingRewardsKey(h, delegator)
		var amt *balance.Amount
		amt, err = drs.get(key)
		if err != nil {
			return
		}
		if !amt.Equals(balance.AmtZero) {
			pdRewards.Rewards = append(pdRewards.Rewards, &PendingRewards{
				Height: h,
				Amount: *amt,
			})
		}
	}
	return
}

// Mature, if any, all delegators' pending rewards at a specific height
func (drs *DelegRewardStore) MaturePendingRewards(height int64) {
	drs.iteratePD(height, func(delegator keys.Address, amt *balance.Amount) bool {
		// clear pending amount
		key := drs.getPendingRewardsKey(height, delegator)
		err := drs.set(key, balance.AmtZero)
		if err != nil {
			return true
		}
		// increase matured amount
		if !amt.Equals(balance.AmtZero) {
			err = drs.addMaturedRewards(delegator, amt)
		}
		return err != nil
	})
}

// Get matured(finalizable) rewards
func (drs *DelegRewardStore) GetMaturedRewards(delegator keys.Address) (amt *balance.Amount, err error) {
	key := drs.getMaturedRewardsKey(delegator)
	amt, err = drs.get(key)
	return
}

// Finalize(deduct) an 'amount' of matured rewards
func (drs *DelegRewardStore) Finalize(delegator keys.Address, amount *balance.Amount) error {
	key := drs.getMaturedRewardsKey(delegator)
	amt, err := drs.get(key)
	if err != nil {
		return err
	}

	result, err := amt.Minus(*amount)
	if err != nil {
		return err
	}

	err = drs.set(key, result)
	return err
}

//-----------------------------helper functions
//
// Set object by key
func (drs *DelegRewardStore) set(key storage.StoreKey, obj interface{}) error {
	dat, err := drs.szlr.Serialize(obj)
	if err != nil {
		return err
	}
	err = drs.state.Set(key, dat)
	return err
}

// Get amount by key
func (drs *DelegRewardStore) get(key storage.StoreKey) (amt *balance.Amount, err error) {
	dat, err := drs.state.Get(key)
	if err != nil {
		return
	}
	amt = balance.NewAmount(0)
	if len(dat) == 0 {
		return
	}
	err = drs.szlr.Deserialize(dat, amt)
	return
}

// iterate rewards by 'balance', 'matured'
func (drs *DelegRewardStore) iterate(subkey string, fn func(delegator keys.Address, amt *balance.Amount) bool) (stopped bool) {
	prefix := append(drs.prefix, subkey...)
	return drs.state.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			amt := balance.NewAmount(0)
			err := drs.szlr.Deserialize(value, amt)
			if err != nil {
				logger.Error("failed to deserialize delegator rewards amount")
				return false
			}
			addr := keys.Address{}
			bytesText := key[len(prefix):]
			err = addr.UnmarshalText(bytesText)
			if err != nil {
				logger.Error("failed to deserialize delegator address")
				return false
			}
			return fn(addr, amt)
		},
	)
}

// iterate pending rewards by height
func (drs *DelegRewardStore) iteratePD(height int64, fn func(delegator keys.Address, amt *balance.Amount) bool) (stopped bool) {
	pfxStr := fmt.Sprintf("%spending_%d_", string(drs.prefix), height)
	prefix := storage.StoreKey(pfxStr)
	return drs.state.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			amt := balance.NewAmount(0)
			err := drs.szlr.Deserialize(value, amt)
			if err != nil {
				logger.Error("failed to deserialize delegator pending rewards amount")
				return true
			}
			addr := keys.Address{}
			bytesText := key[len(prefix):]
			err = addr.UnmarshalText(bytesText)
			if err != nil {
				logger.Error("failed to deserialize delegator address")
				return true
			}
			return fn(addr, amt)
		},
	)
}

// Key for delegator rewards balance
func (drs *DelegRewardStore) getRewardsBalanceKey(delegator keys.Address) []byte {
	key := fmt.Sprintf("%sbalance_%s", string(drs.prefix), delegator)
	return storage.StoreKey(key)
}

// Key for delegator withdrawn but unmatured amount
func (drs *DelegRewardStore) getPendingRewardsKey(height int64, delegator keys.Address) []byte {
	key := fmt.Sprintf("%spending_%d_%s", string(drs.prefix), height, delegator)
	return storage.StoreKey(key)
}

// Key for delegator withdrawn and matured amount
func (drs *DelegRewardStore) getMaturedRewardsKey(delegator keys.Address) []byte {
	key := fmt.Sprintf("%smatured_%s", string(drs.prefix), delegator)
	return storage.StoreKey(key)
}

// Deducts an 'amount' of rewards from rewards balance
func (drs *DelegRewardStore) minusRewardsBalance(delegator keys.Address, amount *balance.Amount) error {
	key := drs.getRewardsBalanceKey(delegator)
	amt, err := drs.get(key)
	if err != nil {
		return err
	}

	result, err := amt.Minus(*amount)
	if err != nil {
		return err
	}

	err = drs.set(key, result)
	return err
}

// Add an 'amount' of rewards as pending
func (drs *DelegRewardStore) addPendingRewards(delegator keys.Address, amount *balance.Amount, height int64) error {
	key := drs.getPendingRewardsKey(height, delegator)
	amt, err := drs.get(key)
	if err != nil {
		return err
	}

	err = drs.set(key, amt.Plus(*amount))
	return err
}

// Add an 'amount' of rewards as matured
func (drs *DelegRewardStore) addMaturedRewards(delegator keys.Address, amount *balance.Amount) error {
	key := drs.getMaturedRewardsKey(delegator)
	amt, err := drs.get(key)
	if err != nil {
		return err
	}

	err = drs.set(key, amt.Plus(*amount))
	return err
}
