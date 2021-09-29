package vm

import (
	"fmt"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (s *CommitStateDB) createObject(addr ethcmn.Address) (newObj, prevObj *stateObject) {
	prevObj = s.getStateObject(addr)

	acc, err := s.accountKeeper.NewAccountWithAddress(keys.Address(addr.Bytes()))
	if err != nil {
		s.setError(fmt.Errorf("failed to create a new address for %s", addr.String()))
		return nil, prevObj
	}
	newObj = newStateObject(s, acc)
	newObj.setNonce(0) // sets the object to dirty

	if prevObj == nil {
		s.journal.append(createObjectChange{account: &addr})
	} else {
		s.journal.append(resetObjectChange{prev: prevObj})
	}

	s.setStateObject(newObj)
	return newObj, prevObj
}

// getStateObject attempts to retrieve a state object given by the address.
// Returns nil and sets an error if not found.
func (s *CommitStateDB) getStateObject(addr ethcmn.Address) (stateObject *stateObject) {
	if idx, found := s.addressToObjectIndex[addr]; found {
		// prefer 'live' (cached) objects
		if so := s.stateObjects[idx].stateObject; so != nil {
			if so.deleted {
				return nil
			}

			return so
		}
	}

	// otherwise, attempt to fetch the account from the account mapper
	acc, err := s.accountKeeper.GetAccount(keys.Address(addr.Bytes()))
	if err != nil {
		if err != balance.ErrAccountNotFound {
			// if not this error, means something wrong
			s.setError(fmt.Errorf("failed to get account %s, error: %s", addr.String(), err))
		}
		return nil
	}

	// insert the state object into the live set
	so := newStateObject(s, acc)
	s.setStateObject(so)

	return so
}

func (s *CommitStateDB) setStateObject(so *stateObject) {
	if idx, found := s.addressToObjectIndex[so.Address()]; found {
		// update the existing object
		s.stateObjects[idx].stateObject = so
		return
	}

	// append the new state object to the stateObjects slice
	se := stateEntry{
		address:     so.Address(),
		stateObject: so,
	}

	s.stateObjects = append(s.stateObjects, se)
	s.addressToObjectIndex[se.address] = len(s.stateObjects) - 1
}

// updateStateObject writes the given state object to the store.
func (s *CommitStateDB) updateStateObject(so *stateObject) error {
	s.logger.Detailf("VM: update state object for address '%s' with nonce: '%d' and balance: '%d' \n", so.Address(), so.account.Sequence, so.account.Balance())
	return s.accountKeeper.SetAccount(*so.account)
}

// setError remembers the first non-nil error it is called with.
func (s *CommitStateDB) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *CommitStateDB) clearJournalAndRefund() {
	s.journal = newJournal()
	s.refund = 0
	s.validRevisions = s.validRevisions[:0]
}

// deleteStateObject removes the given state object from the state store.
func (s *CommitStateDB) deleteStateObject(so *stateObject) {
	so.deleted = true
	s.logger.Detailf("VM: delete state object for address '%s' with nonce: '%d' and balance: '%d' \n", so.Address(), so.account.Sequence, so.account.Balance())
	s.accountKeeper.RemoveAccount(*so.account)
}
