package main

import (
	"./vm"
	"log"
)

func main() {
	reply, err := vm.AutoRun("0x0", "samples://helloworld", "", "", 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply)
}
