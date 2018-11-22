package main

import (
	"github.com/Oneledger/protocol/node/log"

	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
)

func main() {
	request := runner.OLVMRequest{"0x0","samples://helloworld", "",0, ""}
	reply, err := vm.AutoRun(&request)
	if err != nil {
		log.Fatal("Failed",err)
	}
	log.Info("get the result","reply",reply)
}
