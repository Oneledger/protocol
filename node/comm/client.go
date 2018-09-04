/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package comm

import (
	"time"

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

// Set an option in the ABCi app directly
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
func GetClient() (client *rpcclient.HTTP) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Ignoring Client Panic", "r", r)
			client = nil
		}
	}()

	/*
		if cachedClient != nil {
			log.Debug("Cached RpcClient", "address", global.Current.RpcAddress)
			return cachedClient
		}
	*/

	// TODO: Try multiple times before giving up

	for i := 0; i < 3; i++ {
		cachedClient = rpcclient.NewHTTP(global.Current.RpcAddress, "/websocket")

		log.Debug("RPC Client", "address", global.Current.RpcAddress, "client", cachedClient)
		if cachedClient != nil {
			break
		}

		log.Warn("Retrying RPC Client", "address", global.Current.RpcAddress)
		time.Sleep(1 * time.Second)
	}

	return cachedClient
}

// An async Broadcast to the chain
func BroadcastAsync(packet []byte) *ctypes.ResultBroadcastTx {

	client := GetClient()

	result, err := client.BroadcastTxAsync(packet)
	if err != nil {
		log.Error("Broadcast Error", "err", err)
	}

	log.Debug("Broadcast", "packet", packet, "result", result)

	return result
}

// A sync'ed broadcast to the chain that waits for the commit to happen
func Broadcast(packet []byte) *ctypes.ResultBroadcastTxCommit {
	client := GetClient()

	result, err := client.BroadcastTxCommit(packet)
	if err != nil {
		log.Error("Error", "err", err)
	}

	log.Debug("Broadcast", "packet", packet, "result", result)

	return result
}

// Send a very specific query
func Query(path string, packet []byte) (res *ctypes.ResultABCIQuery) {
	client := GetClient()

	result, err := client.ABCIQuery(path, packet)
	if err != nil {
		log.Error("ABCi Query Error", "err", err)
		return nil
	}

	log.Debug("ABCi Query", "path", path, "packet", packet, "result", result)

	return result
}
