/*
	Copyright 2017-2018 OneLedger

	Handle Smart Contarct Install and Execute functions
*/
package action

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/tendermint/tendermint/libs/common"

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
	GetReference() []byte
}

type Install struct {
	//ToDo: all the data you need to install contract
	Owner   id.AccountKey
	Name    string
	Version version.Version
	Script  []byte
	UUID    uuid.UUID //This will be part of the contract address stored into the database
}

func (i Install) GetReference() []byte {
	return _hashToStringBytes(i)
}

type Execute struct {
	//ToDo: all the data you need to execute contract
	Owner      id.AccountKey
	Name       string
	Address    string
	CallString string
	Version    version.Version
	UUID       uuid.UUID //This will be the UUID for one execution
}

func (e Execute) GetReference() []byte {
	return _hashToStringBytes(e)
}

type Compare struct {
	//ToDo: all the data you need to compare contract
	Owner      id.AccountKey
	Name       string
	Address    string
	CallString string
	Version    version.Version
	Result     OLVMResult
	Reference  []byte
}

type PersistContractData = data.ContractData

type ContractReturnData struct {
	Status    status.Code
	Out       string
	Ret       string
	Reference []byte
}

func (crd ContractReturnData) GetReference() []byte {
	return crd.Reference
}

func (c Compare) GetReference() []byte {
	return c.Reference
}

func init() {
	serial.Register(Contract{})
	serial.Register(Install{})
	serial.Register(Execute{})
	serial.Register(Compare{})
	serial.Register(PersistContractData{})

	var prototype ContractData
	serial.RegisterInterface(&prototype)

  serial.Register(uuid.UUID{})

}

func (transaction *Contract) GetData() interface{} {
	return transaction.Data
}

func (transaction *Contract) SetData(data interface{}) bool {
	transaction.Data = data.(ContractData)
	return true
}

func (transaction Contract) TransactionTags(app interface{}) Tags {
	tags := transaction.Base.TransactionTags(app)

	tagReference := transaction.Data.GetReference()
	tag1 := common.KVPair{
		Key:   []byte("tx.contract"),
		Value: []byte(tagReference),
	}

	tags = append(tags, tag1)

	return tags
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
  if transaction.Function == EXECUTE {
    return PreRunSmartContract(app, transaction)
  }
	return status.SUCCESS
}

func (transaction *Contract) ShouldProcess(app interface{}) bool {
	return true
}

func Convert(installData Install) ([]byte, string, version.Version, data.Script) {
	ref := installData.GetReference()
	name := installData.Name
	version := installData.Version
	script := data.Script{
		Script: installData.Script,
	}
	return ref, name, version, script
}

func (transaction *Contract) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Smart Contract Transaction for DeliverTx")

	switch transaction.Function {
	case INSTALL:
		transaction.Install(app)

	case EXECUTE:

		return transaction.Execute(app)

	case COMPARE:
		return transaction.Compare(app)

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
	ref := SaveSmartContractCode(app, installData)
  transaction.SetData(ContractReturnData{Status: status.SUCCESS, Out: string(ref), Reference:ref}) //Send the reference back to the commandline
}

//Execute find the selected validator - who runs the script and keeps the result,
//then creates a Compare transaction then broadcast, get the results, check the results if they're the same
func (transaction *Contract) Execute(app interface{}) status.Code {
	validatorList := id.GetValidators(app)
	selectedValidatorIdentity := validatorList.SelectedValidator

	//if global.Current.NodeName == selectedValidatorIdentity.NodeName {
	accounts := GetAccounts(app)
	nodeAccount, error := accounts.FindName(global.Current.NodeAccountName)
	if error != status.SUCCESS {
		log.Warn("Missing NodeAccount for Contracts", "name", global.Current.NodeAccountName)
		return status.INVALID_SIGNATURE
	}
	//analyze the smart contract code
	executeData := transaction.Data.(Execute)
	script, err := GetSmartContractCode(app, []byte(executeData.Address))
  if err != nil {
		log.Warn("Error when try to fetch the smart contract code", "err", err)
		return status.MISSING_DATA
	}
  scriptBytes := script.Script
	last := GetContext(app, executeData.Address)
	request := NewOLVMResultWithCallString(scriptBytes, executeData.CallString, last)
	analyzeResult := AnalyzeScript(app, request).(OLVMResult)
	if analyzeResult.Out == "__NO_ACTION__" {
    //The __NO_ACTION__ execution should not be looped here
		return status.EXECUTE_ERROR
	}
  ref := executeData.GetReference()
  transaction.SetData(ContractReturnData{Status: status.SUCCESS, Out: string(ref), Reference:ref})
  SaveContractRef(app, ref, data.PENDING)

	if bytes.Compare(nodeAccount.AccountKey(), selectedValidatorIdentity.AccountKey) == 0 {
		ret, err := RunScript(app, request)
    if err != nil {
      return status.EXECUTE_ERROR
    }
    result := ret.(OLVMResult)
    if result.Status == status.SUCCESS {
        go func() {
          resultCompare := transaction.CreateCompareRequest(app, executeData.Owner, executeData.Name,
            executeData.Address, executeData.Version, executeData.GetReference(), executeData.CallString, result)
          if resultCompare != nil {
            comm.Broadcast(resultCompare)
          }
        }()
    }
	}
	return status.SUCCESS
}

