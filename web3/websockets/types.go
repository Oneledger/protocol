package websockets

import (
	rpcfilters "github.com/Oneledger/protocol/web3/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"math/big"
)

type ErrorResponseJSON struct {
	Jsonrpc string            `json:"jsonrpc"`
	Error   *ErrorMessageJSON `json:"error"`
	ID      *big.Int          `json:"id"`
}

type ErrorMessageJSON struct {
	Code    *big.Int `json:"code"`
	Message string   `json:"message"`
}

type SubscriptionResponseJSON struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      float64     `json:"id"`
}

type wsSubscription struct {
	sub          *rpcfilters.Subscription
	unsubscribed chan struct{} // closed when unsubscribing
	conn         *websocket.Conn
}

type SubscriptionNotification struct {
	Jsonrpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  *SubscriptionResult `json:"params"`
}

type SubscriptionResult struct {
	Subscription rpc.ID      `json:"subscription"`
	Result       interface{} `json:"result"`
}
