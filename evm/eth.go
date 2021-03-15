package evm

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// EthAccount implements the keys.Account interface and embeds with code hash for
// the contract
type EthAccount struct {
	Address  keys.Address
	CodeHash []byte
	Coins    map[string]balance.Coin
	Sequence uint64
}

func NewEthAccount(addr keys.Address) *EthAccount {
	return &EthAccount{
		Address:  addr,
		CodeHash: ethcrypto.Keccak256(nil),
		Coins:    make(map[string]balance.Coin),
	}
}

// EthAddress returns the account address ethereum format.
func (acc EthAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

func (acc EthAccount) Balance(currency string) balance.Coin {
	return acc.Coins[currency]
}

func (acc *EthAccount) AddBalance(coin balance.Coin) {
	balance := acc.Coins[coin.Currency.Name]
	acc.Coins[coin.Currency.Name] = balance.Plus(coin)
}

func (acc *EthAccount) SubBalance(coin balance.Coin) {
	balance := acc.Coins[coin.Currency.Name]
	newCoin, _ := balance.Minus(coin)
	acc.Coins[coin.Currency.Name] = newCoin
}

func (acc *EthAccount) SetBalance(coin balance.Coin) {
	acc.Coins[coin.Currency.Name] = coin
}
