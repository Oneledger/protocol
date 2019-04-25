package key

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type PublicKeyHandler interface {
	Address() []byte
	Bytes() []byte
	VerifyBytes(msg []byte, sig []byte) bool
	Equals(PublicKey) bool
}

type PrivateKeyHandler interface {
	Bytes() []byte
	Sign([]byte) ([]byte, error)
	PubKey() PublicKey
	Equals(PrivateKey) bool
}

type PrivateKey struct {
	keytype Algorithm
	data []byte
}

type PublicKey struct {
	Type Algorithm
	Data []byte
}

type Address struct {
	Type AddressType
}

func NewPrivateKeyFromBytes(k []byte, algorithm Algorithm) PrivateKey{
	return PrivateKey{algorithm, k}
}

// Get the public key handler
func (pubkey PublicKey) GetHandler() (PublicKeyHandler, error) {
	switch pubkey.Type {
	case ED25519:
		size := ed25519.PubKeyEd25519Size
		if len(pubkey.Data) != size {
			return new(PublicKeyED25519),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s", pubkey.Type)
		}
		var key [ED25519_PUB_SIZE]byte
		copy(key[:], pubkey.Data)
		return PublicKeyED25519{key}, nil
	case SECP256K1:
		size := SECP256K1_PUB_SIZE
		if len(pubkey.Data) != size {
			return new(PublicKeySECP256K1),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s", pubkey.Type)
		}
		var key [SECP256K1_PUB_SIZE]byte
		copy(key[:], pubkey.Data)
		return PublicKeySECP256K1{key}, nil
	default:
		// Shouldn't reach here
		return nil, errors.New("provided invalid key algorithm")
	}
}

// get the private key handler
func (privkey PrivateKey) GetHandler() (PrivateKeyHandler, error) {
	switch privkey.keytype {
	case ED25519:

		if len(privkey.data) != 64 {
			return new(PrivateKeyED25519),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s", privkey.keytype)
		}
		var key [64]byte
		copy(key[:], privkey.data)
		return PrivateKeyED25519(key), nil
	case SECP256K1:
		size := SECP256K1_PUB_SIZE
		if len(privkey.data) != size {
			return new(PrivateKeySECP256K1),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s", privkey.keytype)
		}
		var key [32]byte
		copy(key[:], privkey.data)
		return PrivateKeySECP256K1(key), nil
	default:
		// Shouldn't reach here
		return nil, errors.New("provided invalid key algorithm")
	}
}



//====================== ED25519 ======================

var _  PublicKeyHandler = PublicKeyED25519{}
var _  PrivateKeyHandler = PrivateKeyED25519{}

type PublicKeyED25519 struct {
	key ed25519.PubKeyEd25519
}

func (k PublicKeyED25519) Bytes() []byte {
	return k.key[:]
}

// Address hashes the key with a RIPEMD-160 hash
func (k PublicKeyED25519) Address() []byte {
	return k.key.Address()
}

func (k PublicKeyED25519) VerifyBytes(msg []byte, sig []byte) bool {
	return k.key.VerifyBytes(msg, sig)
}

func (k PublicKeyED25519) Equals(pubkey PublicKey) bool {
	return pubkey.Type == ED25519 && bytes.Equal(k.Bytes(), pubkey.Data)
}

type PrivateKeyED25519 ed25519.PrivKeyEd25519

func (k PrivateKeyED25519) Bytes() []byte {
	return k[:]
}

func (k PrivateKeyED25519) Sign(msg []byte) ([]byte, error) {
	return ed25519.PrivKeyEd25519(k).Sign(msg)
}

func (k PrivateKeyED25519) PubKey() PublicKey {
	p := ed25519.PrivKeyEd25519(k).PubKey()
	return PublicKey{
		Type:   ED25519,
		Data:   p.Bytes(),
	}
}

func (k PrivateKeyED25519) Equals(privkey PrivateKey) bool {
	return privkey.keytype == ED25519 && bytes.Equal(k.Bytes(), privkey.data)
}

//====================== ED25519 ======================



//====================== SECP256K1 ======================
var _  PublicKeyHandler = PublicKeySECP256K1{}
var _  PrivateKeyHandler = PrivateKeySECP256K1{}

type PublicKeySECP256K1 struct {
	key secp256k1.PubKeySecp256k1
}

func (k PublicKeySECP256K1) Bytes() []byte {
	return k.key[:]
}

// Address hashes the key with a RIPEMD-160 hash
func (k PublicKeySECP256K1) Address() []byte {
	return k.key.Address()
}

func (k PublicKeySECP256K1) VerifyBytes(msg []byte, sig []byte) bool {
	return k.key.VerifyBytes(msg, sig)
}

func (k PublicKeySECP256K1) Equals(PubkeySECP256K1 PublicKey) bool {
	return PubkeySECP256K1.Type == SECP256K1 && bytes.Equal(k.Bytes(), PubkeySECP256K1.Data)
}

type PrivateKeySECP256K1 secp256k1.PrivKeySecp256k1

func (k PrivateKeySECP256K1) Bytes() []byte {
	return k[:]
}

func (k PrivateKeySECP256K1) Sign(msg []byte) ([]byte, error) {
	return secp256k1.PrivKeySecp256k1(k).Sign(msg)
}

func (k PrivateKeySECP256K1) PubKey() PublicKey {
	p := secp256k1.PrivKeySecp256k1(k).PubKey()
	return PublicKey{
		Type:   SECP256K1,
		Data:   p.Bytes(),
	}
}

func (k PrivateKeySECP256K1) Equals(privkey PrivateKey) bool {
	return privkey.keytype == SECP256K1 && bytes.Equal(k.Bytes(), privkey.data)
}


//====================== SECP256K1 ======================