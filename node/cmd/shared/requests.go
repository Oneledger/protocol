/*
	Copyright 2017 - 2018 OneLedger

	Structures and functions for getting command line arguments, and functions
	to convert these into specific requests.
*/
package shared

import (
	"github.com/Oneledger/protocol/node/version"
	"os"
	"regexp"
	"strconv"
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/app"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

// Prepare a transaction to be issued.
func SignAndPack(transaction action.Transaction) []byte {
	return action.SignAndPack(transaction)
}

// Registration
type AccountArguments struct {
	Account     string
	Chain       string
	PublicKey   string
	PrivateKey  string
	NodeAccount bool
}

func UpdateAccountRequest(args *AccountArguments) interface{} {
	return &app.SDKSet{
		Path: "/account",
		Arguments: map[string]string{
			"Account":     args.Account,
			"Chain":       args.Chain,
			"PublicKey":   args.PublicKey,
			"PrivateKey":  args.PrivateKey,
			"NodeAccount": "true",
		},
	}
}

// Registration
type RegisterArguments struct {
	Identity string
	Account  string
	NodeName string
}

// Create a request to register a new identity with the chain
func RegisterIdentityRequest(args *RegisterArguments) interface{} {
	//signers := GetSigners()

	// TODO: Need to check errors here
	//accountKey := GetAccountKey(args.Account)

	// TODO: Need to have access to this data?
	//app.LoadPrivValidatorFile()

	return &app.SDKSet{
		Path: "/register",
		Arguments: map[string]string{
			"Identity": args.Identity,
			"Account":  args.Account,
			"NodeName": args.NodeName,
		},
	}
}

type BalanceArguments struct {
}

func CreateBalanceRequest(args *BalanceArguments) []byte {
	return []byte(nil)
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateApplyValidatorRequest(args *comm.ApplyValidatorArguments) []byte {
	request, err := serial.Serialize(args, serial.CLIENT)
	if err != nil {
		log.Error("Failed to Serialize arguments: ", err)
		return nil
	}

	response := comm.Query("/applyValidators", request)

	if response == nil {
		log.Warn("Query returned no response", "request", request)
		return nil
	}

	return response.([]byte)
}

type InstallArguments struct {
	Owner    string
	Name     string
	Version  string
	File     string
	Currency string
	Gas      string
	Fee      string
}

type ExecuteArguments struct {
	Owner    string
	Name     string
	Version  string
	Currency string
	Gas      string
	Fee      string
}

type ContractArguments struct {
    Address string
    CallString string
    CallFrom string
    SourceCode string
    Value int
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateSendRequest(args *comm.SendArguments) []byte {
	request, err := serial.Serialize(args, serial.CLIENT)

	if err != nil {
		log.Error("Failed to Serialize arguments: ", err)
		return nil
	}

	response := comm.Query("/createSendRequest", request)

	if response == nil {
		log.Warn("Query returned no response", "request", request)
		return nil
	}

	return response.([]byte)
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateMintRequest(args *comm.SendArguments) []byte {
	request, err := serial.Serialize(args, serial.CLIENT)

	if err != nil {
		log.Error("Failed to Serialize arguments: ", err)
		return nil
	}

	response := comm.Query("/createMintRequest", request)

	if response == nil {
		log.Warn("Query returned no response", "request", request)
		return nil
	}

	return response.([]byte)
}

// Create a swap request
func CreateSwapRequest(args *comm.SwapArguments) []byte {
	request, err := serial.Serialize(args, serial.CLIENT)

	if err != nil {
		log.Error("Failed to Serialize arguments: ", err)
		return nil
	}

	response := comm.Query("/createSwapRequest", request)

	if response == nil {
		log.Warn("Query returned no response", "request", request)
		return nil
	}

	return response.([]byte)
}

func CreateExSendRequest(args *comm.ExSendArguments) []byte {
	request, err := serial.Serialize(args, serial.CLIENT)

	if err != nil {
		log.Error("Failed to Serialize arguments: ", err)
		return nil
	}

	response := comm.Query("/createExSendRequest", request)

	if response == nil {
		log.Warn("Query returned no response", "request", request)
		return nil
	}

	return response.([]byte)
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateInstallRequest(args *InstallArguments, script []byte) []byte {
	conv := convert.NewConvert()

	if args.Owner == "" {
		log.Error("Missing Owner argument")
		return nil
	}

	if args.Name == "" {
		log.Error("Missing Name argument")
		return nil
	}

	if args.Version == "" {
		log.Error("Missing Version argument")
		return nil
	}

	if script == nil {
		log.Error("Missing Script")
		return nil
	}

	owner := GetAccountKey(args.Owner)
	if owner == nil {
		log.Fatal("System doesn't recognize the owner", "args", args,
			"owner", owner)
		return nil
	}

	version := ParseVersion(args.Version)
	if version == nil {
		log.Debug("Version error", "version", args.Version)
		return nil
	}

	fee := conv.GetCoin(args.Fee, args.Currency)
	gas := conv.GetCoin(args.Gas, args.Currency)

	sequence := GetSequenceNumber(owner)

	inputs := action.Install{
		Owner:   owner,
		Name:    args.Name,
		Version: *version,
		Script:  script,
	}

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}

	signers := action.GetSigners(owner)
	if signers == nil {
		log.Debug("Missing Signers")
		return nil
	}

	// Create base transaction
	install := &action.Contract{
		Base: action.Base{
			Type:     action.SMART_CONTRACT,
			ChainId:  app.ChainId,
			Owner:    owner,
			Signers:  signers, //action.GetSigners(owner),
			Sequence: sequence,
		},
		Data:     inputs,
		Function: action.INSTALL,
		Fee:      fee,
		Gas:      gas,
	}
	return SignAndPack(action.Transaction(install))
}

func ParseVersion(argsVersion string) *version.Version {
	automata := regexp.MustCompile(`v([0-9]*)\.([0-9]*)\.([0-9]*)`)
	groups := automata.FindStringSubmatch(argsVersion)

	//log.Dump("VersionGroups", groups)
	if groups == nil || len(groups) != 4 {
		log.Debug("ParseVersion", "groups", groups, "groupsLen", len(groups))
		return nil
	}

	major, err := strconv.Atoi(groups[1])
	if err != nil {
		log.Debug("ParseVersion", "major", major, "err", err)
	}

	minor, err := strconv.Atoi(groups[2])
	if err != nil {
		log.Debug("ParseVersion", "minor", minor, "err", err)
	}

	patch, err := strconv.Atoi(groups[3])
	if err != nil {
		log.Debug("ParseVersion", "patch", patch, "err", err)
	}

	return &version.Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

// CreateRequest builds and signs the transaction based on the arguments
func CreateExecuteRequest(args *ExecuteArguments) []byte {
	conv := convert.NewConvert()

	if args.Owner == "" {
		log.Error("Missing Owner argument")
		return nil
	}

	if args.Name == "" {
		log.Error("Missing Name argument")
		return nil
	}

	if args.Version == "" {
		log.Error("Missing Version argument")
		return nil
	}

	owner := GetAccountKey(args.Owner)
	if owner == nil {
		log.Fatal("System doesn't recognize the owner", "args", args,
			"owner", owner)
		return nil
	}

	version := ParseVersion(args.Version)
	if version == nil {
		log.Debug("Version error", "version", args.Version)
		return nil
	}

	fee := conv.GetCoin(args.Fee, args.Currency)
	gas := conv.GetCoin(args.Gas, args.Currency)

	sequence := GetSequenceNumber(owner)

	inputs := action.Execute{
		Owner:   owner,
		Name:    args.Name,
		Version: *version,
	}

	if conv.HasErrors() {
		Console.Error(conv.GetErrors())
		os.Exit(-1)
	}

	// Create base transaction
	execute := &action.Contract{
		Base: action.Base{
			Type:     action.SMART_CONTRACT,
			ChainId:  app.ChainId,
			Owner:    owner,
			Signers:  action.GetSigners(owner),
			Sequence: sequence,
		},
		Data:     inputs,
		Function: action.EXECUTE,
		Fee:      fee,
		Gas:      gas,
	}
	return SignAndPack(action.Transaction(execute))
}
