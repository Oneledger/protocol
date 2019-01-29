/*
	Copyright 2017-2018 OneLedger
*/
package app

import (
	"io/ioutil"
	"os"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
)

func GetSourceCode(name string) []byte {

	path := os.Getenv("OLROOT") + "/protocol/node/olvm/interpreter/samples"

	var filePath string

	// TODO: Just a few hardcoded examples
	switch name {
	case "deadloop":
		filePath = path + "/deadloop.js"
	case "hello":
		filePath = path + "/helloworld.js"
	default:
		filePath = path + "/helloworld.js"
	}

	log.Debug("Loading Contract Script", "filePath", filePath)

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Can get Source File", "err", err)
	}

	return data
}

func StartOLVM() {
	log.Debug("Starting up Smart Contract Engine")
	vm.InitializeClient()
}

/*
func NewOLVMRequest(script []byte) {
	request := &action.OLVMRequest{
		From:       "0x0",
		Address:    "embed://",
		CallString: "",
		Value:      0,
		SourceCode: script,
	}
	return request
}
*/

func (app Application) RunScript(request interface{}) interface{} {
	return RunTestScript(request.(*action.OLVMRequest))
}

func (app Application) AnalyzeScript(request interface{}) interface{} {
	return RunAnalyze(request.(*action.OLVMRequest))
}

func RunTestScriptName(name string) interface{} {
	request := &action.OLVMRequest{
		From:       "0x0",
		Address:    "embed://",
		CallString: "",
		Value:      0,
		SourceCode: GetSourceCode(name),
	}
	return RunTestScript(request)
}

func RunAnalyze(request *action.OLVMRequest)  interface{} {
  reply, err := vm.Analyze(request)
	if err != nil {
		log.Warn("Contract Engine Failed to Start", "err", err)
	}
	log.Dump("Engine output", reply)
	return *reply
}

// Take the engine for a test spin
func RunTestScript(request *action.OLVMRequest) interface{} {
	log.Dump("Engine input", request)
	reply, err := vm.AutoRun(request)
	if err != nil {
		log.Warn("Contract Engine Failed to Start", "err", err)
	}
	log.Dump("Engine output", reply)
	return *reply
}
