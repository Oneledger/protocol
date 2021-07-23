package balance

import (
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/data/keys"
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
	NewAccountWithAddress(addr keys.Address, setAcc bool) (*EthAccount, error)
	GetOrCreateAccount(addr keys.Address) (*EthAccount, error)
	GetAccount(addr keys.Address) (*EthAccount, error)
	GetVersionedAccount(addr keys.Address, height int64) (*EthAccount, error)
	SetAccount(account EthAccount) error
	RemoveAccount(account EthAccount)
	GetVersionedBalance(addr keys.Address, height int64) (*big.Int, error)
	GetNonce(addr keys.Address) uint64
	GetBalance(addr keys.Address) *big.Int
	GetState() *storage.State
	WithState(state *storage.State) AccountKeeper
}

var _ AccountKeeper = (*NesterAccountKeeper)(nil)

// NesterAccountKeeper is used to combine two stores - balance and nonces
type NesterAccountKeeper struct {
	balances   *Store
	currencies *CurrencySet
	state      *storage.State
	prefix     []byte
}

func NewNesterAccountKeeper(state *storage.State, balances *Store, currencies *CurrencySet) AccountKeeper {
	return &NesterAccountKeeper{
		balances:   balances,
		currencies: currencies,
		state:      state,
		prefix:     storage.Prefix("keeper"),
	}
}

func (nak *NesterAccountKeeper) GetState() *storage.State {
	return nak.state
}

func (nak *NesterAccountKeeper) WithState(state *storage.State) AccountKeeper {
	nak.balances.WithState(state)
	nak.state = state
	return nak
}

func (nak *NesterAccountKeeper) NewAccountWithAddress(addr keys.Address, setAcc bool) (*EthAccount, error) {
	coin, err := nak.getOrCreateCurrencyBalance(addr)
	if err != nil {
		return nil, errors.Errorf("Failed to get balance: %s", err)
	}
	acc := NewEthAccount(addr, coin)
	if setAcc {
		err = nak.SetAccount(*acc)
		if err != nil {
			return nil, errors.Errorf("Failed to set account: %s", err)
		}
		acc, err = nak.GetAccount(addr)
		if err != nil {
			return nil, errors.Errorf("Failed to get account: %s", err)
		}
	}
	return acc, nil
}

func (nak *NesterAccountKeeper) GetVersionedBalance(addr keys.Address, height int64) (*big.Int, error) {
	key := nak.balances.BuildKey(addr, nil)
	data := nak.balances.State.GetVersioned(height, key)
	if len(data) == 0 {
		return nil, errors.Errorf("Balance for address '%s' not found", addr)
	}
	amt := NewAmount(0)
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(data, amt)
	if err != nil {
		return nil, err
	}
	return amt.BigInt(), nil
}

func (nak *NesterAccountKeeper) getOrCreateCurrencyBalance(addr keys.Address) (Coin, error) {
	currency, ok := nak.currencies.GetCurrencyByName("OLT")
	if !ok {
		return Coin{}, errors.Errorf("Failed to get currency OLT")
	}
	coin, _ := nak.balances.GetBalanceForCurr(addr, &currency)
	if coin == (Coin{}) {
		coin = Coin{
			Currency: currency,
			Amount:   NewAmountFromBigInt(big.NewInt(0)),
		}
	}
	return coin, nil
}

func (nak *NesterAccountKeeper) GetOrCreateAccount(addr keys.Address) (*EthAccount, error) {
	eoa, err := nak.GetAccount(addr)
	if err != nil {
		eoa, err = nak.NewAccountWithAddress(addr, true)
		if err != nil {
			return nil, err
		}
	}
	return eoa, nil
}

func (nak *NesterAccountKeeper) GetAccount(addr keys.Address) (*EthAccount, error) {
	prefixKey := append(nak.prefix, addr.Bytes()...)

	dat, err := nak.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil, err
	}

	ea := &EthAccount{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ea)
	if err != nil {
		return nil, err
	}

	coin, err := nak.getOrCreateCurrencyBalance(addr)
	if err != nil {
		return nil, err
	}
	ea.Coins = coin
	return ea, nil
}

func (nak *NesterAccountKeeper) GetVersionedAccount(addr keys.Address, height int64) (*EthAccount, error) {
	prefixKey := append(nak.prefix, addr.Bytes()...)

	dat := nak.state.GetVersioned(height, storage.StoreKey(prefixKey))
	if len(dat) == 0 {
		return nil, errors.New(fmt.Sprintf("Previous state on height '%d' for addr '%s' not found", height, addr))
	}

	ea := &EthAccount{}
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ea)
	if err != nil {
		return nil, err
	}

	coin, err := nak.getOrCreateCurrencyBalance(addr)
	if err != nil {
		return nil, err
	}
	ea.Coins = coin
	return ea, nil
}

func (nak *NesterAccountKeeper) SetAccount(account EthAccount) error {
	prefixKey := append(nak.prefix, account.Address.Bytes()...)
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(&account)
	if err != nil {
		return errors.Errorf("Failed to serialize: %s", err)
	}
	err = nak.state.Set(storage.StoreKey(prefixKey), dat)
	if err != nil {
		return errors.Errorf("Failed to update storage for account: %s", err)
	}
	err = nak.balances.SetBalance(account.Address, account.Coins)
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
	coin, err := nak.getOrCreateCurrencyBalance(addr)
	if err != nil {
		return big.NewInt(0)
	}
	return coin.Amount.BigInt()
}
