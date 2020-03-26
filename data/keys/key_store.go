package keys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type AccountKey struct {
	Version    int    `json:"version"`
	CipherText string `json:"cipher_text"`
	Address    string `json:"address"`
}

const (
	EmptyStr = ""
	VERSION  = 1
)

type keystore interface {
	//Hash passphrase so we can get a key that is the correct length for the block cipher.
	createHash(key string) string

	//encrypt data using the given passphrase.
	encrypt(data []byte, passphrase string) ([]byte, error)

	//decrypt data using the given passphrase.
	decrypt(data []byte, passphrase string) ([]byte, error)

	//Encrypt data and store in the given path.
	SaveKeyData(path string, address Address, data []byte, passphrase string) error

	//Open encrypted file and decrypt data using the key generated from the passphrase.
	GetKeyData(path string, address Address, passphrase string) ([]byte, error)

	//Delete data associated with given address
	DeleteKey(path string, address Address, passphrase string) error

	//Verify client has correct password for the given address.
	VerifyPassphrase(path string, address Address, passphrase string) (bool, error)

	//Check if Key exists
	KeyExists(path string, address Address) bool

	//Get address of account given a file name.
	GetAddress(path string, filename string) (Address, error)
}

type KeyStore struct {
}

var _ keystore = &KeyStore{}

func (ks *KeyStore) createHash(key string) string {
	hasher := md5.New()

	_, err := hasher.Write([]byte(key))
	if err != nil {
		return EmptyStr
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func (ks *KeyStore) encrypt(data []byte, passphrase string) ([]byte, error) {

	//Create a key based on the provided passphrase, 32 bytes required for AES.
	key := ks.createHash(passphrase)
	if key == EmptyStr {
		return nil, errors.New("key store.encrypt: hash failed")
	}
	//Create a new AES block cipher based on the key above
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errors.New("Error creating AES block cipher: " + err.Error())
	}
	//Wrap block cipher in Galois Counter Mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("error using Galois Counter Mode: " + err.Error())
	}
	//Creat a nonce compatible with GCM
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("Error Creating nonce: " + err.Error())
	}

	//Prepend nonce to the cipher text to ensure the decrypt function uses the same nonce.
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (ks *KeyStore) decrypt(data []byte, passphrase string) ([]byte, error) {

	//Create a key based on the provided passphrase, 32 bytes required for AES.
	key := ks.createHash(passphrase)
	if key == EmptyStr {
		return nil, errors.New("KeyStore.encrypt: Hash failed")
	}
	//Create a new AES block cipher based on the key above
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, errors.New("Error creating AES block cipher: " + err.Error())
	}
	//Wrap block cipher in Galois Counter Mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("Error using Galois Counter Mode: " + err.Error())
	}

	//Get nonce from prefix of cipher text.
	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	//Decrypt file and return plain text.
	plaintext, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errors.New("Error decrypting: " + err.Error())
	}

	return plaintext, nil
}

func (ks *KeyStore) SaveKeyData(path string, address Address, data []byte, passphrase string) error {
	filename, _ := filepath.Abs(path + buildFileName(address))
	f, _ := os.Create(filename)

	defer f.Close()

	data, err := ks.encrypt(data, passphrase)
	if err != nil {
		return errors.New("error writing encrypted data to file:" + err.Error())
	}

	account := AccountKey{
		Version:    VERSION,
		CipherText: hex.EncodeToString(data),
		Address:    address.Humanize(),
	}

	accData, err := json.Marshal(account)
	if err != nil {
		return err
	}

	f.Write(accData)

	return nil
}

func (ks *KeyStore) GetKeyData(path string, address Address, passphrase string) ([]byte, error) {
	accData := &AccountKey{}
	filename, err := GetFileName(path, address)
	if err != nil {
		return nil, err
	}
	if !ks.KeyExists(path, address) {
		return nil, errors.New("keystore: file doesn't exist")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, accData)
	if err != nil {
		return nil, err
	}

	cipherData, err := hex.DecodeString(accData.CipherText)
	if err != nil {
		return nil, err
	}

	return ks.decrypt(cipherData, passphrase)
}

func (ks *KeyStore) DeleteKey(path string, address Address, passphrase string) error {
	if res, err := ks.VerifyPassphrase(path, address, passphrase); !res {
		if err != nil {
			return err
		}
		return errors.New("error: invalid password")
	}

	filename, err := GetFileName(path, address)
	if err != nil {
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}

func (ks *KeyStore) VerifyPassphrase(path string, address Address, passphrase string) (bool, error) {
	_, err := ks.GetKeyData(path, address, passphrase)
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewKeyStore() *KeyStore {
	return &KeyStore{}
}

func (ks *KeyStore) KeyExists(path string, address Address) bool {
	filename, err := GetFileName(path, address)
	if err != nil {
		return false
	}
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (ks *KeyStore) GetAddress(path string, filename string) (Address, error) {
	accData := &AccountKey{}
	address := Address{}

	absPath, err := filepath.Abs(path + filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, accData)
	if err != nil {
		return nil, err
	}

	err = address.UnmarshalText([]byte(accData.Address))
	if err != nil {
		return nil, err
	}

	return address, nil
}

func buildFileName(address Address) string {
	t := time.Now()
	zone, _ := t.Zone()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%03d%s",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		zone)

	return "OL_" + formatted + "_" + address.Humanize()
}

func GetFileName(path string, address Address) (string, error) {
	pattern, _ := filepath.Abs(path + "*" + address.Humanize())
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(matches) != 1 {
		return "", errors.New("invalid address")
	}

	return matches[0], nil
}
