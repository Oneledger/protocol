package main

import (
	"log"

	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
)

func main() {
	log.Print("Up running vm")
	vm.RunService()
}
