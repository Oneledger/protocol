package accounts

import (
	"encoding/json"
	"errors"
	"github.com/Oneledger/protocol/data/keys"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	STATUS_OPEN   = 44
	STATUS_CLOSED = 45
)

var (
	errorWalletClosed = errors.New("error: wallet is closed")
)

// WalletKeyStore keeps a session storage of accounts on the Full Node
type WalletKeyStore struct {
	keyStorePath string
	keyStore     *keys.KeyStore
	accounts     []keys.Address
	status       int
	sessionKey   string
}

var _ Wallet = &WalletKeyStore{}

func (wks *WalletKeyStore) Accounts() []Account {
	return []Account{}
}

func (wks *WalletKeyStore) ListAddresses() ([]keys.Address, error) {
	files, err := ioutil.ReadDir(wks.keyStorePath)
	if err != nil {
		return []keys.Address{}, err
	}

	for _, f := range files {
		address := keys.Address{}
		err := address.UnmarshalText([]byte(f.Name()))
		if err != nil {
			return []keys.Address{}, err
		}
		wks.accounts = append(wks.accounts, address)
	}

	return wks.accounts, nil
}

func (wks *WalletKeyStore) Add(account Account) error {
	if wks.status == STATUS_OPEN {
		data, err := json.Marshal(account)
		if err != nil {
			return err
		}
		err = wks.keyStore.SaveKeyData(wks.keyStorePath, account.Address(), data, wks.sessionKey)
		if err != nil {
			return err
		}

		return nil
	} else {
		return errorWalletClosed
	}
}

func (wks *WalletKeyStore) Delete(account Account) error {
	if wks.status == STATUS_OPEN {
		return wks.keyStore.DeleteKey(wks.keyStorePath, account.Address(), wks.sessionKey)
	} else {
		return errorWalletClosed
	}
}

func (wks *WalletKeyStore) GetAccount(address keys.Address) (Account, error) {
	account := Account{}
	if wks.status == STATUS_OPEN {
		data, err := wks.keyStore.GetKeyData(wks.keyStorePath, address, wks.sessionKey)
		if err != nil {
			return account, err
		}
		err = json.Unmarshal(data, &account)

		return account, err
	} else {
		return account, errorWalletClosed
	}
}

//Function does not apply to Secure Wallet.
func (wks *WalletKeyStore) SignWithAccountIndex([]byte, int) (keys.PublicKey, []byte, error) {
	return keys.PublicKey{}, nil, nil
}

func (wks *WalletKeyStore) SignWithAddress(data []byte, address keys.Address) (keys.PublicKey, []byte, error) {
	if wks.status == STATUS_OPEN {
		account, err := wks.GetAccount(address)
		if err != nil {
			return keys.PublicKey{}, nil, err
		}
		signature, err := account.Sign(data)
		if err != nil {
			return keys.PublicKey{}, nil, err
		}

		return *account.PublicKey, signature, err
	} else {
		return keys.PublicKey{}, nil, errorWalletClosed
	}
}

func (wks *WalletKeyStore) VerifyPassphrase(address keys.Address, passphrase string) (bool, error) {
	return wks.keyStore.VerifyPassphrase(wks.keyStorePath, address, passphrase)
}

func (wks *WalletKeyStore) Open(address keys.Address, passphrase string) {
	if wks.keyStore.KeyExists(wks.keyStorePath, address) {
		if res, _ := wks.VerifyPassphrase(address, passphrase); res {
			wks.status = STATUS_OPEN
			wks.sessionKey = passphrase
		}
	} else {
		wks.status = STATUS_OPEN
		wks.sessionKey = passphrase
	}
}

func (wks *WalletKeyStore) Close() {
	wks.status = STATUS_CLOSED
	wks.sessionKey = ""
}

func NewWalletKeyStore(path string) (*WalletKeyStore, error) {
	if path == "" {
		return nil, errors.New("walletKeystore: invalid path")
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, errors.New("walletKeystore: " + err.Error())
	}

	//If the path doesn't already exist, then create it.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return nil, errors.New("walletKeystore: error creating path" + err.Error())
		}
	}

	return &WalletKeyStore{
		keyStorePath: absPath + "/",
		keyStore:     keys.NewKeyStore(),
		accounts:     nil,
		status:       STATUS_CLOSED,
		sessionKey:   "",
	}, nil
}
