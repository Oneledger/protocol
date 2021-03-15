package storage

import (
	"fmt"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/ethereum/go-ethereum/common"
)

// stateEntry represents a single key value pair from the StateDB's stateObject mappindg.
// This is to prevent non determinism at genesis initialization or export.
type stateEntry struct {
	// address key of the state object
	address     common.Address
	stateObject *storage.State
}

type CommitStateDB struct {
	// TODO: We need to store the context as part of the structure itself opposed
	// to being passed as a parameter (as it should be) in order to implement the
	// StateDB interface. Perhaps there is a better way.
	ctx action.Context

	storeKey storage.StoreKey

	// array that hold 'live' objects, which will get modified while processing a
	// state transition
	stateObjects         []stateEntry
	addressToObjectIndex map[common.Address]int // map from address to the index of the state objects slice

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (s *CommitStateDB) CreateAccount(addr common.Address) {
	newObj, prev := s.createObject(addr)
	if prev != nil {
		newObj.setBalance(prev.data.Balance)
	}
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (s *CommitStateDB) createObject(addr common.Address) (newObj, prevObj *stateEntry) {
	prevObj = s.getStateObject(addr)

	acc, _ := s.ctx.Accounts.GetAccount(keys.Address(addr.Bytes()))

	newObj = newStateObject(s, acc)
	newObj.setNonce(0) // sets the object to dirty

	s.setStateObject(newObj)
	return newObj, prevObj
}

// getStateObject attempts to retrieve a state object given by the address.
// Returns nil and sets an error if not found.
func (s *CommitStateDB) getStateObject(addr common.Address) (stateObject *stateEntry) {
	// otherwise, attempt to fetch the account from the account mapper
	acc, _ := s.ctx.Accounts.GetAccount(keys.Address(addr.Bytes()))
	if &acc == nil {
		s.setError(fmt.Errorf("no account found for address: %s", addr.String()))
		return nil
	}

	// insert the state object into the live set
	so := newStateObject(s, acc)
	s.setStateObject(so)

	return so
}

// setError remembers the first non-nil error it is called with.
func (s *CommitStateDB) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

// SubBalance(common.Address, *big.Int)
// AddBalance(common.Address, *big.Int)
// GetBalance(common.Address) *big.Int

// GetNonce(common.Address) uint64
// SetNonce(common.Address, uint64)

// GetCodeHash(common.Address) common.Hash
// GetCode(common.Address) []byte
// SetCode(common.Address, []byte)
// GetCodeSize(common.Address) int

// AddRefund(uint64)
// SubRefund(uint64)
// GetRefund() uint64

// GetCommittedState(common.Address, common.Hash) common.Hash
// GetState(common.Address, common.Hash) common.Hash
// SetState(common.Address, common.Hash, common.Hash)

// Suicide(common.Address) bool
// HasSuicided(common.Address) bool

// // Exist reports whether the given account exists in state.
// // Notably this should also return true for suicided accounts.
// Exist(common.Address) bool
// // Empty returns whether the given account is empty. Empty
// // is defined according to EIP161 (balance = nonce = code = 0).
// Empty(common.Address) bool

// PrepareAccessList(sender common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)
// AddressInAccessList(addr common.Address) bool
// SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool)
// // AddAddressToAccessList adds the given address to the access list. This operation is safe to perform
// // even if the feature/fork is not active yet
// AddAddressToAccessList(addr common.Address)
// // AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform
// // even if the feature/fork is not active yet
// AddSlotToAccessList(addr common.Address, slot common.Hash)

// RevertToSnapshot(int)
// Snapshot() int

// AddLog(*types.Log)
// AddPreimage(common.Hash, []byte)

// ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) error
