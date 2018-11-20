package vm

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"runtime/debug"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/monitor"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/node/sdk"
)

var DefaultOLVMService *OLVMService

func InitializeService() {
	protocol := global.Current.OLVMProtocol
	address := global.Current.OLVMAddress
	DefaultOLVMService = NewOLVMService(protocol, address)
}

func (c *Container) Echo(request *OLVMRequest, result *OLVMResult) error {
	// TODO: Do something useful here
	return nil
}

func (c *Container) Exec(request *OLVMRequest, result *OLVMResult) error {

	// TODO: Isn't this just a timer?
	mo := monitor.CreateMonitor(10, monitor.DEFAULT_MODE, "./ovm.pid")

	status_ch := make(chan monitor.Status)
	runner_ch := make(chan bool)

	defer func() {
		if r := recover(); r != nil {
			log.Error("OLVM Panicked", "status", r)
			log.Dump("Details", "request", request, "result", result)
			debug.PrintStack()
		}
	}()

	go mo.CheckStatus(status_ch)
	go func() {
		runner := runner.CreateRunner()

		from, address, callString, value := parseArgs(request)

		out, ret, error := runner.Call(from, address, callString, value)
		if error != nil {
			status_ch <- monitor.Status{"Runtime error", monitor.STATUS_ERROR}
		} else {
			result.Out = out
			result.Ret = ret
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
	log.Debug("Service is running", "protocol", ol.Protocol, "address", ol.Address)

	container := new(Container)
	rpc.Register(container)
	rpc.HandleHTTP()

	listen, err := net.Listen(ol.Protocol, ":"+sdk.GetPort(ol.Address))
	if err != nil {
		log.Fatal("listen error:", "err", err)
	}

	http.Serve(listen, nil)
}

// TODO: Make sure call is not before viper args are handled.
func NewOLVMService(protocol string, address string) *OLVMService {
	return &OLVMService{
		Protocol: protocol,
		Address:  address,
	}
}

func RunService() {
	DefaultOLVMService.Run()
}

func parseArgs(request *OLVMRequest) (string, string, string, int) {

	from := request.From
	address := request.Address
	callString := request.CallString
	value := request.Value

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
