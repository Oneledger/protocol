package networkDelegation

import (
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"sync"
)

const (
	MATURE_KEY  = "m"
	PENDING_KEY = "p"
	ACTIVE_KEY  = "a"
)

type Store struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte
	mux    sync.Mutex
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		state:  state,
		prefix: storage.Prefix(prefix),
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (st *Store) WithState(state *storage.State) *Store {
	st.state = state
	return st
}

//func (st *Store) Get(key []byte) *balance.Coin {
//	prefixKey := append(st.prefix, key...)
//	dat, _ := st.state.Get(storage.StoreKey(prefixKey))
//}

//Set

//------------------------------- Active key store -------------------------------
//Set New Delegator
//Get Delegator amount

//------------------------------- Pending key store -------------------------------
//build pending key
//Set pending amount with height and address
//iterate addresses for height
//iterate all pending amounts

//------------------------------- Mature key store -------------------------------

//build mature key
//get mature value for address
//set mature value for address
//iterate all matured amounts
