package rewards

import (
	"fmt"
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
	numPrivateKeys        = 5
	rewardInterval        = 15
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

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	rewardStore = NewRewardStore(rewardsPrefix, rewardsIntervalPrefix, cs)
	rewardStore.SetOptions(&rewardOptions)
	rewardStore.WithState(cs)
}

func TestRewardStore_Set(t *testing.T) {
	rewardStore.Set(validatorList[0], 1, balance.NewAmount(100))
	amt, _ := rewardStore.Get(validatorList[0], rewardInterval-1)
	assert.Equal(t, amt, balance.NewAmount(100))

	rewardStore.Set(validatorList[1], 1, balance.NewAmount(200))
	amt, _ = rewardStore.Get(validatorList[1], rewardInterval-1)
	assert.Equal(t, amt, balance.NewAmount(200))

	rewardStore.Set(validatorList[2], 1, balance.NewAmount(300))
	amt, _ = rewardStore.Get(validatorList[2], rewardInterval-1)
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
		amt, _ := rewardStore.Get(validatorList[0], int64(i))
		assert.Equal(t, amt, balance.NewAmount(115))
	}

	//Add amounts
	for i := rewardInterval; i < 2*rewardInterval; i++ {
		err := rewardStore.AddToAddress(validatorList[0], int64(i), balance.NewAmount(1))
		assert.Equal(t, err, nil)
	}
	//Verify amount is the same for every version in the interval
	for i := rewardInterval; i < 2*rewardInterval; i++ {
		amt, _ := rewardStore.Get(validatorList[0], int64(i))
		assert.Equal(t, amt, balance.NewAmount(15))
	}
}

func TestRewardStore_Iterate(t *testing.T) {
	rewardStore.State.Commit()
	var amts []*balance.Amount
	rewardStore.Iterate(validatorList[0], func(c string, amt balance.Amount) bool {
		amts = append(amts, &amt)
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
			reward, _ := rewardStore.Get(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(150))
		}
		if i < 60 && i >= 30 {
			rewardStore.SetOptions(&Options{
				RewardInterval:    10,
				RewardPoolAddress: "",
			})
			reward, _ := rewardStore.Get(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(100))
		}
		if i >= 60 && i < 90 {
			rewardStore.SetOptions(&Options{
				RewardInterval:    5,
				RewardPoolAddress: "",
			})
			reward, _ := rewardStore.Get(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(50))
		}
		if i >= 90 {
			rewardStore.SetOptions(&Options{
				RewardInterval:    1,
				RewardPoolAddress: "",
			})
			reward, _ := rewardStore.Get(validatorList[3], int64(i))
			assert.Equal(t, reward, balance.NewAmount(10))
		}
	}
}

func TestRewardStore_GetMaturedAmount(t *testing.T) {
	var matured *balance.Amount

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	rewardStore = NewRewardStore(rewardsPrefix, rewardsIntervalPrefix, cs)
	rewardStore.SetOptions(&rewardOptions)
	rewardStore.WithState(cs)

	for i := 0; i < 100; i++ {
		_ = rewardStore.AddToAddress(validatorList[4], int64(i), balance.NewAmount(10))
		matured, _ = rewardStore.GetMaturedAmount(validatorList[4], int64(i))

		fmt.Println("matured amount: ", matured, "height: ", i)

		rewardStore.State.Commit()
	}
}
