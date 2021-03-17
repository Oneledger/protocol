package evm

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/Oneledger/protocol/action"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var _ ethvm.StateDB = (*CommitStateDB)(nil)

type revision struct {
	id           int
	journalIndex int
}

type CommitStateDB struct {
	// TODO: We need to store the context as part of the structure itself opposed
	// to being passed as a parameter (as it should be) in order to implement the
	// StateDB interface. Perhaps there is a better way.
	ctx *action.Context

	// The refund counter, also used by state transitioning.
	refund uint64

	// keeper interface
	accountKeeper AccountKeeper

	thash, bhash ethcmn.Hash
	txIndex      int
	logs         map[ethcmn.Hash][]*ethtypes.Log
	logSize      uint

	preimages map[common.Hash][]byte

	// array that hold 'live' objects, which will get modified while processing a
	// state transition
	stateObjects         []stateEntry
	addressToObjectIndex map[common.Address]int // map from address to the index of the state objects slice

	// Per-transaction access list
	accessList *accessList

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	// journal        *journal
	validRevisions []revision
	nextRevisionID int
}

// NewCommitStateDB returns a reference to a newly initialized CommitStateDB
// which implements Geth's state.StateDB interface.
//
// CONTRACT: Stores used for state must be cache-wrapped as the ordering of the
// key/value space matters in determining the merkle root.
func NewCommitStateDB(ctx *action.Context, ak AccountKeeper) *CommitStateDB {
	return &CommitStateDB{
		ctx:                  ctx,
		stateObjects:         []stateEntry{},
		accountKeeper:        ak,
		logs:                 make(map[ethcmn.Hash][]*ethtypes.Log),
		preimages:            make(map[common.Hash][]byte),
		addressToObjectIndex: make(map[common.Address]int),
		accessList:           newAccessList(),
		validRevisions:       []revision{},
	}
}

// WithContext returns a Database with an updated protocol context
func (s *CommitStateDB) WithContext(ctx *action.Context) *CommitStateDB {
	s.ctx = ctx
	return s
}

