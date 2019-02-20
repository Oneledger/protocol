/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import (
	"encoding/hex"
	"runtime/debug"
	"strings"

	"github.com/Oneledger/protocol/node/comm"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/chains/driver"

	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/data"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/status"

	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/version"
)

// Top-level list of all query types
func HandleQuery(app Application, path string, arguments map[string]interface{}) (buffer []byte) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Query Panic", "status", r)
			debug.PrintStack()
			buffer, _ = serial.Serialize("Internal Query Error", serial.CLIENT)
		}
	}()

	var result interface{}

	switch path {
	case "/applyValidators":
		result = HandleApplyValidatorQuery(app, arguments)

	case "/createExSendRequest":
		result = HandleCreateExSendRequest(app, arguments)

	case "/createSendRequest":
		result = HandleCreateSendRequest(app, arguments)

	case "/createMintRequest":
		result = HandleCreateMintRequest(app, arguments)

	case "/createSwapRequest":
		result = HandleSwapRequest(app, arguments)

	case "/nodeName":
		result = HandleNodeNameQuery(app, arguments)

	case "/identity":
		result = HandleIdentityQuery(app, arguments)

	case "/signTransaction":
		result = HandleSignTransaction(app, arguments)

	case "/accountPublicKey":
		result = HandleAccountPublicKeyQuery(app, arguments)

	case "/accountKey":
		result = HandleAccountKeyQuery(app, arguments)

	case "/account":
		result = HandleAccountQuery(app, arguments)

	case "/validator":
		result = HandleValidatorQuery(app, arguments)

	case "/balance":
		result = HandleBalanceQuery(app, arguments)

	case "/version":
		result = HandleVersionQuery(app, arguments)

	case "/currencyAddress":
		result = HandleCurrencyAddressQuery(app, arguments)

	case "/sequenceNumber":
		result = HandleSequenceNumberQuery(app, arguments)

	case "/testScript":
		result = HandleTestScript(app, arguments)

	default:
		result = HandleError("Unknown Query", path, arguments)
	}

	buffer, err := serial.Serialize(result, serial.CLIENT)
	if err != nil {
		log.Debug("Failed to serialize query")
	}

	return buffer
}

