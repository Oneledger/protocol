package accounts

import (
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
)


type Account struct {
	Type chain.Type `json:"type"`
	Name string      `json:"name"`

	PublicKey  keys.PublicKey  `json:"publicKey"`
	PrivateKey keys.PrivateKey `json:"privateKey"`
}

func (acc Account) Address() keys.Address {
	handler, err := acc.PublicKey.GetHandler()
	if err != nil {
		logger.Fatal("PublicKey format for account is wrong", err)
	}

	return keys.Address(handler.Address())
}




