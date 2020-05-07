package governance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
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
func (pf *ProposalFundStore) set(key storage.StoreKey, amt ProposalAmount) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(amt)
	if err != nil {
		return err
	}

	prefixed := append(pf.prefix, key...)
	err = pf.State.Set(prefixed, dat)
	return err
}
func (pf *ProposalFundStore) get(key storage.StoreKey) (amt ProposalAmount, err error) {
	prefixed := append(pf.prefix, key...)
	dat, _ := pf.State.Get(prefixed)
	amt = 0
	if len(dat) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, amt)
	return
}

func (pf *ProposalFundStore) getIterable() storage.Iterable {
	return pf.State.GetIterable()
}

func (pf *ProposalFundStore) AddProposalFund(proposalId ProposalID, fundingAddress keys.Address, amount ProposalAmount) error {
	key := storage.StoreKey(string(proposalId) + storage.DB_PREFIX + fundingAddress.String())
	amt, err := pf.get(key)
	if err != nil {
		return errors.Wrapf(err, "failed to get address balance %s", fundingAddress.String())
	}
	newamount := amt + amount
	return pf.set(key, newamount)
}

func (pf *ProposalFundStore) AddNewPropososer(proposalId ProposalID, fundingAddress keys.Address, amount ProposalAmount) error {
	key := storage.StoreKey(string(proposalId) + storage.DB_PREFIX + fundingAddress.String())
	_, err := pf.get(key)
	if err == nil {
		return errors.Wrapf(err, "proposer already exists %s", fundingAddress.String())
	}
	return pf.set(key, amount)
}

func (pf *ProposalFundStore) Iterate(fn func(proposalID ProposalID, addr keys.Address, amt ProposalAmount) bool) bool {
	return pf.State.IterateRange(
		pf.prefix,
		storage.Rangefix(string(pf.prefix)),
		true,
		func(key, value []byte) bool {
			amt := 0
			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, amt)
			if err != nil {
				return true
			}
			arr := strings.Split(string(key), storage.DB_PREFIX)
			proposalID, err := strconv.Atoi(arr[1])
			if err != nil {
				fmt.Println(err)
				return true
			}
			fundingAddress := keys.Address(arr[len(arr)-1])
			return fn(ProposalID(proposalID), fundingAddress, ProposalAmount(amt))
		},
	)
}
