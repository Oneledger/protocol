/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package comm

import (
	"reflect"
	"time"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	client "github.com/tendermint/tendermint/abci/client"
	"github.com/tendermint/tendermint/abci/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// TODO: Why?
var _ *client.Client

// Generic Client interface, allows SetOption
func NewAppClient() client.Client {
	//log.Debug("New Client", "address", global.Current.AppAddress, "transport", global.Current.Transport)

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

	if cachedClient != nil {
		//log.Debug("Cached RpcClient", "address", global.Current.RpcAddress)
		return cachedClient
	}

	// TODO: Try multiple times before giving up

	for i := 0; i < 10; i++ {
		cachedClient = rpcclient.NewHTTP(global.Current.RpcAddress, "/websocket")

		if cachedClient != nil {
			log.Debug("RPC Client", "address", global.Current.RpcAddress, "client", cachedClient)
			break
		}

		log.Warn("Retrying RPC Client", "address", global.Current.RpcAddress)
		time.Sleep(1 * time.Second)
	}

	for i := 0; i < 200; i++ {
		result, err := cachedClient.Status()
		if err == nil {
			log.Debug("Connected", "result", result)
			break
		}
		log.Warn("Waiting for RPC Client", "address", global.Current.RpcAddress)
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

	log.Debug("Start Synced Broadcast", "packet", packet)

	// TODO: result, err := client.BroadcastTxSync(packet)
	result, err := client.BroadcastTxCommit(packet)
	if err != nil {
		log.Error("Error", "err", err)
	}

	log.Debug("Finished Synced Broadcast", "packet", packet, "result", result)

	return result
}

// A sync'ed broadcast to the chain that waits for the commit to happen
func BroadcastSync(packet []byte) *ctypes.ResultBroadcastTx {
	client := GetClient()

	log.Debug("Start Synced Broadcast", "packet", packet)

	result, err := client.BroadcastTxSync(packet)
	if err != nil {
		log.Error("Error", "err", err)
	}

	log.Debug("Finished Synced Broadcast", "packet", packet, "result", result)

	return result
}

func IsError(result interface{}) *string {
	if reflect.TypeOf(result).Kind() == reflect.String {
		final := result.(string)
		return &final
	}
	return nil
}

// Send a very specific query
func Query(path string, packet []byte) interface{} {

	var response *ctypes.ResultABCIQuery
	var err error

	for i := 0; i < 20; i++ {
		client := GetClient()
		response, err = client.ABCIQuery(path, packet)
		if err != nil {
			log.Error("ABCi Query Error", "path", path, "err", err)
			return nil
		}
		if response != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if response == nil {
		return "No results for " + path + " and " + string(packet)
	}

	var prototype interface{}
	result, err := serial.Deserialize(response.Response.Value, prototype, serial.CLIENT)
	if err != nil {
		log.Error("Failed to deserialize Query:", response.Response.Value)
		return "Failed"
	}
	return result
}

func Tx(hash []byte, prove bool) (res *ctypes.ResultTx) {
	client := GetClient()

	result, err := client.Tx(hash, prove)
	if err != nil {
		log.Error("TxSearch Error", "err", err)
	}

	log.Debug("TxSearch", "hash", hash, "prove", prove, "result", result)

	return result
}

func Search(query string, prove bool, page, perPage int) (res *ctypes.ResultTxSearch) {
	client := GetClient()

	result, err := client.TxSearch(query, prove, page, perPage)
	if err != nil {
		log.Error("TxSearch Error", "err", err)
	}

	log.Debug("TxSearch", "query", query, "prove", prove, "result", result)

	return result
}