// GetHeightHash returns the block header hash associated with a given block height and chain epoch number.
func (s *CommitStateDB) GetHeightHash(height uint64) ethcmn.Hash {
	ctx := s.ctx
	bz, _ := ctx.Contracts.Get(HeightHashKey(height))
	if len(bz) == 0 {
		return ethcmn.Hash{}
	}
	return ethcmn.BytesToHash(bz)
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
		newObj.SetBalance(prev.account.Balance())
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (s *CommitStateDB) SubBalance(addr ethcmn.Address, amount *big.Int) {
	so := s.GetOrNewStateObject(addr)
	if so != nil {
		so.SubBalance(amount)
	}
}

// AddBalance adds amount to the account associated with addr.
func (s *CommitStateDB) AddBalance(addr ethcmn.Address, amount *big.Int) {
	so := s.GetOrNewStateObject(addr)
	if so != nil {
		so.AddBalance(amount)
	}
}

// GetBalance retrieves the balance from the given address or 0 if object not
// found.
func (s *CommitStateDB) GetBalance(addr ethcmn.Address) *big.Int {
	so := s.getStateObject(addr)
	if so != nil {
		return so.Balance()
	}
	return big.NewInt(0)
}

// GetNonce returns the nonce (sequence number) for a given account.
func (s *CommitStateDB) GetNonce(addr ethcmn.Address) uint64 {
	so := s.getStateObject(addr)
	if so != nil {
		return so.Nonce()
	}
	return 0
}

// SetNonce sets the nonce (sequence number) of an account.
func (s *CommitStateDB) SetNonce(addr ethcmn.Address, nonce uint64) {
	so := s.GetOrNewStateObject(addr)
	if so != nil {
		so.SetNonce(nonce)
	}
}

// GetCodeHash returns the code hash for a given account.
func (s *CommitStateDB) GetCodeHash(addr ethcmn.Address) ethcmn.Hash {
	so := s.getStateObject(addr)
	if so == nil {
		return ethcmn.Hash{}
	}
	return ethcmn.BytesToHash(so.CodeHash())
}

// GetCode returns the code for a given account.
func (s *CommitStateDB) GetCode(addr ethcmn.Address) []byte {
	so := s.getStateObject(addr)
	if so != nil {
		return so.Code(nil)
	}
	return nil
}

// SetCode sets the code for a given account.
func (s *CommitStateDB) SetCode(addr ethcmn.Address, code []byte) {
	so := s.GetOrNewStateObject(addr)
	if so != nil {
		so.SetCode(ethcrypto.Keccak256Hash(code), code)
	}
}

// GetCodeSize returns the code size for a given account.
func (s *CommitStateDB) GetCodeSize(addr ethcmn.Address) int {
	so := s.getStateObject(addr)
	if so == nil {
		return 0
	}
	if so.code != nil {
		return len(so.code)
	}
	return len(so.Code(nil))
}

// AddRefund adds gas to the refund counter.
func (s *CommitStateDB) AddRefund(gas uint64) {
	// TODO: Add journal
	s.refund += gas
}

// SubRefund removes gas from the refund counter. It will panic if the refund
// counter goes below zero.
func (s *CommitStateDB) SubRefund(gas uint64) {
	// TODO: Add journal
	if gas > s.refund {
		panic("refund counter below zero")
	}
	s.refund -= gas
}

// GetRefund returns the current value of the refund counter.
func (s *CommitStateDB) GetRefund() uint64 {
	return s.refund
}

// GetCommittedState retrieves a value from the given account's committed
// storage.
func (s *CommitStateDB) GetCommittedState(addr ethcmn.Address, hash ethcmn.Hash) ethcmn.Hash {
	so := s.getStateObject(addr)
	if so != nil {
		return so.GetCommittedState(nil, hash)
	}
	return ethcmn.Hash{}
}

// GetState retrieves a value from the given account's storage store.
func (s *CommitStateDB) GetState(addr ethcmn.Address, hash ethcmn.Hash) ethcmn.Hash {
	so := s.getStateObject(addr)
	if so != nil {
		return so.GetState(nil, hash)
	}
	return ethcmn.Hash{}
}

// SetState sets the storage state with a key, value pair for an account.
func (s *CommitStateDB) SetState(addr ethcmn.Address, key, value ethcmn.Hash) {
	so := s.GetOrNewStateObject(addr)
	if so != nil {
		so.SetState(nil, key, value)
	}
}

// Suicide marks the given account as suicided and clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (s *CommitStateDB) Suicide(addr ethcmn.Address) bool {
	so := s.getStateObject(addr)
	if so == nil {
		return false
	}

	// TODO: Add journal
	so.markSuicided()
	so.SetBalance(new(big.Int))

	return true
}

// HasSuicided returns if the given account for the specified address has been
// killed.
func (s *CommitStateDB) HasSuicided(addr ethcmn.Address) bool {
	so := s.getStateObject(addr)
	if so != nil {
		return so.suicided
	}
	return false
}

// Exist reports whether the given account address exists in the state. Notably,
// this also returns true for suicided accounts.
func (s *CommitStateDB) Exist(addr ethcmn.Address) bool {
	return s.getStateObject(addr) != nil
}

// Empty returns whether the state object is either non-existent or empty
// according to the EIP161 specification (balance = nonce = code = 0).
func (s *CommitStateDB) Empty(addr ethcmn.Address) bool {
	so := s.getStateObject(addr)
	return so == nil || so.empty()
}

