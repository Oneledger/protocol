package main

import (
	"log"

	"github.com/Oneledger/protocol/node/olvm/interpreter/vm"
)

func main() {
	reply, err := vm.AutoRun("0x0", "samples://helloworld", "", "", 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply)
}
