package vm

import (
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var _ ethvm.StateDB = (*CommitStateDB)(nil)

type CommitStateDB struct {
	contractStore *evm.ContractStore
	accountKeeper balance.AccountKeeper
	logger        *log.Logger
	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash ethcmn.Hash
	txIndex      int

	logs    map[ethcmn.Hash][]*ethtypes.Log
	logSize uint

	preimages           []preimageEntry
	hashToPreimageIndex map[ethcmn.Hash]int // map from hash to the index of the preimages slice

	// array that hold 'live' objects, which will get modified while processing a
	// state transition
	stateObjects         []stateEntry
	addressToObjectIndex map[ethcmn.Address]int // map from address to the index of the state objects slice
	stateObjectsDirty    map[ethcmn.Address]struct{}

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
	journal        *journal
	validRevisions []revision
	nextRevisionID int

	// mutex for state deep copying
	lock sync.Mutex

	// Transaction counter in a block. Used on StateSB's Prepare function.
	// It is reset to 0 every block on EndBlock so there's no point in storing the counter
	// to store or adding it as a field on the EVM genesis state.
	txCount int

	// Bloom bytes generation from the logs
	bloom *big.Int

	// Used for validation in case of fork
	IsEnabled bool
}

// NewCommitStateDB returns a reference to a newly initialized CommitStateDB
// which implements Geth's state.StateDB interface.
//
// CONTRACT: Stores used for state must be cache-wrapped as the ordering of the
// key/value space matters in determining the merkle root.
func NewCommitStateDB(cs *evm.ContractStore, ak balance.AccountKeeper, logger *log.Logger) *CommitStateDB {
	return &CommitStateDB{
		contractStore:        cs,
		accountKeeper:        ak,
		logger:               logger,
		stateObjects:         []stateEntry{},
		preimages:            []preimageEntry{},
		hashToPreimageIndex:  make(map[ethcmn.Hash]int),
		addressToObjectIndex: make(map[ethcmn.Address]int),
		logs:                 make(map[ethcmn.Hash][]*ethtypes.Log),
		stateObjectsDirty:    make(map[ethcmn.Address]struct{}),
		accessList:           newAccessList(),
		journal:              newJournal(),
		validRevisions:       []revision{},
		txCount:              0,
		bloom:                big.NewInt(0),
	}
}

func (s *CommitStateDB) PrintState(height uint64) {
	s.logger.Detail("sbhash", s.bhash)
	s.logger.Detail("sbheight", height)
	s.logger.Detail("stateObjects", s.stateObjects)
	s.logger.Detail("addressToObjectIndex", s.addressToObjectIndex)
	s.logger.Detail("stateObjectsDirty", s.stateObjectsDirty)

	bl := &BlockLogs{
		BlockHash:   s.bhash,
		BlockNumber: height,
		Logs:        s.logs,
	}
	bz, _ := bl.MarshalLogs()
	s.logger.Detail("logs", string(bz))
	s.logger.Detail("journal", s.journal)
	s.logger.Detail("validRevisions", s.validRevisions)
	s.logger.Detail("txCount", s.txCount)
	s.logger.Detail("bloom", s.bloom)
}

func (s *CommitStateDB) GetContractStore() *evm.ContractStore {
	return s.contractStore
}

func (s *CommitStateDB) WithState(state *storage.State) *CommitStateDB {
	s.contractStore.WithState(state)
	s.accountKeeper.WithState(state)
	return s
}

func (s *CommitStateDB) GetAccountKeeper() balance.AccountKeeper {
	return s.accountKeeper
}

func (s *CommitStateDB) GetAvailableGas() uint64 {
	gc := s.contractStore.State.GetCalculator()
	return gc.GetLeft()
}

// Prepare sets the current transaction hash which is
// used when the EVM emits new state logs. (no state change)
func (s *CommitStateDB) Prepare(thash ethcmn.Hash) {
	s.thash = thash
	s.txIndex = s.txCount
	// refreshing it
	s.accessList = newAccessList()
}

// Finality used to increment tx counter, clear staff and etc. (no state change)
func (s *CommitStateDB) Finality() {
	s.txCount++
}

// Error returns the first non-nil error the StateDB encountered.
func (s *CommitStateDB) Error() error {
	return s.dbErr
}

// Finalise finalizes the state objects (accounts) state by setting their state,
// removing the s destructed objects and clearing the journal as well as the
// refunds.
func (s *CommitStateDB) Finalise(deleteEmptyObjects bool) error {
	defer func() {
		s.dbErr = nil
		s.hashToPreimageIndex = make(map[ethcmn.Hash]int)
	}()

	if s.dbErr != nil {
		return fmt.Errorf("commit aborted due to earlier error: %v", s.dbErr)
	}

	for _, dirty := range s.journal.dirties {
		_, exist := s.addressToObjectIndex[dirty.address]
		if !exist {
			continue
		}
		s.logger.Detail("VM: found dirty, add to update queue", dirty.address)
		s.stateObjectsDirty[dirty.address] = struct{}{}
	}

	// set the state objects
	for _, stateEntry := range s.stateObjects {
		_, isDirty := s.stateObjectsDirty[stateEntry.address]

		s.logger.Detail("VM: starting to process finalise entry", stateEntry.address)

		switch {
		case stateEntry.stateObject.suicided || (isDirty && deleteEmptyObjects && stateEntry.stateObject.empty()):
			// If the state object has been removed, don't bother syncing it and just
			// remove it from the store.
			s.deleteStateObject(stateEntry.stateObject)

		case isDirty:
			// Set all the dirty state storage items for the state object in the
			// protocol and finally set the account in the account mapper.
			stateEntry.stateObject.commitState()

			// write any contract code associated with the state object
			if stateEntry.stateObject.code != nil && stateEntry.stateObject.dirtyCode {
				stateEntry.stateObject.commitCode()
				stateEntry.stateObject.dirtyCode = false
			}

			// update the object in the store
			if err := s.updateStateObject(stateEntry.stateObject); err != nil {
				return err
			}
		}
		delete(s.stateObjectsDirty, stateEntry.address)
	}
	// updating bloom filter
	s.UpdateBloom()
	// clear all state objects in order to have a fresh states at the new tx
	s.stateObjects = make([]stateEntry, 0)
	s.addressToObjectIndex = make(map[ethcmn.Address]int)
	// No need as we delete at previous iteration
	// s.stateObjectsDirty = make(map[ethcmn.Address]struct{})
	// invalidate journal because reverting across transactions is not allowed
	s.clearJournalAndRefund()
	return nil
}

// GetHeightHash returns the block header hash associated with a given block height and chain epoch number.
func (s *CommitStateDB) GetHeightHash(height uint64) ethcmn.Hash {
	bz, _ := s.contractStore.Get(evm.KeyPrefixHeightHash, evm.HeightHashKey(height))
	if len(bz) == 0 {
		return ethcmn.Hash{}
	}
	return ethcmn.BytesToHash(bz)
}

// SetHeightHash set hash and height of the block
func (s *CommitStateDB) SetHeightHash(height uint64, hash ethcmn.Hash) {
	s.bhash = hash
	s.IsEnabled = true
	s.contractStore.Set(evm.KeyPrefixHeightHash, evm.HeightHashKey(height), hash.Bytes())
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying account mapper and store keys to avoid reloading data for the
// next operations.
func (s *CommitStateDB) Reset() error {
	s.stateObjects = make([]stateEntry, 0)
	s.addressToObjectIndex = make(map[ethcmn.Address]int)
	s.stateObjectsDirty = make(map[ethcmn.Address]struct{})
	s.thash = ethcmn.Hash{}
	s.bhash = ethcmn.Hash{}
	s.txIndex = 0
	s.logSize = 0
	s.preimages = make([]preimageEntry, 0)
	s.hashToPreimageIndex = make(map[ethcmn.Hash]int)
	s.accessList = newAccessList()
	s.txCount = 0
	s.logs = make(map[ethcmn.Hash][]*ethtypes.Log)
	s.bloom = big.NewInt(0)
	s.dbErr = nil
	s.nextRevisionID = 0
	s.clearJournalAndRefund()
	return nil
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
func (s *CommitStateDB) CreateAccount(addr ethcmn.Address) {
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
	s.journal.append(refundChange{prev: s.refund})
	s.refund += gas
}

// SubRefund removes gas from the refund counter. It will panic if the refund
// counter goes below zero.
func (s *CommitStateDB) SubRefund(gas uint64) {
	s.journal.append(refundChange{prev: s.refund})
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

	s.journal.append(suicideChange{
		account:     &addr,
		prev:        so.suicided,
		prevBalance: new(big.Int).Set(so.Balance()),
	})

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
func (s *CommitStateDB) PrepareAccessList(sender ethcmn.Address, dst *ethcmn.Address, precompiles []ethcmn.Address, list ethtypes.AccessList) {
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
		s.journal.append(accessListAddAccountChange{&addr})
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
		s.journal.append(accessListAddAccountChange{&addr})
	}
	if slotMod {
		s.journal.append(accessListAddSlotChange{
			address: &addr,
			slot:    &slot,
		})
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

	snapshot := s.validRevisions[idx].journalIndex

	// replay the journal to undo changes and remove invalidated snapshots
	s.journal.revert(s, snapshot)
	s.validRevisions = s.validRevisions[:idx]
}

// Snapshot returns an identifier for the current revision of the state.
func (s *CommitStateDB) Snapshot() int {
	id := s.nextRevisionID
	s.nextRevisionID++

	s.validRevisions = append(
		s.validRevisions,
		revision{
			id:           id,
			journalIndex: s.journal.length(),
		},
	)
	return id
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (s *CommitStateDB) AddPreimage(hash ethcmn.Hash, preimage []byte) {
	if _, ok := s.hashToPreimageIndex[hash]; !ok {
		s.journal.append(addPreimageChange{hash: hash})

		pi := make([]byte, len(preimage))
		copy(pi, preimage)

		s.preimages = append(s.preimages, preimageEntry{hash: hash, preimage: pi})
		s.hashToPreimageIndex[hash] = len(s.preimages) - 1
	}
}

// ForEachStorage iterates over each storage items, all invoke the provided
// callback on each key, value pair.
// Only used in tests https://github.com/ethereum/go-ethereum/search?q=ForEachStorage
func (s *CommitStateDB) ForEachStorage(addr ethcmn.Address, cb func(key, value ethcmn.Hash) bool) error {
	so := s.getStateObject(addr)
	if so == nil {
		return nil
	}

	prefixStore := evm.AddressStoragePrefix(so.Address())
	s.contractStore.Iterate(prefixStore, func(keyD []byte, valueD []byte) bool {
		key := ethcmn.BytesToHash(keyD)
		value := ethcmn.BytesToHash(valueD)

		if idx, dirty := so.keyToDirtyStorageIndex[key]; dirty {
			// check if iteration stops
			if cb(key, ethcmn.HexToHash(so.dirtyStorage[idx].Value)) {
				return true
			}
		} else if cb(key, value) {
			return true
		}
		return false
	})
	return nil
}

// GetOrNewStateObject retrieves a state object or create a new state object if
// nil.
func (s *CommitStateDB) GetOrNewStateObject(addr ethcmn.Address) StateObject {
	so := s.getStateObject(addr)
	if so == nil || so.deleted {
		so, _ = s.createObject(addr)
		s.logger.Detail("VM: account created", addr)
	}
	return so
}

// Copy creates a deep, independent copy of the state.
//
// NOTE: Snapshots of the copied state cannot be applied to the copy.
func (s *CommitStateDB) Copy() *CommitStateDB {
	// copy all the basic fields, initialize the memory ones
	state := &CommitStateDB{}
	CopyCommitStateDB(s, state)
	return state
}

func CopyCommitStateDB(from, to *CommitStateDB) {
	to.IsEnabled = from.IsEnabled
	// safe as store states will be changed during process check and deliver
	to.contractStore = from.contractStore
	to.accountKeeper = from.accountKeeper
	to.logger = from.logger
	to.refund = from.refund

	to.thash = from.thash
	to.bhash = from.bhash
	to.txIndex = from.txIndex
	to.logs = make(map[ethcmn.Hash][]*ethtypes.Log, len(from.logs))
	to.logSize = from.logSize

	to.preimages = make([]preimageEntry, len(from.preimages))
	to.hashToPreimageIndex = make(map[ethcmn.Hash]int, len(from.hashToPreimageIndex))

	to.stateObjects = make([]stateEntry, len(from.journal.dirties))
	to.addressToObjectIndex = make(map[ethcmn.Address]int)
	to.stateObjectsDirty = make(map[ethcmn.Address]struct{}, len(from.journal.dirties))

	// no need to copy it as it will be filled by tx, not took too much for execution
	to.accessList = newAccessList()

	to.journal = newJournal()

	// snapshots ignoring
	to.validRevisions = make([]revision, 0)
	to.nextRevisionID = 0

	to.txCount = from.txCount
	to.bloom = new(big.Int).Set(from.bloom)

	// copy the dirty states, logs, and preimages
	for _, dirty := range from.journal.dirties {
		// There is a case where an object is in the journal but not in the
		// stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we
		// need to check for nil.
		//
		// Ref: https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527
		if idx, exist := from.addressToObjectIndex[dirty.address]; exist {
			to.stateObjects = append(to.stateObjects, stateEntry{
				address:     dirty.address,
				stateObject: from.stateObjects[idx].stateObject.deepCopy(to),
			})
			to.addressToObjectIndex[dirty.address] = len(to.stateObjects) - 1
			to.stateObjectsDirty[dirty.address] = struct{}{}
		}
	}

	// Above, we don't copy the actual journal. This means that if the copy is
	// copied, the loop above will be a no-op, since the copy's journal is empty.
	// Thus, here we iterate over stateObjects, to enable copies of copies.
	for addr := range from.stateObjectsDirty {
		if idx, exist := to.addressToObjectIndex[addr]; !exist {
			to.setStateObject(from.stateObjects[idx].stateObject.deepCopy(to))
		}
		to.stateObjectsDirty[addr] = struct{}{}
	}

	for hash, logs := range from.logs {
		cpy := make([]*types.Log, len(logs))
		for i, l := range logs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		to.logs[hash] = cpy
	}

	// copy pre-images
	for i, preimage := range from.preimages {
		to.preimages[i] = preimage
		to.hashToPreimageIndex[preimage.hash] = i
	}
}
