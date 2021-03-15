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

func NewEthAccount(addr keys.Address) *EthAccount {
	return &EthAccount{
		Address:  addr,
		CodeHash: ethcrypto.Keccak256(nil),
		Coins:    Coin{},
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
	GetVersionedAccount(height int64, addr keys.Address) (*EthAccount, error)
	SetAccount(account EthAccount) error
	RemoveAccount(account EthAccount)
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

func (nak *NesterAccountKeeper) WithState(state *storage.State) AccountKeeper {
	nak.balances.WithState(state)
	nak.state = state
	return nak
}

func (nak *NesterAccountKeeper) NewAccountWithAddress(addr keys.Address) (*EthAccount, error) {
	acc := NewEthAccount(addr)
	err := nak.SetAccount(*acc)
	if err != nil {
		return nil, errors.Errorf("Failed to set account: %s", err)
	}
	acc, err = nak.GetAccount(addr)
	if err != nil {
		return nil, errors.Errorf("Failed to get account: %s", err)
	}
	return acc, nil
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

	balance, err := nak.balances.GetBalance(addr, nak.currencies)
	if err != nil {
		return nil, err
	}
	ea.Coins = balance.Amounts["OLT"]
	return ea, nil
}

func (nak *NesterAccountKeeper) GetVersionedAccount(height int64, addr keys.Address) (*EthAccount, error) {
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

	balance, err := nak.balances.GetBalance(addr, nak.currencies)
	if err != nil {
		return nil, err
	}
	ea.Coins = balance.Amounts["OLT"]
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

	if account.Coins != (Coin{}) {
		err = nak.balances.SetBalance(account.Address, account.Coins)
		if err != nil {
			return errors.Errorf("Failed to set balance: %s", err)
		}
		// mark as zero as we do not use balances here as storage
		account.Coins = Coin{}
	}
	return nil
}

func (nak *NesterAccountKeeper) RemoveAccount(account EthAccount) {
	prefixKey := append(nak.prefix, account.Address.Bytes()...)
	nak.state.Delete(prefixKey)
}
