/*
	Copyright 2017 - 2018 OneLedger

	Handle setting any options for the node.
*/
package app

import (
	"github.com/Oneledger/protocol/node/comm"
	"strings"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/status"
)

// Arguments for registration
type RegisterArguments struct {
	Identity   string
	Chain      string
	PublicKey  string
	PrivateKey string
}

func HandleSet(app Application, path string, arguments map[string]interface{}) []byte {
	var result interface{}

	switch path {
	case "/account":
		result = HandleSetAccount(app, arguments)

	case "/register":
		result = HandleRegisterIdentity(app, arguments)
	}

	if result == nil {
		return nil
	}

	buffer, err := clSerializer.Serialize(result)
	if err != nil {
		log.Fatal("Failed to serialize query", "err", err)
	}
	return buffer
}

func GetChain(chainName string) data.ChainType {
	switch chainName {
	case "OneLedger":
		return data.ONELEDGER
	case "BitCoin":
		return data.BITCOIN
	case "Ethereum":
		return data.ETHEREUM
	}
	return data.UNKNOWN
}

// TODO: The datatype for Key, depends on Chain
func GetKeys(chain data.ChainType, name string, publicKey string, privateKey string) (id.PublicKeyED25519, id.PrivateKeyED25519) {

	// TODO: Need to push the passphrase back through the CLI
	priv, public := id.GenerateKeys([]byte(name+"as password"), true)
	return public, priv
}

// TODO: Should be in common library
func GetBool(boolean string) bool {
	if strings.EqualFold(boolean, "true") {
		return true
	}
	if strings.EqualFold(boolean, "false") {
		return false
	}

	// TODO: matches golang?
	return false
}

// TODO: Pass in App pointer?
func HandleSetAccount(app Application, arguments map[string]interface{}) interface{} {
	chain := GetChain(arguments["Chain"].(string))
	accountName := arguments["Account"].(string)
	nodeAccount := arguments["NodeAccount"].(bool)

	publicKey, privateKey := GetKeys(chain, accountName,
		arguments["PublicKey"].(string), arguments["PrivateKey"].(string))

	chainKey := arguments["ChainKey"].(string)

	account, err := app.Accounts.FindNameOnChain(accountName, chain)
	if err == status.SUCCESS {
		return "Account has already been set up."
	}

	AddAccount(&app, accountName, chain, publicKey, privateKey, chainKey, nodeAccount)

	account, err = app.Accounts.FindNameOnChain(accountName, chain)
	if err == status.SUCCESS {
		return account
	}

	return "Error in Setting up Account"
}

func HandleRegisterIdentity(app Application, arguments map[string]interface{}) interface{} {

	identity := arguments["Identity"].(string)
	accountName := arguments["Account"].(string)
	fee := arguments["Fee"].(float64)

	account, ok := app.Accounts.FindName(accountName)
	if ok != status.SUCCESS {
		log.Warn("Missing Registration Account", "name", accountName)
		return "ERROR: Account Not Found"
	}

	// TODO Broadcast the transaction
	transaction := CreateRegisterRequest(app.ClientContext, identity, account.AccountKey(), fee)
	action.BroadcastTransaction(app.ClientContext, transaction, false)

	return "Broadcast Identity"
}

// TODO: Called by olfullnode, not olclient?
func CreateRegisterRequest(ctx comm.ClientContext, identityName string, accountKey id.AccountKey, fee float64) action.Transaction {
	LoadPrivValidatorFile()

	reg := &action.Register{
		Base: action.Base{
			Type:     action.REGISTER,
			ChainId:  ChainId,
			Owner:    accountKey,
			Signers:  action.GetSigners(ctx, accountKey), // TODO: Server-side? Then this is wrong
			Sequence: global.Current.Sequence,
		},
		Identity:          identityName,
		NodeName:          global.Current.NodeName,
		AccountKey:        accountKey,
		TendermintAddress: global.Current.TendermintAddress,
		TendermintPubKey:  global.Current.TendermintPubKey,
		Fee:               data.NewCoinFromFloat(fee, "OLT"),
	}
	return reg
}

// TODO: This probably doesn't work. It was replaced by the SDK direct connection
// Handle a SetOption ABCi reqeust
func SetOption(app *Application, key string, value string) bool {
	log.Debug("Setting Application Options", "key", key, "value", value)

	switch key {

	case "Register":
		var args = &RegisterArguments{}
		// I am not sure where this comes from; doing client deser for now
		err := clSerializer.Deserialize([]byte(value), args)
		//result, err := serial.Deserialize([]byte(value), &arguments, serial.NETWORK)

		if err != nil {
			log.Error("Can't set options", "status", err)
			return false
		}

		privateKey, publicKey := id.GenerateKeys([]byte(args.Identity), true) // TODO: Switch with passphrase
		AddAccount(app, args.Identity, id.ParseAccountType(args.Chain), publicKey, privateKey, nil, false)

	default:
		log.Warn("Unknown Option", "key", key)
		return false
	}
	return true
}
