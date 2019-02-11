/*
	Copyright 2017-2018 OneLedger
*/
package vm

import (
	"net/rpc"
	"strings"
	"time"

	"github.com/Oneledger/protocol/node/action"
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
func InitializeClient() {
	protocol := global.Current.OLVMProtocol
	address := global.Current.OLVMAddress

	defaultClient = NewClient(protocol, address)
}

func AutoRun(request *action.OLVMRequest) (result *action.OLVMResult, err error) {

	log.Debug("Trying to Run")
	result, err = defaultClient.Run(request)
	var count int = 10

	// TODO: Should be based on error code, not text...
	if err != nil {
		log.Dump("Failed to run", err, result)
		if strings.HasSuffix(err.Error(), "connection refused") {

			// Pause for a bit, might be a race condition
			time.Sleep(time.Second)

			for strings.HasSuffix(err.Error(), "connection refused") {

				// Always bound loops with a fixed count
				if count < 0 {
					log.Fatal("Can't connect", "err", err)
				}

				log.Dump("Failed Again", err, result)
				time.Sleep(time.Second)
				log.Debug("Trying to ReRun")
				result, err = defaultClient.Run(request)
				count--
			}
		} else {
			log.Error("Run Failed", "err", err)
		}
		return
	}
	return
}

func Analyze(request *action.OLVMRequest) (result *action.OLVMResult, err error) {

	log.Debug("Analyze the smart contract")
	result, err = defaultClient.RunAnalyze(request)
	var count int = 10

	// TODO: Should be based on error code, not text...
	if err != nil {
		log.Dump("Failed to run", err, result)
		if strings.HasSuffix(err.Error(), "connection refused") {

			// Pause for a bit, might be a race condition
			time.Sleep(time.Second)

			for strings.HasSuffix(err.Error(), "connection refused") {

				// Always bound loops with a fixed count
				if count < 0 {
					log.Fatal("Can't connect", "err", err)
				}

				log.Dump("Failed Again", err, result)
				time.Sleep(time.Second)
				log.Debug("Trying to ReRun")
				result, err = defaultClient.RunAnalyze(request)
				count--
			}
		} else {
			log.Error("Run Failed", "err", err)
		}
		return
	}
	return
}

// Run a smart contract
func (c OLVMClient) RunAnalyze(request *action.OLVMRequest) (*action.OLVMResult, error) {

	log.Info("Dialing service...", "protocol", c.Protocol, "service", c.ServicePath)

	client, err := rpc.DialHTTP(c.Protocol, ":"+sdk.GetPort(c.ServicePath))
	if err != nil {
		log.Dump("Failded to Connect", err, client)
		return nil, err
	}

	// TODO: Shouldn't pass by address for the result
	result := &action.OLVMResult{}
	err = client.Call("Container.Analyze", request, result)
	if err != nil {
		log.Dump("Failded to Exec", err, result)
		return nil, err
	}

	client.Close()

	log.Dump("Have a Result", result)
	return result, nil
}

// Run a smart contract
func (c OLVMClient) Run(request *action.OLVMRequest) (*action.OLVMResult, error) {

	log.Info("Dialing service...", "protocol", c.Protocol, "service", c.ServicePath)

	client, err := rpc.DialHTTP(c.Protocol, ":"+sdk.GetPort(c.ServicePath))
	if err != nil {
		log.Dump("Failded to Connect", err, client)
		return nil, err
	}

	// TODO: Shouldn't pass by address for the result
	result := &action.OLVMResult{}
	err = client.Call("Container.Exec", request, result)
	if err != nil {
		log.Dump("Failded to Exec", err, result)
		return result, err
	}

	client.Close()

	log.Dump("Have a Result", result)
	return result, nil
}

func Run(request *action.OLVMRequest) (*action.OLVMResult, error) {
	return defaultClient.Run(request)
}
