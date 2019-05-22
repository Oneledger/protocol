/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package data

import (
	"github.com/Oneledger/protocol/serialize"
)

const (
	TagOLVMResult  = "action_olvm_result"
	TagOLVMRequest = "action_olvm_request"
	TagOLVMContext = "action_olvm_context"
)

type TransactionData []byte

type OLVMContext struct {
	Data map[string]interface{}
}

func (context OLVMContext) GetValue(key string) interface{} {

	ret := context.Data[key]
	if ret == nil {
		return nil
	}

	return string(ret.([]byte))
}

// All of the input necessary to perform a computation on a transaction
type OLVMRequest struct {
	From        string
	Address     string
	CallString  string
	Value       int
	SourceCode  []byte
	Reference   []byte
	Transaction TransactionData
	Context     OLVMContext
	// TODO: Data Handle (some way to call out for large data requests)
}

// All of the output received from the computation
type OLVMResult struct {
	Status       string
	Out          string
	Ret          string // TODO: Should be a real name
	Elapsed      string
	Reference    []byte
	Transactions []TransactionData
	Context      OLVMContext
}

func init() {

	serialize.RegisterConcrete(new(OLVMRequest), TagOLVMRequest)
	serialize.RegisterConcrete(new(OLVMResult), TagOLVMResult)
	serialize.RegisterConcrete(new(OLVMContext), TagOLVMContext)
}

func NewOLVMResultWithCallString(script []byte, callString string, context OLVMContext) *OLVMRequest {
	request := &OLVMRequest{
		From:       "0x0",
		Address:    "embed://",
		CallString: callString,
		Value:      0,
		SourceCode: script,
		Context:    context,
	}
	return request
}

func NewOLVMRequest(script []byte, context OLVMContext) *OLVMRequest {
	return NewOLVMResultWithCallString(script, "", context)
}

func NewOLVMResult() *OLVMResult {
	result := &OLVMResult{
		Status: "MISSING DATA",
	}
	return result
}
