/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package accounts

import (
	"testing"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/stretchr/testify/assert"
)

func TestNewWallet(t *testing.T) {
	cs := config.Server{}
	cs.Node = config.DefaultNodeConfig()

	w := NewWallet(cs, "/tmp/")

	pubkey, privkey, err := keys.NewKeyPairFromTendermint()
	assert.NoError(t, err)

	acc, err := NewAccount(0, "accountName", &privkey, &pubkey)

	accs := w.Accounts()
	assert.Len(t, accs, 0)

	err = w.Add(acc)
	assert.NoError(t, err)

	err = w.Add(acc)
	assert.Error(t, err)

	acc.PrivateKey = nil
	accs = w.Accounts()
	assert.Len(t, accs, 1)
	assert.Equal(t, accs[0], acc)

	accNew, err := w.GetAccount(acc.Address())
	assert.NoError(t, err)
	assert.Equal(t, accNew, acc)

	accNew, err = w.GetAccount(make([]byte, 20))
	assert.Error(t, err)
	assert.Equal(t, accNew, Account{})

	err = w.Delete(acc)
	assert.NoError(t, err)

	err = w.Delete(acc)
	assert.Error(t, err)

	w.Close()
}

func TestWalletSign(t *testing.T) {
	cs := config.Server{}
	cs.Node = config.DefaultNodeConfig()

	w := NewWallet(cs, "/tmp")

	pubkey, privkey, err := keys.NewKeyPairFromTendermint()
	assert.NoError(t, err)

	acc, err := NewAccount(0, "accountName", &privkey, &pubkey)
	err = w.Add(acc)
	assert.NoError(t, err)

	msg := []byte("iah89230h23uiehd923hg9283hd93h82h3dh238h")
	pubKey, signed, err := w.SignWithAddress(msg, acc.Address())

	assert.NoError(t, err)
	assert.Equal(t, &pubKey, acc.PublicKey)

	// here the index is 1 because there is anaother account present created in previous test function
	pubKey, signedNew, err := w.SignWithAccountIndex(msg, 0)
	assert.NoError(t, err)
	assert.Equal(t, &pubKey, acc.PublicKey)

	assert.Equal(t, signed, signedNew)

	_, _, err = w.SignWithAccountIndex(msg, 1)
	assert.Error(t, err)
	_, _, err = w.SignWithAccountIndex(msg, 2)
	assert.Error(t, err)

	_, _, err = w.SignWithAddress(msg, make([]byte, 20))
	assert.Error(t, err)
}
