package keys

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/privval"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/p2p"
)

/*
	Interfaces
*/
type PublicKeyHandler interface {
	Address() Address
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

/*
	Key Types
*/
type PublicKey struct {
	KeyType Algorithm `json:"keyType"`
	Data    []byte    `json:"data"`
}

type PrivateKey struct {
	Keytype Algorithm `json:"keyType"`
	Data    []byte    `json:"data"`
}

type gobkey struct {
	K int
	D []byte
}

// tmPrefixSize is the size of the prefix returned by tendermint/amino-encoded byte arrays
const tmPrefixSize = 5

func (pubKey *PublicKey) GobEncode() ([]byte, error) {
	//a := map[string]interface{}{
	//	"K": int(pubKey.Keytype),
	//	"D":    pubKey.Data,
	//}
	a := gobkey{int(pubKey.KeyType), pubKey.Data}
	return json.Marshal(&a)
}

func (pubKey *PublicKey) GobDecode(buf []byte) error {
	a := gobkey{}
	err := json.Unmarshal(buf, &a)
	if err != nil {
		return err
	}

	//pubKey.Keytype = Algorithm(a["K"].(float64))
	pubKey.KeyType = Algorithm(a.K)
	pubKey.Data = a.D

	return nil
}

func (privKey *PrivateKey) GobEncode() ([]byte, error) {
	a := gobkey{int(privKey.Keytype), privKey.Data}
	return json.Marshal(&a)
}

func (privKey *PrivateKey) GobDecode(buf []byte) error {
	a := gobkey{}
	err := json.Unmarshal(buf, &a)
	if err != nil {
		return err
	}

	//pubKey.Keytype = Algorithm(a["K"].(float64))
	privKey.Keytype = Algorithm(a.K)
	privKey.Data = a.D

	return nil
}

func (pubKey PublicKey) Equal(pkey PublicKey) bool {
	if pubKey.KeyType != pkey.KeyType {
		return false
	}
	if !bytes.Equal(pubKey.Data, pkey.Data) {
		return false
	}
	return true
}

func GetPrivateKeyFromBytes(k []byte, algorithm Algorithm) (PrivateKey, error) {
	key := PrivateKey{algorithm, k}

	if _, err := key.GetHandler(); err != nil {
		return PrivateKey{}, err
	}

	return key, nil
}

func GetPublicKeyFromBytes(k []byte, algorithm Algorithm) (PublicKey, error) {
	key := PublicKey{algorithm, k}

	if _, err := key.GetHandler(); err != nil {
		return PublicKey{}, err
	}

	return key, nil
}

func NewKeyPairFromTendermint() (PublicKey, PrivateKey, error) {
	tmPrivKey := ed25519.GenPrivKey()
	tmPublicKey := tmPrivKey.PubKey()

	pubKey, err := GetPublicKeyFromBytes(tmPublicKey.Bytes()[tmPrefixSize:], ED25519)
	if err != nil {
		return PublicKey{}, PrivateKey{}, errors.Wrap(err, "error creating public key")
	}

	privKey, err := GetPrivateKeyFromBytes(tmPrivKey.Bytes()[tmPrefixSize:], ED25519)
	if err != nil {
		return PublicKey{}, PrivateKey{}, errors.Wrap(err, "error in cr3eating private key")
	}

	return pubKey, privKey, nil
}

// NodeKeyFromTendermint returns a PrivateKey from a tendermint NodeKey.
// The input key must be a ED25519 key.
func NodeKeyFromTendermint(key *p2p.NodeKey) (PrivateKey, error) {
	if key == nil {
		return PrivateKey{}, errors.New("NodeKeyFromTendermint: got nil argument")
	}
	bz := key.PrivKey.Bytes()[tmPrefixSize:]
	return GetPrivateKeyFromBytes(bz, ED25519)
}

// NodeKeyFromTendermint returns a PrivateKey from a tendermint NodeKey.
// The input key must be a ED25519 key.
func PVKeyFromTendermint(key *privval.FilePVKey) (PrivateKey, error) {
	if key == nil {
		return PrivateKey{}, errors.New("PVKeyFromTendermint: got nil argument")
	}
	bz := key.PrivKey.Bytes()[tmPrefixSize:]
	return GetPrivateKeyFromBytes(bz, ED25519)
}

func PubKeyFromTendermint(key []byte) (PublicKey, error) {
	if key == nil {
		return PublicKey{}, errors.New("PubKeyFromTendermint: got nil argument")
	}
	bz := key[tmPrefixSize:]
	return GetPublicKeyFromBytes(bz, ED25519)
}

func (pubKey PublicKey) GetABCIPubKey() types.PubKey {
	return types.PubKey{
		Type: pubKey.KeyType.Name(),
		Data: pubKey.Data,
	}
}

// Get the public key handler
func (pubKey PublicKey) GetHandler() (PublicKeyHandler, error) {
	switch pubKey.KeyType {
	case ED25519:
		size := ed25519.PubKeyEd25519Size
		if len(pubKey.Data) != size {
			return new(PublicKeyED25519),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s length %d", pubKey.KeyType.Name(), len(pubKey.Data))
		}
		var key [ED25519_PUB_SIZE]byte
		copy(key[:], pubKey.Data)
		return PublicKeyED25519{key}, nil
	case SECP256K1:
		size := SECP256K1_PUB_SIZE
		if len(pubKey.Data) != size {
			return new(PublicKeySECP256K1),
				fmt.Errorf("given key doesn't match the size of the key algorithm %s length %d", pubKey.KeyType.Name(), len(pubKey.Data))
		}
		var key [SECP256K1_PUB_SIZE]byte
		copy(key[:], pubKey.Data)
		return PublicKeySECP256K1{key}, nil
	default:
		// Shouldn't reach here
		return nil, errors.New("provided invalid key algorithm")
	}
}

// get the private key handler
func (privKey PrivateKey) GetHandler() (PrivateKeyHandler, error) {
	switch privKey.Keytype {
	case ED25519:

		if len(privKey.Data) != ED25519_PRIV_SIZE {
			return new(PrivateKeyED25519),
				fmt.Errorf("given key doesn't match the size of the key algorithm %d", privKey.Keytype)
		}
		var key [64]byte
		copy(key[:], privKey.Data)
		return PrivateKeyED25519(key), nil
	case SECP256K1:
		size := SECP256K1_PRIV_SIZE
		if len(privKey.Data) != size {
			return new(PrivateKeySECP256K1),
				fmt.Errorf("given key doesn't match the size of the key algorithm %d", privKey.Keytype)
		}
		var key [32]byte
		copy(key[:], privKey.Data)
		return PrivateKeySECP256K1(key), nil
	default:
		// Shouldn't reach here
		return nil, errors.New("provided invalid key algorithm")
	}
}

//====================== ED25519 ======================

var _ PublicKeyHandler = PublicKeyED25519{}
var _ PrivateKeyHandler = PrivateKeyED25519{}

type PublicKeyED25519 struct {
	key ed25519.PubKeyEd25519
}

func (k PublicKeyED25519) Bytes() []byte {
	return k.key[:]
}

// Address hashes the key with a RIPEMD-160 hash
func (k PublicKeyED25519) Address() Address {
	return k.key.Address().Bytes()
}

func (k PublicKeyED25519) VerifyBytes(msg []byte, sig []byte) bool {
	return k.key.VerifyBytes(msg, sig)
}

func (k PublicKeyED25519) Equals(pubkey PublicKey) bool {
	return pubkey.KeyType == ED25519 && bytes.Equal(k.Bytes(), pubkey.Data)
}

func (k PublicKeyED25519) String() string {
	return k.key.String()
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
		KeyType: ED25519,
		Data:    p.Bytes()[tmPrefixSize:],
	}
}

