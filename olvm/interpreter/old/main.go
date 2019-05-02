package main

import (
	"flag"
	"os"

	logger "github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/olvm/interpreter/committor"
	"github.com/Oneledger/protocol/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/olvm/interpreter/vm"
)

var log = logger.NewDefaultLogger(os.Stdout).WithPrefix("olvm/interpreter/old")

func main() {
	log.Info("starting OVM")

	address := flag.String("address", "samples://helloworld", "the address of your smart contract")

	call_string := flag.String("call_string", "default__('hello,world from Oneledger')", "the call string on that contract address")

	call_from := flag.String("from", "0x0", "the public address of the caller")

	code := flag.String("sourceCode", "", "the source code of the smart contract(optional)")

	call_value := flag.Int("value", 0, "number of OLT put on this call")

	flag.Parse()

	log.Infof("\nfrom:\t%s\naddress:\t%s\ncall string:\t%s\ncode:\t%s\nvalue:\t%x\n",
		*call_from,
		*address,
		*call_string,
		*code,
		*call_value)

	reply, err := vm.AutoRun(&runner.OLVMRequest{*call_from, *address, *call_string, *call_value, *code})
	if err != nil {
		log.Fatal(err)
	}

	c := committor.Create()
	log.Info("return value:", reply.Ret)
	log.Info("transaction out:", reply.Out)

	s, err := c.Commit(reply.Ret, reply.Out)
	if err != nil {
		log.Error("error in committing transaction", err)
	}

	log.Info(s)
}
