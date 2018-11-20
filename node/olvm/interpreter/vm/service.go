package vm

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"runtime/debug"
	"strconv"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/monitor"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
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
			log.Error("OLVM Panicked", "status", r)
			log.Dump("Details", "args", args, "reply", reply)
			debug.PrintStack()
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
	log.Debug("Service is running", "protocol", ol.Protocol, "port", ol.Port)

	container := new(Container)
	rpc.Register(container)
	rpc.HandleHTTP()

	listen, err := net.Listen(ol.Protocol, ":"+strconv.Itoa(ol.Port))
	if err != nil {
		log.Fatal("listen error:", "err", err)
	}

	http.Serve(listen, nil)
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
