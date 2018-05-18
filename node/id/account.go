/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains

	TODO: Need to pick a system key for identities. Is a hash of pubkey reasonable?
*/
package id

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	crypto "github.com/tendermint/go-crypto"
	wdata "github.com/tendermint/go-wire/data"
	"golang.org/x/crypto/ripemd160"
)

// Aliases to hide some of the basic underlying types.

type Address = wdata.Bytes // OneLedger address, like Tendermint the hash of the associated PubKey

type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey
type Signature = crypto.Signature

// enum for type
type AccountType int

const (
	ONELEDGER AccountType = iota
	BITCOIN
	ETHEREUM
)

// The persistent collection of all accounts known by this node
type Accounts struct {
	data *data.Datastore
}

func NewAccounts(name string) *Accounts {
	data := data.NewDatastore(name, data.PERSISTENT)

	return &Accounts{
		data: data,
	}
}

type AccountKey []byte

// Polymorphism
type Account interface {
	AddPrivateKey(PrivateKey)
	Name() string
	Key() []byte
}

type AccountBase struct {
	Type AccountType

	Key        AccountKey
	Name       string
	PublicKey  PublicKey
	PrivateKey PrivateKey
}

// Hash the public key to get a unqiue hash that can act as a key
func NewAccountKey(key PublicKey) AccountKey {
	hasher := ripemd160.New()

	bytes, err := key.MarshalJSON()
	if err != nil {
		panic("Unable to Marshal the key into bytes")
	}

	hasher.Write(bytes)

	return hasher.Sum(nil)
}

func NewAccount(newType AccountType, name string, Key PublicKey) Account {
	switch newType {

	case ONELEDGER:
		return &AccountOneLedger{}

	case BITCOIN:
		return &AccountBitcoin{}

	case ETHEREUM:
		return &AccountEthereum{}

	default:
		panic("Unknown Type")
	}
}

// TODO: really should be part of the enum, as a map...
func FindAccountType(typeName string) (AccountType, err.Code) {
	switch typeName {
	case "OneLedger":
		return ONELEDGER, err.SUCCESS

	case "Ethereum":
		return ETHEREUM, err.SUCCESS

	case "Bitcoin":
		return BITCOIN, err.SUCCESS
	}
	return 0, 42
}

func FindAccount(name string) (Account, err.Code) {
	// TODO: Lookup the identity in the node's database
	return &AccountOneLedger{AccountBase: AccountBase{Name: name}}, 0
}

// OneLedger

// Information we need about our own fullnode identities
type AccountOneLedger struct {
	AccountBase
}

func (account *AccountOneLedger) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountOneLedger) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}

func (account *AccountOneLedger) Name() string {
	return account.AccountBase.Name
}

func (account *AccountOneLedger) Key() []byte {
	return []byte(account.AccountBase.Name)
}

// Bitcoin

// Information we need for a Bitcoin account
type AccountBitcoin struct {
	AccountBase
}

func (account *AccountBitcoin) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountBitcoin) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}

func (account *AccountBitcoin) Name() string {
	return account.AccountBase.Name
}

func (account *AccountBitcoin) Key() []byte {
	return []byte(account.AccountBase.Name)
}

// Ethereum

// Information we need for an Ethereum account
type AccountEthereum struct {
	AccountBase
}

func (account *AccountEthereum) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountEthereum) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}

func (account *AccountEthereum) Name() string {
	return account.AccountBase.Name
}

func (account *AccountEthereum) Key() []byte {
	return []byte(account.AccountBase.Name)
}
