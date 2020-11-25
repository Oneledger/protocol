package rewards

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/magiconair/properties/assert"
	db "github.com/tendermint/tm-db"
)

const (
	rewardsPrefix         = "rew"
	rewardsIntervalPrefix = "rewInt"
	rewardsAddrList       = "rewAddr"
	numPrivateKeys        = 6
	rewardInterval        = 15
	rewardIntervalNew     = 10
)

var (
	validatorList []keys.Address
	rewardStore   *RewardStore
	rewardOptions Options
)

func init() {
	//Setup Options
	rewardOptions.RewardInterval = rewardInterval

	//Create Validator Keys
	for i := 0; i < numPrivateKeys; i++ {
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()
		validatorList = append(validatorList, h.Address())
	}

	// Sort validator address ANSC
	sort.Slice(validatorList, func(i, j int) bool {
		return validatorList[i].String() < validatorList[j].String()
	})

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	rewardStore = NewRewardStore(rewardsPrefix, rewardsIntervalPrefix, rewardsAddrList, cs)
	rewardStore.SetOptions(&rewardOptions)
	rewardStore.WithState(cs)
}

func TestRewardStore_Set(t *testing.T) {
	rewardStore.SetWithHeight(validatorList[0], 1, balance.NewAmount(100))
	amt, _ := rewardStore.GetWithHeight(validatorList[0], rewardInterval-1)
	assert.Equal(t, amt, balance.NewAmount(100))

	rewardStore.SetWithHeight(validatorList[1], 1, balance.NewAmount(200))
	amt, _ = rewardStore.GetWithHeight(validatorList[1], rewardInterval-1)
	assert.Equal(t, amt, balance.NewAmount(200))

	rewardStore.SetWithHeight(validatorList[2], 1, balance.NewAmount(300))
	amt, _ = rewardStore.GetWithHeight(validatorList[2], rewardInterval-1)
	assert.Equal(t, amt, balance.NewAmount(300))
}

func TestRewardStore_AddToAddress(t *testing.T) {
	//Add amounts
	for i := 0; i < rewardInterval; i++ {
		err := rewardStore.AddToAddress(validatorList[0], int64(i), balance.NewAmount(1))
		assert.Equal(t, err, nil)
	}
	//Verify amount is the same for every version in the interval
	for i := 0; i < rewardInterval; i++ {
		amt, _ := rewardStore.GetWithHeight(validatorList[0], int64(i))
		assert.Equal(t, amt, balance.NewAmount(115))
	}

	//Add amounts
	for i := rewardInterval; i < 2*rewardInterval; i++ {
		err := rewardStore.AddToAddress(validatorList[0], int64(i), balance.NewAmount(1))
		assert.Equal(t, err, nil)
	}
	//Verify amount is the same for every version in the interval
	for i := rewardInterval; i < 2*rewardInterval; i++ {
		amt, _ := rewardStore.GetWithHeight(validatorList[0], int64(i))
		assert.Equal(t, amt, balance.NewAmount(15))
	}
}

func TestRewardStore_Iterate(t *testing.T) {
	rewardStore.State.Commit()
	var amts []*balance.Amount
	rewardStore.Iterate(validatorList[0], func(addr keys.Address, index int64, amt *balance.Amount) bool {
		amts = append(amts, amt)
		return false
	})
	assert.Equal(t, len(amts), 2)
	assert.Equal(t, amts[0], balance.NewAmount(115))
	assert.Equal(t, amts[1], balance.NewAmount(15))
}

func TestRewardStore_GetOptions(t *testing.T) {
	options := rewardStore.GetOptions()
	assert.Equal(t, *options, rewardOptions)
}

