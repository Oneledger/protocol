/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains
*/
package id

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Oneledger/protocol/node/data"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

// Temporary typing for signatures
type Signature = []byte

// The persistent collection of all accounts known by this node
type Accounts struct {
	store data.Datastore
}

func NewAccounts(name string) *Accounts {
	store := data.NewDatastore(name, data.PERSISTENT)
	return &Accounts{
		store: store,
	}
}

func (acc *Accounts) Add(account Account) {
	if value := acc.store.Get(account.AccountKey()); value != nil {
		//log.Fatal("Key exists already", "key", account.AccountKey())
		log.Debug("Key is being updated", "key", account.AccountKey())
	}

	session := acc.store.Begin()
	session.Set(account.AccountKey(), account)
	session.Commit()
}

func (acc *Accounts) Delete(account Account) {
}

func (acc *Accounts) Exists(newType data.ChainType, name string) bool {
	// TODO: Probably shouldn't need to create a fake account here...
	account := NewAccount(newType, name, NilPublicKey(), NilPrivateKey())
	return acc.store.Exists(account.AccountKey())
}

func (acc *Accounts) Find(account Account) (Account, status.Code) {
	return acc.FindKey(account.AccountKey())
}

func (acc *Accounts) FindIdentity(identity Identity) (Account, status.Code) {
	// TODO: Should have better name mapping between identities and accounts
	account := NewAccount(data.ONELEDGER, identity.Name+"-OneLedger", NilPublicKey(), NilPrivateKey())
	return acc.Find(account)
}

func (acc *Accounts) FindNameOnChain(name string, chain data.ChainType) (Account, status.Code) {
	// TODO: Should be replaced with a real index
	for _, entry := range acc.FindAll() {
		if Matches(entry, name, chain) {
			return entry, status.SUCCESS
		}
	}
	return nil, status.MISSING_DATA
}

func (acc *Accounts) FindName(name string) (Account, status.Code) {
	return acc.FindNameOnChain(name, data.ONELEDGER)
}

func Matches(account Account, name string, chain data.ChainType) bool {
	// TODO: Incorrect, all names for all chains are in the same scope...
	if strings.EqualFold(account.Name(), name) {
		if account.Chain() == chain {
			return true
		}
	}
	return false
}

func (acc *Accounts) FindKey(key AccountKey) (Account, status.Code) {
	interim := acc.store.Get(key)

	if interim == nil {
		return nil, status.MISSING_DATA
	}

	result := interim.(Account)
	return result.(Account), status.SUCCESS

}

func (acc *Accounts) FindAll() []Account {
	keys := acc.store.FindAll()

	size := len(keys)
	results := make([]Account, size, size)

	for i := 0; i < size; i++ {

		account, ok := acc.FindKey(keys[i])
		if ok != status.SUCCESS {
			log.Fatal("Missing Account", "status", ok, "account", account)
		}
		results[i] = account
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

func (acc *Accounts) Close() {
	acc.store.Close()
}

// Polymorphism
type Account interface {
	Name() string
	Chain() data.ChainType

	AccountKey() AccountKey
	PublicKey() PublicKey
	PrivateKey() PrivateKey

	String() string

	//AddPublicKey(PublicKey)
	//AddPrivateKey(PrivateKey)
}

func init() {
	var prototype Account
	serial.RegisterInterface(&prototype)
}

func getAccountType(chain data.ChainType) string {
	switch chain {
	case data.ONELEDGER:
		return "OneLedger"

	case data.ETHEREUM:
		return "Ethereum"

	case data.BITCOIN:
		return "Bitcoin"
	}
	return "Unknown"
}

type AccountBase struct {
	Type data.ChainType `json:"type"`
	Key  AccountKey     `json:"key"`
	Name string         `json:"name"`

	// TODO: Should handle key polymorphism properly..
	PublicKey  PublicKeyED25519  `json:"publicKey"`
	PrivateKey PrivateKeyED25519 `json:"privateKey"`
}

func init() {
	serial.Register(AccountBase{})
	serial.Register(AccountOneLedger{})
	serial.Register(AccountBitcoin{})
	serial.Register(AccountEthereum{})
}

// Create a new account for a given chain
func NewAccount(newType data.ChainType, name string, key PublicKeyED25519, priv PrivateKeyED25519) Account {
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
			//todo: change to olt wallet auth
			//NewAccountKey(key),
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
			//todo: change to bitcoin auth
			//NewAccountKey(key),
		}

	case data.ETHEREUM:
		return &AccountEthereum{
			AccountBase{
				Type:       newType,
				Key:        NewAccountKey(key),
				Name:       name,
				PublicKey:  NilPublicKey(),
				PrivateKey: NilPrivateKey(),
			},
			//ethereum.GetAuth(),
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

// Information we need about our own olfullnode identities
type AccountOneLedger struct {
	AccountBase
	//todo: need to be change to the right type
	//ChainAuth AccountKey
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

//String used in fmt and Dump
func (account *AccountOneLedger) String() string {
	buffer := fmt.Sprintf("%x", account.AccountKey())
	return "OneLedger:" + account.AccountBase.Name + ":" + buffer
}

func (account *AccountOneLedger) Chain() data.ChainType {
	return data.ONELEDGER
}

// Bitcoin

// Information we need for a Bitcoin account
type AccountBitcoin struct {
	AccountBase
	//todo: need to be change to the right type
	//ChainAuth AccountKey
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

//String used in fmt and Dump
func (account *AccountBitcoin) String() string {
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

	//ChainAuth *bind.TransactOpts
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

//String used in fmt and Dump
func (account *AccountEthereum) String() string {
	buffer := fmt.Sprintf("%x", account.AccountKey())
	return "ETH:" + account.AccountBase.Name + ":" + buffer
}

func (account *AccountEthereum) Chain() data.ChainType {
	return data.ETHEREUM
}
