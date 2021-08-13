package websockets

import (
	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

// PubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
type PubSubAPI struct {
	clientCtx rpctypes.Web3Context
	logger    log.Logger
}

// NewAPI creates an instance of the ethereum PubSub API.
func NewAPI(clientCtx rpctypes.Web3Context) *PubSubAPI {
	return &PubSubAPI{
		clientCtx: clientCtx,
		logger:    log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "websocket-client"),
	}
}
