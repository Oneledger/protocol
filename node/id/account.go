/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains
*/
package id

import (
	"reflect"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
	crypto "github.com/tendermint/go-crypto"
	wdata "github.com/tendermint/go-wire/data"
	"golang.org/x/crypto/ripemd160"
)

// Aliases to hide some of the basic underlying types.
type Address = wdata.Bytes // OneLedger address, like Tendermint the hash of the associated PubKey

type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey
type Signature = crypto.Signature

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

func (acc *Accounts) Add(account Account) {

	if value := acc.data.Load(account.Key()); value != nil {
		log.Debug("Key is being updated")
	}

	buffer, _ := comm.Serialize(account)

	acc.data.Store(account.Key(), buffer)
	acc.data.Commit()
}

func (acc *Accounts) Delete(account Account) {
}

func (acc *Accounts) Exists(newType data.ChainType, name string) bool {
	account := NewAccount(newType, name, PublicKey{})
	value := acc.data.Load(account.Key())
	if value != nil {
		return true
	}
	return false
}

func (acc *Accounts) FindAll() []Account {
	keys := acc.data.List()
	size := len(keys)
	results := make([]Account, size, size)

	for i := 0; i < size; i++ {
		// TODO: This is dangerous...
		account := &AccountOneLedger{}
		base, _ := comm.Deserialize(acc.data.Load(keys[i]), account)
		results[i] = base.(Account)
	}

	return results
}

func (acc *Accounts) Dump() {
	list := acc.FindAll()
	size := len(list)

	for i := 0; i < size; i++ {
		account := list[i]
		log.Info("Account", "Name", account.Name())
		log.Info("Type", "Type", reflect.TypeOf(account))
	}
}

func (acc *Accounts) Find(name string) (Account, err.Code) {
	// TODO: Lookup the identity in the node's database
	return &AccountOneLedger{AccountBase: AccountBase{Name: name}}, 0
}

type AccountKey []byte

// Polymorphism
type Account interface {
	Key() data.DatabaseKey
	Name() string

	AddPublicKey(PublicKey)
	AddPrivateKey(PrivateKey)

	AsString() string
}

type AccountBase struct {
	Type data.ChainType

	Key AccountKey

	Name       string
	PublicKey  PublicKey
	PrivateKey PrivateKey
}

// Hash the public key to get a unqiue hash that can act as a key
func NewAccountKey(key PublicKey) AccountKey {
	hasher := ripemd160.New()

	// TODO: This deosn't seem right?
	bytes, err := key.MarshalJSON()
	if err != nil {
		panic("Unable to Marshal the key into bytes")
	}

	hasher.Write(bytes)

	return hasher.Sum(nil)
}

// Create a new account for a given chain
func NewAccount(newType data.ChainType, name string, key PublicKey) Account {
	switch newType {

	case data.ONELEDGER:
		return &AccountOneLedger{
			AccountBase{
				Type:      newType,
				Key:       NewAccountKey(key),
				Name:      name,
				PublicKey: key,
			},
		}

	case data.BITCOIN:
		return &AccountBitcoin{
			AccountBase{
				Type:      newType,
				Key:       NewAccountKey(key),
				Name:      name,
				PublicKey: key,
			},
		}

	case data.ETHEREUM:
		return &AccountEthereum{
			AccountBase{
				Type:      newType,
				Key:       NewAccountKey(key),
				Name:      name,
				PublicKey: key,
			},
		}

	default:
		panic("Unknown Type")
	}
}

// Map type to string
func ParseAccountType(typeName string) data.ChainType {
	switch typeName {
	case "OneLedger":
		return data.ONELEDGER

	case "Ethereum":
		return data.ETHEREUM

	case "Bitcoin":
		return data.BITCOIN
	}
	return data.UNKNOWN
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

func (account *AccountOneLedger) Key() data.DatabaseKey {
	return data.DatabaseKey(account.AccountBase.Name)
}

func (account *AccountOneLedger) AsString() string {

	// TODO: Add in UTXO entry
	return "- " + account.AccountBase.Name
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

func (account *AccountBitcoin) Key() data.DatabaseKey {
	return data.DatabaseKey(account.AccountBase.Name)
}

func (account *AccountBitcoin) AsString() string {
	return "- " + account.AccountBase.Name
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

func (account *AccountEthereum) Key() data.DatabaseKey {
	return data.DatabaseKey(account.AccountBase.Name)
}

func (account *AccountEthereum) AsString() string {
	return "- " + account.AccountBase.Name
}
