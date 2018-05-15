/*
	Copyright 2017-2018 OneLedger

	Cover over the Tendermint client handling.

	TODO: Make this generic to handle HTTP and local clients
*/
package main

import (
	"os"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var cachedClient *rpcclient.HTTP

// TODO: Want to switch client type, based on config or cli args.
func GetClient() *rpcclient.HTTP {
	//cachedClient = rpcclient.NewHTTP("127.0.0.1:46657", "/websocket")
	cachedClient = rpcclient.NewHTTP(global.Current.Address, "/websocket")
	return cachedClient
}

// Broadcast packet to the chain
func Broadcast(packet []byte) *ctypes.ResultBroadcastTxCommit {
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
	client := GetClient()

	result, err := client.ABCIQuery(path, packet)
	if err != nil {
		log.Error("Error", "err", err)
		os.Exit(-1)
	}
	return result
}
