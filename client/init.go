package client

import (
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
	"os"
)

var clSerializer serialize.Serializer
var logger *log.Logger

func init() {
	clSerializer = serialize.GetSerializer(serialize.CLIENT)
	logger = log.NewLoggerWithPrefix(os.Stdout, "client")
}

const (

	// Broadcast a tx and wait for being commit to block
	BroadcastCommit = "commit"

	// Broadcast a tx and wait for CheckTx response
	// Recommend for smart contract
	BroadcastSync = "sync"

	// Broadcast a tax and return immediately
	BroadcastAsync = "async"



	RPC_ADDRESS = "127.0.0.1:6019"
)