func TestRewardStore_Interval(t *testing.T) {
	//Loop through different heights and add rewards
	for i := 0; i < 100; i++ {

		_ = rewardStore.AddToAddress(validatorList[3], int64(i), balance.NewAmount(10))

		switch i {
		case 30:
			_ = rewardStore.SetInterval(30)
			rewardStore.SetOptions(&Options{
				RewardInterval:    10,
				RewardPoolAddress: "",
			})

		case 60:
			_ = rewardStore.SetInterval(60)
			rewardStore.SetOptions(&Options{
				RewardInterval:    5,
				RewardPoolAddress: "",
			})
		case 90:
			_ = rewardStore.SetInterval(90)
			rewardStore.SetOptions(&Options{
				RewardInterval:    1,
				RewardPoolAddress: "",
			})
		}

		rewardStore.State.Commit()
	}

	//Reset Interval back to 15 for validation.
	rewardStore.SetOptions(&Options{
		RewardInterval:    15,
		RewardPoolAddress: "",
	})
	//Validate whether the different intervals are being followed
	for i := 0; i < 100; i++ {
		if i < 30 {
			reward, _ := rewardStore.GetWithHeight(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(150))
		}
		if i < 60 && i >= 30 {
			rewardStore.SetOptions(&Options{
				RewardInterval:    10,
				RewardPoolAddress: "",
			})
			reward, _ := rewardStore.GetWithHeight(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(100))
		}
		if i >= 60 && i < 90 {
			rewardStore.SetOptions(&Options{
				RewardInterval:    5,
				RewardPoolAddress: "",
			})
			reward, _ := rewardStore.GetWithHeight(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(50))
		}
		if i >= 90 {
			rewardStore.SetOptions(&Options{
				RewardInterval:    1,
				RewardPoolAddress: "",
			})
			reward, _ := rewardStore.GetWithHeight(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(10))
		}
	}
}

func TestRewardStore_GetMaturedAmount(t *testing.T) {
	var matured *balance.Amount

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	rewardStore = NewRewardStore(rewardsPrefix, rewardsIntervalPrefix, rewardsAddrList, cs)
	rewardStore.SetOptions(&rewardOptions)
	rewardStore.WithState(cs)

	_ = rewardStore.AddToAddress(validatorList[5], int64(1), balance.NewAmount(10))

	for i := 0; i < 100; i++ {
		_ = rewardStore.AddToAddress(validatorList[4], int64(i), balance.NewAmount(10))
		matured, _ = rewardStore.GetMaturedAmount(validatorList[4], int64(i))

		if i >= 30 {
			assert.Equal(t, matured, balance.NewAmount(150))
		}

		rewardStore.State.Commit()
	}
}

func TestRewardStore_GetLastTwoChunks(t *testing.T) {
	amount, err := rewardStore.GetLastTwoChunks(validatorList[4])
	assert.Equal(t, err, nil)
	assert.Equal(t, amount, balance.NewAmount(250))

	amount, err = rewardStore.GetMaturedAmount(validatorList[5], 35)
	assert.Equal(t, err, nil)
	assert.Equal(t, amount, balance.NewAmount(10))
}

func TestRewardStore_GetLastTwoChunks2(t *testing.T) {
	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	rewardStore = NewRewardStore(rewardsPrefix, rewardsIntervalPrefix, rewardsAddrList, cs)
	rewardStore.SetOptions(&Options{
		RewardInterval:    1,
		RewardPoolAddress: "",
	})
	rewardStore.WithState(cs)

	for i := 0; i < 4; i++ {
		_ = rewardStore.AddToAddress(validatorList[5], int64(i), balance.NewAmount(10))

		if i < 1 {
			amount, _ := rewardStore.GetLastTwoChunks(validatorList[5])
			assert.Equal(t, amount, balance.NewAmount(10))
		}

		if i >= 2 {
			amount, _ := rewardStore.GetLastTwoChunks(validatorList[5])
			assert.Equal(t, amount, balance.NewAmount(20))
		}

		rewardStore.State.Commit()
	}
}

