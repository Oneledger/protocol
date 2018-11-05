package main

import (
	"./committor"
	"./monitor"
	"./runner"
	"flag"
	"log"
	"os"
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

func main() {
	log.Println("starting OVM")

	address := flag.String("address", "samples://helloworld", "the address of your smart contract")

	call_transaction := flag.String("transaction", "default__('hello,world from Oneledger')", "the transaction on that address, either string or encrypted bytes")

	call_method_type := flag.String("type", "plain", "plain text call or encrypted call. [plain|encrypted]")

	call_from := flag.String("from", "0x0", "the public address of the caller")

	call_value := flag.Int("value", 0, "number of OLT put on this call")

	flag.Parse()

	log.Printf("\nfrom:\t%s\naddress:\t%s\ntransaction:\t%s\ntype:\t%s\nvalue:\t%x\n",
		*call_from,
		*address,
		*call_transaction,
		*call_method_type,
		*call_value)

	transaction_ch := make(chan string)
	returnValue_ch := make(chan string)
	status_ch := make(chan monitor.Status)
	monitor := monitor.CreateMonitor(10, monitor.DEFAULT_MODE, "./ovm.pid")

	status, err := monitor.CheckUnique()

	defer func() {
		r := recover()
		if r != nil {
			log.Println(status)
			os.Remove(monitor.GetPidFilePath())
			log.Fatal(r)
		} else {
			os.Remove(monitor.GetPidFilePath())
		}

	}()

	if err == true {
		panic(status)
	} else {
		log.Println("VM Initialized finished, with status:", status.Details, ",  and code:", status.Code)
	}

	os.Create(monitor.GetPidFilePath())

	go monitor.CheckStatus(status_ch)
	go run(transaction_ch, returnValue_ch, status_ch, *call_from, *address, *call_transaction, *call_value)
	ready := 0
	var transaction string
	var returnValue string
	for {
		select {
		case transaction = <-transaction_ch:
			ready++
		case returnValue = <-returnValue_ch:
			ready++
		case status := <-status_ch:
			log.Println("retuning: ", status.Details, "with code", status.Code)
			panic("exit with error")
		}
		if ready == 2 {
			commit(returnValue, transaction)
			log.Println("ending OVM")
			return
		}
	}

}
