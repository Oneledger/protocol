package main

import (
	"github.com/Oneledger/protocol/olvm/interpreter/vm"
)

func main() {
	log.Info("Up running vm")
	vm.RunService()
}
