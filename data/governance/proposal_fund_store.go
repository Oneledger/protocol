package governance

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type ProposalFundStore struct {
	State  *storage.State
	prefix []byte
}

func (pf *ProposalFundStore) set(key storage.StoreKey, amt balance.Amount) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(amt)
	if err != nil {
		return errors.Wrap(err, errorSerialization)
	}
	prefixed := append(pf.prefix, key...)
	err = pf.State.Set(prefixed, dat)
	return errors.Wrap(err, errorSettingRecord)
}

func (pf *ProposalFundStore) get(key storage.StoreKey) (amt *balance.Amount, err error) {
	prefixed := append(pf.prefix, storage.StoreKey(key)...)
	dat, err := pf.State.Get(prefixed)
	//fmt.Println("dat :", dat, "err", err)
	if err != nil {
		return nil, errors.Wrap(err, errorGettingRecord)
	}
	amt = balance.NewAmount(0)
	if len(dat) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, amt)
	if err != nil {
		err = errors.Wrap(err, errorDeSerialization)
	}
	return
}

func (pf *ProposalFundStore) delete(key storage.StoreKey) (bool, error) {
	prefixed := append(pf.prefix, key...)
	res, err := pf.State.Delete(prefixed)
	if err != nil {
		return false, errors.Wrap(err, errorDeletingRecord)
	}
	return res, err
}

func (pf *ProposalFundStore) iterate(fn func(proposalID ProposalID, addr keys.Address, amt *balance.Amount) bool) bool {
	return pf.State.IterateRange(
		pf.prefix,
		storage.Rangefix(string(pf.prefix)),
		true,
		func(key, value []byte) bool {
			amt := balance.NewAmount(0)
			err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(value, amt)
			if err != nil {
				return true
			}
			arr := strings.Split(string(key), storage.DB_PREFIX)
			proposalID := arr[1]
			fundingAddress := keys.Address(arr[len(arr)-1])
			err = fundingAddress.UnmarshalText([]byte(arr[len(arr)-1]))
			if err != nil {
				fmt.Println("Error Unmarshalling ", err)
				return true
			}
			return fn(ProposalID(proposalID), fundingAddress, amt)
		},
	)
}

func (pf *ProposalFundStore) WithState(state *storage.State) *ProposalFundStore {
	pf.State = state
	return pf
}

// Store Function Called my external Layers
func NewProposalFundStore(prefix string, state *storage.State) *ProposalFundStore {
	return &ProposalFundStore{
		State:  state,
		prefix: storage.Prefix(prefix),
	}
}

func (pf *ProposalFundStore) GetFundsForProposalID(id ProposalID, fn func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund) []ProposalFund {
	var foundProposals []ProposalFund
	pf.iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) bool {
		if proposalID == id {
			foundProposals = append(foundProposals, fn(proposalID, fundingAddr, amt))
		}
		return false
	})
	return foundProposals
}

func (pf *ProposalFundStore) GetProposalsForFunder(funderAddress keys.Address, fn func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund) []ProposalFund {
	var foundProposals []ProposalFund
	pf.iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) bool {
		if bytes.Equal(funderAddress, fundingAddr) {
			foundProposals = append(foundProposals, fn(proposalID, fundingAddr, amt))
		}
		return false
	})
	return foundProposals
}

func (store *ProposalFundStore) IsFundedByFunder(id ProposalID, funder keys.Address) bool {
	haveFunderAddress := false
	store.GetFundsForProposalID(id, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		if fundingAddr.Equal(funder) {
			haveFunderAddress = true
		}
		return ProposalFund{}
	})
	return haveFunderAddress
}

func (pf *ProposalFundStore) AddFunds(proposalId ProposalID, fundingAddress keys.Address, amount *balance.Amount) error {
	key := storage.StoreKey(string(proposalId) + storage.DB_PREFIX + fundingAddress.String())
	amt, err := pf.get(key)
	if err != nil {
		return errors.Wrap(err, errorGettingRecord)
	}
	return pf.set(key, *amt.Plus(*amount))
}

func (pf *ProposalFundStore) DeductFunds(proposalId ProposalID, fundingAddress keys.Address, amount *balance.Amount) error {
	key := storage.StoreKey(string(proposalId) + storage.DB_PREFIX + fundingAddress.String())
	amt, err := pf.get(key)
	if err != nil {
		return errors.Wrap(err, errorGettingRecord)
	}
	result, err := amt.Minus(*amount)
	if err != nil {
		return errors.Wrap(err, errorGettingRecord)
	}
	return pf.set(key, *result)
}

func (pf *ProposalFundStore) DeleteFunds(proposalId ProposalID, fundingAddress keys.Address) (bool, error) {
	key := storage.StoreKey(string(proposalId) + storage.DB_PREFIX + fundingAddress.String())

	_, err := pf.get(key)
	if err != nil {

		return false, errors.Wrap(err, errorGettingRecord)
	}
	ok, err := pf.delete(key)
	if err != nil {
		return false, errors.Wrap(err, errorDeletingRecord)
	}
	return ok, nil
}
