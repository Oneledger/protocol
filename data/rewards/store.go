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
	prefixAddrList  []byte

	rewardOptions *Options
}

func NewRewardStore(prefix string, intervalPrefix string, addrListPrefix string, state *storage.State) *RewardStore {
	return &RewardStore{
		State:           state,
		szlr:            serialize.GetSerializer(serialize.PERSISTENT),
		prefix:          storage.Prefix(prefix),
		prefixIntervals: storage.Prefix(intervalPrefix),
		prefixAddrList:  storage.Prefix(addrListPrefix),
	}
}

func (rs *RewardStore) GetState() *storage.State {
	return rs.State
}

func (rs *RewardStore) WithState(state *storage.State) *RewardStore {
	rs.State = state
	return rs
}

func (rs *RewardStore) generateMaturedKey(address keys.Address, height int64, interval int64) (Key storage.StoreKey) {
	Key = nil

	lastInterval := rs.GetInterval(height)
	index := lastInterval.LastIndex + int64((height-lastInterval.LastHeight)/interval) + 1
	if index >= 2 {
		Key = storage.StoreKey(address.String() + storage.DB_PREFIX + strconv.FormatInt(index-2, 10))
	}
	return
}

func (rs *RewardStore) generatePreviousKey(address keys.Address, height int64, interval int64) (Key storage.StoreKey) {
	Key = nil
	lastInterval := rs.GetInterval(height)
	index := lastInterval.LastIndex + int64((height-lastInterval.LastHeight)/interval)
	Key = storage.StoreKey(address.String() + storage.DB_PREFIX + strconv.FormatInt(index, 10))
	return
}

func (rs *RewardStore) generateKey(address keys.Address, height int64, interval int64) (Key storage.StoreKey) {
	lastInterval := rs.GetInterval(height)
	index := lastInterval.LastIndex + int64((height-lastInterval.LastHeight)/interval) + 1
	Key = storage.StoreKey(address.String() + storage.DB_PREFIX + strconv.FormatInt(index, 10))
	return
}

func (rs *RewardStore) Get(key storage.StoreKey) (amount *balance.Amount, err error) {
	data, err := rs.State.Get(key)
	amount = balance.NewAmount(0)
	if len(data) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(data, amount)
	return
}

func (rs *RewardStore) GetWithHeight(address keys.Address, height int64) (amount *balance.Amount, err error) {
	key := append(rs.prefix, rs.generateKey(address, height, rs.rewardOptions.RewardInterval)...)
	data, err := rs.State.Get(key)
	amount = balance.NewAmount(0)
	if len(data) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(data, amount)
	return
}

