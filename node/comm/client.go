/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package comm

import (
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	client "github.com/tendermint/abci/client"
	"github.com/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// TODO: Why?
var _ *client.Client

// Generic Client interface, allows SetOption
func NewAppClient() client.Client {
	log.Debug("New Client", "address", global.Current.AppAddress, "transport", global.Current.Transport)

	// TODO: Try multiple times before giving up
	client, err := client.NewClient(global.Current.AppAddress, global.Current.Transport, true)
	if err != nil {
		log.Fatal("Can't start client", "err", err)
	}
	log.Debug("Have Client", "client", client)

	return client
}

func SetOption(key string, value string) {
	log.Debug("Setting Option")

	client := NewAppClient()
	options := types.RequestSetOption{
		Key:   key,
		Value: value,
	}

	// response := client.SetOptionAsync(options)
	response, err := client.SetOptionSync(options)
	log.Debug("Have Set Option")

	if err != nil {
		log.Error("SetOption Failed", "err", err, "response", response)
	}
}

var cachedClient *rpcclient.HTTP

// HTTP interface, allows Broadcast?
// TODO: Want to switch client type, based on config or cli args.
func GetClient() *rpcclient.HTTP {

	if cachedClient != nil {
		log.Debug("Cached RpcClient", "address", global.Current.RpcAddress)
		return cachedClient
	}

	log.Debug("Initializing RpcClient", "address", global.Current.RpcAddress)

	// TODO: Try multiple times before giving up
	cachedClient = rpcclient.NewHTTP(global.Current.RpcAddress, "/websocket")

	return cachedClient
}

// Broadcast packet to the chain
func Broadcast(packet []byte) *ctypes.ResultBroadcastTx {
	log.Debug("Broadcast")

	client := GetClient()

	//result, err := client.BroadcastTxCommit(packet)
	result, err := client.BroadcastTxAsync(packet)
	if err != nil {
		log.Error("Error", "err", err)
		os.Exit(-1)
	}

	return result
}

// Send a very specific query
func Query(path string, packet []byte) *ctypes.ResultABCIQuery {
	log.Debug("sending ABCi Query", "path", path, "packet", packet)

	client := GetClient()

	result, err := client.ABCIQuery(path, packet)
	if err != nil {
		log.Error("Error", "err", err)
		os.Exit(-1)
	}

	log.Debug("ABCi Query", "result", result)

	return result
}
