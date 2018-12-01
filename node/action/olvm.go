/*
	Copyright 2017-2018 OneLedger
*/

package action

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

type OLVMContext struct {
	Data interface{}
}

// All of the input necessary to perform a computation on a transaction
type OLVMRequest struct {
	From        string
	Address     string
	CallString  string
	Value       int
	SourceCode  string
	Transaction Transaction
	Context     OLVMContext

	// TODO: Data Handle (some way to call out for large data requests)
}

// All of the output received from the computation
type OLVMResult struct {
	Status  status.Code
	Out     string
	Ret     string // TODO: Should be a real name
	Elapsed string

	Transactions []Transaction
	Context      OLVMContext
}

func init() {
	serial.Register(OLVMRequest{})
	serial.Register(OLVMResult{})
	serial.Register(OLVMContext{})

	// TODO: Doesn't work in serial?
	//var prototype time.Time
	//serial.Register(prototype)
	//var prototype2 time.Duration
	//serial.Register(prototype2)
}
