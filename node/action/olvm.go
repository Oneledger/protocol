/*
	Copyright 2017-2018 OneLedger
*/

package action

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/serialize"
	"github.com/Oneledger/protocol/node/status"
)

type OLVMContext struct {
	Data map[string]interface{}
}

func (context OLVMContext) GetValue(key string) interface{} {
	ret := context.Data[key]
	if ret == nil {
		return nil
	} else {
		return string(ret.([]byte))
	}
}

// All of the input necessary to perform a computation on a transaction
type OLVMRequest struct {
	From        string
	Address     string
	CallString  string
	Value       int
	SourceCode  []byte
	Reference   []byte
	Transaction Transaction
	Context     OLVMContext
	// TODO: Data Handle (some way to call out for large data requests)
}

// All of the output received from the computation
type OLVMResult struct {
	Status       status.Code
	Out          string
	Ret          string // TODO: Should be a real name
	Elapsed      string
	Reference    []byte
	Transactions []Transaction
	Context      OLVMContext
}

func init() {
	serial.Register(OLVMRequest{})
	serial.Register(OLVMResult{})
	serial.Register(OLVMContext{})


	serialize.RegisterConcrete(new(OLVMResult), "action_olvmresult")
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
		Status: status.MISSING_DATA,
	}
	return result
}
