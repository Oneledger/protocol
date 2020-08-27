package data

import (
	"testing"

	"github.com/magiconair/properties/assert"
	db2 "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
)

const (
	balanceType = "b"
	governType  = "g"
	feeType     = "f"
)

var (
	chainstate      *storage.ChainState
	balanceStore    *balance.Store
	governanceStore *governance.Store
	feeStore        *fees.Store

	stores Router
)

func init() {
	//Create a few stores
	db := db2.NewDB("testDB", db2.MemDBBackend, "")
	chainstate = storage.NewChainState("chainstate", db)

	balanceStore = balance.NewStore(balanceType, storage.NewState(chainstate))
	governanceStore = governance.NewStore(governType, storage.NewState(chainstate))
	feeStore = fees.NewStore(feeType, storage.NewState(chainstate))

	feeOpt := fees.FeeOption{
		FeeCurrency: balance.Currency{
			Name: "olt", Chain: chain.Type(0), Decimal: 18, Unit: "ones",
		},
		MinFeeDecimal: 1,
	}

	feeStore.SetupOpt(&feeOpt)

	//Create data router
	stores = NewStorageRouter()

}

func TestStorageRouter_Add(t *testing.T) {

	//Add stores to the data router
	_ = stores.Add(balanceType, balanceStore)
	_ = stores.Add(governType, governanceStore)
	_ = stores.Add(feeType, feeStore)

	db, _ := stores.Get(balanceType)
	balanceDB := db.(*balance.Store)

	_ = balanceDB.State.Set(storage.StoreKey("bob"), []byte("1000"))

	db, _ = stores.Get(governType)
	governanceDB := db.(*governance.Store)

	_ = governanceDB.WithHeight(0).SetEpoch(1024)
	governanceDB.WithHeight(0).SetAllLUH()
	db, _ = stores.Get(feeType)
	feeDB := db.(*fees.Store)

	coin := balance.Coin{
		Currency: balance.Currency{
			Name: "olt", Chain: chain.Type(0), Decimal: 18, Unit: "ones",
		},
		Amount: balance.NewAmount(17),
	}

	_ = feeDB.Set(keys.Address("Address1"), coin)
}

func TestStorageRouter_Get(t *testing.T) {
	db, _ := stores.Get(balanceType)
	balanceDB := db.(*balance.Store)

	val, _ := balanceDB.State.Get(storage.StoreKey("bob"))
	assert.Equal(t, val, []byte("1000"))

	db, _ = stores.Get(governType)
	governanceDB := db.(*governance.Store)

	epoch, _ := governanceDB.GetEpoch()
	assert.Equal(t, epoch, int64(1024))

	db, _ = stores.Get(feeType)
	feeDB := db.(*fees.Store)

	coin, _ := feeDB.Get(keys.Address("Address1"))
	assert.Equal(t, coin.Amount, balance.NewAmount(17))
	assert.Equal(t, coin.Currency.Name, "olt")
}
