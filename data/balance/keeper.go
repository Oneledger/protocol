package balance

import (
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type EthAccount struct {
	Address  keys.Address
	Nonce    uint64
	CodeHash []byte
	Coins    Coin
	Sequence uint64
}

func NewEthAccount(addr keys.Address) *EthAccount {
	return &EthAccount{
		Address:  addr,
		Nonce:    0,
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
	NewAccountWithAddress(addr keys.Address) *EthAccount
	GetAccount(addr keys.Address) *EthAccount
	SetAccount(account EthAccount)
	RemoveAccount(account EthAccount)
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

func (nak *NesterAccountKeeper) NewAccountWithAddress(addr keys.Address) *EthAccount {
	acc := NewEthAccount(addr)
	fmt.Printf("New account: %+v", acc)
	nak.SetAccount(*acc)
	return nak.GetAccount(addr)
}

func (nak *NesterAccountKeeper) GetAccount(addr keys.Address) *EthAccount {
	prefixKey := append(nak.prefix, addr.Bytes()...)

	dat, err := nak.state.Get(storage.StoreKey(prefixKey))
	if err != nil {
		return nil
	}

	ea := &EthAccount{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, ea)
	if err != nil {
		return nil
	}

	balance, err := nak.balances.GetBalance(addr, nak.currencies)
	if err != nil {
		panic(fmt.Sprintf("Failed to get balance: %s", err))
	}
	ea.Coins = balance.Amounts["OLT"]
	return ea
}

func (nak *NesterAccountKeeper) SetAccount(account EthAccount) {
	prefixKey := append(nak.prefix, account.Address.Bytes()...)
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(account)
	if err != nil {
		panic(fmt.Sprintf("Failed to serialize: %s", err))
	}
	err = nak.state.Set(storage.StoreKey(prefixKey), dat)
	if err != nil {
		panic(fmt.Sprintf("Failed to set account: %s", err))
	}

	if account.Coins != (Coin{}) {
		err = nak.balances.SetBalance(account.Address, account.Coins)
		if err != nil {
			panic(fmt.Sprintf("Failed to set balance: %s", err))
		}
		// mark as zero as we do not use balances here as storage
		account.Coins = Coin{}
	}
}

func (nak *NesterAccountKeeper) RemoveAccount(account EthAccount) {
	prefixed := append(nak.prefix, account.Address.Bytes()...)
	nak.state.Delete(prefixed)
}