//Execute calls this and compare the results
// TODO: Maybe the compare should be in CheckTx, so that it can fail?
func (transaction *Contract) Compare(app interface{}) status.Code {
	//todo: FindSaveContract(app, reference) use this to compare input, if input is different, then ignore,
	// if the input is the same, execute the olvm run, and delete the saved reference.
	compareData := transaction.Data.(Compare)
  log.Debug("Start comparation")
  ref := compareData.GetReference()
  log.Debug("reference", "ref", ref)
  contractRefStatus := GetContractRef(app,ref)
  log.Debug("contract reference status", "contractRefStatus", contractRefStatus)
  switch  contractRefStatus {
    case data.NOT_FOUND:
      log.Warn("No reference was found", "contractRefStatus", contractRefStatus)
      return status.INVALID
    case data.ERROR:
      log.Warn("Execution Error", "contractRefStatus", contractRefStatus)
      return status.INVALID
    case data.COMPLETED:
      log.Warn("Execution completed", "contractRefStatus", contractRefStatus)
      return status.INVALID
  }

  //contractRefStatus == data.PENDING
  log.Debug("Executing comparation", "contractRefStatus", contractRefStatus)
	scriptBytes, err := GetSmartContractCode(app, []byte(compareData.Address))
	if err != nil {
		log.Warn("Error when try to fetch the smart contract code", "err", err)
    SaveContractRef(app,ref, data.ERROR)
		return status.INVALID
	}

	last := GetContext(app, compareData.Address)
	request := NewOLVMResultWithCallString(scriptBytes.Script, compareData.CallString, last)
	ret, err := RunScript(app, request)
  if err != nil {
    SaveContractRef(app,ref, data.ERROR)
    return status.EXECUTE_ERROR
  }
  result := ret.(OLVMResult)
	transaction.SetData(GenerateContractReturnData(result))
	// TODO: Comparison should be on the structure, not a string
	if CompareResults(result, compareData.Result) {
		SaveContext(app, compareData.Address, []byte(result.Out))
    SaveContractRef(app,ref, data.COMPLETED)
		return status.SUCCESS
	}
  SaveContractRef(app, ref, data.ERROR)
	return status.INVALID
}

func GetContext(app interface{}, address string) OLVMContext {
  log.Debug("GetContext")
	var olvmContext OLVMContext
	context := GetExecutionContext(app)
	raw := context.Get(data.DatabaseKey(address))
	if raw == nil {
		return olvmContext
	}
	contractData := raw.(*PersistContractData)
	olvmContext.Data = contractData.Data
	return olvmContext
}

func GenerateContractReturnData(data interface{}) ContractReturnData{
  contractData := data.(OLVMResult)
  return ContractReturnData{contractData.Status, contractData.Out, contractData.Ret, contractData.Reference}
}


func SaveContext(app interface{}, address string, resultOut []byte) {
  log.Debug("Save Context Data")
  addressBytes := []byte(address)
	context := GetExecutionContext(app)
	var contractData *data.ContractData
	raw := context.Get(data.DatabaseKey(address))
	if raw == nil {
		log.Debug("No data create new data into the database")
		contractData = data.NewContractData(addressBytes)
	} else {
		contractData = raw.(*PersistContractData)
	}
	contractData.UpdateByJSONData(resultOut)
	session := context.Begin()
	session.Set(data.DatabaseKey(address), contractData)
	session.Commit()
}

func CompareResults(recent OLVMResult, original OLVMResult) bool {
	return true
}

// Create a comparison request

func (transaction *Contract) CreateCompareRequest(app interface{}, owner id.AccountKey, name string, address string,
	version version.Version, reference []byte, callString string, result OLVMResult) []byte {
	chainId := GetChainID(app)

	// Costs have been taken care of already
	fee := data.NewCoinFromInt(0, "OLT")
	gas := data.NewCoinFromInt(0, "OLT")

	next := id.NextSequence(app, transaction.Owner)

	inputs := Compare{
		Owner:      owner,
		Name:       name,
		Address:    address,
		CallString: callString,
		Version:    version,
		Result:     result,
		Reference:  reference,
	}

	account := GetNodeAccount(app)

	// Create base transaction
	log.Debug("create compare transaction")
	log.Debug("accountkey", "account.AccountKey", account.AccountKey())
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
	log.Debug("end create compare transaction")
	signed := SignAndPack(Transaction(compare))
	log.Debug("get transaction signed")
	return signed
}

func PreRunSmartContract(app interface{}, transaction *Contract) status.Code {
  executeData := transaction.Data.(Execute)
	script, err := GetSmartContractCode(app, []byte(executeData.Address))
  if err != nil {
		log.Warn("Error when try to fetch the smart contract code", "err", err)
		return status.MISSING_DATA
	}
  scriptBytes := script.Script
	last := GetContext(app, executeData.Address)
	request := NewOLVMResultWithCallString(scriptBytes, executeData.CallString, last)
	analyzeResult := AnalyzeScript(app, request).(OLVMResult)
	if analyzeResult.Out == "__NO_ACTION__" {
		result, err := RunScript(app, request)
		transaction.SetData(GenerateContractReturnData(result))
    if err == nil {
      return status.QUERY_COMPLETED
    } else {
      return status.EXECUTE_ERROR
    }
	}
  return status.SUCCESS
}
