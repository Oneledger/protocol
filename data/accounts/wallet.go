package accounts

import (
	"bytes"
	"fmt"
	"github.com/Oneledger/protocol/serialize"
	"sync"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
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

const rootkey = "rootkey"

// WalletStore keeps a session storage of accounts on the Full Node
type WalletStore struct {
	store storage.SessionedStorage

	accounts []storage.StoreKey

	sync.Mutex
}

var ErrGetAccountByAddress = errors.New("account not found for address")

// WalletStore satisfies the Wallet interface
var _ Wallet = &WalletStore{}

// Accounts returns all the accounts in the wallet store
func (ws WalletStore) Accounts() []Account {

	accounts := make([]Account, len(ws.accounts))
	for i, key := range ws.accounts {
		account, err := ws.GetAccount(key.Bytes())
		if err != nil {
			logger.Error(ErrGetAccountByAddress.Error())
			continue
		}
		//for safety reason, don't pass the private key
		account.PrivateKey = nil
		accounts[i] = account
	}
	return accounts

}

func (ws *WalletStore) Add(account Account) error {
	ws.Lock()
	defer ws.Unlock()

	session := ws.store.BeginSession()

	exist := session.Exists(account.Address().Bytes())
	if exist {
		return errors.New("account already exist: " + account.Address().String())
	}
	value := account.Bytes()
	err := session.Set(account.Address().Bytes(), value)
	if err != nil {
		return errors.Wrap(err, "failed to set the new account")
	}

	ws.accounts = append(ws.accounts, account.Address().Bytes())
	data, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(ws.accounts)
	if err != nil {
		return err
	}
	err = session.Set([]byte(rootkey), data)
	if err != nil {
		return err
	}
	session.Commit()
	return nil
}

func (ws *WalletStore) Delete(account Account) error {
	ws.Lock()
	defer ws.Unlock()

	session := ws.store.BeginSession()

	exist := session.Exists(account.Address().Bytes())
	if !exist {
		return errors.New("account not exist: " + account.Address().String())
	}

	_, err := session.Delete(account.Address().Bytes())
	if err == nil {
		for i, addr := range ws.accounts {
			if bytes.Equal(addr, account.Address().Bytes()) {
				ws.accounts = append(ws.accounts[:i], ws.accounts[i+1:]...)
				break
			}
		}
		data, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(ws.accounts)
		if err != nil {
			return err
		}
		err = session.Set([]byte(rootkey), data)
		if err != nil {
			return err
		}
		session.Commit()
	}
	return err
}

func (ws WalletStore) getAccount(address keys.Address) (Account, error) {
	value, err := ws.store.Get(address.Bytes())
	if err != nil || len(value) == 0 {
		return Account{}, fmt.Errorf("failed to get account by address: %s", address.String())
	}

	var account = &Account{}
	account = account.FromBytes(value)
	return *account, nil
}

func (ws WalletStore) GetAccount(address keys.Address) (Account, error) {
	account, err := ws.getAccount(address)
	if err != nil {
		return Account{}, err
	}
	//for safety reason, don't pass the private key
	account.PrivateKey = nil
	return account, nil
}

func (ws WalletStore) SignWithAccountIndex(msg []byte, index int) (keys.PublicKey, []byte, error) {
	if index >= len(ws.accounts) {
		return keys.PublicKey{}, nil, fmt.Errorf("account index out of range")
	}
	return ws.SignWithAddress(msg, ws.accounts[index].Bytes())
}

func (ws WalletStore) SignWithAddress(msg []byte, address keys.Address) (keys.PublicKey, []byte, error) {
	account, err := ws.getAccount(address)
	if err != nil {
		return keys.PublicKey{}, nil, err
	}
	signed, err := account.Sign(msg)
	return *account.PublicKey, signed, err
}

func NewWallet(config config.Server, dbDir string) Wallet {
	store := storage.NewStorageDB(storage.KEYVALUE, "accounts", dbDir, config.Node.DB)

	accounts := make([]storage.StoreKey, 0, 10)
	data, err := store.Get([]byte(rootkey))
	if err == nil {
		_ = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(data, &accounts)
	}
	return &WalletStore{
		store,
		accounts,
		sync.Mutex{},
	}
}

func (ws WalletStore) Close() {
	ws.store.Close()
}
