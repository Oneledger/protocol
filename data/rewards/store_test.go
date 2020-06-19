package rewards

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/magiconair/properties/assert"
	db "github.com/tendermint/tm-db"
	"testing"
)

const (
	rewardsPrefix  = "rew"
	numPrivateKeys = 3
	rewardInterval = 15
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

	rewardStore = NewRewardStore(rewardsPrefix, cs)
	rewardStore.SetOptions(rewardOptions)
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
	assert.Equal(t, options, rewardOptions)
}
