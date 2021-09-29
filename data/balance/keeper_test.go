package balance

import (
	"math/big"
	"testing"

	"github.com/Oneledger/protocol/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func prepareKeeperEnv(t *testing.T) AccountKeeper {
	olt := Currency{
		Id:      0,
		Name:    "OLT",
		Chain:   0,
		Decimal: 18,
		Unit:    "nue",
	}
	db := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("balance", db))
	store := NewStore("b", cs)
	currencies := NewCurrencySet()
	err := currencies.Register(olt)
	assert.NoError(t, err)
	keeper := NewNesterAccountKeeper(
		storage.NewState(storage.NewChainState("keeper", db)),
		store,
		currencies,
	)
	return keeper
}

func TestKeeper(t *testing.T) {

	t.Run("test new account creation and it is OK", func(t *testing.T) {
		keeper := prepareKeeperEnv(t)

		addr := common.Address{}
		_, err := keeper.GetAccount(addr.Bytes())
		assert.Error(t, err)

		acc, err := keeper.NewAccountWithAddress(addr.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(0), acc.Coins.Amount.BigInt())
		assert.Equal(t, addr.Bytes(), acc.Address.Bytes())
		assert.Equal(t, uint64(0), acc.Sequence)

		_, err = keeper.GetAccount(addr.Bytes())
		assert.Error(t, err)

		err = keeper.SetAccount(*acc)
		assert.NoError(t, err)

		acc, err = keeper.GetAccount(addr.Bytes())
		assert.NoError(t, err)

		assert.Equal(t, big.NewInt(0), acc.Coins.Amount.BigInt())
		assert.Equal(t, addr.Bytes(), acc.Address.Bytes())
		assert.Equal(t, uint64(0), acc.Sequence)
	})

	t.Run("test balance add and it is OK", func(t *testing.T) {
		keeper := prepareKeeperEnv(t)

		addr := common.Address{}
		acc, err := keeper.NewAccountWithAddress(addr.Bytes())
		assert.NoError(t, err)

		value := big.NewInt(100)
		acc.AddBalance(value)
		keeper.SetAccount(*acc)

		acc, err = keeper.GetAccount(addr.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, value, acc.Coins.Amount.BigInt())
		assert.Equal(t, addr.Bytes(), acc.Address.Bytes())
		assert.Equal(t, uint64(0), acc.Sequence)
	})
}
