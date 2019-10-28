// ecdsa.go contains some utility functions for reading and saving ECDSA keys
package eth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

const DefaultKeyDumpFilename = "ecdsa.json"

// privateKeyDump holds the raw binary blob of an ECDSA key
type privateKeyDump struct {
	// Binary contains the raw ECDSA key in a binary format. Pass this to
	// crypto.
	Binary []byte `json:"binary"`
}

// Save marshals the privateKeyDump into JSON and then saves the file to the given location.
// If the location includes a filen
func (dump privateKeyDump) Save(location string) error {
	raw, err := json.Marshal(dump)
	if err != nil {
		return err
	}

	// Give it the default filename if the path doesn't include the filename
	hasFilename := filepath.Ext(location) != ""
	if !hasFilename {
		location = filepath.Join(location, DefaultKeyDumpFilename)
	}

	return ioutil.WriteFile(location, raw, 0644)
}

func (dump privateKeyDump) ECDSA() (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(dump.Binary)
}

// ReadUnsafeKeyDump reads a privateKeyDump from a given JSON file, and unmarshals
// the contents into an ecdsa.PrivateKey
func ReadUnsafeKeyDump(location string) (*ecdsa.PrivateKey, error) {
	blob, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}

	var dump privateKeyDump
	err = json.Unmarshal(blob, &dump)
	if err != nil {
		return nil, err
	}

	return dump.ECDSA()
}

// SaveUnsafeKeyDump saves a given ecdsa.PrivateKey into a JSON file at the given location
func SaveUnsafeKeyDump(key *ecdsa.PrivateKey, location string) error {
	dump := unsafeKeyDump(key)
	return dump.Save(location)
}

// unsafeKeyDump takes a raw ecdsa key and returns a privateKeyDump
func unsafeKeyDump(key *ecdsa.PrivateKey) privateKeyDump {
	return privateKeyDump{crypto.FromECDSA(key)}
}

// GenerateKey returns an secp256k1 key for use with Ethereum
func GenerateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
}
