package keys

import (
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// Code is account Code type alias
type Code []byte

func (c Code) String() string {
	return string(c)
}

type EthAccount struct {
	Address  Address
	Nonce    uint64
	CodeHash []byte
	Coins    *big.Int
	Sequence uint64
}

func NewEthAccount(addr Address) *EthAccount {
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
