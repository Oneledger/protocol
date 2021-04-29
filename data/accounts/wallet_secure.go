package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Oneledger/protocol/data/keys"
)

const (
	STATUS_OPEN     = 44
	STATUS_CLOSED   = 45
	SESSION_TIMEOUT = 30
)

var (
	errorWalletClosed = errors.New("error: wallet is closed")
)

// WalletKeyStore keeps a session storage of accounts on the Full Node
type WalletKeyStore struct {
	keyStorePath   string
	keyStore       *keys.KeyStore
	accounts       []keys.Address
	status         int
	sessionKey     string
	sessionClose   chan error
	sessionTimeout chan error
}

func (wks *WalletKeyStore) ListAddresses() ([]keys.Address, error) {
	files, err := ioutil.ReadDir(wks.keyStorePath)
	if err != nil {
		return []keys.Address{}, err
	}

	//Clear current slice of accounts.
	wks.accounts = nil

	//Need to range over all files in directory, Cannot depend on number of files to determine if list needs to be updated.
	for _, f := range files {
		address, err := wks.keyStore.GetAddress(wks.keyStorePath, f.Name())
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
	}

	return errorWalletClosed
}

func (wks *WalletKeyStore) Delete(address keys.Address) error {
	if wks.status == STATUS_OPEN {
		return wks.keyStore.DeleteKey(wks.keyStorePath, address, wks.sessionKey)
	}
	return errorWalletClosed
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
	}
	return account, errorWalletClosed
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
	}

	return keys.PublicKey{}, nil, errorWalletClosed
}

func (wks *WalletKeyStore) VerifyPassphrase(address keys.Address, passphrase string) (bool, error) {
	return wks.keyStore.VerifyPassphrase(wks.keyStorePath, address, passphrase)
}

func (wks *WalletKeyStore) Open(address keys.Address, passphrase string) bool {
	if wks.keyStore.KeyExists(wks.keyStorePath, address) {
		if res, _ := wks.VerifyPassphrase(address, passphrase); !res {
			return false
		}
	}

	//Validate input
	if address == nil || passphrase == "" {
		return false
	}

	//Can't open wallet if its already open.
	if wks.status == STATUS_OPEN {
		return false
	}

	wks.status = STATUS_OPEN
	wks.sessionKey = passphrase

	//Go routine to trigger timeout if session is opened too long.
	go func(wks *WalletKeyStore) {
		time.Sleep(SESSION_TIMEOUT * time.Second)
		if wks.status == STATUS_CLOSED {
			return
		}
		wks.sessionTimeout <- nil
	}(wks)

	//Go routine to handle any channel signals.
	go func(wks *WalletKeyStore) {
		select {
		case <-wks.sessionTimeout:
			if wks.status == STATUS_OPEN {
				wks.status = STATUS_CLOSED
				wks.sessionKey = ""
				fmt.Println("WALLET SESSION TIMEOUT!! CLOSING WALLET...")
			}
			break
		case <-wks.sessionClose:
		}
	}(wks)

	return true
}

func (wks *WalletKeyStore) KeyExists(address keys.Address) bool {
	return wks.keyStore.KeyExists(wks.keyStorePath, address)
}

func (wks *WalletKeyStore) Close() {
	if wks.status == STATUS_OPEN {
		wks.sessionClose <- nil
	}
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
		keyStorePath:   absPath + "/",
		keyStore:       keys.NewKeyStore(),
		accounts:       nil,
		status:         STATUS_CLOSED,
		sessionKey:     "",
		sessionClose:   make(chan error),
		sessionTimeout: make(chan error),
	}, nil
}
