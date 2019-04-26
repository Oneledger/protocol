package main

import (
	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
)

func main() {
	log.Info("Up running vm")
	vm.RunService()
}