func (k PrivateKeyED25519) Equals(privkey PrivateKey) bool {
	return privkey.Keytype == ED25519 && bytes.Equal(k.Bytes(), privkey.Data)
}

//====================== ED25519 ======================

//====================== SECP256K1 ======================
var _ PublicKeyHandler = PublicKeySECP256K1{}
var _ PrivateKeyHandler = PrivateKeySECP256K1{}

type PublicKeySECP256K1 struct {
	key secp256k1.PubKeySecp256k1
}

func (k PublicKeySECP256K1) Bytes() []byte {
	return k.key[:]
}

// Address hashes the key with a RIPEMD-160 hash
func (k PublicKeySECP256K1) Address() Address {
	return k.key.Address().Bytes()
}

func (k PublicKeySECP256K1) VerifyBytes(msg []byte, sig []byte) bool {
	return k.key.VerifyBytes(msg, sig)
}

func (k PublicKeySECP256K1) Equals(PubkeySECP256K1 PublicKey) bool {
	return PubkeySECP256K1.KeyType == SECP256K1 && bytes.Equal(k.Bytes(), PubkeySECP256K1.Data)
}

func (k PublicKeySECP256K1) String() string {
	return k.key.String()
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
		KeyType: SECP256K1,
		Data:    p.Bytes(),
	}
}

func (k PrivateKeySECP256K1) Equals(privkey PrivateKey) bool {
	return privkey.Keytype == SECP256K1 && bytes.Equal(k.Bytes(), privkey.Data)
}

//====================== SECP256K1 ======================