func TestRewardStore_IterateAddrList(t *testing.T) {
	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	rewardStore = NewRewardStore(rewardsPrefix, rewardsIntervalPrefix, rewardsAddrList, cs)
	rewardStore.SetOptions(&rewardOptions)
	rewardStore.WithState(cs)

	for _, addr := range validatorList {
		_ = rewardStore.AddToAddress(addr, 5, balance.NewAmount(19))
	}
	rewardStore.State.Commit()

	count := 0
	rewardStore.IterateAddrList(func(key keys.Address) bool {
		count++
		return false
	})
	assert.Equal(t, count, 6)
}

func TestRewardStore_DumpLoadState(t *testing.T) {
	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	rewardStore = NewRewardStore(rewardsPrefix, rewardsIntervalPrefix, rewardsAddrList, cs)
	rewardStore.SetOptions(&rewardOptions)
	rewardStore.WithState(cs)

	//Add amounts
	for i := 0; i < rewardInterval; i++ {
		err := rewardStore.AddToAddress(validatorList[0], int64(i), balance.NewAmount(1))
		assert.Equal(t, err, nil)
		err = rewardStore.AddToAddress(validatorList[1], int64(i), balance.NewAmount(2))
		assert.Equal(t, err, nil)
	}
	//Add amounts
	for i := rewardInterval; i < 2*rewardInterval; i++ {
		err := rewardStore.AddToAddress(validatorList[0], int64(i), balance.NewAmount(2))
		assert.Equal(t, err, nil)
		err = rewardStore.AddToAddress(validatorList[1], int64(i), balance.NewAmount(2))
		assert.Equal(t, err, nil)
	}
	// set new interval
	newOptions := Options{RewardInterval: rewardIntervalNew}
	err := rewardStore.UpdateOptions(2*rewardInterval-1, &newOptions)
	assert.Equal(t, err, nil)
	rewardStore.State.Commit()

	// state expected
	expected := NewRewardState()
	expected.Rewards = append(expected.Rewards, IntervalReward{
		Address: validatorList[0],
		Index:   1,
		Amount:  balance.NewAmount(1 * rewardInterval),
	})
	expected.Rewards = append(expected.Rewards, IntervalReward{
		Address: validatorList[0],
		Index:   2,
		Amount:  balance.NewAmount(2 * rewardInterval),
	})
	expected.Rewards = append(expected.Rewards, IntervalReward{
		Address: validatorList[1],
		Index:   1,
		Amount:  balance.NewAmount(2 * rewardInterval),
	})
	expected.Rewards = append(expected.Rewards, IntervalReward{
		Address: validatorList[1],
		Index:   2,
		Amount:  balance.NewAmount(2 * rewardInterval),
	})
	expected.Intervals = append(expected.Intervals, Interval{LastIndex: 1, LastHeight: 2})
	expected.AddrList = append(expected.AddrList, validatorList[0])
	expected.AddrList = append(expected.AddrList, validatorList[1])

	// prepare to dump
	dir, _ := os.Getwd()
	file := filepath.Join(dir, "genesis.json")
	writer, err := os.Create(file)
	assert.Equal(t, err, nil)
	defer func() { _ = os.Remove(file) }()
	state, err := rewardStore.dumpState()
	assert.Equal(t, err, nil)
	assert.Equal(t, state, expected)

	// dump to Genesis
	str, err := json.MarshalIndent(state, "", " ")
	assert.Equal(t, err, nil)
	_, err = writer.Write(str)
	assert.Equal(t, err, nil)
	err = writer.Close()
	assert.Equal(t, err, nil)

	// load from Genesis
	reader, err := os.Open(file)
	stateBytes, _ := ioutil.ReadAll(reader)
	assert.Equal(t, err, nil)
	stateLoaded := NewRewardState()
	err = json.Unmarshal(stateBytes, stateLoaded)
	assert.Equal(t, err, nil)
	assert.Equal(t, stateLoaded, expected)
}
