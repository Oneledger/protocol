/*
	Copyright 2017-2018 OneLedger

	Declare basic types used by the Application

	If a type requires functions or a few types are intertwinded, then should be in their own file.
*/
package action

import (
	"bytes"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/serialize"
	"github.com/pkg/errors"

	"strconv"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

func init() {
	serial.Register(Event{})
	serial.Register(data.Script{})

	clSerializer = serialize.GetSerializer(serialize.CLIENT)
}

var clSerializer serialize.Serializer

/*
func NewSendInput(accountKey id.AccountKey, amount data.Coin) SendInput {

	if bytes.Equal(accountKey, []byte("")) {
		log.Fatal("Missing AccountKey", "key", accountKey, "amount", amount)
		// TODO: Error handling should be better
		return SendInput{}
	}

	if !amount.IsValid() {
		log.Fatal("Missing Amount", "key", accountKey, "amount", amount)
	}

	return SendInput{
		AccountKey: accountKey,
		Amount:     amount,
	}
}
*/

/*
// outputs for a send transaction (similar to Bitcoin)
type SendOutput struct {
	AccountKey id.AccountKey `json:"account_key"`
	Amount     data.Coin     `json:"coin"`
}
*/

/*
func NewSendOutput(accountKey id.AccountKey, amount data.Coin) SendOutput {

	if bytes.Equal(accountKey, []byte("")) {
		log.Fatal("Missing AccountKey", "key", accountKey, "amount", amount)
		// TODO: Error handling should be better
		return SendOutput{}
	}

	if !amount.IsValid() {
		log.Fatal("Missing Amount", "key", accountKey, "amount", amount)
	}

	return SendOutput{
		AccountKey: accountKey,
		Amount:     amount,
	}
}
*/

/*
func CheckBalance(app interface{}, accountKey id.AccountKey, amount data.Coin) bool {
	balances := GetBalances(app)

	balance := balances.Get(accountKey)
	if balance == nil {
		// New accounts don't have a balance until the first transaction
		log.Debug("New Balance", "key", accountKey, "amount", amount, "balance", balance)
		interim := data.NewBalanceFromInt(0, amount.Currency.Name)
		balance = interim
		if !balance.GetAmountByName(amount.Currency.Name).Equals(amount) {
			return false
		}
		return true
	}

	if !balance.GetAmountByName(amount.Currency.Name).Equals(amount) {
		log.Warn("Balance Mismatch", "key", accountKey, "amount", amount, "balance", balance)
		return false
	}
	return true
}
*/

func GetHeight(app interface{}) int64 {
	balances := GetBalances(app)

	height := int64(balances.Version)
	return height
}

func CheckAmounts(app interface{}, inputs []SendInput, outputs []SendOutput) bool {
	total := data.NewCoinFromInt(0, "OLT")
	for _, input := range inputs {
		if input.Amount.LessThan(0) {
			log.Debug("FAILED: Less Than 0", "input", input)
			return false
		}

		if !input.Amount.IsCurrency("OLT") {
			log.Debug("FAILED: Send on Currency isn't implement yet")
			return false
		}

		if bytes.Compare(input.AccountKey, []byte("")) == 0 {
			log.Debug("FAILED: Key is Empty", "input", input)
			return false
		}
		if !CheckBalance(app, input.AccountKey, input.Amount) {
			log.Warn("Balance is incorrect", "input", input)
			//return false
		}
		total.Plus(input.Amount)
	}
	for _, output := range outputs {

		if output.Amount.LessThan(0) {
			log.Debug("FAILED: Less Than 0", "output", output)
			return false
		}

		if !output.Amount.IsCurrency("OLT") {
			log.Debug("FAILED: Send on Currency isn't implement yet")
			return false
		}

		if bytes.Compare(output.AccountKey, []byte("")) == 0 {
			log.Debug("FAILED: Key is Empty", "output", output)
			return false
		}
		total.Minus(output.Amount)
	}
	if !total.Equals(data.NewCoinFromInt(0, "OLT")) {
		log.Debug("FAILED: Doesn't add up", "inputs", inputs, "outputs", outputs)
		return false
	}
	return true
}

type Event struct {
	Type        Type   `json:"type"`
	SwapKeyHash []byte `json:"swapkeyhash"`
	Step        int    `json:"step"`
}

func (e Event) ToKey() []byte {

	buffer, err := clSerializer.Serialize(e)
	if err != nil {
		log.Error("Failed to Serialize event key")
	}
	return buffer
}

func SaveEvent(app interface{}, eventKey Event, status bool) {
	events := GetEvent(app)
	s := "0"
	if status {
		s = "1"
	}
	log.Debug("Save Event", "key", eventKey)

	session := events.Begin()
	session.Set(eventKey.ToKey(), []byte(s))
	session.Commit()
}

func FindEvent(app interface{}, eventKey Event) bool {
	log.Debug("Load Event", "key", eventKey)
	events := GetEvent(app)
	result := events.Get(eventKey.ToKey())
	if result == nil {
		return false
	}

	if bytes.Equal(result.([]byte), []byte("1")) {
		return true
	}

	return false
}

func SaveContract(app interface{}, contractKey []byte, nonce int64, contract []byte) {
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	log.Debug("Save contract", "key", contractKey, "afterNonce", n)
	session := contracts.Begin()
	session.Set(n, contract)
	session.Commit()
}

func FindContract(app interface{}, contractKey []byte, nonce int64) []byte {
	log.Debug("Load Contract", "key", contractKey)
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	result := contracts.Get(n)
	if result == nil {
		return nil
	}
	return result.([]byte)
}

func DeleteContract(app interface{}, contractKey []byte, nonce int64) {
	contracts := GetContracts(app)
	n := strconv.AppendInt(contractKey, nonce, 10)
	log.Debug("Delete contract", "key", contractKey, "afterNonce", n)
	session := contracts.Begin()
	session.Delete(n)
	session.Commit()
}

func SaveContractRef(app interface{}, key []byte, contractRefStatus data.ContractRefStatus) {
	datastore := GetContracts(app)
	session := datastore.Begin()
	session.Set(key, data.ContractRefRecord{contractRefStatus})
	session.Commit()
}

func GetContractRef(app interface{}, key []byte) data.ContractRefStatus {
	datastore := GetContracts(app)
	log.Debug("getdatastore")
	raw := datastore.Get(key)
	log.Debug("get raw", "raw", raw)
	if raw == nil {
		return data.NOT_FOUND
	} else {
		contractRefRecord := raw.(data.ContractRefRecord)
		log.Debug("contractRefRecord", "contractRefRecord", contractRefRecord)
		log.Debug("status", "contractRefRecord.Status", contractRefRecord.Status)
		return contractRefRecord.Status
	}
}

func SaveSmartContractCode(app interface{}, installData Install) []byte {
	ref, name, version, script := Convert(installData)
	smartContracts := GetSmartContracts(app)
	var scriptRecords *data.ScriptRecords

	raw := smartContracts.Get(ref)
	if raw == nil {
		scriptRecords = data.NewScriptRecords()
	} else {
		scriptRecords = raw.(*data.ScriptRecords)
	}

	scriptRecords.Set(name, version, script)
	session := smartContracts.Begin()
	session.Set(ref, scriptRecords)
	session.Commit()
	log.Debug("save the smart contract with key", "key", string(ref))
	return ref
}

func GetSmartContractCode(app interface{}, ref []byte) (script data.Script, err error) {
	log.Debug("load the smart contract with key", "key", ref)
	smartContracts := GetSmartContracts(app)
	raw := smartContracts.Get(ref)
	if raw == nil {
		err = errors.New("Not code found")
		return
	}
	scriptRecords := raw.(*data.ScriptRecords)
	script = scriptRecords.Script
	return
}

func FindContractCode(hash []byte, proof bool) (script data.Script, err error) {
	result := comm.Tx(hash, false)
	if result == nil {
		err = errors.New("Failed to find the result from hash")
		return
	}
	var signedTx SignedTransaction
	out, err := serial.Deserialize(result.Tx, signedTx, serial.CLIENT)
	if err != nil {
		err = errors.New("Unable to serialize transaction")
		return
	}
	// This includes the signature
	// tx := out.(SignedTransaction)
	// This returns the wrapped Transaction
	tx := out.(SignedTransaction).Transaction

	if tx.GetType() != SMART_CONTRACT {
		err = errors.New("The transaction is not for the Contract Type")
		return
	}
	contractTx := tx.(*Contract)

	if contractTx.Function != INSTALL {
		err = errors.New("The contract doesn't have scripts")
		return
	}

	_, _, _, script = Convert(contractTx.Data.(Install))

	return script, nil

}
