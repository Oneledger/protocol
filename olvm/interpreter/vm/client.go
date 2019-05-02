/*
	Copyright 2017-2018 OneLedger
*/
package vm

import (
	"net/rpc"
	"strings"
	"time"

	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/utils"
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
func InitializeClient(protocol, address string) {

	defaultClient = NewClient(protocol, address)
}

func AutoRun(request *data.OLVMRequest) (result *data.OLVMResult, err error) {

	log.Debug("Trying to Run")
	result, err = defaultClient.Run(request)
	var count int = 10

	// TODO: Should be based on error code, not text...
	if err != nil {
		log.Errorf("Failed to run %s %#v", err, *result)
		if strings.HasSuffix(err.Error(), "connection refused") {

			// Pause for a bit, might be a race condition
			time.Sleep(time.Second)

			for strings.HasSuffix(err.Error(), "connection refused") {

				// Always bound loops with a fixed count
				if count < 0 {
					log.Fatal("Can't connect", "err", err)
				}

				log.Errorf("Failed Again %s %#v", err, result)
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

func Analyze(request *data.OLVMRequest) (result *data.OLVMResult, err error) {

	log.Debug("Analyze the smart contract")
	result, err = defaultClient.RunAnalyze(request)
	var count int = 10

	// TODO: Should be based on error code, not text...
	if err != nil {

		log.Errorf("Failed to run %s %#v", err, result)
		if strings.HasSuffix(err.Error(), "connection refused") {

			// Pause for a bit, might be a race condition
			time.Sleep(time.Second)

			for strings.HasSuffix(err.Error(), "connection refused") {

				// Always bound loops with a fixed count
				if count < 0 {
					log.Fatal("Can't connect", "err", err)
				}

				log.Error("Failed Again", err, *result)
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
func (c OLVMClient) RunAnalyze(request *data.OLVMRequest) (*data.OLVMResult, error) {

	log.Info("Dialing service...", "protocol", c.Protocol, "service", c.ServicePath)

	port, err := utils.GetPort(c.ServicePath)
	if err != nil {
		log.Fatal("parsing error", "err", err, "address:", c.ServicePath)
	}

	client, err := rpc.DialHTTP(c.Protocol, ":"+port)
	if err != nil {
		log.Errorf("Failded to Connect error:%s %#v ", err, client)
		return nil, err
	}

	// TODO: Shouldn't pass by address for the result
	result := &data.OLVMResult{}
	err = client.Call("Container.Analyze", request, result)
	if err != nil {
		log.Errorf("Failded to Exec error:%s result:%#v", err, result)
		return nil, err
	}

	client.Close()

	log.Error("Have a Result", result)
	return result, nil
}

// Run a smart contract
func (c OLVMClient) Run(request *data.OLVMRequest) (*data.OLVMResult, error) {

	log.Info("Dialing service...", "protocol", c.Protocol, "service", c.ServicePath)

	port, err := utils.GetPort(c.ServicePath)
	if err != nil {
		log.Fatal("parsing error", "err", err, "address:", c.ServicePath)
	}

	client, err := rpc.DialHTTP(c.Protocol, ":"+port)
	if err != nil {
		log.Errorf("Failed to Connect error:%s result:%#v", err, client)
		return nil, err
	}

	// TODO: Shouldn't pass by address for the result
	result := &data.OLVMResult{}
	err = client.Call("Container.Exec", request, result)
	if err != nil {
		log.Errorf("Failed to Exec error:%s result:%#v", err, result)
		return result, err
	}

	client.Close()

	log.Debug("Have a Result", *result)
	return result, nil
}

func Run(request *data.OLVMRequest) (*data.OLVMResult, error) {
	return defaultClient.Run(request)
}
