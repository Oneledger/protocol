package rewards

import (
	"strings"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type RewardStore struct {
	State  *storage.State
	szlr   serialize.Serializer
	prefix []byte

	rewardOptions Options
}

func NewRewardStore(prefix string, state *storage.State) *RewardStore {
	return &RewardStore{
		State:  state,
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
		prefix: storage.Prefix(prefix),
	}
}

func (rs *RewardStore) WithState(state *storage.State) {
	rs.State = state
}

func generateKey(address keys.Address, height int64, interval int64) (Key storage.StoreKey) {
	index := int64(height / interval)
	Key = storage.StoreKey(address.String() + storage.DB_PREFIX + string(index))
	return
}

func (rs *RewardStore) Get(address keys.Address, height int64) (amount *balance.Amount, err error) {
	key := append(rs.prefix, generateKey(address, height, rs.rewardOptions.RewardInterval)...)

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

	key := append(rs.prefix, generateKey(address, height, rs.rewardOptions.RewardInterval)...)
	err = rs.State.Set(key, data)
	return err
}

func (rs *RewardStore) AddToAddress(address keys.Address, height int64, amount *balance.Amount) error {
	prevAmount, err := rs.Get(address, height)
	if err != nil {
		return err
	}
	newAmount := prevAmount.Plus(*amount)
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

func (rs *RewardStore) SetOptions(options Options) {
	rs.rewardOptions = options
}

func (rs *RewardStore) GetOptions() Options {
	return rs.rewardOptions
}
