package vm

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"runtime/debug"
	//"time"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/monitor"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/node/sdk"
)

var DefaultOLVMService *OLVMService

func InitializeService() {
	log.Debug("Initialize Service")
	protocol := global.Current.OLVMProtocol
	address := global.Current.OLVMAddress
	DefaultOLVMService = NewOLVMService(protocol, address)
}

func (c *Container) Echo(request *runner.OLVMRequest, result *runner.OLVMResult) error {
	// TODO: Do something useful here
	return nil
}

func (c *Container) Exec(request *runner.OLVMRequest, result *runner.OLVMResult) (err error) {
	log.Dump("Exec a Contract", request)

	// TODO: Isn't this just a timer? There is an otto based timer as well, that is nicer
	mo := monitor.CreateMonitor(4, monitor.DEFAULT_MODE, "./ovm.pid")

	status_ch := make(chan monitor.Status)
	runner_ch := make(chan bool)

	defer func() {
		if r := recover(); r != nil {
			go func() {
				debug.PrintStack()
				log.Dump("Details", "request", request, "result", result)
				log.Warn("OLVM Panicked", "status", r)
			}()
		}
	}()

	go mo.CheckStatus(status_ch)

	go func() {
		log.Debug("Creating a Runner")
		runner := runner.CreateRunner()

		log.Debug("Calling the Runner")
		error := runner.Call(request, result)
		if error != nil {
			status_ch <- monitor.Status{error.Error(), monitor.STATUS_ERROR}
		} else {
			runner_ch <- true
		}
	}()

	for {
		log.Debug("Checking Status")
		select {

		case <-runner_ch:
			log.Debug("Finished Processing Normally")
			err = nil
			return

		case status := <-status_ch:
			err = errors.New(fmt.Sprintf("%s : %d", status.Details, status.Code))
			log.Dump("Have an error, exiting process", err, request, result)
			if result.Ret == "HALT" {
				// SOFT EXIT
				log.Debug("Picked up a Soft Exit")
				break
			} else {
				// HARD EXIT
				//panic(status) // TODO: the panic is caught elsewhere, so this doesn't work
				os.Exit(0)
			}
		}
		log.Debug("Redoing Select")
	}
	return
}

func (ol OLVMService) Run() {
	defer func() {
		if r := recover(); r != nil {
			go func() {
				debug.PrintStack()
				log.Warn("OLVM Panicked", "status", r)
			}()
		}
	}()

	log.Debug("Starting Service", "protocol", ol.Protocol, "address", ol.Address)

	container := new(Container)
	rpc.Register(container)
	rpc.HandleHTTP()

	log.Debug("Listening on the port")
	listen, err := net.Listen(ol.Protocol, ":"+sdk.GetPort(ol.Address))
	if err != nil {
		log.Fatal("listen error:", "err", err)
	}

	log.Debug("Waiting for a request")
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
