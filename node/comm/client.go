/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package comm

import (
	"fmt"
	"github.com/tendermint/tendermint/node"

	"github.com/Oneledger/protocol/node/serial"
	"github.com/pkg/errors"
	"reflect"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type ClientContext struct {
	Client        rpcclient.Client
	BroadcastMode string
	Proof         bool
}

func NewLocalClientContext(node *node.Node) (cliCtx ClientContext) {
	rpc := rpcclient.NewLocal(node)

	cliCtx = ClientContext{
		Client:        rpc,
		BroadcastMode: BroadcastAsync,
		Proof:         false,
	}
	return
}

func NewClientContext() (cliCtx ClientContext) {
	var rpc rpcclient.Client

	defer func() {
		if r := recover(); r != nil {
			log.Debug("Ignoring Client Panic", "r", r)
		}
	}()

	rpc = rpcclient.NewHTTP(global.Current.Config.Network.RPCAddress, "/websocket")

	if _, err := rpc.Status(); err == nil {
		log.Debug("Client is running")

		cliCtx = ClientContext{
			Client:        rpc,
			BroadcastMode: global.Current.ClientConfig.BroadcastMode,
		}
		return
	}

	if err := rpc.Start(); err != nil {
		log.Fatal("Client is unavailable", "address", global.Current.Config.Network.RPCAddress)
		return
	}
	return
}

func (ctx ClientContext) GetClient() (client rpcclient.Client, err error) {
	if ctx.Client == nil {
		return nil, errors.New("no rpc comm initialized")
	}

	if _, err := ctx.Client.Status(); err != nil {
		err = ctx.Client.Start()
		if err != nil {
			return nil, fmt.Errorf("rpc comm not available: %v", err)
		}
	}

	return ctx.Client, nil
}

func (ctx ClientContext) BroadcastTxSync(packet []byte) (res *ctypes.ResultBroadcastTx, err error) {

	client, err := ctx.GetClient()
	if err != nil {
		return res, err
	}

	if len(packet) < 1 {
		return res, errors.New("empty transaction")
	}

	result, err := client.BroadcastTxSync(packet)

	if err != nil {
		return res, err
	}

	return result, nil
}

func (ctx ClientContext) BroadcastTxAsync(packet []byte) (res *ctypes.ResultBroadcastTx, err error) {

	client, err := ctx.GetClient()
	if err != nil {
		return res, err
	}

	if len(packet) < 1 {
		return nil, errors.New("empty transaction")
	}

	result, err := client.BroadcastTxAsync(packet)

	if err != nil {
		return res, err
	}

	return result, nil
}

func (ctx ClientContext) BroadcastTxCommit(packet []byte) (res *ctypes.ResultBroadcastTxCommit, err error) {

	client, err := ctx.GetClient()
	if err != nil {
		return res, err
	}

	if len(packet) < 1 {
		return nil, errors.New("empty transaction")
	}

	result, err := client.BroadcastTxCommit(packet)

	if err != nil {
		return nil, err
	}

	return result, nil
}

//query to return abci response
func (ctx ClientContext) ABCIQuery(path string, packet []byte) (res *ctypes.ResultABCIQuery, err error) {

	if len(path) < 1 {
		return nil, errors.New("empty query path")
	}

	var response *ctypes.ResultABCIQuery

	response, err = ctx.Client.ABCIQuery(path, packet)

	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("empty response")
	}

	return response, nil
}

// Send a very specific query to return parse response.value
func (ctx ClientContext) Query(path string, packet []byte) (res interface{}) {

	var response *ctypes.ResultABCIQuery

	response, err := ctx.ABCIQuery(path, packet)
	if err != nil {
		log.Debug("query response is empty", "request", packet, "err", err)
		return res
	}

	var result interface{}

	_, isTransactionPath := transactionPathsMap[path]
	if isTransactionPath {
		// we continue to use old serializer for query handlers who
		// return transaction interface, which is yet to be moved to
		// the new serializer
		var proto interface{}
		result, err = serial.Deserialize(response.Response.Value, proto, serial.CLIENT)
	} else {

		err = clSerializer.Deserialize(response.Response.Value, &result)
	}

	if err != nil {
		log.Error("Failed to deserialize Query:", "response", response.Response.Value)
		return res
	}

	return result
}

func (ctx ClientContext) Tx(hash []byte, prove bool) (res *ctypes.ResultTx) {

	result, err := ctx.Client.Tx(hash, prove)
	if err != nil {
		log.Error("TxSearch Error", "err", err)
		return nil
	}

	log.Debug("TxSearch", "hash", hash, "prove", prove, "result", result)
	return result
}

func (ctx ClientContext) Block(height int64) (res *ctypes.ResultBlock) {

	// Pass nil if given 0 to return the latest block
	var h *int64
	if height != 0 {
		h = &height
	}
	result, err := ctx.Client.Block(h)
	if err != nil {
		return nil
	}
	return result
}

func (ctx ClientContext) Search(query string, prove bool, page, perPage int) (res *ctypes.ResultTxSearch) {
	client, err := NewClientContext().GetClient()
	if err != nil {
		return
	}

	result, err := client.TxSearch(query, prove, page, perPage)
	if err != nil {
		log.Error("TxSearch Error", "err", err)
		return nil
	}

	log.Debug("TxSearch", "query", query, "prove", prove, "result", result)

	return result
}

var transactionPathsMap = map[string]bool{
	"/applyValidators":     true,
	"/createExSendRequest": true,
	"/createSendRequest":   true,
	"/createMintRequest":   true,
	"/createSwapRequest":   true,
	"/nodeName":            true,
	"/signTransaction":     true,
}

func IsError(result interface{}) *string {
	if reflect.TypeOf(result).Kind() == reflect.String {
		final := result.(string)
		return &final
	}
	return nil
}
