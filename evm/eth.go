package evm

import (
	"math/big"

	"github.com/Oneledger/protocol/data/keys"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type EthAccount struct {
	Address  keys.Address
	Nonce    uint64
	CodeHash []byte
	Coins    *big.Int
	Sequence uint64
}

func NewEthAccount(addr keys.Address) *EthAccount {
	return &EthAccount{
		Address:  addr,
		Nonce:    0,
		CodeHash: ethcrypto.Keccak256(nil),
		Coins:    big.NewInt(0),
	}
}

// EthAddress returns the account address ethereum format.
func (acc EthAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

func (acc EthAccount) Balance() *big.Int {
	return acc.Coins
}

func (acc *EthAccount) AddBalance(amount *big.Int) {
	acc.Coins = big.NewInt(0).Add(acc.Coins, amount)
}

func (acc *EthAccount) SubBalance(amount *big.Int) {
	acc.Coins = big.NewInt(0).Sub(acc.Coins, amount)
}

func (acc *EthAccount) SetBalance(amount *big.Int) {
	acc.Coins = amount
}
