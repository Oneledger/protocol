package accounts

import (
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"

)

type Wallet interface {
	//returns the account that the wallet holds
	Accounts() []Account

	Add(Account) bool

	Delete(Account) bool

	Sign ([]byte, Account) []byte
}


type WalletStore struct {
	store storage.Storage

	accounts []keys.Address
}

func (WalletStore) Accounts() []Account {
	panic("implement me")
}

func (WalletStore) Add(Account) bool {
	panic("implement me")
}

func (WalletStore) Delete(Account) bool {
	panic("implement me")
}

func (WalletStore) Sign([]byte, Account) []byte {
	panic("implement me")
}

func NewWallet(config config.Server) WalletStore {
	ctx := storage.Context{
		//todo: get the database path
		DbDir: "dbpath",
		ConfigDB: config.Node.DB,
	}

	store := storage.NewSessionStorage(storage.KEYVALUE, "accounts", ctx)

	accountKeys := store.Begin().FindAll()

	accounts := make([]keys.Address, len(accountKeys))
	for i, key := range accountKeys {
		accounts[i] = key.Bytes()
	}

	return WalletStore{
		store,
		accounts,
	}
}




