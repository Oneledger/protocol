package network_delegation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
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

func (drs *DelegRewardStore) GetState() *storage.State {
	return drs.state
}

// Add rewards balance
func (drs *DelegRewardStore) AddRewardsBalance(delegator keys.Address, amount *balance.Amount) error {
	key := drs.getRewardsBalanceKey(delegator)
	amt, err := drs.get(key)
	if err != nil {
		return err
	}

	keyTotal := drs.getTotalRewardsKey()
	amtTotal, err := drs.get(keyTotal)
	if err != nil {
		return err
	}

	err = drs.set(key, amt.Plus(*amount))
	err = drs.set(keyTotal, amtTotal.Plus(*amount))
	return err
}

// Get rewards balance
func (drs *DelegRewardStore) GetRewardsBalance(delegator keys.Address) (amt *balance.Amount, err error) {
	key := drs.getRewardsBalanceKey(delegator)
	amt, err = drs.get(key)
	return
}

// Get total rewards
func (drs *DelegRewardStore) GetTotalRewards() (amt *balance.Amount, err error) {
	key := drs.getTotalRewardsKey()
	amt, err = drs.get(key)
	return
}

// Deducts an 'amount' of rewards from rewards balance
func (drs *DelegRewardStore) MinusRewardsBalance(delegator keys.Address, amount *balance.Amount) error {
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

// Initiate a withdrawal of an 'amount' of rewards
func (drs *DelegRewardStore) Withdraw(delegator keys.Address, amount *balance.Amount, matureHeight int64) error {
	err := drs.MinusRewardsBalance(delegator, amount)
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
	pdRewards = &DelegPendingRewards{Address: delegator, Rewards: []*PendingRewards{}}
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

// Set pending rewards for a certain height
func (drs *DelegRewardStore) SetPendingRewards(delegator keys.Address, amount *balance.Amount, height int64) error {
	key := drs.getPendingRewardsKey(height, delegator)

	err := drs.set(key, amount)
	return err
}

func (drs *DelegRewardStore) IterateActiveRewards(fn func(addr *keys.Address, amt *balance.Amount) bool) bool {
	balanceKey := string(storage.Prefix("balance"))
	return drs.iterate(balanceKey, func(delegator keys.Address, amt *balance.Amount) bool {
		return fn(&delegator, amt)
	})
}

// below is removed since finalize withdraw rewards logic is moved to block beginner, OLP-1266
//// Mature, if any, all delegators' pending withdrawal at a specific height
//func (drs *DelegRewardStore) MaturePendingRewards(height int64) (event abciTypes.Event, any bool) {
//	event.Type = "deleg_rewards"
//	event.Attributes = append(event.Attributes, kv.Pair{
//		Key:   []byte("height"),
//		Value: []byte(strconv.FormatInt(height, 10)),
//	})
//
//	drs.IteratePD(height, func(delegator keys.Address, amt *balance.Amount) bool {
//		// clear pending amount
//		key := drs.getPendingRewardsKey(height, delegator)
//		err := drs.set(key, balance.AmtZero)
//		if err != nil {
//			return true
//		}
//		// increase matured amount
//		if !amt.Equals(balance.AmtZero) {
//			err = drs.addMaturedRewards(delegator, amt)
//			event.Attributes = append(event.Attributes, kv.Pair{
//				Key:   []byte(delegator.String()),
//				Value: []byte(amt.String()),
//			})
//		}
//		return err != nil
//	})
//
//	any = len(event.Attributes) > 1
//	return
//}
//
//// Get matured(finalizable) rewards
//func (drs *DelegRewardStore) GetMaturedRewards(delegator keys.Address) (amt *balance.Amount, err error) {
//	key := drs.getMaturedRewardsKey(delegator)
//	amt, err = drs.get(key)
//	return
//}
//
//// Finalize(deduct) an 'amount' of matured rewards
//func (drs *DelegRewardStore) Finalize(delegator keys.Address, amount *balance.Amount) error {
//	key := drs.getMaturedRewardsKey(delegator)
//	amt, err := drs.get(key)
//	if err != nil {
//		return err
//	}
//
//	result, err := amt.Minus(*amount)
//	if err != nil {
//		return err
//	}
//
//	err = drs.set(key, result)
//	return err
//}
//
//// Mature, if any, all delegators' pending withdrawal at a specific height
//func (drs *DelegRewardStore) MaturePendingRewards(height int64) (event abciTypes.Event, any bool) {
//	event.Type = "deleg_rewards"
//	event.Attributes = append(event.Attributes, kv.Pair{
//		Key:   []byte("height"),
//		Value: []byte(strconv.FormatInt(height, 10)),
//	})
//
//	drs.IteratePD(height, func(delegator keys.Address, amt *balance.Amount) bool {
//		// clear pending amount
//		key := drs.getPendingRewardsKey(height, delegator)
//		err := drs.set(key, balance.AmtZero)
//		if err != nil {
//			return true
//		}
//		// increase matured amount
//		if !amt.Equals(balance.AmtZero) {
//			err = drs.addMaturedRewards(delegator, amt)
//			event.Attributes = append(event.Attributes, kv.Pair{
//				Key:   []byte(delegator.String()),
//				Value: []byte(amt.String()),
//			})
//		}
//		return err != nil
//	})
//
//	any = len(event.Attributes) > 1
//	return
//}
//
//// Get matured(finalizable) rewards
//func (drs *DelegRewardStore) GetMaturedRewards(delegator keys.Address) (amt *balance.Amount, err error) {
//	key := drs.getMaturedRewardsKey(delegator)
//	amt, err = drs.get(key)
//	return
//}
//
//// Finalize(deduct) an 'amount' of matured rewards
//func (drs *DelegRewardStore) Finalize(delegator keys.Address, amount *balance.Amount) error {
//	key := drs.getMaturedRewardsKey(delegator)
//	amt, err := drs.get(key)
//	if err != nil {
//		return err
//	}
//
//	result, err := amt.Minus(*amount)
//	if err != nil {
//		return err
//	}
//
//	err = drs.set(key, result)
//	return err
//}

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
func (drs *DelegRewardStore) IteratePD(height int64, fn func(delegator keys.Address, amt *balance.Amount) bool) (stopped bool) {
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

func (drs *DelegRewardStore) IterateAllPD(fn func(height int64, delegator keys.Address, amt *balance.Amount) bool) bool {
	pfxStr := fmt.Sprintf("%spending_", string(drs.prefix))
	prefix := storage.StoreKey(pfxStr)
	return drs.state.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			keyArr := strings.Split(string(key), storage.DB_PREFIX)
			addr := &keys.Address{}
			err := addr.UnmarshalText([]byte(keyArr[len(keyArr)-1]))
			if err != nil {
				return true
			}
			height, err := strconv.ParseInt(keyArr[len(keyArr)-2], 10, 64)
			if err != nil {
				return true
			}

			amt := balance.NewAmount(0)
			err = drs.szlr.Deserialize(value, amt)
			if err != nil {
				logger.Error("failed to deserialize delegator pending rewards amount")
				return true
			}
			return fn(height, *addr, amt)
		},
	)
}

// Key for delegator rewards balance
func (drs *DelegRewardStore) getRewardsBalanceKey(delegator keys.Address) []byte {
	key := fmt.Sprintf("%sbalance_%s", string(drs.prefix), delegator)
	return storage.StoreKey(key)
}

// Key for total rewards
func (drs *DelegRewardStore) getTotalRewardsKey() []byte {
	key := fmt.Sprintf("%stotal_rewards", string(drs.prefix))
	return storage.StoreKey(key)
}

// Key for delegator withdrawn but unmatured amount
func (drs *DelegRewardStore) getPendingRewardsKey(height int64, delegator keys.Address) []byte {
	key := fmt.Sprintf("%spending_%d_%s", string(drs.prefix), height, delegator)
	return storage.StoreKey(key)
}

//// Key for delegator withdrawn and matured amount
//func (drs *DelegRewardStore) getMaturedRewardsKey(delegator keys.Address) []byte {
//	key := fmt.Sprintf("%smatured_%s", string(drs.prefix), delegator)
//	return storage.StoreKey(key)
//}

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

//// Add an 'amount' of rewards as matured
//func (drs *DelegRewardStore) addMaturedRewards(delegator keys.Address, amount *balance.Amount) error {
//	key := drs.getMaturedRewardsKey(delegator)
//	amt, err := drs.get(key)
//	if err != nil {
//		return err
//	}
//
//	err = drs.set(key, amt.Plus(*amount))
//	return err
//}

//---------------------------------- Save Delegation Reward Store -------------------------------------
func (drs *DelegRewardStore) SaveState() (*RewardState, bool) {
	balanceKey := string(storage.Prefix("balance"))
	//matureKey := string(storage.Prefix("mature"))

	//Populate Current Balances
	var balanceList []Reward
	drs.iterate(balanceKey, func(delegator keys.Address, amt *balance.Amount) bool {
		reward := Reward{
			Amount:  amt,
			Address: delegator,
		}
		balanceList = append(balanceList, reward)
		return false
	})

	//var matureList []Reward
	////Populate Mature Balances
	//drs.iterate(matureKey, func(delegator keys.Address, amt *balance.Amount) bool {
	//	reward := Reward{
	//		Amount:  amt,
	//		Address: delegator,
	//	}
	//	matureList = append(matureList, reward)
	//	return false
	//})

	var pendingList []PendingReward
	//Populate Pending Balances
	drs.IterateAllPD(func(height int64, delegator keys.Address, amt *balance.Amount) bool {
		pendingRew := PendingReward{
			Amount:  amt,
			Address: delegator,
			Height:  height,
		}
		pendingList = append(pendingList, pendingRew)
		return false
	})

	return &RewardState{
		BalanceList: balanceList,
		//MatureList:  matureList,
		PendingList: pendingList,
	}, true
}

//__________________________________ Load Delegation Reward Store -------------------------------------
func (drs *DelegRewardStore) LoadState(state *RewardState) error {
	for _, v := range state.BalanceList {
		err := drs.AddRewardsBalance(v.Address, v.Amount)
		if err != nil {
			return err
		}
	}
	//for _, v := range state.MatureList {
	//	err := drs.addMaturedRewards(v.Address, v.Amount)
	//	if err != nil {
	//		return err
	//	}
	//}
	for _, v := range state.PendingList {
		height := v.Height - drs.state.Version()
		if height <= 0 {
			continue
		}
		err := drs.addPendingRewards(v.Address, v.Amount, height)
		if err != nil {
			return err
		}
	}
	return nil
}
