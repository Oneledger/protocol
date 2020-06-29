package rewards

import (
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type RewardStore struct {
	State           *storage.State
	szlr            serialize.Serializer
	prefix          []byte
	prefixIntervals []byte

	rewardOptions *Options
}

func NewRewardStore(prefix string, intervalPrefix string, state *storage.State) *RewardStore {
	return &RewardStore{
		State:           state,
		szlr:            serialize.GetSerializer(serialize.PERSISTENT),
		prefix:          storage.Prefix(prefix),
		prefixIntervals: storage.Prefix(intervalPrefix),
	}
}

func (rs *RewardStore) WithState(state *storage.State) *RewardStore {
	rs.State = state
	return rs
}

func (rs *RewardStore) generateKey(address keys.Address, height int64, interval int64) (Key storage.StoreKey) {
	lastInterval := rs.GetInterval(height)
	index := lastInterval.LastIndex + int64((height-lastInterval.LastHeight)/interval) + 1
	Key = storage.StoreKey(address.String() + storage.DB_PREFIX + string(index))
	return
}

func (rs *RewardStore) Get(address keys.Address, height int64) (amount *balance.Amount, err error) {
	key := append(rs.prefix, rs.generateKey(address, height, rs.rewardOptions.RewardInterval)...)
	data, err := rs.State.Get(key)
	amount = balance.NewAmount(0)
	if len(data) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(data, amount)
	return
}

func (rs *RewardStore) Set(address keys.Address, height int64, amount *balance.Amount) error {
	data, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(amount)
	if err != nil {
		return err
	}
	key := append(rs.prefix, rs.generateKey(address, height, rs.rewardOptions.RewardInterval)...)
	err = rs.State.Set(key, data)
	return err
}

func (rs *RewardStore) SetInterval(height int64) error {
	//Get last height to find the difference
	LastInterval := rs.GetInterval(height)
	diff := height - LastInterval.LastHeight

	key := append(rs.prefixIntervals, storage.StoreKey(storage.DB_PREFIX+strconv.FormatInt(height, 10))...)
	interval := &Interval{
		LastIndex:  LastInterval.LastIndex + int64(diff/rs.rewardOptions.RewardInterval),
		LastHeight: height,
	}

	data, err := rs.szlr.Serialize(interval)
	if err != nil {
		return err
	}
	return rs.State.Set(key, data)
}

func (rs *RewardStore) GetInterval(height int64) *Interval {
	maxHeight := int64(0)
	lastInterval := &Interval{
		LastIndex:  0,
		LastHeight: 0,
	}

	//Iterate to find closest Interval where LastHeight <= height
	rs.State.IterateRange(
		rs.prefixIntervals,
		storage.Rangefix(string(rs.prefixIntervals)),
		true,
		func(key, value []byte) bool {
			interval := &Interval{}

			err := rs.szlr.Deserialize(value, interval)
			if err != nil {
				return true
			}

			if interval.LastHeight > maxHeight && interval.LastHeight <= height {
				lastInterval = interval
				maxHeight = interval.LastHeight
			}

			return false
		},
	)

	//If there aren't any stored intervals then return default value
	return lastInterval
}

func (rs *RewardStore) AddToAddress(address keys.Address, height int64, amount *balance.Amount) error {
	prevAmount, err := rs.Get(address, height)
	newAmount := amount
	if err == nil {
		newAmount = prevAmount.Plus(*amount)
	}

	return rs.Set(address, height, newAmount)
}

//Iterate through all reward records for a given Address
func (rs *RewardStore) Iterate(addr keys.Address, fn func(c string, amt balance.Amount) bool) bool {
	return rs.State.IterateRange(
		append(rs.prefix, addr.String()...),
		storage.Rangefix(string(append(rs.prefix, addr.String()...))),
		true,
		func(key, value []byte) bool {
			amt := balance.NewAmount(0)

			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, amt)
			if err != nil {
				return true
			}

			arr := strings.Split(string(key), storage.DB_PREFIX)
			return fn(arr[len(arr)-1], *amt)
		},
	)
}

func (rs *RewardStore) SetOptions(options *Options) {
	rs.rewardOptions = options
}

func (rs *RewardStore) GetOptions() *Options {
	return rs.rewardOptions
}
