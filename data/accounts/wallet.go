package accounts

import (
	"fmt"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
)

type Wallet interface {
	//returns the account that the wallet holds
	Accounts() []Account

	Add(Account) error

	Delete(Account) error

	GetAccount(address keys.Address) (Account, error)

	SignWithAccountIndex([]byte, int) (keys.PublicKey, []byte, error)

	SignWithAddress([]byte, keys.Address) (keys.PublicKey, []byte, error)

	Close()
}

// WalletStore keeps a session storage of accounts on the Full Node
type WalletStore struct {
	store storage.SessionedStorage

	accounts []storage.StoreKey
}

// WalletStore satisfies the Wallet interface
var _ Wallet = &WalletStore{}

// Accounts returns all the accounts in the wallet store
func (ws WalletStore) Accounts() []Account {

	accounts := make([]Account, len(ws.accounts))
	for i, key := range ws.accounts {

		acc, err := ws.store.Get(key)
		if err != nil {
			logger.Error("account not exist anymore")
		}

		var account Account
		err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(acc, &account)
		if err != nil {
			logger.Error("failed to deserialize account")
		}

		accounts[i] = account
	}
	return accounts

}

func (ws *WalletStore) Add(account Account) error {
	session := ws.store.BeginSession()

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
	session := ws.store.BeginSession()

	exist := session.Exists(account.Address().Bytes())
	if !exist {
		return errors.New("account already exist: " + string(account.Address()))
	}

	_, err := session.Delete(account.Address().Bytes())
	if err == nil {
		session.Commit()
	}

	return err
}

func (ws WalletStore) GetAccount(address keys.Address) (Account, error) {
	value, err := ws.store.Get(address.Bytes())
	if err != nil || len(value) == 0 {
		return Account{}, fmt.Errorf("failed to get account by address: %s", address.String())
	}

	var account = &Account{}
	account = account.FromBytes(value)

	return *account, nil
}

func (ws WalletStore) SignWithAccountIndex(msg []byte, index int) (keys.PublicKey, []byte, error) {
	if index >= len(ws.accounts) {
		return keys.PublicKey{}, nil, fmt.Errorf("account index out of range")
	}
	return ws.SignWithAddress(msg, ws.accounts[index].Bytes())
}

func (ws WalletStore) SignWithAddress(msg []byte, address keys.Address) (keys.PublicKey, []byte, error) {
	account, err := ws.GetAccount(address)
	if err != nil {
		return keys.PublicKey{}, nil, errors.Wrap(err, "sign with address")
	}

	signed, err := account.Sign(msg)
	return *account.PublicKey, signed, err
}

func NewWallet(config config.Server, dbDir string) Wallet {
	store := storage.NewStorageDB(storage.KEYVALUE, "accounts", dbDir, config.Node.DB)

	accountKeys := store.BeginSession().FindAll()

	accounts := make([]storage.StoreKey, len(accountKeys))
	for i, key := range accountKeys {
		accounts[i] = key
	}

	return &WalletStore{
		store,
		accounts,
	}
}

func (ws WalletStore) Close() {
	ws.store.Close()
}
