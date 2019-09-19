/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package balance

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

var ErrNoBalanceFoundForThisAddress = errors.New("no balance found for the address")

type Store struct {
	State  *storage.State
	prefix []byte
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		State:  state,
		prefix: storage.Prefix(prefix),
	}
}

func (st *Store) WithState(state *storage.State) *Store {
	st.State = state
	return st
}

func (st *Store) Get(address []byte) (bal *Balance, err error) {
	key := append(st.prefix, storage.StoreKey(address)...)
	dat, _ := st.State.Get(key)

	if len(dat) == 0 {
		err = ErrNoBalanceFoundForThisAddress
		bal = NewBalance()
		return
	}
	bal = NewBalance()
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, bal)
	return
}

func (st *Store) Set(address keys.Address, balance Balance) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(&balance)
	if err != nil {
		return err
	}

	key := append(st.prefix, storage.StoreKey(address)...)
	err = st.State.Set(key, dat)
	return err
}

//func (st *Store) FindAll() map[string]*Balance {
//	balMap := make(map[string]*Balance)
//
//	pSzlr := serialize.GetSerializer(serialize.PERSISTENT)
//	for key, dat := range st.State.FindAll() {
//		bal := &Balance{}
//		var err error
//		err = pSzlr.Deserialize(dat, bal)
//		if err != nil {
//			logger.Error("error in deserializing balance", "key", key)
//		}
//
//		balMap[key] = bal
//	}
//
//	return balMap
//}

func (st *Store) Exists(address keys.Address) bool {
	key := append(st.prefix, storage.StoreKey(address)...)
	return st.State.Exists(key)
}
