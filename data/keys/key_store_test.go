package keys

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var (
	keyStore        KeyStore
	path            string
	secretData      string
	filename        string
	address         Address
	address1        Address
	passphrase      string
	wrongPassPhrase string
)

func init() {
	keyStore = KeyStore{}
	path = "keystore/"

	pub, _, _ := NewKeyPairFromTendermint()
	h, _ := pub.GetHandler()
	address = h.Address()

	pub1, _, _ := NewKeyPairFromTendermint()
	h1, _ := pub1.GetHandler()
	address1 = h1.Address()

	secretData = "My Secret Data"
	passphrase = "password"
	wrongPassPhrase = "wrong password"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return
		}
	}

	filename, _ = filepath.Abs(path + address.Humanize())
}

func TestKeyStore_SaveKeyData(t *testing.T) {
	err := keyStore.SaveKeyData(path, address, []byte(secretData), passphrase)
	assert.Equal(t, nil, err)
}

func TestKeyStore_GetKeyData(t *testing.T) {
	plainText, err := keyStore.GetKeyData(path, address, passphrase)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, plainText)
}

func TestKeyStore_VerifyPassphrase(t *testing.T) {
	result, err := keyStore.VerifyPassphrase(path, address, passphrase)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, result)
}

func TestKeyStore_VerifyPassphrase2(t *testing.T) {
	result, err := keyStore.VerifyPassphrase(path, address, wrongPassPhrase)
	assert.NotEqual(t, nil, err)
	assert.NotEqual(t, true, result)
}

func TestKeyStore_DeleteKey(t *testing.T) {
	err := keyStore.SaveKeyData(path, address1, []byte(secretData), passphrase)
	assert.Equal(t, nil, err)

	err = keyStore.DeleteKey(path, address1, passphrase)
	assert.Equal(t, nil, err)
}

func TestKeyStore_TearDown(t *testing.T) {
	os.Remove(filename)
	os.RemoveAll(path)
}
