/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package comm

import (
	"github.com/Oneledger/protocol/node/status"
	"reflect"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type ClientContext struct {
	Client rpcclient.Client
	Async bool
}

type ClientInterface interface {
	BroadcastTxSync(packet []byte) (*ctypes.ResultBroadcastTx, status.Code)
	BroadcastTxAsync(packet []byte) (*ctypes.ResultBroadcastTx, status.Code)
	BroadcastTxCommit(packet []byte) (*ctypes.ResultBroadcastTxCommit, status.Code)
	ABCIQuery(path string, packet []byte) (*ctypes.ResultABCIQuery, error)
	Tx(hash []byte, prove bool) (*ctypes.ResultTx, error)
	TxSearch(query string, prove bool, page, perPage int) (*ctypes.ResultTxSearch, error)
	Block(height *int64) (*ctypes.ResultBlock, error)
	Status() (*ctypes.ResultStatus, error)
	Start() error
	IsRunning() bool
}

func (client ClientContext) BroadcastTxSync(packet []byte) (*ctypes.ResultBroadcastTx, status.Code) {

	if len(packet) < 1 {
		log.Debug("Empty Transaction")
		return nil, status.MISSING_DATA
	}

	log.Debug("Start Synced Broadcast", "packet", packet)

	result, err := client.Client.BroadcastTxSync(packet)

	StopClient()

	if err != nil {
		log.Error("Error", "err", err)
		return nil, status.EXECUTE_ERROR
	}

	log.Debug("Finished Synced Broadcast", "packet", packet, "result", result)

	return result, status.SUCCESS
}

func (client ClientContext) BroadcastTxAsync(packet []byte) (*ctypes.ResultBroadcastTx, status.Code) {
	if len(packet) < 1 {
		log.Debug("Empty Transaction")
		return nil, status.MISSING_DATA
	}

	result, err := client.Client.BroadcastTxAsync(packet)

	// @todo Do we need to stop Client?
	StopClient()

	if err != nil {
		log.Error("Broadcast Error", "err", err)
		return nil, status.EXECUTE_ERROR
	}

	log.Debug("Broadcast", "packet", packet, "result", result)

	return result, status.SUCCESS
}

func (client ClientContext) BroadcastTxCommit(packet []byte) (*ctypes.ResultBroadcastTxCommit, status.Code) {
	if len(packet) < 1 {
		log.Debug("Empty Transaction")
		return nil, status.MISSING_DATA
	}

	log.Debug("Start Synced Broadcast", "packet", packet)

	result, err := client.Client.BroadcastTxCommit(packet)

	// @todo Do we need to stop Client?
	StopClient()

	if err != nil {
		log.Error("Error", "err", err)
		return nil, status.EXECUTE_ERROR
	}

	log.Debug("Finished Synced Broadcast", "packet", packet, "result", result)

	return result, status.SUCCESS
}

func (client ClientContext) ABCIQuery(path string, packet []byte) (*ctypes.ResultABCIQuery, error) {
	return client.Client.ABCIQuery(path, packet)
}

func (client ClientContext) Tx(hash []byte, prove bool) (*ctypes.ResultTx, error) {
	return client.Client.Tx(hash, prove)
}

func (client ClientContext) TxSearch(query string, prove bool, page, perPage int) (*ctypes.ResultTxSearch, error) {
	return client.Client.TxSearch(query, prove, page, perPage)
}

func (client ClientContext) Block(height *int64) (*ctypes.ResultBlock, error) {
	return client.Client.Block(height)
}

func (client ClientContext) Status() (*ctypes.ResultStatus, error) {
	return client.Client.Status()
}

func (client ClientContext) Start() error {
	return client.Client.Start()
}

func (client ClientContext) IsRunning() bool {
	return client.Client.IsRunning()
}

var cachedClient ClientInterface

// HTTP interface, allows Broadcast?
// TODO: Want to switch client type, based on config or cli args.
func GetClient() (client ClientInterface) {

	var rpc rpcclient.Client

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
		rpc = rpcclient.NewLocal(global.Current.ConsensusNode)

	} else {
		log.Debug("Using new HTTP ABCI Client")
		rpc = rpcclient.NewHTTP(global.Current.RpcAddress, "/websocket")
	}

	client = ClientContext {
		Client: rpc,
		Async:  false,
	}

	if _, err := client.Status(); err == nil {
		log.Debug("Client is running")
		cachedClient = client
		return
	}

	if err := client.Start(); err != nil {
		log.Fatal("Client is unavailable", "address", global.Current.RpcAddress)
		client = nil
		return
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

	return
}

func StopClient() {
	if cachedClient != nil && cachedClient.IsRunning() {
		//cachedClient.Stop()
	}
}

// An async Broadcast to the chain
func BroadcastAsync(packet []byte) *ctypes.ResultBroadcastTx {

	client := GetClient()

	result, _ := client.BroadcastTxAsync(packet)

	return result
}

// A sync'ed broadcast to the chain that waits for the commit to happen
func Broadcast(packet []byte) *ctypes.ResultBroadcastTxCommit {

	client := GetClient()

	result, _ := client.BroadcastTxCommit(packet)

	return result
}

// A sync'ed broadcast to the chain that waits for the commit to happen
func BroadcastSync(packet []byte) *ctypes.ResultBroadcastTx {

	client := GetClient()

	result, _ := client.BroadcastTxSync(packet)

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
	if client == nil {
		log.Debug("Client Unavailable")
		return nil
	}

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
