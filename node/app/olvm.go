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

func GetSourceCode() string {

	// TODO: Just a hardcoded example
	path := os.Getenv("OLROOT") + "/protocol/node/olvm/interpreter/samples"
	filePath := path + "/deadloop.js"
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

// Take the engine for a test spin
func RunTestScript() {
	log.Debug("########### TESTING OLVM EXECUTION ###########")

	request := &runner.OLVMRequest{
		From:       "0x0",
		Address:    "embed://",
		CallString: "",
		Value:      0,
		SourceCode: GetSourceCode(),
	}

	log.Dump("Engine input", request)

	reply, err := vm.AutoRun(request)
	if err != nil {
		log.Warn("Contract Engine Failed to Start", "err", err)
	}

	log.Dump("Engine output", reply)
}