// PrepareAccessList handles the preparatory steps for executing a state transition with
// regards to both EIP-2929 and EIP-2930:
//
// - Add sender to access list (2929)
// - Add destination to access list (2929)
// - Add precompiles to access list (2929)
// - Add the contents of the optional tx access list (2930)
//
// This method should only be called if Yolov3/Berlin/2929+2930 is applicable at the current number.
func (s *CommitStateDB) PrepareAccessList(sender common.Address, dst *common.Address, precompiles []common.Address, list types.AccessList) {
	s.AddAddressToAccessList(sender)
	if dst != nil {
		s.AddAddressToAccessList(*dst)
		// If it's a create-tx, the destination will be added inside evm.create
	}
	for _, addr := range precompiles {
		s.AddAddressToAccessList(addr)
	}
	for _, el := range list {
		s.AddAddressToAccessList(el.Address)
		for _, key := range el.StorageKeys {
			s.AddSlotToAccessList(el.Address, key)
		}
	}
}

// AddressInAccessList returns true if the given address is in the access list.
func (s *CommitStateDB) AddressInAccessList(addr ethcmn.Address) bool {
	return s.accessList.ContainsAddress(addr)
}

// SlotInAccessList returns true if the given (address, slot)-tuple is in the access list.
func (s *CommitStateDB) SlotInAccessList(addr ethcmn.Address, slot ethcmn.Hash) (bool, bool) {
	return s.accessList.Contains(addr, slot)
}

// AddAddressToAccessList adds the given address to the access list
func (s *CommitStateDB) AddAddressToAccessList(addr ethcmn.Address) {
	if s.accessList.AddAddress(addr) {
		// TODO: Add journal
	}
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (s *CommitStateDB) AddSlotToAccessList(addr ethcmn.Address, slot ethcmn.Hash) {
	addrMod, slotMod := s.accessList.AddSlot(addr, slot)
	if addrMod {
		// In practice, this should not happen, since there is no way to enter the
		// scope of 'address' without having the 'address' become already added
		// to the access list (via call-variant, create, etc).
		// Better safe than sorry, though

		// TODO: Add journal
	}
	if slotMod {
		// TODO: Add journal
	}
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (s *CommitStateDB) RevertToSnapshot(revID int) {
	// find the snapshot in the stack of valid snapshots
	idx := sort.Search(len(s.validRevisions), func(i int) bool {
		return s.validRevisions[i].id >= revID
	})

	if idx == len(s.validRevisions) || s.validRevisions[idx].id != revID {
		panic(fmt.Errorf("revision ID %v cannot be reverted", revID))
	}

	// snapshot := s.validRevisions[idx].journalIndex

	// replay the journal to undo changes and remove invalidated snapshots
	// csdb.journal.revert(csdb, snapshot)
	// TODO: Add journal revert
	s.validRevisions = s.validRevisions[:idx]
}

// Snapshot returns an identifier for the current revision of the state.
func (s *CommitStateDB) Snapshot() int {
	id := s.nextRevisionID
	s.nextRevisionID++

	s.validRevisions = append(
		s.validRevisions,
		revision{
			id: id,
			// TODO: Add journal
			// journalIndex: csdb.journal.length(),
		},
	)
	return id
}

// AddLog adds a new log to the state and sets the log metadata from the state.
func (s *CommitStateDB) AddLog(log *ethtypes.Log) {
	// TODO: Add journal
	// s.journal.append(addLogChange{txhash: s.thash})

	log.TxHash = s.thash
	log.BlockHash = s.bhash
	log.TxIndex = uint(s.txIndex)
	log.Index = s.logSize
	s.logs[s.thash] = append(s.logs[s.thash], log)
	s.logSize++
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (s *CommitStateDB) AddPreimage(hash common.Hash, preimage []byte) {
	if _, ok := s.preimages[hash]; !ok {
		// TODO: Add journal
		// s.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		s.preimages[hash] = pi
	}
}

// ForEachStorage iterates over each storage items, all invoke the provided
// callback on each key, value pair.
// Only used in tests https://github.com/ethereum/go-ethereum/search?q=ForEachStorage
func (s *CommitStateDB) ForEachStorage(addr common.Address, cb func(key, value common.Hash) bool) error {
	return nil
}
