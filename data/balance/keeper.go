package balance

import (
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type EthAccount struct {
	Address  keys.Address `json:"address"`
	CodeHash []byte       `json:"codeHash"`
	Coins    Coin         `json:"coins"`
	Sequence uint64       `json:"sequence"`
}

func NewEthAccount(addr keys.Address, coin Coin) *EthAccount {
	return &EthAccount{
		Address:  addr,
		CodeHash: ethcrypto.Keccak256(nil),
		Coins:    coin,
	}
}

// EthAddress returns the account address ethereum format.
func (acc EthAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

func (acc EthAccount) Balance() *big.Int {
	if acc.Coins == (Coin{}) {
		return big.NewInt(0)
	}
	return acc.Coins.Amount.BigInt()
}

func (acc *EthAccount) AddBalance(amount *big.Int) {
	coin := acc.Coins.Plus(Coin{
		Currency: acc.Coins.Currency,
		Amount:   NewAmountFromBigInt(amount),
	})
	acc.Coins = coin
}

func (acc *EthAccount) SubBalance(amount *big.Int) {
	coin, err := acc.Coins.Minus(Coin{
		Currency: acc.Coins.Currency,
		Amount:   NewAmountFromBigInt(amount),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to minus balance: %s", err))
	}
	acc.Coins = coin
}

func (acc *EthAccount) SetBalance(amount *big.Int) {
	acc.Coins = Coin{
		Currency: acc.Coins.Currency,
		Amount:   NewAmountFromBigInt(amount),
	}
}

type AccountKeeper interface {
	NewAccountWithAddress(addr keys.Address) (*EthAccount, error)
	GetAccount(addr keys.Address) (*EthAccount, error)
	GetVersionedAccount(addr keys.Address, height int64) (*EthAccount, error)
	SetAccount(account EthAccount) error
	RemoveAccount(account EthAccount)
	GetNonce(addr keys.Address) uint64
	GetBalance(addr keys.Address) *big.Int
	WithState(state *storage.State) AccountKeeper
}

var _ AccountKeeper = (*NesterAccountKeeper)(nil)

// NesterAccountKeeper is used to combine two stores - balance and nonces
type NesterAccountKeeper struct {
	balances   *Store
	currencies *CurrencySet
	state      *storage.State
	prefix     []byte
	logger     *log.Logger

	// as we use two different prefixes, for preventing race conditions when:
	// 1) Creating, fetching accounts;
	// 2) Balance versioned check as hit map which is not concurrent;
	mu  sync.Mutex
	rmu sync.RWMutex
}

func NewNesterAccountKeeper(state *storage.State, balances *Store, currencies *CurrencySet) AccountKeeper {
	return &NesterAccountKeeper{
		balances:   balances,
		currencies: currencies,
		state:      state,
		prefix:     storage.Prefix("keeper"),
		logger:     log.NewLoggerWithPrefix(os.Stdout, "account_keeper"),
	}
}

func (nak *NesterAccountKeeper) WithState(state *storage.State) AccountKeeper {
	nak.balances.WithState(state)
	nak.state = state
	return nak
}

func (nak *NesterAccountKeeper) NewAccountWithAddress(addr keys.Address) (*EthAccount, error) {
	coin, err := nak.getOrCreateCurrencyBalance(addr, nil)
	if err != nil {
		return nil, errors.Errorf("Failed to get balance: %s", err)
	}
	return NewEthAccount(addr, coin), nil
}

func (nak *NesterAccountKeeper) getOrCreateCurrencyBalance(addr keys.Address, height *int64) (Coin, error) {
	nak.rmu.RLock()
	defer nak.rmu.RUnlock()

	currency, ok := nak.currencies.GetCurrencyByName("OLT")
	if !ok {
		return Coin{}, errors.Errorf("Failed to get currency OLT")
	}
	var coin Coin
	if height != nil {
		coin, _ = nak.balances.GetVersionedBalanceForCurr(addr, *height, &currency)
	} else {
		coin, _ = nak.balances.GetBalanceForCurr(addr, &currency)
	}
	if coin == (Coin{}) {
		coin = Coin{
			Currency: currency,
			Amount:   NewAmountFromBigInt(big.NewInt(0)),
		}
	}
	return coin, nil
}

// legacyFix for old accounts without balance
func (nak *NesterAccountKeeper) legacyFix(addr keys.Address, coin Coin) (*EthAccount, error) {
	ea := NewEthAccount(addr, coin)
	// if balance is exists like in genesis or other place set, assume that we already have an account as legacy
	// and we need to update it with new one
	if len(ea.Balance().Bits()) != 0 {
		nak.logger.Debug("Legacy acc processing", ea.Address)
		return ea, nil
	}
	return nil, ErrAccountNotFound
}

func (nak *NesterAccountKeeper) GetAccount(addr keys.Address) (*EthAccount, error) {
	nak.mu.Lock()
	defer nak.mu.Unlock()

	prefixKey := append(nak.prefix, addr.Bytes()...)

	dat, err := nak.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil, err
	}

	coin, err := nak.getOrCreateCurrencyBalance(addr, nil)
	if err != nil {
		return nil, err
	}

	if len(dat) == 0 {
		return nak.legacyFix(addr, coin)
	}

	ea := &EthAccount{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ea)
	if err != nil {
		return nil, err
	}
	ea.Coins = coin
	return ea, nil
}

func (nak *NesterAccountKeeper) GetVersionedAccount(addr keys.Address, height int64) (*EthAccount, error) {
	nak.mu.Lock()
	defer nak.mu.Unlock()

	prefixKey := append(nak.prefix, addr.Bytes()...)

	dat := nak.state.GetVersioned(height, storage.StoreKey(prefixKey))

	coin, err := nak.getOrCreateCurrencyBalance(addr, &height)
	if err != nil {
		return nil, err
	}

	if len(dat) == 0 {
		return nak.legacyFix(addr, coin)
	}

	ea := &EthAccount{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ea)
	if err != nil {
		return nil, err
	}

	ea.Coins = coin
	return ea, nil
}

func (nak *NesterAccountKeeper) SetAccount(account EthAccount) error {
	nak.mu.Lock()
	defer nak.mu.Unlock()

	// cache and not store this balance in account as we have another store
	coins := account.Coins
	account.Coins = Coin{}

	prefixKey := append(nak.prefix, account.Address.Bytes()...)
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(&account)
	if err != nil {
		return errors.Errorf("Failed to serialize: %s", err)
	}
	err = nak.state.Set(storage.StoreKey(prefixKey), dat)
	if err != nil {
		return errors.Errorf("Failed to update storage for account: %s", err)
	}
	err = nak.balances.SetBalance(account.Address, coins)
	if err != nil {
		return errors.Errorf("Failed to set balance: %s", err)
	}
	return nil
}

func (nak *NesterAccountKeeper) RemoveAccount(account EthAccount) {
	prefixKey := append(nak.prefix, account.Address.Bytes()...)
	nak.state.Delete(prefixKey)
}

func (nak *NesterAccountKeeper) GetNonce(addr keys.Address) uint64 {
	acc, err := nak.GetAccount(addr)
	if err != nil {
		return 0
	}
	return acc.Sequence
}
func (nak *NesterAccountKeeper) GetBalance(addr keys.Address) *big.Int {
	coin, err := nak.getOrCreateCurrencyBalance(addr, nil)
	if err != nil {
		return big.NewInt(0)
	}
	return coin.Amount.BigInt()
}
