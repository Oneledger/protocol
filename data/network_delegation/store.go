package network_delegation

import (
	"strconv"
	"strings"
	"sync"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

const (
	//MatureKey  = "m"
	PendingKey = "p"
	ActiveKey  = "a"
)

type Store struct {
	State         *storage.State
	szlr          serialize.Serializer
	prefix        []byte
	currentPrefix []byte
	mux           sync.Mutex
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		State:         state,
		prefix:        storage.StoreKey(prefix),
		currentPrefix: storage.StoreKey(prefix + storage.DB_PREFIX + ActiveKey),
		szlr:          serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (st *Store) WithState(state *storage.State) *Store {
	st.State = state
	return st
}

func (st *Store) GetState() *storage.State {
	return st.State
}

func (st *Store) Exists(addr *keys.Address) bool {
	key := append(st.currentPrefix, addr.String()...)
	return st.State.Exists(key)
}

//Set coin to specific key
func (st *Store) set(key []byte, coin *balance.Coin) (err error) {
	dat, err := st.szlr.Serialize(coin)
	if err != nil {
		return
	}
	err = st.State.Set(storage.StoreKey(key), dat)
	return
}

//get coin from specific key
func (st *Store) get(key []byte) (coin *balance.Coin, err error) {
	coin = &balance.Coin{
		Currency: balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"},
		Amount:   balance.NewAmount(0),
	}
	dat, err := st.State.Get(storage.StoreKey(key))
	if dat == nil || err != nil {
		return
	}
	err = st.szlr.Deserialize(dat, coin)
	return
}

//Set coin to specific address
func (st *Store) Set(address keys.Address, coin *balance.Coin) (err error) {
	prefixKey := append(st.currentPrefix, address.String()...)
	err = st.set(storage.StoreKey(prefixKey), coin)
	return
}

//Get coin from address
func (st *Store) Get(address keys.Address) (coin *balance.Coin, err error) {
	prefixKey := append(st.currentPrefix, address.String()...)
	coin, err = st.get(storage.StoreKey(prefixKey))
	return
}

func (st *Store) WithPrefix(prefix DelegationPrefixType) *Store {
	switch prefix {
	//case MatureType:
	//	st.currentPrefix = st.buildMatureKey()
	case PendingType:
		st.currentPrefix = st.buildPendingKey()
	case ActiveType:
		st.currentPrefix = st.buildActiveKey()
	}
	return st
}

func (st *Store) iterate(prefix storage.StoreKey, fn func(key []byte, coin *balance.Coin) bool) bool {
	return st.State.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			coin := &balance.Coin{}
			err := st.szlr.Deserialize(value, coin)
			if err != nil {
				return true
			}
			return fn(key, coin)
		},
	)
}

func (st *Store) iterateAddresses(prefix storage.StoreKey, fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
	return st.iterate(prefix, func(key []byte, coin *balance.Coin) bool {
		arr := strings.Split(string(key), storage.DB_PREFIX)
		addr := &keys.Address{}
		err := addr.UnmarshalText([]byte(arr[len(arr)-1]))
		if err != nil {
			return true
		}
		return fn(addr, coin)
	})
}

//------------------------------- Active key store -------------------------------
//build Active key
func (st *Store) buildActiveKey() storage.StoreKey {
	return storage.Prefix(string(st.prefix) + storage.DB_PREFIX + ActiveKey)
}

func (st *Store) IterateActiveAmounts(fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
	prefix := st.buildActiveKey()
	//fmt.Println("Active prefix: ", prefix)
	return st.iterateAddresses(prefix, func(addr *keys.Address, coin *balance.Coin) bool {
		return fn(addr, coin)
	})
}

//------------------------------- Pending key store -------------------------------
//build pending key
func (st *Store) buildPendingKey() storage.StoreKey {
	return storage.Prefix(string(st.prefix) + storage.DB_PREFIX + PendingKey)
}

//check existence of pending amount on specific height
func (st *Store) PendingExists(addr keys.Address, height int64) bool {
	prefix := st.buildPendingKey()
	pendingKey := strconv.FormatInt(height, 10) + storage.DB_PREFIX + addr.String()
	return st.State.Exists(append(prefix, pendingKey...))
}

//check existence of pending amount on specific height
func (st *Store) GetPendingAmount(addr keys.Address, height int64) (coin *balance.Coin, err error) {
	prefix := st.buildPendingKey()
	pendingKey := strconv.FormatInt(height, 10) + storage.DB_PREFIX + addr.String()
	coin, err = st.get(append(prefix, pendingKey...))
	return
}

//Set pending amount with height and address
func (st *Store) SetPendingAmount(addr keys.Address, height int64, coin *balance.Coin) error {
	prefix := st.buildPendingKey()
	pendingKey := strconv.FormatInt(height, 10) + storage.DB_PREFIX + addr.String()
	return st.set(append(prefix, pendingKey...), coin)
}

