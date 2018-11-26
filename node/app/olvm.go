/*
	Copyright 2017-2018 OneLedger
*/
package app

import (
	"io/ioutil"
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
	"github.com/spf13/viper"
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

func RunVM() {
	log.Debug("Starting up Smart Contract Engine")

	// TODO: Temporary until fullnodes correctly integrate with viper
	global.Current.OLVMAddress = viper.Get("OLVMAddress").(string)
	global.Current.OLVMProtocol = viper.Get("OLVMProtocol").(string)

	//vm.InitializeService()
	vm.InitializeClient()

	request := &runner.OLVMRequest{
		From:       "0x0",
		Address:    "embed://",
		CallString: "",
		Value:      0,
		SourceCode: GetSourceCode(),
	}

	// TODO: Take the engine for a test spin

	log.Dump("Engine input", request)
	reply, err := vm.AutoRun(request)
	if err != nil {
		log.Warn("Contract Engine Failed to Start", "err", err)
	}

	log.Dump("Engine output", reply)
}
