/*
	Copyright 2017 - 2018 OneLedger

	Key Management
*/
package id

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ripemd160"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type KeyAlgorithm int

const (
	ED25519 KeyAlgorithm = iota
	SECP256K1
	ED25519_PUB_SIZE   int = ed25519.PubKeyEd25519Size
	SECP256K1_PUB_SIZE int = secp256k1.PubKeySecp256k1Size
)

func (algo KeyAlgorithm) String() string {
	switch algo {
	case ED25519:
		return "ED25519"
	case SECP256K1:
		return "SECP256K1"
	}
	return "Unknown algorithm"
}

func (algo KeyAlgorithm) Size() int {
	switch algo {
	case ED25519:
		return ED25519_PUB_SIZE
	case SECP256K1:
		return SECP256K1_PUB_SIZE
	default:
		log.Error("asked for size of unknown algorithm", "algo", algo)
		return 0
	}
}

func init() {
	serial.Register(AccountKey(""))
}

// Aliases to hide some of the basic underlying types.
type AccountKey []byte // OneLedger address, like Tendermint the hash of the associated PubKey

func (accountKey AccountKey) String() string {
	return hex.EncodeToString(accountKey)
}

func (accountKey AccountKey) Bytes() []byte {
	return accountKey
}

// NewAccountKey hashes the public key to get a unique hash that can act as a key
func NewAccountKey(key PublicKeyED25519) AccountKey {
	return key.Address()
}

type PublicKey interface {
	Address() []byte
	Bytes() []byte
	VerifyBytes(msg []byte, sig []byte) bool
	Equals(PublicKey) bool
	Hex() string
}

type PrivateKey interface {
	Bytes() []byte
	Sign([]byte) ([]byte, error)
	PubKey() PublicKey
	Equals(PrivateKey) bool
}

type PublicKeyED25519 struct {
	Key ed25519.PubKeyEd25519
}

type PrivateKeyED25519 ed25519.PrivKeyEd25519

type PublicKeySECP256K1 struct {
	Key secp256k1.PubKeySecp256k1
}

type PrivateKeySECP256K1 secp256k1.PrivKeySecp256k1

// Ensure these key types implement PublicKey and PrivateKey
var _ PublicKey = new(PublicKeyED25519)
var _ PublicKey = new(PublicKeySECP256K1)

var _ PrivateKey = new(PrivateKeyED25519)
var _ PrivateKey = new(PrivateKeySECP256K1)

func init() {
	serial.Register(NilPublicKey())
	serial.Register(NilPrivateKey())
	serial.Register(ed25519.PubKeyEd25519{})
	serial.Register(secp256k1.PubKeySecp256k1{})
	serial.Register(PublicKeyED25519{})
	serial.Register(PrivateKeyED25519{})
	serial.Register(PublicKeySECP256K1{})
	serial.Register(PrivateKeySECP256K1{})
	var prototypePublicKey PublicKey
	var prototypePrivateKey PrivateKey
	serial.RegisterInterface(&prototypePublicKey)
	serial.RegisterInterface(&prototypePrivateKey)
}

func NilPublicKey() PublicKeyED25519 {
	return PublicKeyED25519{}
}

func NilPrivateKey() PrivateKeyED25519 {
	return PrivateKeyED25519{}
}

// Public Keys
// --------------------------------------------------

func (k PublicKeyED25519) Bytes() []byte {
	return k.Key[:]
}

// Address hashes the key with a RIPEMD-160 hash
func (k PublicKeyED25519) Address() []byte {
	return hash(k)
}

func (k PublicKeyED25519) VerifyBytes(msg []byte, sig []byte) bool {
	return ed25519.PubKeyEd25519(k.Key).VerifyBytes(msg, sig)
}

func (k PublicKeyED25519) Equals(key PublicKey) bool {
	return bytes.Equal(k.Bytes(), key.Bytes())
}

func (k PublicKeyED25519) Hex() string {
	return hex.EncodeToString(k.Bytes())
}

func (k PublicKeySECP256K1) Bytes() []byte {
	return k.Key[:]
}

func (k PublicKeySECP256K1) Address() []byte {
	return hash(k)
}

func (k PublicKeySECP256K1) VerifyBytes(msg []byte, sig []byte) bool {
	return secp256k1.PubKeySecp256k1(k.Key).VerifyBytes(msg, sig)
}

func (k PublicKeySECP256K1) Equals(key PublicKey) bool {
	return bytes.Equal(k.Bytes(), key.Bytes())
}

func OnePublicKey() PublicKeyED25519 {
	return PublicKeyED25519{ed25519.PubKeyEd25519{1}}
}

func OnePrivateKey() PrivateKeyED25519 {
	return PrivateKeyED25519{1}
}

func (k PublicKeySECP256K1) Hex() string {
	return hex.EncodeToString(k.Bytes())
}

// Private keys
//--------------------------------------------------
func (k PrivateKeyED25519) Bytes() []byte {
	return k[:]
}