func (rs *RewardStore) SetWithHeight(address keys.Address, height int64, amount *balance.Amount) error {
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

	key := append(rs.prefixIntervals, storage.StoreKey(strconv.FormatInt(height, 10))...)
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

func (rs *RewardStore) IterateAddrList(fn func(key keys.Address) bool) bool {
	return rs.State.IterateRange(
		append(rs.prefixAddrList),
		storage.Rangefix(string(rs.prefixAddrList)),
		true,
		func(key, value []byte) bool {
			address := keys.Address{}
			arr := strings.Split(string(key), storage.DB_PREFIX)
			keyStr := arr[len(arr)-1]

			err := address.UnmarshalText([]byte(keyStr))
			if err != nil {
				return true
			}
			return fn(address)
		},
	)
}

func (rs *RewardStore) AddAddressToList(address keys.Address) {
	key := append(rs.prefixAddrList, storage.StoreKey(address.String())...)
	if !rs.State.Exists(key) {
		_ = rs.State.Set(key, []byte("active"))
	}
}

func (rs *RewardStore) AddToAddress(address keys.Address, height int64, amount *balance.Amount) error {
	prevAmount, err := rs.GetWithHeight(address, height)
	newAmount := amount
	if err == nil {
		newAmount = prevAmount.Plus(*amount)
	}
	rs.AddAddressToList(address)

	return rs.SetWithHeight(address, height, newAmount)
}

//Iterate through all reward records for a given Address
func (rs *RewardStore) Iterate(addr keys.Address, fn func(addr keys.Address, index int64, amt *balance.Amount) bool) bool {
	return rs.State.IterateRange(
		append(rs.prefix, addr.String()...),
		storage.Rangefix(string(append(rs.prefix, addr.String()...))),
		true,
		func(key, value []byte) bool {
			amt := balance.NewAmount(0)
			err := rs.szlr.Deserialize(value, amt)
			if err != nil {
				return true
			}

			// key in format "validator_index"
			keyStr := string(key[len(rs.prefix):])
			vi := strings.Split(keyStr, storage.DB_PREFIX)
			if len(vi) != 2 {
				return true
			}
			// parse address and index
			address := keys.Address{}
			err = address.UnmarshalText([]byte(vi[0]))
			if err != nil {
				return true
			}
			index, err := strconv.ParseInt(vi[1], 10, 64)
			if err != nil {
				return true
			}
			return fn(address, index, amt)
		},
	)
}

func (rs *RewardStore) GetMaturedAmount(address keys.Address, height int64) (*balance.Amount, error) {
	key := append(rs.prefix, rs.generateMaturedKey(address, height, rs.rewardOptions.RewardInterval)...)
	return rs.Get(key)
}

func (rs *RewardStore) SetOptions(options *Options) {
	rs.rewardOptions = options
}

func (rs *RewardStore) GetOptions() *Options {
	return rs.rewardOptions
}
func (rs *RewardStore) UpdateOptions(height int64, options *Options) error {
	if rs.rewardOptions.RewardInterval != options.RewardInterval {
		err := rs.SetInterval(height)
		if err != nil {
			return err
		}
	}
	rs.SetOptions(options)
	return nil
}

func (rs *RewardStore) GetLastTwoChunks(address keys.Address) (*balance.Amount, error) {
	amount := balance.NewAmount(0)
	previousAmount := balance.NewAmount(0)
	currentKey := append(rs.prefix, rs.generateKey(address, rs.State.Version(), rs.rewardOptions.RewardInterval)...)
	previousKey := append(rs.prefix, rs.generatePreviousKey(address, rs.State.Version(), rs.rewardOptions.RewardInterval)...)
	amount, err := rs.Get(currentKey)
	if err != nil {
		return nil, err
	}

	previousAmount, err = rs.Get(previousKey)
	if err != nil {
		return amount, nil
	}

	return amount.Plus(*previousAmount), nil
}

//-----------------------------Dump/Load chain state
//
type IntervalReward struct {
	Address keys.Address    `json:"address"`
	Index   int64           `json:"index"`
	Amount  *balance.Amount `json:"amount"`
}

type RewardState struct {
	Rewards   []IntervalReward `json:"rewards"`
	Intervals []Interval       `json:"intervals"`
	AddrList  []keys.Address   `json:"addrList"`
}

func NewRewardState() *RewardState {
	return &RewardState{
		Rewards:   []IntervalReward{},
		Intervals: []Interval{},
		AddrList:  []keys.Address{},
	}
}

func (rs *RewardStore) dumpState() (state *RewardState, err error) {
	// dump rewards
	state = NewRewardState()
	rs.Iterate(keys.Address{}, func(addr keys.Address, index int64, amt *balance.Amount) bool {
		reward := IntervalReward{
			Address: addr,
			Index:   index,
			Amount:  amt,
		}
		state.Rewards = append(state.Rewards, reward)
		return false
	})
	// dump intervals
	rs.State.IterateRange(
		rs.prefixIntervals,
		storage.Rangefix(string(rs.prefixIntervals)),
		true,
		func(key, value []byte) bool {
			interval := &Interval{}
			err = rs.szlr.Deserialize(value, interval)
			if err != nil {
				return true
			}
			state.Intervals = append(state.Intervals, *interval)
			return false
		},
	)
	// dump validator address list
	rs.IterateAddrList(func(addr keys.Address) bool {
		state.AddrList = append(state.AddrList, addr)
		return false
	})
	return
}

func (rs *RewardStore) loadState(state *RewardState) error {
	// load rewards
	for _, reward := range state.Rewards {
		key := append(rs.prefix, storage.StoreKey(reward.Address.String()+storage.DB_PREFIX+string(reward.Index))...)
		data, err := rs.szlr.Serialize(reward.Amount)
		if err != nil {
			return err
		}
		err = rs.State.Set(key, data)
		if err != nil {
			return err
		}
	}
	// load intervals
	for _, interval := range state.Intervals {
		key := append(rs.prefixIntervals, storage.StoreKey(strconv.FormatInt(interval.LastHeight, 10))...)
		data, err := rs.szlr.Serialize(interval)
		if err != nil {
			return err
		}
		err = rs.State.Set(key, data)
		if err != nil {
			return err
		}
	}
	// dump validator address list
	for _, addr := range state.AddrList {
		rs.AddAddressToList(addr)
	}
	return nil
}
