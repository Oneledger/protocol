/*
	Copyright 2017-2018 OneLedger

	Identities management for any of the associated chains
*/
package id

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

// Temporary typing for signatures
type Signature = []byte

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
		log.Fatal("Key exists already", "key", account.AccountKey())
		log.Debug("Key is being updated", "key", account.AccountKey())
	}

	buffer, err := serial.Serialize(account, serial.PERSISTENT)
	if err != nil {
		log.Fatal("Failed to Deserialize account: ", err)
	}

	acc.data.Store(account.AccountKey(), buffer)
	acc.data.Commit()
}

func (acc *Accounts) Delete(account Account) {
}

func (acc *Accounts) Exists(newType data.ChainType, name string) bool {
	account := NewAccount(newType, name, NilPublicKey(), NilPrivateKey())

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
	account := NewAccount(data.ONELEDGER, identity.Name+"-OneLedger", NilPublicKey(), NilPrivateKey())
	return acc.Find(account)
}

func (acc *Accounts) FindNameOnChain(name string, chain data.ChainType) (Account, err.Code) {
	log.Debug("FindNameOnChain", "name", name, "chain", chain)

	// TODO: Should be replaced with a real index
	for _, entry := range acc.FindAll() {
		if Matches(entry, name, chain) {
			return entry, err.SUCCESS
		}
	}
	return nil, err.MISSING_VALUE
}

func (acc *Accounts) FindName(name string) (Account, err.Code) {
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

func (acc *Accounts) FindKey(key AccountKey) (Account, err.Code) {

	value := acc.data.Load(key)
	account := Account(nil)

	result, status := serial.Deserialize(value, account, serial.PERSISTENT)

	if status != nil {
		log.Fatal("Failed to Deserialize Account", "status", status)
	}

	//log.Dump("Deserialized", value, result, account)
	return result.(Account), err.SUCCESS

	/*

			if value != nil {

				// TODO: Should be switchable
				accountOneLedger := &AccountOneLedger{}
				base, _ := serial.Deserialize(value, accountOneLedger, serial.PERSISTENT)
				if base != nil {
					return base.(Account), err.SUCCESS
				}

				accountEthereum := &AccountEthereum{}
				base, _ = serial.Deserialize(value, accountEthereum, serial.PERSISTENT)
				if base != nil {
					return base.(Account), err.SUCCESS
				}

				accountBitcoin := &AccountBitcoin{}
				base, _ = serial.Deserialize(value, accountBitcoin, serial.PERSISTENT)
				if base != nil {
					return base.(Account), err.SUCCESS
				}
				log.Fatal("Can't deserialize", "value", value)
			}
		return nil, err.SUCCESS
	*/
}

func (acc *Accounts) FindAll() []Account {
	log.Debug("Begin Account FindAll")
	keys := acc.data.List()

	size := len(keys)
	results := make([]Account, size, size)

	for i := 0; i < size; i++ {
		account, status := acc.FindKey(keys[i])
		if status != err.SUCCESS {
			log.Fatal("Missing Account", "status", status, "account", account)
		}

		/*
			// TODO: This is dangerous...
			account := &AccountOneLedger{}

			base, err := serial.Deserialize(acc.data.Load(keys[i]), account, serial.PERSISTENT)
			if err != nil {
				log.Fatal("Failed to Deserialize Account at index ", "i", i, "err", err)
			}
		*/

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
	acc.data.Close()
}

// Polymorphism
type Account interface {
	Name() string
	Chain() data.ChainType

	AccountKey() AccountKey
	PublicKey() PublicKey
	PrivateKey() PrivateKey

	Export() AccountExport
	AsString() string

	//AddPublicKey(PublicKey)
	//AddPrivateKey(PrivateKey)
}

// AccountExport struct holds important account info in a
type AccountExport struct {
	Type       string
	AccountKey string
	Name       string

	// Balance must come from utxo database, fill when needed
	Balance  string
	NodeName string
}

func init() {
	serial.Register(AccountExport{})
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

// Information we need about our own fullnode identities
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

func (account *AccountOneLedger) AsString() string {
	buffer := fmt.Sprintf("%x", account.AccountKey())
	return "OneLedger:" + account.AccountBase.Name + ":" + buffer
}

func (account *AccountOneLedger) Chain() data.ChainType {
	return data.ONELEDGER
}

func (account *AccountOneLedger) Export() AccountExport {
	accountType := getAccountType(account.AccountBase.Type)
	name := account.Name()
	key := hex.EncodeToString(account.AccountKey())
	return AccountExport{
		Type:       accountType,
		AccountKey: key,
		Name:       name,
		NodeName:   global.Current.NodeName,
	}
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

func (account *AccountBitcoin) AsString() string {
	buffer := fmt.Sprintf("%x", account.AccountKey())
	return "BTC:" + account.AccountBase.Name + ":" + buffer
}

func (account *AccountBitcoin) Chain() data.ChainType {
	return data.BITCOIN
}

func (account *AccountBitcoin) Export() AccountExport {
	accountType := getAccountType(account.AccountBase.Type)
	name := account.Name()
	key := hex.EncodeToString(account.AccountKey())
	return AccountExport{
		Type:       accountType,
		AccountKey: key,
		Name:       name,
		NodeName:   global.Current.NodeName,
	}
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

func (account *AccountEthereum) AsString() string {
	buffer := fmt.Sprintf("%x", account.AccountKey())
	return "ETH:" + account.AccountBase.Name + ":" + buffer
}

func (account *AccountEthereum) Chain() data.ChainType {
	return data.ETHEREUM
}

func (account *AccountEthereum) Export() AccountExport {
	accountType := getAccountType(account.AccountBase.Type)
	name := account.Name()
	key := hex.EncodeToString(account.AccountKey())
	return AccountExport{
		Type:       accountType,
		AccountKey: key,
		Name:       name,
		NodeName:   global.Current.NodeName,
	}
}
