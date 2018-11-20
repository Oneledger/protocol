package main

import (
	"./committor"
	"./vm"
	"flag"
	"log"
)

func main() {
	log.Println("starting OVM")

	address := flag.String("address", "samples://helloworld", "the address of your smart contract")

	call_string := flag.String("call_string", "default__('hello,world from Oneledger')", "the call string on that contract address")

	call_from := flag.String("from", "0x0", "the public address of the caller")

	code := flag.String("sourceCode", "", "the source code of the smart contract(optional)")

	call_value := flag.Int("value", 0, "number of OLT put on this call")

	flag.Parse()

	log.Printf("\nfrom:\t%s\naddress:\t%s\ncall string:\t%s\ncode:\t%s\nvalue:\t%x\n",
		*call_from,
		*address,
		*call_string,
		*code,
		*call_value)

	reply, err := vm.AutoRun(*call_from, *address, *call_string, *code, *call_value)
	if err != nil {
		log.Fatal(err)
	}
	c := committor.Create()
	log.Println("return value:", reply.Ret)
	log.Println("transaction out:", reply.Out)
	s, _ := c.Commit(reply.Ret, reply.Out)
	log.Println(s)
}
