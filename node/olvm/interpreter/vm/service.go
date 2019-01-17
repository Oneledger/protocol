package vm

import (
	"net"
	"net/http"
	"net/rpc"
	"runtime/debug"
	//"time"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/node/sdk"
)

// TODO: Make sure call is not before viper args are handled.
func NewOLVMService() *OLVMService {
	return &OLVMService{
		Protocol: global.Current.OLVMProtocol,
		Address:  global.Current.OLVMAddress,
	}
}

// Start up the service
func (ol OLVMService) StartService() {
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

// Echo as defined by RPC
func (c *Container) Echo(request *action.OLVMRequest, result *action.OLVMResult) error {
	// TODO: Do something useful here
	return nil
}

// Exec as defined by RPC
func (c *Container) Exec(request *action.OLVMRequest, result *action.OLVMResult) error {
	log.Dump("Exec a Contract", request)
	runner := runner.CreateRunner()
	return runner.Call(request, result)
}
