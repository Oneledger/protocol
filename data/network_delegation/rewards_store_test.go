package network_delegation

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	memDb      db.DB
	storeRwz   *DelegRewardStore
	cs         *storage.State
	delegators []keys.Address
	zero       *balance.Amount
	amt1       *balance.Amount
	amt2       *balance.Amount
	amt3       *balance.Amount
	draw1      *balance.Amount
	draw2      *balance.Amount
	draw3      *balance.Amount
	draw4      *balance.Amount
	draw5      *balance.Amount

	delegOpt = &Options{
		RewardsMaturityTime: 2,
	}
)

func setup() {
	fmt.Println("####### Testing delegator rewards store #######")
	memDb = db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("chainstate", memDb))
	storeRwz = NewDelegRewardStore("delegRwz", cs)
	setupVariables()
}

func genDelegAddress() keys.Address {
	pub, _, _ := keys.NewKeyPairFromTendermint()
	h, _ := pub.GetHandler()
	return h.Address()
}

func setupVariables() {
	// generates and sorts delegator addresses ANSC
	for i := 0; i < 2; i++ {
		delegators = append(delegators, genDelegAddress())
	}
	sort.Slice(delegators, func(i, j int) bool {
		return delegators[i].String() < delegators[j].String()
	})

	// some amounts
	zero = balance.NewAmount(0)
	amt1 = balance.NewAmount(100)
	amt2 = balance.NewAmount(200)
	amt3 = balance.NewAmount(377)
	draw1 = balance.NewAmount(163)
	draw2 = balance.NewAmount(17)
	draw3 = balance.NewAmount(20)
	draw4 = balance.NewAmount(77)
	draw5 = balance.NewAmount(50)
}

func TestNewDelegRewardStore(t *testing.T) {
	setup()
	balance, err := storeRwz.GetRewardsBalance(delegators[0])
	assert.Nil(t, err)
	pending, err := storeRwz.GetPendingRewards(delegators[0], 1, delegOpt.RewardsMaturityTime)
	assert.Nil(t, err)
	matured, err := storeRwz.GetMaturedRewards(delegators[0])
	assert.Nil(t, err)
	assert.Equal(t, zero, balance)
	assert.EqualValues(t, DelegPendingRewards{Address: delegators[0], Rewards: []*PendingRewards{}}, *pending)
	assert.Equal(t, zero, matured)
}

func TestDelegRewardStore_AddGetRewardsBalance(t *testing.T) {
	setup()
	err := storeRwz.AddRewardsBalance(delegators[0], amt1)
	assert.Nil(t, err)
	balance, err := storeRwz.GetRewardsBalance(delegators[0])
	assert.Nil(t, err)
	assert.Equal(t, balance, amt1)

	err = storeRwz.AddRewardsBalance(delegators[0], amt2)
	assert.Nil(t, err)
	balance, err = storeRwz.GetRewardsBalance(delegators[0])
	assert.Nil(t, err)
	assert.Equal(t, balance, amt1.Plus(*amt2))
}

func TestDelegRewardStore_Withdraw(t *testing.T) {
	setup()
	curHeight := int64(8)
	storeRwz.AddRewardsBalance(delegators[0], amt1)
	storeRwz.AddRewardsBalance(delegators[0], amt2)

	err := storeRwz.Withdraw(delegators[0], draw1, curHeight+delegOpt.RewardsMaturityTime)
	assert.Nil(t, err)
	balance, err := storeRwz.GetRewardsBalance(delegators[0])
	assert.Nil(t, err)
	expected, _ := amt1.Plus(*amt2).Minus(*draw1)
	assert.Equal(t, balance, expected)
}

func TestDelegRewardStore_GetPendingRewards(t *testing.T) {
	setup()
	curHeight := int64(8)
	storeRwz.AddRewardsBalance(delegators[0], amt1)
	storeRwz.AddRewardsBalance(delegators[0], amt2)

	storeRwz.Withdraw(delegators[0], draw1, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw2, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw3, curHeight+delegOpt.RewardsMaturityTime+1)
	pending, err := storeRwz.GetPendingRewards(delegators[0], curHeight+1, delegOpt.RewardsMaturityTime+1)
	assert.Nil(t, err)
	expected := &DelegPendingRewards{Address: delegators[0]}
	expected.Rewards = append(expected.Rewards, &PendingRewards{
		Height: curHeight + delegOpt.RewardsMaturityTime,
		Amount: *draw1.Plus(*draw2),
	})
	expected.Rewards = append(expected.Rewards, &PendingRewards{
		Height: curHeight + delegOpt.RewardsMaturityTime + 1,
		Amount: *draw3,
	})
	assert.Equal(t, *pending, *expected)
}

