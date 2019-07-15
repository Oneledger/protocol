package client

import (
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
)

var clSerializer serialize.Serializer
var logger *log.Logger

func init() {
	clSerializer = serialize.GetSerializer(serialize.CLIENT)
	logger = log.NewLoggerWithPrefix(os.Stdout, "client")
}

type BroadcastMode string

const (
	BROADCASTASYNC  BroadcastMode = "async"
	BROADCASTSYNC   BroadcastMode = "sync"
	BROADCASTCOMMIT BroadcastMode = "commit"
)
