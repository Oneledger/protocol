/*
	Copyright 2017 - 2018 OneLedger

	Key Management
*/
package id

import (
	"errors"
	"github.com/tendermint/tendermint/crypto"
	"golang.org/x/crypto/bcrypt"

	//"github.com/tendermint/go-crypto/keys/bcrypt"
	"github.com/Oneledger/protocol/node/log"
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

func (accountKey AccountKey) String() string {
	return string(accountKey)
}

func (accountKey AccountKey) Bytes() []byte {
	return accountKey
}

type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey

type ED25519PublicKey = ed25519.PubKeyEd25519
type ED25519PrivateKey = ed25519.PrivKeyEd25519

type SECP256K1PublicKey = secp256k1.PubKeySecp256k1
type SECP256K1PrivateKey = secp256k1.PrivKeySecp256k1

func NilPublicKey() ED25519PublicKey {
	return ED25519PublicKey{}
}

func NilPrivateKey() ED25519PrivateKey {
	return ED25519PrivateKey{}
}

func GenerateKeys(secret []byte) (ED25519PrivateKey, ED25519PublicKey) {
	// TODO: Should be configurable
	private, public, err := generateKeys(secret, ED25519)
	if err != nil {
		log.Fatal("Key Generation Failed")
	}
	return private.(ED25519PrivateKey), public.(ED25519PublicKey)
}

func generateKeys(secret []byte, algorithm KeyAlgorithm) (PrivateKey, PublicKey, error) {
	hash, err := bcrypt.GenerateFromPassword(secret, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to generate bcrypt hash from secret", "secret", secret)
	}
	switch algorithm {
	case ED25519:
		private := ed25519.GenPrivKeyFromSecret(hash)

		log.Info("Generate", "private", private)

		public := private.PubKey()
		return ED25519PrivateKey(private), public, nil
	case SECP256K1:
		// NOTE: secret should be the output of a KDF like bcrypt,
		// if it's derived from user input.
		private := secp256k1.GenPrivKeySecp256k1(hash)
		public := private.PubKey()
		return SECP256K1PrivateKey(private), public, nil
	}
	return NilPrivateKey(), NilPublicKey(), errors.New("Unknown Algorithm: " + string(algorithm))
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
