/*
	Copyright 2017-2018 OneLedger

	Handle Smart Contarct Install and Execute functions
*/
package action

import (
	"bytes"
  "encoding/hex"
  "errors"

	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/Oneledger/protocol/node/version"
)

type ContractFunction int

const (
	INSTALL ContractFunction = iota
	EXECUTE
	COMPARE
)

// Synchronize a swap between two users
type Contract struct {
	Base
	Data     ContractData     `json:"data"`
	Function ContractFunction `json:"function"`
	Gas      data.Coin        `json:"gas"`
	Fee      data.Coin        `json:"fee"`
}

type ContractData interface {
}

type Install struct {
	//ToDo: all the data you need to install contract
	Owner   id.AccountKey
	Name    string
	Version version.Version
	Script  []byte
}

type Execute struct {
	//ToDo: all the data you need to execute contract
	Owner   id.AccountKey
	Name    string
  Address string
  CallString string
	Version version.Version
}

type Compare struct {
	//ToDo: all the data you need to compare contract
	Owner   id.AccountKey
	Name    string
  Address string
  CallString string
	Version version.Version
	Result  OLVMResult
}

func init() {
	serial.Register(Contract{})
	serial.Register(Install{})
	serial.Register(Execute{})
	serial.Register(Compare{})

	var prototype ContractData
	serial.RegisterInterface(&prototype)
}

func getSmartContractCode(app interface{}, address string, owner []byte, name string,version version.Version) (scriptBytes []byte, err error) {
  if address != "" {
    addressBytes, err_ := hex.DecodeString(address)
    if err_ != nil {
      log.Warn("Failed to decode address", "err", err_)
      err = err_
      return
    }
    script, err_ := FindContractCode(addressBytes, false) //TODO: support proof to true?
    if err_ != nil {
      log.Warn("Failed to decode smart contract", "err", err_)
      err = err_
      return
    }
    scriptBytes =  script.Script
  }else{
    smartContracts := GetSmartContracts(app)
    raw := smartContracts.Get(owner)
    if raw == nil {
      log.Warn("Failed to get smart contract record", "owner", owner)
      err = errors.New("Failed to get smart contract rcord")
      return
    }
    scriptRecords := raw.(*data.ScriptRecords)
		versions := scriptRecords.Name[name]
		script := versions.Version[version.String()]
    scriptBytes = script.Script
  }
  return
}

func (transaction *Contract) TransactionType() Type {
	return transaction.Base.Type
}

func (transaction *Contract) Validate() status.Code {
	log.Debug("Validating Smart Contract Transaction")

	baseValidate := transaction.Base.Validate()

	if baseValidate != status.SUCCESS {
		return baseValidate
	}

	//check that the data supplied is valid and no security problems
	if transaction.Function == INSTALL {
		if transaction.Data == nil {
			log.Debug("Smart Contract Missing Data", "transaction data", transaction)
			return status.MISSING_DATA
		}

		installData := transaction.Data.(Install)
		if installData.Name == "" {
			log.Debug("Smart Contract Missing Data", "name", installData.Name)
			return status.MISSING_DATA
		}

		if installData.Version.String() == "" {
			log.Debug("Smart Contract Missing Data", "version", installData.Version)
			return status.MISSING_DATA
		}

		if installData.Script == nil {
			log.Debug("Smart Contract Missing Data", "script", installData.Script)
			return status.MISSING_DATA
		}

		if !transaction.Fee.IsCurrency("OLT") {
			log.Debug("Wrong Fee token", "fee", transaction.Fee)
			return status.INVALID
		}

		if transaction.Fee.LessThan(global.Current.MinSendFee) {
			log.Debug("Missing Fee", "fee", transaction.Fee)
			return status.MISSING_DATA
		}
	}

	if transaction.Function == EXECUTE {
		if transaction.Data == nil {
			log.Debug("Smart Contract Missing Data", "transaction data", transaction)
			return status.MISSING_DATA
		}

		executeData := transaction.Data.(Execute)
		if executeData.Name == "" {
			log.Debug("Smart Contract Missing Data", "name", executeData.Name)
			return status.MISSING_DATA
		}

		if executeData.Version.String() == "" {
			log.Debug("Smart Contract Missing Data", "version", executeData.Version)
			return status.MISSING_DATA
		}

		if !transaction.Fee.IsCurrency("OLT") {
			log.Debug("Wrong Fee token", "fee", transaction.Fee)
			return status.INVALID
		}

		if transaction.Fee.LessThan(global.Current.MinSendFee) {
			log.Debug("Missing Fee", "fee", transaction.Fee)
			return status.MISSING_DATA
		}
	}

	return status.SUCCESS
}

func (transaction *Contract) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Smart Contract Transaction for CheckTx")

	return status.SUCCESS
}

func (transaction *Contract) ShouldProcess(app interface{}) bool {
	return true
}

func Convert(installData Install) (id.AccountKey, string, version.Version, data.Script) {
	owner := installData.Owner
	name := installData.Name
	version := installData.Version
	script := data.Script{
		Script: installData.Script,
	}
	return owner, name, version, script
}

func (transaction *Contract) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Smart Contract Transaction for DeliverTx")

	switch transaction.Function {
	case INSTALL:
		transaction.Install(app)

	case EXECUTE:
		//TODO: Need to fix this after
		go transaction.Execute(app)

	case COMPARE:
		status := transaction.Compare(app)
		return status

	default:
		return status.INVALID
	}

	return status.SUCCESS
}

