/*
	Copyright 2017 - 2018 OneLedger

	Key Management
*/
package id

import (
	"errors"

	//"github.com/tendermint/go-crypto/keys/bcrypt"
	"github.com/tendermint/go-crypto"
)

type KeyAlgorithm int

const (
	ED25519 KeyAlgorithm = iota
	SECP256K1
)

func GenerateKeys() (PublicKey, PrivateKey) {
	return PublicKey{}, PrivateKey{}
}

func Generate(secret []byte, algorithm KeyAlgorithm) (PrivateKey, error) {
	switch algorithm {
	case ED25519:
		return crypto.GenPrivKeyEd25519FromSecret(secret).Wrap(), nil
	case SECP256K1:
		return crypto.GenPrivKeySecp256k1FromSecret(secret).Wrap(), nil
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

	bytes := privateKey.Bytes()
	return crypto.EncryptSymmetric(bytes, base), nil
}

func Dearmour(buffer []byte, passphrase string, salt []byte) (PrivateKey, error) {
	key, err := []byte(passphrase), error(nil)
	//key, err := bcrypt.GenerateFromPassword(salt, []byte(passphrse), 16)
	if err != nil {
		return PrivateKey{}, errors.New("Failed Bcrypt")
	}
	base := crypto.Sha256(key) // Is this necessary?
	result, err := crypto.DecryptSymmetric(buffer, base)
	if err != nil {
		return PrivateKey{}, errors.New("Failed Symmetric Decrypt")
	}
	return crypto.PrivKeyFromBytes(result)
}
