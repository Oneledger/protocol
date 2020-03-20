package accounts

import (
	"fmt"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	numAddresses = 6
	path         = "keystore/"
)

var (
	accounts       []Account
	passwords      []string
	walletKeyStore *WalletKeyStore
	err            error
)

func init() {
	for i := 0; i < numAddresses; i++ {
		pub, priv, err := keys.NewKeyPairFromTendermint()
		if err != nil {
			break
		}
		account, err := NewAccount(chain.ETHEREUM, string(i), &priv, &pub)
		if err != nil {
			break
		}
		accounts = append(accounts, account)
	}

	for k := 0; k < numAddresses; k++ {
		passwords = append(passwords, "password"+string(k))
	}

	walletKeyStore, err = NewWalletKeyStore(path)
	if err != nil {
		return
	}

	walletKeyStore.keyStore = keys.NewKeyStore()
}

func TestWalletKeyStore_Open(t *testing.T) {
	walletKeyStore.Open(accounts[0].Address(), passwords[0])
	walletKeyStore.Close()
}

func TestWalletKeyStore_Add(t *testing.T) {
	for i := 0; i < numAddresses; i++ {
		walletKeyStore.Open(accounts[i].Address(), passwords[i])

		err := walletKeyStore.Add(accounts[i])
		assert.Equal(t, nil, err)

		walletKeyStore.Close()
	}
}

func TestWalletKeyStore_ListAddresses(t *testing.T) {
	addresses, err := walletKeyStore.ListAddresses()
	assert.Equal(t, nil, err)

	for i := 0; i < numAddresses; i++ {
		fmt.Println("address" + string(i) + ": " + addresses[i].Humanize())
	}
}

func TestWalletKeyStore_GetAccount(t *testing.T) {
	for i := 0; i < numAddresses; i++ {
		walletKeyStore.Open(accounts[i].Address(), passwords[i])

		_, err := walletKeyStore.GetAccount(accounts[i].Address())
		assert.Equal(t, nil, err)
		//fmt.Println(account)

		walletKeyStore.Close()
	}
}

func TestWalletKeyStore_SignWithAddress(t *testing.T) {
	walletKeyStore.Open(accounts[0].Address(), passwords[0])
	pub, sig, err := walletKeyStore.SignWithAddress([]byte("MY TRANSACTION DATA"), accounts[0].Address())
	assert.Equal(t, nil, err)
	assert.NotEqual(t, keys.PublicKey{}, pub)
	assert.NotEqual(t, nil, sig)
	walletKeyStore.Close()
}

func TestWalletKeyStore_VerifyPassphrase(t *testing.T) {
	walletKeyStore.Open(accounts[1].Address(), passwords[1])
	res, err := walletKeyStore.VerifyPassphrase(accounts[1].Address(), passwords[1])
	assert.Equal(t, true, res)
	assert.Equal(t, nil, err)

	res, err = walletKeyStore.VerifyPassphrase(accounts[1].Address(), passwords[0])
	assert.Equal(t, false, res)
	assert.NotEqual(t, nil, err)

	walletKeyStore.Close()
}

func TestWalletKeyStore_Delete(t *testing.T) {
	for i := 0; i < numAddresses; i++ {
		walletKeyStore.Open(accounts[i].Address(), passwords[i])

		err := walletKeyStore.Delete(accounts[i].Address())
		assert.Equal(t, nil, err)

		walletKeyStore.Close()
	}

	_ = os.RemoveAll(walletKeyStore.keyStorePath)
}
