package ethereum

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Contract is our main access point to the Ethereum smart contract we use to lock and redeem ethereum tokens

// newKeyTransactor returns a transactror
func newKeyTransactor(key *ecdsa.PrivateKey) *bind.TransactOpts {
	return bind.NewKeyedTransactor(key)
}

// TODO: embed and implement
// SContract is the stub implementation thereof the
type LRContract interface {
	Lock(*big.Int) (*Transaction, error)
	Redeem(*big.Int) (*Transaction, error)
	ProposeAddValidator(Address) (*Transaction, error)
	ProposeRemoveValidator(Address) (*Transaction, error)
	IsValidator(Address) (bool, error)
	NumValidators() (uint, *Transaction, error)
	Epoch() (uint, *Transaction, error)
	DoNewEpoch() (uint, *Transaction, error)
}
