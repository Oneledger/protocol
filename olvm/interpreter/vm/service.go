package vm

import (
	"net"
	"net/http"
	"net/rpc"
	"runtime/debug"
	//"time"

	"github.com/Oneledger/protocol/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/node/sdk"
	"github.com/Oneledger/protocol/data"
)

// TODO: Make sure call is not before viper args are handled.
func NewOLVMService(protocol, address string) *OLVMService {
	return &OLVMService{
		Protocol: protocol,
		Address:  address,
	}
}

// Start up the service
func (ol OLVMService) StartService() {
	defer func() {
		if r := recover(); r != nil {
			go func() {
				debug.PrintStack()
				log.Error("OLVM Panicked", "status", r)
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
func (c *Container) Echo(request *data.OLVMRequest, result *data.OLVMResult) error {
	// TODO: Do something useful here
	return nil
}

// Exec as defined by RPC
func (c *Container) Exec(request *data.OLVMRequest, result *data.OLVMResult) error {
	log.Info("Exec a Contract", *request)
	runner := runner.CreateRunner()
	return runner.Call(request, result)
}

// Exec as defined by RPC
func (c *Container) Analyze(request *data.OLVMRequest, result *data.OLVMResult) error {
	log.Info("Analyze a Contract", *request)
	runner := runner.CreateRunner()
	return runner.Analyze(request, result)
}
