/*
	Copyright 2017 - 2018 OneLedger

	Key Management
*/
package id

import (
	"errors"
	//"github.com/tendermint/go-crypto/keys/bcrypt"
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/xsalsa20symmetric"
	"github.com/tendermint/tendermint/crypto"
)

type KeyAlgorithm int

const (
	ED25519 KeyAlgorithm = iota
	SECP256K1
)

func GenerateKeys(secret []byte) (PublicKey, PrivateKey) {
	private, err := Generate(secret, ED25519) // TODO: Should be configurable
	if err != nil {
		log.Fatal("Key Generation Failed")
	}
	public := PublicKey{private.Interface.PubKey()}

	return public, private
}

func Generate(secret []byte, algorithm KeyAlgorithm) (PrivateKey, error) {
	// go-crypto doesn't work
	switch algorithm {
	case ED25519:
		return PrivateKey{ed25519.GenPrivKeyFromSecret(secret)}, nil
	case SECP256K1:
		// NOTE: secret should be the output of a KDF like bcrypt,
		// if it's derived from user input.
		return PrivateKey{secp256k1.GenPrivKeySecp256k1(secret)}, nil
	}
	return PrivateKey{}, errors.New("Unknown Algorithm: " + string(algorithm))
}

//	salt := cytpo.CRandBytes(16)
func Armour(privateKey PrivateKey, passphrase string, salt []byte) ([]byte, error) {
	key, err := []byte(passphrase), error(nil)
	//key, err := bcrypt.GenerateFromPassword(salt, []byte(passphrse), 16)
	if err != nil {
		return nil, errors.New("Failed Bcrypt")
	}
	base := crypto.Sha256(key) // Is this necessary?

	bytes := privateKey.Interface.Bytes()
	return xsalsa20symmetric.EncryptSymmetric(bytes, base), nil
}

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
