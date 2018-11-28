package runner

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/robertkrimen/otto"
)

type Runner struct {
	vm *otto.Otto
}

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
	Transaction action.Transaction
	Context     OLVMContext

	// TODO: Data Handle (some way to call out for large data requests)
}

// All of the output received from the computation
type OLVMResult struct {
	Out     string
	Ret     string // TODO: Should be a real name
	Elapsed string

	Transactions []action.Transaction
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
