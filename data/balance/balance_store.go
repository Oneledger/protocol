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
	*storage.ChainState
}

func NewStore(name, dbDir, configDB string, typ storage.StorageType) *Store {
	cs := storage.NewChainState(name, dbDir, configDB, typ)

	return &Store{cs}
}

func (st *Store) Get(address []byte, lastCommit bool) (bal *Balance, err error) {
	dat := st.ChainState.Get(storage.StoreKey(address), lastCommit)

	if len(dat) == 0 {
		err = ErrNoBalanceFoundForThisAddress
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

	err = st.ChainState.Set(storage.StoreKey(address), dat)

	return err
}

func (st *Store) FindAll() map[string]*Balance {
	balMap := make(map[string]*Balance)

	pSzlr := serialize.GetSerializer(serialize.PERSISTENT)
	for key, dat := range st.ChainState.FindAll() {
		bal := &Balance{}
		var err error
		err = pSzlr.Deserialize(dat, bal)
		if err != nil {
			logger.Error("error in deserializing balance", "key", key)
		}

		balMap[key] = bal
	}

	return balMap
}

func (st *Store) Exists(address keys.Address) bool {
	return st.ChainState.Exists(storage.StoreKey(address))
}
