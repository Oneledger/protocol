/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package comm

import (
	"reflect"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var cachedClient rpcclient.Client

// HTTP interface, allows Broadcast?
// TODO: Want to switch client type, based on config or cli args.
func GetClient() (client rpcclient.Client) {

	defer func() {
		if r := recover(); r != nil {
			log.Debug("Ignoring Client Panic", "r", r)
			client = nil
		}
	}()

	if cachedClient != nil {
		return cachedClient
	}

	if global.Current.ConsensusNode != nil {
		log.Debug("Using local ConsensusNode ABCI Client")
		cachedClient = rpcclient.NewLocal(global.Current.ConsensusNode)

	} else {
		log.Debug("Using new HTTP ABCI Client")
		cachedClient = rpcclient.NewHTTP(global.Current.RpcAddress, "/websocket")
	}

	if _, err := cachedClient.Status(); err == nil {
		log.Debug("Client is running")
		return cachedClient
	}

	if err := cachedClient.Start(); err != nil {
		log.Fatal("Client is unavailable", "address", global.Current.RpcAddress)
		client = nil
	}
	// TODO: Try multiple times before giving up

	//for i := 0; i < 10; i++ {
	//	cachedClient = rpcclient.NewHTTP(global.Current.RpcAddress, "/websocket")
	//
	//	if cachedClient != nil {
	//		log.Debug("RPC Client", "address", global.Current.RpcAddress, "client", cachedClient)
	//		break
	//	}
	//
	//	log.Warn("Retrying RPC Client", "address", global.Current.RpcAddress)
	//	time.Sleep(1 * time.Second)
	//}
	//
	//for i := 0; i < 10; i++ {
	//	result, err := cachedClient.Status()
	//	if err == nil {
	//		log.Debug("Connected", "result", result)
	//		break
	//	}
	//	log.Warn("Waiting for RPC Client", "address", global.Current.RpcAddress)
	//	time.Sleep(2 * time.Second)
	//}

	return cachedClient
}

func StopClient() {
	if cachedClient != nil && cachedClient.IsRunning() {
		//cachedClient.Stop()
	}
}

// An async Broadcast to the chain
func BroadcastAsync(packet []byte) *ctypes.ResultBroadcastTx {

	if len(packet) < 1 {
		log.Debug("Empty Transaction")
		return nil
	}

	client := GetClient()

	result, err := client.BroadcastTxAsync(packet)
	StopClient()

	if err != nil {
		log.Error("Broadcast Error", "err", err)
	}

	log.Debug("Broadcast", "packet", packet, "result", result)

	return result
}

// A sync'ed broadcast to the chain that waits for the commit to happen
func Broadcast(packet []byte) *ctypes.ResultBroadcastTxCommit {
	if len(packet) < 1 {
		log.Debug("Empty Transaction")
		return nil
	}

	client := GetClient()

	log.Debug("Start Synced Broadcast", "packet", packet)

	// TODO: result, err := client.BroadcastTxSync(packet)
	result, err := client.BroadcastTxCommit(packet)
	StopClient()

	if err != nil {
		log.Error("Error", "err", err)
	}

	log.Debug("Finished Synced Broadcast", "packet", packet, "result", result)

	return result
}

// A sync'ed broadcast to the chain that waits for the commit to happen
func BroadcastSync(packet []byte) *ctypes.ResultBroadcastTx {
	if len(packet) < 1 {
		log.Debug("Empty Transaction")
		return nil
	}

	client := GetClient()

	log.Debug("Start Synced Broadcast", "packet", packet)

	result, err := client.BroadcastTxSync(packet)
	StopClient()

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
	if len(path) < 1 {
		log.Debug("Empty Query Path")
		return nil
	}

	var response *ctypes.ResultABCIQuery
	var err error

	//for i := 0; i < 20; i++ {
	client := GetClient()
	response, err = client.ABCIQuery(path, packet)
	StopClient()

	if err != nil {
		log.Debug("ABCi Query Error", "path", path, "err", err)
		return nil
	}
	//if response != nil {
	//	break
	//}
	//time.Sleep(2 * time.Second)
	//}

	if response == nil {
		//return "No results for " + path + " and " + string(packet)
		log.Debug("response is empty")
		return nil
	}

	var prototype interface{}
	result, err := serial.Deserialize(response.Response.Value, prototype, serial.CLIENT)
	if err != nil {
		log.Error("Failed to deserialize Query:", "response", response.Response.Value)
		return nil
	}

	return result
}

func Tx(hash []byte, prove bool) (res *ctypes.ResultTx) {
	client := GetClient()

	result, err := client.Tx(hash, prove)
	if err != nil {
		log.Error("TxSearch Error", "err", err)
		return nil
	}

	log.Debug("TxSearch", "hash", hash, "prove", prove, "result", result)

	return result
}

func Search(query string, prove bool, page, perPage int) (res *ctypes.ResultTxSearch) {
	client := GetClient()

	result, err := client.TxSearch(query, prove, page, perPage)
	if err != nil {
		log.Error("TxSearch Error", "err", err)
		return nil
	}

	log.Debug("TxSearch", "query", query, "prove", prove, "result", result)

	return result
}

func Block(height int64) (res *ctypes.ResultBlock) {
	client := GetClient()

	// Pass nil if given 0 to return the latest block
	var h *int64
	if height != 0 {
		h = &height
	}
	result, err := client.Block(h)
	if err != nil {
		return nil
	}
	return result
}
