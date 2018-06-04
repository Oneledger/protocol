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

var _ *client.Client

// Generic Client interface, allows SetOption
func NewClient() client.Client {
	log.Debug("New Client", "address", global.Current.App, "transport", global.Current.Transport)

	client, err := client.NewClient(global.Current.App, global.Current.Transport, true)
	if err != nil {
		log.Fatal("Can't start client", "err", err)
	}
	log.Debug("Have Client", "client", client)

	return client
}

func SetOption(key string, value string) {
	log.Debug("Setting Option")

	client := NewClient()
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
	log.Debug("RPCClient", "address", global.Current.Address)

	cachedClient = rpcclient.NewHTTP(global.Current.Address, "/websocket")
	return cachedClient
}

// Broadcast packet to the chain
func Broadcast(packet []byte) *ctypes.ResultBroadcastTxCommit {
	log.Debug("Broadcast")

	client := GetClient()

	result, err := client.BroadcastTxCommit(packet)
	if err != nil {
		log.Error("Error", "err", err)
		os.Exit(-1)
	}
	return result
}

// Send a very specific query
func Query(path string, packet []byte) *ctypes.ResultABCIQuery {
	log.Debug("Query")

	client := GetClient()

	result, err := client.ABCIQuery(path, packet)
	if err != nil {
		log.Error("Error", "err", err)
		os.Exit(-1)
	}
	return result
}
