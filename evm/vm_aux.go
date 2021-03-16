package evm

import (
	"fmt"

	"github.com/Oneledger/protocol/data/accounts"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (s *CommitStateDB) createObject(addr ethcmn.Address) (newObj, prevObj *stateObject) {
	prevObj = s.getStateObject(addr)

	acc, _ := s.ctx.Accounts.GetAccount(addr.Bytes())

	newObj = newStateObject(s, acc)
	newObj.setNonce(0) // sets the object to dirty

	s.setStateObject(newObj)
	return newObj, prevObj
}

// getStateObject attempts to retrieve a state object given by the address.
// Returns nil and sets an error if not found.
func (s *CommitStateDB) getStateObject(addr ethcmn.Address) (stateObject *stateObject) {
	// otherwise, attempt to fetch the account from the account mapper
	acc, _ := s.ctx.Accounts.GetAccount(addr.Bytes())
	if (accounts.Account{}) == acc {
		s.setError(fmt.Errorf("no account found for address: %s", addr.String()))
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

// setError remembers the first non-nil error it is called with.
func (s *CommitStateDB) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

// GetOrNewStateObject retrieves a state object or create a new state object if
// nil.
func (s *CommitStateDB) GetOrNewStateObject(addr ethcmn.Address) StateObject {
	so := s.getStateObject(addr)
	if so == nil || so.deleted {
		so, _ = s.createObject(addr)
	}
	return so
}
