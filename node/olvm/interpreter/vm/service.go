package vm

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"runtime/debug"
	//"time"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/monitor"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
)

var DefaultOLVMService = NewOLVMService("tcp", ":1980")

func InitializeService() {
	protocol := global.Current.OLVMProtocol
	address := global.Current.OLVMAddress
	DefaultOLVMService = NewOLVMService(protocol, address)
}

func (c *Container) Echo(request *runner.OLVMRequest, result *runner.OLVMResult) error {
	// TODO: Do something useful here
	return nil
}

func (c *Container) Exec(request *runner.OLVMRequest, result *runner.OLVMResult) (err error) {
	log.Debug("Exec the Contract")

	// TODO: Isn't this just a timer?
	mo := monitor.CreateMonitor(10, monitor.DEFAULT_MODE, "./ovm.pid")

	status_ch := make(chan monitor.Status)
	runner_ch := make(chan bool)

	defer func() {
		if r := recover(); r != nil {
			go func() {
				debug.PrintStack()
				log.Dump("Details", "request", request, "result", result)
				log.Fatal("OLVM Panicked", "status", r)
			}()
		}
	}()

	go mo.CheckStatus(status_ch)
	go func() {
		runner := runner.CreateRunner()

		error := runner.Call(request, result)
		if error != nil {
			status_ch <- monitor.Status{error.Error(), monitor.STATUS_ERROR}
		} else {
			runner_ch <- true
		}
	}()

	for {
		select {
		case <-runner_ch:
			err = nil
			return
		case status := <-status_ch:
			err = errors.New(fmt.Sprintf("%s : %d", status.Details, status.Code))
			panic(status)
		}
	}
	return
}

func (ol OLVMService) Run() {
	log.Debug("Service is running", "protocol", ol.Protocol, "address", ol.Address)

	container := new(Container)
	rpc.Register(container)
	rpc.HandleHTTP()

	listen, err := net.Listen(ol.Protocol, ol.Address)
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
