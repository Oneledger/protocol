package network_delegation

import (
	"testing"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/magiconair/properties/assert"
	db "github.com/tendermint/tm-db"
)

const (
	numPrivateKeys    = 15
	activeStartIndex  = 0
	pendingStartIndex = 5
	matureStartIndex  = 10
)

var (
	activeAddrList  map[string]*balance.Coin
	pendingAddrList map[string]*balance.Coin
	matureAddrList  map[string]*balance.Coin
	coinList        []*balance.Coin
	store           *Store
)

func init() {
	for i := 100; i < 115; i++ {
		coin := &balance.Coin{
			Currency: balance.Currency{
				Id:      0,
				Name:    "cur",
				Chain:   0,
				Decimal: 10,
				Unit:    "u",
			},
			Amount: balance.NewAmount(int64(i)),
		}
		coinList = append(coinList, coin)
	}

	activeAddrList = make(map[string]*balance.Coin)
	pendingAddrList = make(map[string]*balance.Coin)
	matureAddrList = make(map[string]*balance.Coin)

	//Create Validator Keys
	for i := 0; i < numPrivateKeys; i++ {
		pub, _, _ := keys.NewKeyPairFromTendermint()
		h, _ := pub.GetHandler()

		if i >= activeStartIndex && i < pendingStartIndex {
			activeAddrList[h.Address().String()] = coinList[i]
		}
		if i >= pendingStartIndex && i < matureStartIndex {
			pendingAddrList[h.Address().String()] = coinList[i]
		}
		//if i >= matureStartIndex {
		//	matureAddrList[h.Address().String()] = coinList[i]
		//}
	}

	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))

	store = NewStore("nd", cs)
}

func TestStore_Set_Get(t *testing.T) {
	//Test Active Store
	store.WithPrefix(ActiveType)
	for i, v := range activeAddrList {
		addr := keys.Address{}
		_ = addr.UnmarshalText([]byte(i))
		err := store.Set(addr, v)
		assert.Equal(t, err, nil)

		coin, err := store.Get(addr)
		assert.Equal(t, coin, v)
	}

	////Test Mature Store
	//store.WithPrefix(MatureType)
	//for i, v := range matureAddrList {
	//	addr := keys.Address{}
	//	_ = addr.UnmarshalText([]byte(i))
	//	err := store.Set(addr, v)
	//	assert.Equal(t, err, nil)
	//
	//	coin, err := store.Get(addr)
	//	assert.Equal(t, coin, v)
	//}

	store.State.Commit()
}

func TestStore_Exists(t *testing.T) {
	store.WithPrefix(ActiveType)
	for key := range activeAddrList {
		addr := &keys.Address{}
		_ = addr.UnmarshalText([]byte(key))
		res := store.Exists(addr)
		assert.Equal(t, res, true)
	}
}

func TestStore_IterateActiveAmounts(t *testing.T) {
	//iterate all Active amounts
	store.IterateActiveAmounts(func(addr *keys.Address, coin *balance.Coin) bool {
		assert.Equal(t, coin, activeAddrList[addr.String()])
		return false
	})
}

//func TestStore_IterateMatureAmounts(t *testing.T) {
//	//iterate all Active amounts
//	store.IterateMatureAmounts(func(addr *keys.Address, coin *balance.Coin) bool {
//		assert.Equal(t, coin, matureAddrList[addr.String()])
//		return false
//	})
//}

func TestStore_SetPendingAmount(t *testing.T) {
	//Test Set Pending Amount at different heights - pending address list is at 500, active address list is at 600
	store.WithPrefix(PendingType)
	for i, v := range pendingAddrList {
		addr := keys.Address{}
		_ = addr.UnmarshalText([]byte(i))
		err := store.SetPendingAmount(addr, 500, v)
		assert.Equal(t, err, nil)
	}

	for i, v := range activeAddrList {
		addr := keys.Address{}
		_ = addr.UnmarshalText([]byte(i))
		err := store.SetPendingAmount(addr, 600, v)
		assert.Equal(t, err, nil)
	}
	store.State.Commit()

	//Test Get Pending Amount for one address at height 500
	for i, v := range pendingAddrList {
		addr := keys.Address{}
		_ = addr.UnmarshalText([]byte(i))
		amount, err := store.GetPendingAmount(addr, 500)
		assert.Equal(t, err, nil)
		assert.Equal(t, amount, v)
	}

	//Test Check If Pending Amount Exists
	for i := range pendingAddrList {
		addr := keys.Address{}
		_ = addr.UnmarshalText([]byte(i))
		exist := store.PendingExists(addr, 500)
		assert.Equal(t, exist, true)
	}

	//Test iterate Pending Amounts at different heights
	//Iterate pending amounts at height 500
	store.IteratePendingAmounts(500, func(addr *keys.Address, coin *balance.Coin) bool {
		assert.Equal(t, coin, pendingAddrList[addr.String()])
		return false
	})
	//Iterate pending amounts at height 600
	store.IteratePendingAmounts(600, func(addr *keys.Address, coin *balance.Coin) bool {
		assert.Equal(t, coin, activeAddrList[addr.String()])
		return false
	})

	//Test iterate All Pending Amounts
	//Should be 10 records in total within Pending amounts
	count := 0
	store.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
		count++
		return false
	})
	assert.Equal(t, count, 10)
}
