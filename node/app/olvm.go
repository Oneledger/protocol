/*
	Copyright 2017-2018 OneLedger
*/
package app

import (
	"io/ioutil"
	"os"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
)

func GetSourceCode(name string) string {

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

	return string(data)
}

func StartVM() {
	log.Debug("Starting up Smart Contract Engine")
	vm.InitializeClient()
}

func (app Application) RunScript(script string) interface{} {
	return RunTestScript(script)
}

func RunTestScriptName(name string) interface{} {
	return RunTestScript(GetSourceCode(name))
}

// Take the engine for a test spin
func RunTestScript(script string) interface{} {
	log.Debug("########### TESTING OLVM EXECUTION ###########")

	request := &runner.OLVMRequest{
		From:       "0x0",
		Address:    "embed://",
		CallString: "",
		Value:      0,
		SourceCode: script,
	}

	log.Dump("Engine input", request)

	reply, err := vm.AutoRun(request)
	if err != nil {
		log.Warn("Contract Engine Failed to Start", "err", err)
	}

	log.Dump("Engine output", reply)

	return reply.Out
}
