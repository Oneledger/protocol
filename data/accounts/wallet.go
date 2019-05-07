package accounts

import (
	"errors"
	"fmt"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

var _ Wallet = &WalletStore{}

type Wallet interface {
	//returns the account that the wallet holds
	Accounts() []Account

	Add(Account) error

	Delete(Account) error

	SignWithAccountIndex([]byte, int) ([]byte, error)

	SignWithAddress([]byte, keys.Address) ([]byte, error)
}

type WalletStore struct {
	store storage.Storage

	accounts []storage.StoreKey
}

func (ws WalletStore) Accounts() []Account {
	accounts := make([]Account, len(ws.accounts))
	for i, key := range ws.accounts {
		acc, err := ws.store.Get(key)
		if err != nil {
			logger.Error("account not exist anymore")
		}
		var account Account
		err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(acc, account)
		if err != nil {
			logger.Error("failed to deserialize account")
		}
		accounts[i] = account
	}
	return accounts

}

func (ws *WalletStore) Add(account Account) error {
	session := ws.store.Begin()
	exist := session.Exists(account.Address().Bytes())
	if exist {
		return errors.New("account already exist: " + string(account.Address()))
	}
	value := account.Bytes()
	err := session.Set(account.Address().Bytes(), value)
	if err != nil {
		return fmt.Errorf("failed to set the new account: %s", err)
	}
	session.Commit()
	ws.accounts = append(ws.accounts, account.Address().Bytes())
	return nil
}

func (ws *WalletStore) Delete(account Account) error {
	session := ws.store.Begin()
	exist := session.Exists(account.Address().Bytes())
	if !exist {
		return errors.New("account already exist: " + string(account.Address()))
	}
	_, err := session.Delete(account.Address().Bytes())
	return err
}

func (ws WalletStore) SignWithAccountIndex(msg []byte, index int) ([]byte, error) {
	if index > len(ws.accounts) {
		return nil, fmt.Errorf("account index out of range")
	}
	return ws.SignWithAddress(msg, ws.accounts[index].Bytes())
}

func (ws WalletStore) SignWithAddress(msg []byte, address keys.Address) ([]byte, error) {
	value, err := ws.store.Get(address.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to get account by address: %s", err)
	}
	var account = &Account{}
	account = account.FromBytes(value)
	return account.Sign(msg)
}

func NewWallet(config config.Server, dbDir string) WalletStore {
	ctx := storage.Context{
		DbDir: dbDir,
		ConfigDB: config.Node.DB,
	}

	store := storage.NewSessionStorageDB(storage.KEYVALUE, "accounts", ctx)

	accountKeys := store.Begin().FindAll()

	accounts := make([]storage.StoreKey, len(accountKeys))
	for i, key := range accountKeys {
		accounts[i] = key
	}

	return WalletStore{
		store,
		accounts,
	}
}
