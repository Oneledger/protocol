package storage

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

var (
	_ StateObject = (*stateObject)(nil)

	emptyCodeHash = ethcrypto.Keccak256(nil)
)

// StateObject interface for interacting with state object
type StateObject interface {
	// GetCommittedState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash
	// GetState(db ethstate.Database, key ethcmn.Hash) ethcmn.Hash
	// SetState(db ethstate.Database, key, value ethcmn.Hash)

	Code(db ethstate.Database) []byte
	SetCode(codeHash ethcmn.Hash, code []byte)
	CodeHash() []byte

	AddBalance(amount *big.Int)
	SubBalance(amount *big.Int)
	SetBalance(amount *big.Int)

	Balance() *big.Int
	ReturnGas(gas *big.Int)
	Address() ethcmn.Address

	SetNonce(nonce uint64)
	Nonce() uint64
}

// stateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call CommitTrie to write the modified storage trie into a database.
type stateObject struct {
	code keys.Code // contract bytecode, which gets set when code is loaded
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	originStorage Storage // Storage cache of original entries to dedup rewrites
	dirtyStorage  Storage // Storage entries that need to be flushed to disk

	// DB error
	dbErr   error
	stateDB *CommitStateDB
	account *keys.EthAccount

	keyToOriginStorageIndex map[ethcmn.Hash]int
	keyToDirtyStorageIndex  map[ethcmn.Hash]int

	address ethcmn.Address

	// cache flags
	//
	// When an object is marked suicided it will be delete from the trie during
	// the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

func newStateObject(db *CommitStateDB, acc accounts.Account) *stateObject {
	protocolAccount := keys.NewEthAccount(acc.Address())

	return &stateObject{
		stateDB:                 db,
		account:                 protocolAccount,
		address:                 protocolAccount.EthAddress(),
		originStorage:           Storage{},
		dirtyStorage:            Storage{},
		keyToOriginStorageIndex: make(map[ethcmn.Hash]int),
		keyToDirtyStorageIndex:  make(map[ethcmn.Hash]int),
	}
}

// Code returns the contract code associated with this object, if any.
func (so *stateObject) Code(_ ethstate.Database) []byte {
	if len(so.code) > 0 {
		return so.code
	}

	if bytes.Equal(so.CodeHash(), emptyCodeHash) {
		return nil
	}

	ctx := so.stateDB.ctx
	store := ctx.State.Get(so.stateDB.storeKey)
	code := store.Get(so.CodeHash())

	if len(code) == 0 {
		so.setError(fmt.Errorf("failed to get code hash %x for address %s", so.CodeHash(), so.Address().String()))
	}

	return code
}

// SetCode sets the state object's code.
func (so *stateObject) SetCode(codeHash ethcmn.Hash, code []byte) {
	so.setCode(codeHash, code)
}

func (so *stateObject) setCode(codeHash ethcmn.Hash, code []byte) {
	so.code = code
	so.account.CodeHash = codeHash.Bytes()
	so.dirtyCode = true
}

// CodeHash returns the state object's code hash.
func (so *stateObject) CodeHash() []byte {
	if so.account == nil || len(so.account.CodeHash) == 0 {
		return emptyCodeHash
	}
	return so.account.CodeHash
}

// AddBalance adds an amount to a state object's balance. It is used to add
// funds to the destination account of a transfer.
func (so *stateObject) AddBalance(coin balance.Coin) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if coin.IsZero() {
		return
	}
	so.AddBalance(coin)
}

// SubBalance removes an amount from the stateObject's balance. It is used to
// remove funds from the origin account of a transfer.
func (so *stateObject) SubBalance(coin balance.Coin) {
	if coin.IsZero() {
		return
	}
	so.SubBalance(coin)
}

// SubBalance removes an amount from the stateObject's balance. It is used to
// remove funds from the origin account of a transfer.
func (so *stateObject) SetBalance(coin balance.Coin) {
	so.SetBalance(coin)
}

// Balance returns the state object's current balance.
func (so *stateObject) Balance() *balance.Amount {
	return so.account.Balance("OLT")
}

// ReturnGas returns the gas back to the origin. Used by the Virtual machine or
// Closures. It performs a no-op.
func (so *stateObject) ReturnGas(gas *big.Int) {}

// Address returns the address of the state object.
func (so stateObject) Address() ethcmn.Address {
	return so.address
}

// SetNonce sets the state object's nonce (i.e sequence number of the account).
func (so *stateObject) SetNonce(nonce uint64) {
	so.setNonce(nonce)
}

func (so *stateObject) setNonce(nonce uint64) {
	if so.account == nil {
		panic("state object account is empty")
	}
	so.account.Sequence = nonce
}

// Nonce returns the state object's current nonce (sequence number).
func (so *stateObject) Nonce() uint64 {
	if so.account == nil {
		return 0
	}
	return so.account.Sequence
}

// setError remembers the first non-nil error it is called with.
func (so *stateObject) setError(err error) {
	if so.dbErr == nil {
		so.dbErr = err
	}
}