func HandleApplyValidatorQuery(application Application, arguments map[string]interface{}) interface{} {
	log.Debug("HandleApplyValidatorQuery", "arguments", arguments)

	var argsHolder comm.ApplyValidatorArguments

	text := arguments["parameters"].([]byte)

	args, err := serial.Deserialize(text, argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize ApplyValidatorArguments", "error", err)
		return "Could not deserialize ApplyValidatorArguments" + " error=" + err.Error()
	}

	amount := args.(*comm.ApplyValidatorArguments).Amount
	idName := args.(*comm.ApplyValidatorArguments).Id
	purge := args.(*comm.ApplyValidatorArguments).Purge

	nodeAccount := action.GetNodeAccount(application)
	if nodeAccount == nil {
		return "Node account is not found."
	}

	nodeAccountBalance := application.Balances.Get(data.DatabaseKey(nodeAccount.AccountKey()), false)

	if &nodeAccountBalance == nil {
		return "Missing balance for node account."
	}

	identity, ok := application.Identities.FindName(idName)
	if ok != status.SUCCESS {
		return "The identity is not found." + " id:" + idName
	}

	identityBalance := application.Balances.Get(data.DatabaseKey(identity.AccountKey), false)

	if &identityBalance == nil {
		return "Missing balance for identity."
	}

	//check who is running the applyValidator command, if 10 VTs or more, then it's an administrator
	tenVT := data.NewCoinFromFloat(10.0, "VT")
	oneVT := data.NewCoinFromFloat(1.0, "VT")

	stake := data.NewCoinFromFloat(amount, "VT")

	if identity.AccountKey.String() != nodeAccount.AccountKey().String() {
		if nodeAccountBalance.GetAmountByName("VT").LessThanCoin(tenVT) {
			return "Node account has less than 10 VTs. It can not run applyValidator for other identities."
		}
	}

	//when removing validator, the node account doesn't have to have a VT
	if purge == true {

		found := false
		for _, validator := range application.Validators.Approved {
			if validator.AccountKey.String() == identity.AccountKey.String() {
				found = true
			}
		}

		if found == false {
			return "Node account is not a validator."
		}

		stake = identityBalance.GetAmountByName("VT")

	} else {

		if amount == 0.0 {
			return "Missing an amount argument."
		}

		if identityBalance.GetAmountByName("VT").LessThanCoin(stake) {
			return "Validator Token is not enough."
		}

		if stake.LessThanCoin(oneVT) {
			return "Validator Token amount can not be less than 1."
		}

	}

	sequence := SequenceNumber(application, nodeAccount.AccountKey().Bytes())

	validator := &action.ApplyValidator{
		Base: action.Base{
			Type:     action.APPLY_VALIDATOR,
			ChainId:  ChainId,
			Owner:    nodeAccount.AccountKey(),
			Signers:  GetSigners(nodeAccount.AccountKey(), application),
			Sequence: sequence.Sequence,
		},

		AccountKey:        identity.AccountKey,
		Identity:          identity.Name,
		NodeName:          global.Current.NodeName,
		TendermintAddress: identity.TendermintAddress,
		TendermintPubKey:  identity.TendermintPubKey,
		Stake:             stake,
		Purge:             purge,
	}

	signed := SignTransaction(action.Transaction(validator), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleCreateExSendRequest(application Application, arguments map[string]interface{}) interface{} {
	conv := convert.NewConvert()

	result := make([]byte, 0)

	var argsHolder comm.ExSendArguments

	text := arguments["parameters"].([]byte)

	argsDeserialized, err := serial.Deserialize(text, argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize ExSendArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.ExSendArguments)

	partyKey := AccountKey(application, args.SenderId)
	cpartyKey := AccountKey(application, args.ReceiverId)

	fee := conv.GetCoinFromFloat(args.Fee, "OLT")
	gas := conv.GetCoinFromUnits(args.Gas, "OLT")

	amount := conv.GetCoinFromFloat(args.Amount, args.Currency)
	chain := conv.GetChainFromCurrency(args.Chain)

	sender := CurrencyAddress(application, conv.GetCurrency(args.Currency), args.SenderId)
	reciever := CurrencyAddress(application, conv.GetCurrency(args.Currency), args.ReceiverId)

	sequence := SequenceNumber(application, partyKey)

	exSend := &action.ExternalSend{
		Base: action.Base{
			Type:     action.EXTERNAL_SEND,
			ChainId:  ChainId,
			Signers:  GetSigners(sender, application),
			Owner:    partyKey,
			Target:   cpartyKey,
			Sequence: sequence.Sequence,
		},
		Gas:      gas,
		Fee:      fee,
		Chain:    chain,
		Sender:   string(sender),
		Receiver: string(reciever),
		Amount:   amount,
	}

	signed := SignTransaction(action.Transaction(exSend), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleCreateSendRequest(application Application, arguments map[string]interface{}) interface{} {
	conv := convert.NewConvert()

	result := make([]byte, 0)

	var argsHolder comm.SendArguments

	text := arguments["parameters"].([]byte)

	argsDeserialized, err := serial.Deserialize(text, argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize SendArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.SendArguments)

	if args.Party == "" {
		log.Error("Missing Party argument")
		return result
	}

	if args.CounterParty == "" {
		log.Error("Missing CounterParty argument")
		return result
	}

	// TODO: Can't convert identities to accounts, this way!
	party := AccountKey(application, args.Party)
	counterParty := AccountKey(application, args.CounterParty)

	if party == nil || counterParty == nil {
		log.Error("System doesn't recognize the parties", "args", args,
			"party", party, "counterParty", counterParty)
		return result
	}

	if args.Currency == "" || args.Amount == 0.0 {
		log.Error("Missing an amount argument")
		return result
	}

	amount := conv.GetCoinFromFloat(args.Amount, args.Currency)

	fee := conv.GetCoinFromFloat(args.Fee, "OLT")

	sendTo := action.SendTo{
		AccountKey: counterParty,
		Amount:     amount,
	}

	if conv.HasErrors() {
		log.Error("Conversion error", "error", conv.GetErrors())
		return result
	}

	sequence := SequenceNumber(application, party.Bytes())

	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  ChainId,
			Owner:    party,
			Signers:  GetSigners(party, application),
			Sequence: sequence.Sequence,
		},
		SendTo: sendTo,
		Fee:    fee,
	}

	signed := SignTransaction(action.Transaction(send), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleCreateMintRequest(application Application, arguments map[string]interface{}) interface{} {
	result := make([]byte, 0)

	var argsHolder comm.SendArguments

	text := arguments["parameters"].([]byte)

	argsDeserialized, err := serial.Deserialize(text, argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize SendArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.SendArguments)

	conv := convert.NewConvert()

	if args.Party == "" {
		log.Warn("Missing Party arguments", "args", args)
		return result
	}

	zero := AccountKey(application, "Zero")
	party := AccountKey(application, args.Party)

	if party == nil || zero == nil {
		log.Warn("Missing Party information", "args", args, "party", party, "zero", zero)
		return result
	}

	amount := conv.GetCoinFromFloat(args.Amount, args.Currency)

	log.Debug("Getting TestMint Account Balances")
	zeroBalance := Balance(application, zero).(*data.Balance).GetAmountByName(args.Currency)

	if &zeroBalance == nil {
		log.Warn("Missing Balances", "zero", zero)
		return result
	}

	if zeroBalance.LessThanEqual(0) {
		log.Warn("No more money left...")
		return result
	}

	sendTo := action.SendTo{
		AccountKey: party,
		Amount:     amount,
	}

	// Fixed price for fees for minting
	fee := data.NewCoinFromFloat(1, "OLT")

	if conv.HasErrors() {
		log.Error("Conversion error", "error", conv.GetErrors())
		return result
	}

	sequence := SequenceNumber(application, party)

	send := &action.Send{
		Base: action.Base{
			Type:     action.SEND,
			ChainId:  ChainId,
			Signers:  GetSigners(zero, application),
			Owner:    zero,
			Sequence: sequence.Sequence,
		},
		SendTo: sendTo,
		Fee:    fee,
	}

	signed := SignTransaction(action.Transaction(send), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleSwapRequest(application Application, arguments map[string]interface{}) interface{} {
	result := make([]byte, 0)

	var argsHolder comm.SwapArguments

	text := arguments["parameters"].([]byte)

	argsDeserialized, err := serial.Deserialize(text, argsHolder, serial.CLIENT)

	if err != nil {
		log.Error("Could not deserialize SwapArgumentsArguments", "error", err)
		return result
	}

	args := argsDeserialized.(*comm.SwapArguments)

	conv := convert.NewConvert()

	partyKey := AccountKey(application, args.Party)
	counterPartyKey := AccountKey(application, args.CounterParty)

	fee := conv.GetCoinFromFloat(args.Fee, "OLT")
	gas := conv.GetCoinFromUnits(args.Gas, "OLT")

	amount := conv.GetCoinFromFloat(args.Amount, args.Currency)
	exchange := conv.GetCoinFromFloat(args.Exchange, args.Excurrency)

	if conv.HasErrors() {
		log.Error("Conversion error", "error", conv.GetErrors())
		return result
	}

	account := make(map[data.ChainType][]byte)
	counterAccount := make(map[data.ChainType][]byte)

	account[conv.GetChainFromCurrency(args.Currency)] = CurrencyAddress(application, conv.GetCurrency(args.Currency), args.Party)
	account[conv.GetChainFromCurrency(args.Excurrency)] = CurrencyAddress(application, conv.GetCurrency(args.Excurrency), args.Party)

	party := action.Party{Key: partyKey, Accounts: account}
	counterParty := action.Party{Key: counterPartyKey, Accounts: counterAccount}

	swapInit := action.SwapInit{
		Party:        party,
		CounterParty: counterParty,
		Fee:          fee,
		Gas:          gas,
		Amount:       amount,
		Exchange:     exchange,
		Nonce:        args.Nonce,
	}

	sequence := SequenceNumber(application, partyKey.Bytes())

	signers := GetSigners(partyKey, application)
	if signers == nil || len(signers) == 0 {
		return result
	}

	swap := &action.Swap{
		Base: action.Base{
			Type:     action.SWAP,
			ChainId:  ChainId,
			Signers:  signers,
			Owner:    partyKey,
			Target:   counterPartyKey,
			Sequence: sequence.Sequence,
		},
		SwapMessage: swapInit,
		Stage:       action.SWAP_MATCHING,
	}

	signed := SignTransaction(action.Transaction(swap), application)
	packet := action.PackRequest(signed)

	return packet
}

func HandleNodeNameQuery(app Application, arguments map[string]interface{}) interface{} {
	return global.Current.NodeName
}

func HandleSignTransaction(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("SignTransactionQuery", "arguments", arguments)

	var tx action.Transaction

	text := arguments["parameters"].([]byte)
	transaction, transactionErr := serial.Deserialize(text, tx, serial.CLIENT)

	signatures := []action.TransactionSignature{}

	if transactionErr != nil {
		log.Error("Could not deserialize a transaction", "error", transactionErr)
		return signatures
	}

	accountKey := transaction.(action.Transaction).GetOwner()

	if accountKey == nil {
		log.Error("Account key is null", "transaction", transaction)
		return signatures
	}

	account, accountStatus := app.Accounts.FindKey(accountKey)

	if accountStatus != status.SUCCESS {
		log.Error("Could not find an account", "status", accountStatus)
		return signatures
	}

	privateKey := account.PrivateKey()

	signature, signatureError := privateKey.Sign(text)

	if signatureError != nil {
		log.Error("Could not sign a transaction", "error", signatureError)
		return signatures
	}

	signatures = append(signatures, action.TransactionSignature{signature})

	return signatures
}

func HandleAccountPublicKeyQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("AccountPublicKeyQuery", "arguments", arguments)

	text := arguments["parameters"].([]byte)

	account, accountStatus := app.Accounts.FindKey(text)

	if accountStatus != status.SUCCESS {
		log.Error("Could not find an account", "status", accountStatus)
		return "Missing Account: " + string(text)
	}

	privateKey := account.PrivateKey()
	publicKey := privateKey.PubKey()

	if publicKey.Equals(account.PublicKey()) == false {
		log.Warn("Public keys don't match", "derivedPublicKey",
			publicKey, "accountPublicKey", account.PublicKey())
	}
	return publicKey
}

// Get the account information for a given user
func HandleAccountKeyQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("AccountKeyQuery", "arguments", arguments)

	text := arguments["parameters"].([]byte)

	name := ""
	parts := strings.Split(string(text), "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return AccountKey(app, name)
}

func AccountKey(app Application, name string) id.AccountKey {

	// Check Itdentities First
	identity, ok := app.Identities.FindName(name)
	if ok == status.SUCCESS && identity.Name != "" {
		return identity.AccountKey
	}

	// TODO: This is a bit dangerous (can cause confusion)
	// Maybe this is an AccountName, not an identity
	account, ok := app.Accounts.FindName(name)
	if ok == status.SUCCESS && account.Name() != "" {
		return account.AccountKey()
	}

	accountKey, err := hex.DecodeString(name)
	if err != nil {
		return nil
	}

	identity = app.Identities.FindKey(accountKey)
	if identity.Name != "" {
		return identity.AccountKey
	}
	return accountKey
}

// Get the account information for a given user
func HandleIdentityQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("IdentityQuery", "arguments", arguments)

	text := arguments["parameters"].([]byte)

	name := ""
	parts := strings.Split(string(text), "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return IdentityInfo(app, name)
}

func IdentityInfo(app Application, name string) interface{} {
	if name == "" {
		identities := app.Identities.FindAll()
		return identities
	}

	identity, ok := app.Identities.FindName(name)
	if ok == status.SUCCESS {
		return []id.Identity{identity}
	}
	return "Identity " + name + " Not Found" + global.Current.NodeName
}

// Get the validator information for a given user
func HandleValidatorQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("ValidatorQuery", "arguments", arguments)

	text := arguments["parameters"].([]byte)

	name := ""
	parts := strings.Split(string(text), "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return ValidatorInfo(app, name)
}

func ValidatorInfo(app Application, name string) interface{} {
	if name == "" {
		validators := app.Validators.Approved
		return validators
	}

	return "Validator " + name + " Not Found" + global.Current.NodeName
}

// Get the account information for a given user
func HandleAccountQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("AccountQuery", "arguments", arguments)

	text := arguments["parameters"].([]byte)

	name := ""
	parts := strings.Split(string(text), "=")
	if len(parts) > 1 {
		name = parts[1]
	}
	return AccountInfo(app, name)
}

// AccountInfo returns the information for a given account
func AccountInfo(app Application, name string) interface{} {
	if name == "" {
		accounts := app.Accounts.FindAll()
		return accounts
	}

	account, ok := app.Accounts.FindName(name)
	if ok == status.SUCCESS {
		return []id.Account{account}
	}
	return "Account " + name + " Not Found" + global.Current.NodeName
}

func HandleVersionQuery(app Application, arguments map[string]interface{}) interface{} {
	return version.Fullnode.String()
}

// Get the account information for a given user
func HandleBalanceQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("BalanceQuery", "arguments", arguments)

	text := arguments["parameters"].([]byte)

	var key []byte
	parts := strings.Split(string(text), "=")
	if len(parts) > 1 {

		key, _ = hex.DecodeString(parts[1])
	}
	return Balance(app, key)
}

func Balance(app Application, accountKey []byte) interface{} {
	balance := app.Balances.Get(accountKey, true)
	if balance != nil {
		return balance
	}
	result := data.NewBalance()

	return result
}

func HandleCurrencyAddressQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("CurrencyAddressQuery", "arguments", arguments)

	text := arguments["parameters"].(string)
	conv := convert.NewConvert()
	var chain data.ChainType
	var identity id.Identity
	parts := strings.Split(text, "|")
	if len(parts) == 2 {
		currencyString := strings.Split(parts[0], "=")
		if len(currencyString) > 1 {
			chain = conv.GetChainFromCurrency(currencyString[1])
		} else {
			return nil
		}

		identityString := strings.Split(parts[1], "=")
		if len(identityString) > 1 {
			identity, ok := app.Identities.FindName(identityString[1])
			if ok == status.SUCCESS {
				if chain != data.ONELEDGER {
					//todo : right now, the idenetity.Chain do not contain the address on the other chain
					// need to fix register to remove this part
					_ = identity
					return ChainAddress(app, identityString[1], chain)
				}

			}
		}

	}
	return identity.Chain[chain]
}

func CurrencyAddress(application Application, currency string, identityName string) []byte {
	conv := convert.NewConvert()

	chain := conv.GetChainFromCurrency(currency)

	identity, ok := application.Identities.FindName(identityName)

	if ok == status.SUCCESS {
		if chain != data.ONELEDGER {
			//todo : right now, the idenetity.Chain do not contain the address on the other chain, need to fix register to remove this part
			_ = identity
			return ChainAddress(application, identityName, chain).([]byte)
		}
	}

	return identity.Chain[chain]
}

func ChainAddress(application Application, name string, chain data.ChainType) interface{} {
	accountName := name + "-" + chain.String()

	account, stat := application.Accounts.FindNameOnChain(accountName, chain)

	if stat != status.SUCCESS {
		return nil
	}

	return chaindriver.GetDriver(chain).GetChainAddress(account.GetChainKey())
}

// Return a nicely formatted error message
func HandleError(text string, path string, arguments map[string]interface{}) interface{} {
	// TODO: Add in arguments to output
	return "Unknown Query " + text + " " + path
}

func SignTransaction(transaction action.Transaction, application Application) action.SignedTransaction {
	packet, err := serial.Serialize(transaction, serial.CLIENT)

	signed := action.SignedTransaction{transaction, nil}

	if err != nil {
		log.Error("Failed to Serialize packet: ", "error", err)
	} else {
		request := action.Message(packet)
		arguments := map[string]interface{}{
			"parameters": request,
		}

		response := HandleSignTransaction(application, arguments)

		if response == nil {
			log.Warn("Query returned no signature", "request", request)
		} else {
			signed.Signatures = response.([]action.TransactionSignature)
		}
	}

	return signed
}

func GetSigners(owner []byte, application Application) []id.PublicKey {
	log.Debug("GetSigners", "owner", owner)

	arguments := map[string]interface{}{
		"parameters": owner,
	}
	result := HandleAccountPublicKeyQuery(application, arguments)

	if result == nil {
		return nil
	}

	switch value := result.(type) {
	case id.PublicKey:
		return []id.PublicKey{value}
	}

	log.Debug("Missing Signers", "owner", owner)

	return []id.PublicKey{}
}

// Get the account information for a given user
func HandleSequenceNumberQuery(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("SequenceNumberQuery", "arguments", arguments)

	text := arguments["parameters"].([]byte)

	var key []byte
	parts := strings.Split(string(text), "=")
	if len(parts) > 1 {

		// TODO: Encoded because it is dumped into a string
		key, _ = hex.DecodeString(parts[1])
	}
	return SequenceNumber(app, key)
}

func SequenceNumber(app Application, accountKey []byte) id.SequenceRecord {
	sequenceRecord := id.NextSequence(&app, accountKey)
	return sequenceRecord
}

func HandleTestScript(app Application, arguments map[string]interface{}) interface{} {
	log.Debug("TestScript", "arguments", arguments)
	text := arguments["parameters"].(string)

	results, _ := RunTestScriptName(text)

	return results
}