func TestDelegRewardStore_MaturePendingRewards(t *testing.T) {
	setup()
	curHeight := int64(8)
	storeRwz.AddRewardsBalance(delegators[0], amt1)
	storeRwz.AddRewardsBalance(delegators[0], amt2)
	storeRwz.AddRewardsBalance(delegators[1], amt3)

	storeRwz.Withdraw(delegators[0], draw1, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw2, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw3, curHeight+delegOpt.RewardsMaturityTime+1)
	storeRwz.Withdraw(delegators[1], draw4, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[1], draw5, curHeight+delegOpt.RewardsMaturityTime+1)
	storeRwz.state.Commit()
	storeRwz.MaturePendingRewards(curHeight + delegOpt.RewardsMaturityTime)

	// check each delegator's pending
	pending1, err := storeRwz.GetPendingRewards(delegators[0], curHeight+1, delegOpt.RewardsMaturityTime+1)
	assert.Nil(t, err)
	expected1 := &DelegPendingRewards{Address: delegators[0]}
	expected1.Rewards = append(expected1.Rewards, &PendingRewards{
		Height: curHeight + delegOpt.RewardsMaturityTime + 1,
		Amount: *draw3,
	})
	assert.Equal(t, *pending1, *expected1)

	pending2, err := storeRwz.GetPendingRewards(delegators[1], curHeight+1, delegOpt.RewardsMaturityTime+1)
	assert.Nil(t, err)
	expected2 := &DelegPendingRewards{Address: delegators[1]}
	expected2.Rewards = append(expected2.Rewards, &PendingRewards{
		Height: curHeight + delegOpt.RewardsMaturityTime + 1,
		Amount: *draw5,
	})
	assert.Equal(t, *pending2, *expected2)

	// check each delegator's matured rewards
	matured1, err := storeRwz.GetMaturedRewards(delegators[0])
	assert.Nil(t, err)
	maturedExp1 := draw1.Plus(*draw2)
	assert.Equal(t, *matured1, *maturedExp1)

	matured2, err := storeRwz.GetMaturedRewards(delegators[1])
	assert.Nil(t, err)
	maturedExp2 := *draw4
	assert.Equal(t, *matured2, maturedExp2)
}

func TestDelegRewardStore_Finalize(t *testing.T) {
	setup()
	curHeight := int64(8)
	storeRwz.AddRewardsBalance(delegators[0], amt1)
	storeRwz.AddRewardsBalance(delegators[0], amt2)

	storeRwz.Withdraw(delegators[0], draw1, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw2, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw3, curHeight+delegOpt.RewardsMaturityTime+1)
	storeRwz.state.Commit()
	storeRwz.MaturePendingRewards(curHeight + delegOpt.RewardsMaturityTime)

	// finalizes a part of matured rewards
	err := storeRwz.Finalize(delegators[0], draw1)
	assert.Nil(t, err)
	matured, err := storeRwz.GetMaturedRewards(delegators[0])
	assert.Nil(t, err)
	assert.Equal(t, *matured, *draw2)

	// finalizes more than matured rewards
	err = storeRwz.Finalize(delegators[0], draw1.Plus(*draw2))
	assert.NotNil(t, err)

	// finalizes all matured rewards
	err = storeRwz.Finalize(delegators[0], draw2)
	assert.Nil(t, err)
	matured, err = storeRwz.GetMaturedRewards(delegators[0])
	assert.Nil(t, err)
	assert.Equal(t, *matured, balance.AmtZero)
}

func TestGetTotalRewards(t *testing.T) {
	setup()
	curHeight := int64(8)
	storeRwz.AddRewardsBalance(delegators[0], amt1)
	storeRwz.AddRewardsBalance(delegators[0], amt2)
	storeRwz.AddRewardsBalance(delegators[1], amt2)
	storeRwz.AddRewardsBalance(delegators[1], amt3)

	storeRwz.Withdraw(delegators[0], draw1, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw2, curHeight+delegOpt.RewardsMaturityTime)
	storeRwz.Withdraw(delegators[0], draw3, curHeight+delegOpt.RewardsMaturityTime+1)
	storeRwz.Withdraw(delegators[1], draw3, curHeight+delegOpt.RewardsMaturityTime+1)

	storeRwz.state.Commit()
	storeRwz.MaturePendingRewards(curHeight + delegOpt.RewardsMaturityTime)

	// finalizes all matured rewards for delegator[0]
	err := storeRwz.Finalize(delegators[0], draw1.Plus(*draw2))
	assert.Nil(t, err)
	matured, err := storeRwz.GetMaturedRewards(delegators[0])
	assert.Nil(t, err)
	assert.Equal(t, *matured, balance.AmtZero)

	// get total rewards that has been generated
	totalRewards, err := storeRwz.GetTotalRewards()
	assert.Equal(t, amt1.Plus(*amt2).Plus(*amt2).Plus(*amt3), totalRewards)
}

func TestStore_IterateAllPendingAmounts(t *testing.T) {
	setup()
	storeRwz.addPendingRewards(delegators[0], balance.NewAmount(30), 500)
	storeRwz.addPendingRewards(delegators[0], balance.NewAmount(30), 600)
	storeRwz.addPendingRewards(delegators[0], balance.NewAmount(30), 700)
	storeRwz.addPendingRewards(delegators[0], balance.NewAmount(30), 800)

	storeRwz.state.Commit()
	count := 0
	storeRwz.iterateAllPD(func(height int64, addr keys.Address, amt *balance.Amount) bool {
		count++
		return false
	})
	assert.Equal(t, 4, count)
}
