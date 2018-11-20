/*
	Copyright 2017-2018 OneLedger
*/
package vm

import (
	"net/rpc"
	"strings"
	"time"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk"
)

// TODO: Hardcoded port, needs to come from config
//var DefaultClient = NewClient("tcp", "localhost:1980")
var defaultClient *OLVMClient

func NewClient(protocol string, address string) *OLVMClient {

	return &OLVMClient{
		Protocol:    protocol,
		ServicePath: address,
	}
}

// Initialize the vm/daemon/etc.
func Initialize() {
	protocol := global.Current.OLVMProtocol
	address := global.Current.OLVMAddress

	defaultClient = NewClient(protocol, address)
}

func AutoRun(from string, address string, callString string, sourceCode string, value int) (result *OLVMResult, err error) {

	result, err = defaultClient.Run(from, address, callString, value)

	// TODO: Should be based on error code, not text...
	if err != nil && strings.HasSuffix(err.Error(), "connection refused") {
		//try to launch the service
		log.Debug("Launching OLVM")

		// TODO: Not started the first time?
		go RunService()

		for err != nil && strings.HasSuffix(err.Error(), "connection refused") {
			time.Sleep(time.Second)
			result, err = defaultClient.Run(from, address, callString, value)
		}
		return
	}
	return
}

// Run a smart contract
func (c OLVMClient) Run(from string, address string, callString string, value int) (*OLVMResult, error) {

	request := &OLVMRequest{
		From:       from,
		Address:    address,
		CallString: callString,
		Value:      value,
	}

	client, err := rpc.DialHTTP(c.Protocol, ":"+sdk.GetPort(c.ServicePath))
	if err != nil {
		return nil, err
	}

	var result OLVMResult
	// TODO: Shouldn't pass by address for the result
	err = client.Call("Container.Exec", request, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Run(from string, address string, callString string, value int) (*OLVMResult, error) {
	return defaultClient.Run(from, address, callString, value)
}
