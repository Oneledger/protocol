package main

import (
	"flag"
	"log"
	"os"
	"github.com/Oneledger/protocol/node/olvm/interpreter/committor"
	"github.com/Oneledger/protocol/node/olvm/interpreter/monitor"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
)

func run(x chan string, y chan string, status_ch chan monitor.Status, from string, address string, transaction string, olt int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			status_ch <- monitor.Status{"scripting running error", monitor.STATUS_ERROR}
		}
	}()

	runner := runner.CreateRunner()
	transaction, returnValue := runner.Call(from, address, transaction, olt)
	x <- transaction
	y <- returnValue
}

func commit(returnValue string, transaction string) {
	log.Print(returnValue)
	log.Print(transaction)
	c := committor.Create()
	c.Commit(returnValue, transaction)
}

func runAsCommand(monitor monitor.Monitor, x chan string, y chan string, status_ch chan monitor.Status, from string, address string, transaction string, olt int) {
	go monitor.CheckStatus(status_ch)
	go run(x, y, status_ch, from, address, transaction, olt)
}


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
