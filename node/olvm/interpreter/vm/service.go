package vm

import (
	"../monitor"
	"../runner"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

var DefaultOLVMService = NewOLVMService("tcp", 1980)

func (c *Container) Echo(args *Args, reply *Reply) error {
	return nil
}

func (c *Container) Exec(args *Args, reply *Reply) error {
	mo := monitor.CreateMonitor(10, monitor.DEFAULT_MODE, "./ovm.pid")
	status_ch := make(chan monitor.Status)
	runner_ch := make(chan bool)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("error happens %v\n", r)
			os.Exit(1)
		}
	}()

	go mo.CheckStatus(status_ch)
	go func() {
		runner := runner.CreateRunner()
		from, address, callString, value := parseArgs(args)
		out, ret, error := runner.Call(from, address, callString, value)
		if error != nil {
			status_ch <- monitor.Status{"Runtime error", monitor.STATUS_ERROR}
		} else {
			reply.Out = out
			reply.Ret = ret
			runner_ch <- true
		}

	}()
	for {
		select {
		case <-runner_ch:
			return nil
		case status := <-status_ch:
			panic(status)
			return errors.New(fmt.Sprintf("%s : %d", status.Details, status.Code))
		}
	}

	return nil
}

func (ol OLVMService) Run() {
	log.Printf("Service is running with protocol %s on the port %d...\n", ol.Protocol, ol.Port)
	container := new(Container)
	rpc.Register(container)
	rpc.HandleHTTP()
	l, e := net.Listen(ol.Protocol, ":"+strconv.Itoa(ol.Port))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

func NewOLVMService(protocol string, port int) OLVMService {
	return OLVMService{protocol, port}
}

func RunService() {
	DefaultOLVMService.Run()
}

func parseArgs(args *Args) (string, string, string, int) {
	from := args.From
	address := args.Address
	callString := args.CallString
	value := args.Value
	if from == "" {
		from = "0x0"
	}
	if address == "" {
		address = "samples://helloworld"
	}
	if callString == "" {
		callString = "default__('hello,world from Oneledger')"
	}
	return from, address, callString, value
}
