package governance

import (
	"fmt"
	"strings"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"

	"github.com/pkg/errors"
)

type ProposalFundStore struct {
	State  *storage.State
	prefix []byte
}

func NewProposalFundStore(prefix string, state *storage.State) *ProposalFundStore {
	return &ProposalFundStore{
		State:  state,
		prefix: storage.Prefix(prefix),
	}
}

func (st *ProposalFundStore) set(key storage.StoreKey, amt ProposalAmount) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(amt)
	if err != nil {
		return err
	}
	prefixed := append(st.prefix, key...)
	err = st.State.Set(prefixed, dat)
	return err
}

func (st *ProposalFundStore) get(key storage.StoreKey) (amt *ProposalAmount, err error) {
	prefixed := append(st.prefix, storage.StoreKey(key)...)
	dat, _ := st.State.Get(prefixed)
	amt = NewAmount(0)
	if len(dat) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, amt)
	return
}

func (pf *ProposalFundStore) getIterable() storage.Iterable {
	return pf.State.GetIterable()
}

func (pf *ProposalFundStore) AddFunds(proposalId ProposalID, fundingAddress keys.Address, amount *ProposalAmount) error {
	key := storage.StoreKey(string(proposalId) + storage.DB_PREFIX + fundingAddress.String())
	amt, err := pf.get(key)
	if err != nil {
		return errors.Wrapf(err, "proposer already exists %s", fundingAddress.String())
	}
	return pf.set(key, *amt.Plus(amount))
}

func (pf *ProposalFundStore) Iterate(fn func(proposalID ProposalID, addr keys.Address, amt *ProposalAmount) bool) bool {
	return pf.State.IterateRange(
		pf.prefix,
		storage.Rangefix(string(pf.prefix)),
		true,
		func(key, value []byte) bool {

			amt := NewAmount(0)
			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, amt)
			if err != nil {
				fmt.Println("err", err)
				return true
			}
			arr := strings.Split(string(key), storage.DB_PREFIX)
			proposalID := arr[1]
			fmt.Println("Address :", arr[2])
			fundingAddress := keys.Address(arr[2]) //arr[len(arr)-1]
			return fn(ProposalID(proposalID), fundingAddress, amt)
		},
	)
}