//iterate addresses for height
func (st *Store) IteratePendingAmounts(height int64, fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
	prefix := append(st.buildPendingKey(), strconv.FormatInt(height, 10)...)
	return st.iterateAddresses(prefix, func(addr *keys.Address, coin *balance.Coin) bool {
		return fn(addr, coin)
	})
}

//iterate all pending amounts
func (st *Store) IterateAllPendingAmounts(fn func(height int64, addr *keys.Address, coin *balance.Coin) bool) bool {
	prefix := st.buildPendingKey()
	return st.iterate(prefix, func(key []byte, coin *balance.Coin) bool {
		arr := strings.Split(string(key), storage.DB_PREFIX)
		addr := &keys.Address{}
		err := addr.UnmarshalText([]byte(arr[len(arr)-1]))
		if err != nil {
			return true
		}
		height, err := strconv.ParseInt(arr[len(arr)-2], 10, 64)
		if err != nil {
			return true
		}
		return fn(height, addr, coin)
	})
}

// below is removed since withdraw logic is moved to block beginner, OLP-1267
////------------------------------- Mature key store -------------------------------
//
////build mature key
//func (st *Store) buildMatureKey() storage.StoreKey {
//	return storage.Prefix(string(st.prefix) + storage.DB_PREFIX + MatureKey)
//}
//
////iterate all matured amounts
//func (st *Store) IterateMatureAmounts(fn func(addr *keys.Address, coin *balance.Coin) bool) bool {
//	prefix := st.buildMatureKey()
//	return st.iterateAddresses(prefix, func(addr *keys.Address, coin *balance.Coin) bool {
//		return fn(addr, coin)
//	})
//}
//
////mature pending delegator
//func (st *Store) maturePendingDelegator(delegator *PendingDelegator) error {
//	//Get current Matured amount stored for Address
//	currentCoin, err := st.WithPrefix(MatureType).Get(*delegator.Address)
//	if err != nil {
//		return err
//	}
//
//	//Add amount to matured prefix
//	newCoin := currentCoin.Plus(*delegator.Amount)
//	err = st.WithPrefix(MatureType).Set(*delegator.Address, &newCoin)
//	if err != nil {
//		return err
//	}
//
//	//clear pending record amount
//	delegator.Amount.Amount = balance.NewAmount(0)
//	return st.SetPendingAmount(*delegator.Address, delegator.Height, delegator.Amount)
//}
//
//func (st *Store) HandlePendingDelegates(height int64) error {
//	var pendingList []*PendingDelegator
//	st.IteratePendingAmounts(height, func(addr *keys.Address, coin *balance.Coin) bool {
//		pendingDelegator := &PendingDelegator{
//			Address: addr,
//			Amount:  coin,
//			Height:  height,
//		}
//		pendingList = append(pendingList, pendingDelegator)
//		return false
//	})
//
//	for _, delegator := range pendingList {
//		err := st.maturePendingDelegator(delegator)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

//__________________________________ Load State ____________________________________

func (st *Store) LoadDelegators(state State) error {
	//Load Active Delegators
	for _, delegator := range state.ActiveList {
		err := st.WithPrefix(ActiveType).Set(*delegator.Address, delegator.Amount)
		if err != nil {
			return err
		}
	}
	////Load Mature Delegators
	//for _, delegator := range state.MatureList {
	//	err := st.WithPrefix(MatureType).Set(*delegator.Address, delegator.Amount)
	//	if err != nil {
	//		return err
	//	}
	//}
	//Load Pending Delegators
	for _, delegator := range state.PendingList {
		err := st.SetPendingAmount(*delegator.Address, delegator.Height, delegator.Amount)
		if err != nil {
			return err
		}
	}

	//st.LoadTestData()

	return nil
}

//func (st *Store)LoadTestData() {
//	var coinList []*balance.Coin
//
//	for i := 100; i < 115; i++ {
//		coin := &balance.Coin{
//			Currency: balance.Currency{
//				Id:      0,
//				Name:    "cur",
//				Chain:   0,
//				Decimal: 10,
//				Unit:    "u",
//			},
//			Amount: balance.NewAmount(int64(i)),
//		}
//		coinList = append(coinList, coin)
//	}
//
//	var addrList []keys.Address
//	//Create Validator Keys
//	for i := 0; i < 10; i++ {
//		pub, _, _ := keys.NewKeyPairFromTendermint()
//		h, _ := pub.GetHandler()
//		addrList = append(addrList, h.Address())
//	}
//
//	_ = st.WithPrefix(ActiveType).Set(addrList[0], coinList[0])
//	_ = st.WithPrefix(ActiveType).Set(addrList[1], coinList[1])
//
//	_ = st.WithPrefix(MatureType).Set(addrList[2], coinList[2])
//	_ = st.WithPrefix(MatureType).Set(addrList[3], coinList[3])
//
//	_ = st.SetPendingAmount(addrList[4], 100, coinList[4])
//	_ = st.SetPendingAmount(addrList[5], 100, coinList[5])
//}