func (transaction *Contract) Resolve(app interface{}) Commands {
	return []Command{}
}

//Install store the script in the SmartContract database
func (transaction *Contract) Install(app interface{}) {

	installData := transaction.Data.(Install)
	owner, name, version, script := Convert(installData)

	smartContracts := GetSmartContracts(app)
	var scriptRecords *data.ScriptRecords

	raw := smartContracts.Get(owner)
	if raw == nil {
		scriptRecords = data.NewScriptRecords()
	} else {
		scriptRecords = raw.(*data.ScriptRecords)
	}

	scriptRecords.Set(name, version, script)
	session := smartContracts.Begin()
	session.Set(owner, scriptRecords)
	session.Commit()
}

//Execute find the selected validator - who runs the script and keeps the result,
//then creates a Compare transaction then broadcast, get the results, check the results if they're the same
func (transaction *Contract) Execute(app interface{}) Transaction {
	validatorList := id.GetValidators(app)
	selectedValidatorIdentity := validatorList.SelectedValidator

	//if global.Current.NodeName == selectedValidatorIdentity.NodeName {
	accounts := GetAccounts(app)
	nodeAccount, err := accounts.FindName(global.Current.NodeAccountName)
	if err != status.SUCCESS {
		log.Warn("Missing NodeAccount for Contracts", "name", global.Current.NodeAccountName)
		return nil
	}
	if bytes.Compare(nodeAccount.AccountKey(), selectedValidatorIdentity.AccountKey) == 0 {
		executeData := transaction.Data.(Execute)
    scriptBytes, err := getSmartContractCode(app, executeData.Address, executeData.Owner, executeData.Name, executeData.Version)
    if err != nil {
      log.Warn("Error when try to fetch the smart contract code", "err", err)
      return nil
    }
    last := GetContext(app, executeData.Owner, executeData.Name, executeData.Version)
    request := NewOLVMResultWithCallString(scriptBytes, executeData.CallString, last.Context)
    result := RunScript(app, request).(OLVMResult)

    if result.Status == status.SUCCESS {
      resultCompare := transaction.CreateCompareRequest(app, executeData.Owner, executeData.Name, executeData.Version, result)
      if resultCompare != nil {
        //TODO: check this later
        comm.Broadcast(resultCompare)
      }
		}
	}
	return nil
}

//Execute calls this and compare the results
// TODO: Maybe the compare should be in CheckTx, so that it can fail?
func (transaction *Contract) Compare(app interface{}) status.Code {
	compareData := transaction.Data.(Compare)
  scriptBytes, err := getSmartContractCode(app, compareData.Address, compareData.Owner, compareData.Name, compareData.Version)
  if err != nil {
    log.Warn("Error when try to fetch the smart contract code", "err", err)
    return status.INVALID
  }

  last := GetContext(app, compareData.Owner, compareData.Name, compareData.Version)
  request := NewOLVMResultWithCallString(scriptBytes,compareData.CallString, last.Context)
  result := RunScript(app, request).(OLVMResult)

  // TODO: Comparison should be on the structure, not a string
  if CompareResults(result, compareData.Result) {
    SaveContext(app, compareData.Owner, compareData.Name, compareData.Version, result)

    return status.SUCCESS
  }

	return status.INVALID
}

func GetContext(app interface{}, owner id.AccountKey, name string, version version.Version) OLVMResult {
	context := GetExecutionContext(app)
	raw := context.Get(owner)
	if raw == nil {
		// TODO: Should be a NewOLVMResult to initialize properly
		return OLVMResult{}
	}
	mmap := raw.(*data.MultiMap)
	entry := mmap.Get(name, version)
	if entry.Value == nil {
		return OLVMResult{}
	}
	return entry.Value.(OLVMResult)
}

func SaveContext(app interface{}, owner id.AccountKey, name string, version version.Version, result OLVMResult) {
	context := GetExecutionContext(app)

	var mmap *data.MultiMap
	raw := context.Get(owner)
	if raw == nil {
		mmap = data.NewMultiMap()
	} else {
		mmap = raw.(*data.MultiMap)
	}
	mmap.Set(name, version, result)

	session := context.Begin()
	session.Set(owner, mmap)

	// TODO: Wrong, should only commit on block commit, but it is this way now
	// so that inter-block executions work correctly
	session.Commit()
}

func CompareResults(recent OLVMResult, original OLVMResult) bool {
	return true
}

// Create a comparison request
func (transaction *Contract) CreateCompareRequest(app interface{}, owner id.AccountKey, name string,
	version version.Version, result OLVMResult) []byte {

	chainId := GetChainID(app)

	// Costs have been taken care of already
	fee := data.NewCoinFromInt(0, "OLT")
	gas := data.NewCoinFromInt(0, "OLT")

	next := id.NextSequence(app, transaction.Owner)

	inputs := Compare{
		Owner:   owner,
		Name:    name,
		Version: version,
		Result:  result,
	}

	account := GetNodeAccount(app)

	// Create base transaction
	compare := &Contract{
		Base: Base{
			Type:     SMART_CONTRACT,
			ChainId:  chainId,
			Owner:    account.AccountKey(),
			Signers:  GetSigners(account.AccountKey()),
			Sequence: next.Sequence,
		},
		Data:     inputs,
		Function: COMPARE,
		Fee:      fee,
		Gas:      gas,
	}

	signed := SignAndPack(Transaction(compare))
	return signed
}
