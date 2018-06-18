/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains
*/
package id

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
	crypto "github.com/tendermint/go-crypto"
	wdata "github.com/tendermint/go-wire/data"
	"golang.org/x/crypto/ripemd160"
)

// Aliases to hide some of the basic underlying types.
type AccountKey = wdata.Bytes // OneLedger address, like Tendermint the hash of the associated PubKey

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

	if value := acc.data.Load(account.AccountKey()); value != nil {
		log.Debug("Key is being updated", "key", account.AccountKey())
	}

	buffer, _ := comm.Serialize(account)

	acc.data.Store(account.AccountKey(), buffer)
	acc.data.Commit()
}

func (acc *Accounts) Delete(account Account) {
}

func (acc *Accounts) Exists(newType data.ChainType, name string) bool {
	account := NewAccount(newType, name, PublicKey{}, PrivateKey{})

	if value := acc.data.Load(account.AccountKey()); value != nil {
		return true
	}

	return false
}

func (acc *Accounts) Find(account Account) (Account, err.Code) {
	return acc.FindKey(account.AccountKey())
}

func (acc *Accounts) FindIdentity(identity Identity) (Account, err.Code) {
	// TODO: Should have better name mapping between identities and accounts
	account := NewAccount(data.ONELEDGER, identity.Name+"-OneLedger", PublicKey{}, PrivateKey{})
	return acc.Find(account)

}

func (acc *Accounts) FindNameOnChain(name string, chain data.ChainType) (Account, err.Code) {

	// TODO: Should be replaced with a real index
	for _, entry := range acc.FindAll() {
		if Matches(entry, name, chain) {
			return entry, err.SUCCESS
		}
	}
	return nil, err.SUCCESS
}

func (acc *Accounts) FindName(name string) (Account, err.Code) {
	return acc.FindNameOnChain(name, data.ONELEDGER)
}

func Matches(account Account, name string, chain data.ChainType) bool {
	if strings.EqualFold(account.Name(), name) {
		if account.Chain() == chain {
			return true
		}
	}
	return false
}

func (acc *Accounts) FindKey(key AccountKey) (Account, err.Code) {
	value := acc.data.Load(key)
	if value != nil {
		// TODO: Should be switchable
		accountOneLedger := &AccountOneLedger{}
		base, _ := comm.Deserialize(value, accountOneLedger)
		if base != nil {
			return base.(Account), err.SUCCESS
		}
		accountEthereum := &AccountEthereum{}
		base, _ = comm.Deserialize(value, accountEthereum)
		if base != nil {
			return base.(Account), err.SUCCESS
		}
		accountBitcoin := &AccountBitcoin{}
		base, _ = comm.Deserialize(value, accountBitcoin)
		if base != nil {
			return base.(Account), err.SUCCESS
		}
		log.Fatal("Can't deserialize", "value", value)
	}
	return nil, err.SUCCESS
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

// List out all of the accounts
func (acc *Accounts) Dump() {
	list := acc.FindAll()
	size := len(list)

	for i := 0; i < size; i++ {
		account := list[i]
		log.Info("Account", "Name", account.Name(), "Key", account.AccountKey(), "Type", reflect.TypeOf(account))
	}
}

// Polymorphism
type Account interface {
	Name() string
	Chain() data.ChainType

	AccountKey() AccountKey
	PublicKey() PublicKey
	PrivateKey() PrivateKey

	//AddPublicKey(PublicKey)
	//AddPrivateKey(PrivateKey)

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
func NewAccount(newType data.ChainType, name string, key PublicKey, priv PrivateKey) Account {
	switch newType {

	case data.ONELEDGER:
		return &AccountOneLedger{
			AccountBase{
				Type:       newType,
				Key:        NewAccountKey(key),
				Name:       name,
				PublicKey:  key,
				PrivateKey: priv,
			},
		}

	case data.BITCOIN:
		return &AccountBitcoin{
			AccountBase{
				Type:       newType,
				Key:        NewAccountKey(key),
				Name:       name,
				PublicKey:  key,
				PrivateKey: priv,
			},
		}

	case data.ETHEREUM:
		return &AccountEthereum{
			AccountBase{
				Type:       newType,
				Key:        NewAccountKey(key),
				Name:       name,
				PublicKey:  key,
				PrivateKey: priv,
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

/*
func (account *AccountOneLedger) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountOneLedger) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}
*/

func (account *AccountOneLedger) Name() string {
	return account.AccountBase.Name
}

func (account *AccountOneLedger) AccountKey() AccountKey {
	return data.DatabaseKey(account.AccountBase.Key)
}

func (account *AccountOneLedger) PublicKey() PublicKey {
	return account.AccountBase.PublicKey
}

func (account *AccountOneLedger) PrivateKey() PrivateKey {
	return account.AccountBase.PrivateKey
}

func (account *AccountOneLedger) AsString() string {
	return account.AccountBase.Name
}

func (account *AccountOneLedger) Chain() data.ChainType {
	return data.ONELEDGER
}

// Bitcoin

// Information we need for a Bitcoin account
type AccountBitcoin struct {
	AccountBase
}

/*
func (account *AccountBitcoin) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountBitcoin) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}
*/

func (account *AccountBitcoin) Name() string {
	return account.AccountBase.Name
}

func (account *AccountBitcoin) AccountKey() AccountKey {
	return data.DatabaseKey(account.AccountBase.Key)
}

func (account *AccountBitcoin) PublicKey() PublicKey {
	return account.AccountBase.PublicKey
}

func (account *AccountBitcoin) PrivateKey() PrivateKey {
	return account.AccountBase.PrivateKey
}

func (account *AccountBitcoin) AsString() string {
	buffer := fmt.Sprintf("%x", account.AccountKey())
	return "BTC:" + account.AccountBase.Name + ":" + buffer
}

func (account *AccountBitcoin) Chain() data.ChainType {
	return data.BITCOIN
}

// Ethereum

// Information we need for an Ethereum account
type AccountEthereum struct {
	AccountBase
}

/*
func (account *AccountEthereum) AddPublicKey(key PublicKey) {
	account.PublicKey = key
}

func (account *AccountEthereum) AddPrivateKey(key PrivateKey) {
	account.PrivateKey = key
}
*/

func (account *AccountEthereum) Name() string {
	return account.AccountBase.Name
}

func (account *AccountEthereum) AccountKey() AccountKey {
	return data.DatabaseKey(account.AccountBase.Key)
}

func (account *AccountEthereum) PublicKey() PublicKey {
	return account.AccountBase.PublicKey
}

func (account *AccountEthereum) PrivateKey() PrivateKey {
	return account.AccountBase.PrivateKey
}

func (account *AccountEthereum) AsString() string {
	buffer := fmt.Sprintf("%x", account.AccountKey())
	return "ETH:" + account.AccountBase.Name + ":" + buffer
}

func (account *AccountEthereum) Chain() data.ChainType {
	return data.ETHEREUM
}
