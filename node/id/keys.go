/*
	Copyright 2017 - 2018 OneLedger

	Key Management
*/
package id

import (
	"bytes"
	"errors"

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
)

// Aliases to hide some of the basic underlying types.
type AccountKey []byte // OneLedger address, like Tendermint the hash of the associated PubKey

func init() {
	serial.Register(AccountKey(""))
}

func (accountKey AccountKey) String() string {
	return string(accountKey)
}

// TODO: Just use String for all presentation variations....
func (accountKey AccountKey) AsString() string {
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
}

type PrivateKey interface {
	Bytes() []byte
	Sign([]byte) ([]byte, error)
	PubKey() PublicKey
	Equals(PrivateKey) bool
}

type PublicKeyED25519 ed25519.PubKeyEd25519
type PrivateKeyED25519 ed25519.PrivKeyEd25519

type PublicKeySECP256K1 secp256k1.PubKeySecp256k1
type PrivateKeySECP256K1 secp256k1.PrivKeySecp256k1

func init() {
	serial.Register(NilPublicKey())
	serial.Register(NilPrivateKey())
	serial.Register(PublicKeySECP256K1{})
	serial.Register(PrivateKeySECP256K1{})
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
	return k[:]
}

// Address hashes the key with a RIPEMD-160 hash
func (k PublicKeyED25519) Address() []byte {
	return hash(k)
}

func (k PublicKeyED25519) VerifyBytes(msg []byte, sig []byte) bool {
	return ed25519.PubKeyEd25519(k).VerifyBytes(msg, sig)
}

func (k PublicKeyED25519) Equals(key PublicKey) bool {
	return bytes.Equal(k.Bytes(), key.Bytes())
}

func (k PublicKeySECP256K1) Bytes() []byte {
	return k[:]
}

func (k PublicKeySECP256K1) Address() []byte {
	return hash(k)
}

func (k PublicKeySECP256K1) VerifyBytes(msg []byte, sig []byte) bool {
	return secp256k1.PubKeySecp256k1(k).VerifyBytes(msg, sig)
}

func (k PublicKeySECP256K1) Equals(key PublicKey) bool {
	return bytes.Equal(k.Bytes(), key.Bytes())
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
	return PublicKeyED25519(p)
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
	p := secp256k1.PrivKeySecp256k1(k).PubKey().(ed25519.PubKeyEd25519)
	return PublicKeyED25519(p)
}

func (k PrivateKeySECP256K1) Equals(key PrivateKey) bool {
	return bytes.Equal(k.Bytes(), key.Bytes())
}

func GenerateKeys(secret []byte) (PrivateKeyED25519, PublicKeyED25519) {
	// TODO: Should be configurable
	private, public, err := generateKeys(secret, ED25519)
	if err != nil {
		log.Fatal("Key Generation Failed")
	}
	return private.(PrivateKeyED25519), public.(PublicKeyED25519)
}

func generateKeys(secret []byte, algorithm KeyAlgorithm) (PrivateKey, PublicKey, error) {
	hash, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to generate bcrypt hash from secret", "secret", secret)
	}
	switch algorithm {

	case ED25519:
		private := ed25519.GenPrivKeyFromSecret(hash)
		public := private.PubKey().(ed25519.PubKeyEd25519)
		return PrivateKeyED25519(private), PublicKeyED25519(public), nil

	case SECP256K1:
		// NOTE: secret should be the output of a KDF like bcrypt,
		// if it's derived from user input.
		private := secp256k1.GenPrivKeySecp256k1(hash)
		public := private.PubKey().(secp256k1.PubKeySecp256k1)
		return PrivateKeySECP256K1(private), PublicKeySECP256K1(public), nil
	}
	return NilPrivateKey(), NilPublicKey(), errors.New("Unknown Algorithm: " + string(algorithm))
}

func hash(k PublicKey) []byte {
	hasher := ripemd160.New()
	hasher.Write(k.Bytes())

	return hasher.Sum(nil)
}

//	salt := cytpo.CRandBytes(16)
//func Armour(privateKey PrivateKey, passphrase string, salt []byte) ([]byte, error) {
//	key, err := []byte(passphrase), error(nil)
//	//key, err := bcrypt.GenerateFromPassword(salt, []byte(passphrse), 16)
//	if err != nil {
//		return nil, errors.New("Failed Bcrypt")
//	}
//	base := crypto.Sha256(key) // Is this necessary?
//
//	bytes := privateKey.Bytes()
//	return xsalsa20symmetric.EncryptSymmetric(bytes, base), nil
//}

//func Dearmour(buffer []byte, passphrase string, salt []byte) (PrivateKey, error) {
//	key, err := []byte(passphrase), error(nil)
//	//key, err := bcrypt.GenerateFromPassword(salt, []byte(passphrse), 16)
//	if err != nil {
//		return PrivateKey{}, errors.New("Failed Bcrypt")
//	}
//	base := crypto.Sha256(key) // Is this necessary?
//	result, err := crypto.DecryptSymmetric(buffer, base)
//	if err != nil {
//		return PrivateKey{}, errors.New("Failed Symmetric Decrypt")
//	}
//	crypto.PrivKeyFromBytes(result)
//}
