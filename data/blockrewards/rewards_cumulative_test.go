package blockrewards

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

var (
	store      *RewardsCumulativeStore
	cs         *storage.State
	validator1 keys.Address
	validator2 keys.Address
	zero       *balance.Amount
	amt1       *balance.Amount
	amt2       *balance.Amount
	amt3       *balance.Amount
	amt4       *balance.Amount
	amt5       *balance.Amount
	withdraw1  *balance.Amount
	withdraw2  *balance.Amount
)

func setup() {
	fmt.Println("####### Testing  rewards cumulative store #######")
	memDb := db.NewDB("test", db.MemDBBackend, "")
	cs = storage.NewState(storage.NewChainState("rewards", memDb))
	store = NewRewardsCumulativeStore("rws", cs)
	generateAddresses()
}

func generateAddresses() {
	pub1, _, _ := keys.NewKeyPairFromTendermint()
	h1, _ := pub1.GetHandler()
	validator1 = h1.Address()

	pub2, _, _ := keys.NewKeyPairFromTendermint()
	h2, _ := pub2.GetHandler()
	validator2 = h2.Address()

	zero = balance.NewAmount(0)
	amt1 = balance.NewAmount(100)
	amt2 = balance.NewAmount(200)
	amt3 = balance.NewAmount(377)
	withdraw1 = balance.NewAmount(163)
	withdraw2 = balance.NewAmount(499)
}

func TestNewRewardsCumulativeStore(t *testing.T) {
	setup()
	mutured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	withdrawn, err := store.GetWithdrawnRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, zero, mutured)
	assert.Equal(t, zero, balance)
	assert.Equal(t, zero, withdrawn)
}

func TestRewardsCumulativeStore_AddGetMaturedBalance(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected := amt1.Plus(amt2)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, matured, expected)
}

func TestRewardsCumulativeStore_WithdrawRewards(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	store.WithdrawRewards(validator1, withdraw1)
	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected, _ := amt1.Plus(amt2).Minus(withdraw1)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	expected = amt1.Plus(amt2)
	assert.Equal(t, matured, expected)
}

func TestRewardsCumulativeStore_GetWithdrawnRewards(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	store.AddMaturedBalance(validator2, amt1)
	store.WithdrawRewards(validator1, withdraw1)
	store.AddMaturedBalance(validator1, amt3)
	store.AddMaturedBalance(validator2, amt2)
	store.WithdrawRewards(validator1, withdraw2)

	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected, _ := amt1.Plus(amt2).Plus(amt3).Minus(withdraw1.Plus(withdraw2))
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	expected = amt1.Plus(amt2).Plus(amt3)
	assert.Equal(t, matured, expected)

	withdrawn, err := store.GetWithdrawnRewards(validator1)
	assert.Nil(t, err)
	expected = withdraw1.Plus(withdraw2)
	assert.Equal(t, withdrawn, expected)
}

func TestRewardsCumulativeStore_WithdrawOthers(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	err := store.WithdrawRewards(validator2, withdraw1)
	assert.NotNil(t, err)

	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected := amt1.Plus(amt2)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, matured, expected)
}

func TestRewardsCumulativeStore_OverWithdraw(t *testing.T) {
	setup()
	store.AddMaturedBalance(validator1, amt1)
	store.AddMaturedBalance(validator1, amt2)
	err := store.WithdrawRewards(validator1, withdraw2)
	assert.NotNil(t, err)

	balance, err := store.GetMaturedBalance(validator1)
	assert.Nil(t, err)
	expected := amt1.Plus(amt2)
	assert.Equal(t, balance, expected)

	matured, err := store.GetMaturedRewards(validator1)
	assert.Nil(t, err)
	assert.Equal(t, matured, expected)
}