func (k PrivateKeyED25519) Sign(msg []byte) ([]byte, error) {
	return ed25519.PrivKeyEd25519(k).Sign(msg)
}

func (k PrivateKeyED25519) PubKey() PublicKey {
	p := ed25519.PrivKeyEd25519(k).PubKey().(ed25519.PubKeyEd25519)
	return PublicKeyED25519{p}
}

func (k PrivateKeyED25519) Equals(key PrivateKey) bool {
	return bytes.Equal(k.Bytes(), key.Bytes())
}

func (k PrivateKeySECP256K1) Bytes() []byte {
	return k[:]
}

func (k PrivateKeySECP256K1) Sign(msg []byte) ([]byte, error) {
	return secp256k1.PrivKeySecp256k1(k).Sign(msg)
}

func (k PrivateKeySECP256K1) PubKey() PublicKey {
	p := secp256k1.PrivKeySecp256k1(k).PubKey().(secp256k1.PubKeySecp256k1)
	return PublicKeySECP256K1{p}
}

func (k PrivateKeySECP256K1) Equals(key PrivateKey) bool {
	return bytes.Equal(k.Bytes(), key.Bytes())
}

func GenerateKeys(secret []byte, random bool) (PrivateKeyED25519, PublicKeyED25519) {
	// TODO: Should be configurable
	private, public, err := generateKeys(secret, ED25519, random)
	if err != nil {
		log.Fatal("Key Generation Failed")
	}
	return private.(PrivateKeyED25519), public.(PublicKeyED25519)
}

func generateKeys(secret []byte, algorithm KeyAlgorithm, random bool) (PrivateKey, PublicKey, error) {

	var hash []byte

	if random {
		hash_, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to generate bcrypt hash from secret", "secret", secret)
		}
		hash = hash_
	} else {
		hash = secret
	}
	switch algorithm {

	case ED25519:
		private := ed25519.GenPrivKeyFromSecret(hash)
		public := private.PubKey().(ed25519.PubKeyEd25519)
		return PrivateKeyED25519(private), PublicKeyED25519{public}, nil

	case SECP256K1:
		private := secp256k1.GenPrivKeySecp256k1(hash)
		public := private.PubKey().(secp256k1.PubKeySecp256k1)
		return PrivateKeySECP256K1(private), PublicKeySECP256K1{public}, nil
	}
	return NilPrivateKey(), NilPublicKey(), errors.New("Unknown Algorithm: " + string(algorithm))
}

func hash(k PublicKey) []byte {
	hasher := ripemd160.New()
	hasher.Write(k.Bytes())

	return hasher.Sum(nil)
}

// ImportHexKey returns a PublicKey given a hex-encoded string
func ImportHexKey(h string, k KeyAlgorithm) (PublicKey, error) {
	bz, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key %s", err)
	}
	size := k.Size()
	if size == 0 {
		return NilPublicKey(), errors.New("provided invalid key algorithm")
	}
	return ImportBytesKey(bz, k)
}

// ImportBytesKey takes a byteslice and returns a PublicKey
func ImportBytesKey(bz []byte, k KeyAlgorithm) (PublicKey, error) {
	switch k {
	case ED25519:
		size := ED25519_PUB_SIZE
		if len(bz) != size {
			return new(PublicKeyED25519),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s", k)
		}
		var key [ED25519_PUB_SIZE]byte
		copy(key[:], bz)
		return PublicKeyED25519{key}, nil
	case SECP256K1:
		size := SECP256K1_PUB_SIZE
		if len(bz) != size {
			return new(PublicKeySECP256K1),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s", k)
		}
		var key [SECP256K1_PUB_SIZE]byte
		copy(key[:], bz)
		return PublicKeySECP256K1{key}, nil
	default:
		// Shouldn't reach here
		return nil, errors.New("provided invalid key algorithm")
	}
}

//	salt := cytpo.CRandBytes(16)
//func Armour(privateKey PrivateKey, passphrase string, salt []byte) ([]byte, error) {
//	key, status := []byte(passphrase), error(nil)
//	//key, status := bcrypt.GenerateFromPassword(salt, []byte(passphrse), 16)
//	if status != nil {
//		return nil, errors.New("Failed Bcrypt")
//	}
//	base := crypto.Sha256(key) // Is this necessary?
//
//	bytes := privateKey.Bytes()
//	return xsalsa20symmetric.EncryptSymmetric(bytes, base), nil
//}

//func Dearmour(buffer []byte, passphrase string, salt []byte) (PrivateKey, error) {
//	key, status := []byte(passphrase), error(nil)
//	//key, status := bcrypt.GenerateFromPassword(salt, []byte(passphrse), 16)
//	if status != nil {
//		return PrivateKey{}, errors.New("Failed Bcrypt")
//	}
//	base := crypto.Sha256(key) // Is this necessary?
//	result, status := crypto.DecryptSymmetric(buffer, base)
//	if status != nil {
//		return PrivateKey{}, errors.New("Failed Symmetric Decrypt")
//	}
//	crypto.PrivKeyFromBytes(result)
//}
