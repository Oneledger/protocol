package websockets

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	web3types "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
)

// WsServer defines a server that handles Ethereum websockets.
type WsServer struct {
	Address string
	api     *PubSubAPI
	logger  *log.Logger
}

// NewServer creates a new websocket server instance.
func NewServer(ctx web3types.Web3Context, config *config.Server) *WsServer {
	// TODO: make a new web3 websocket address in config.toml
	return &WsServer{
		Address: strings.Replace(config.Network.RPCAddress, "266", "276", 1),
		api:     NewAPI(ctx),
		logger:  log.NewLoggerWithPrefix(os.Stdout, "ws"),
	}
}

func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		err := fmt.Errorf("websocket upgrade failed: %w", err)
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Fatal(err)
		return
	}

	s.readLoop(wsConn)
}

// Name returns the JSON-RPC service name
func (s *WsServer) Name() string {
	return "Ethereum Websocket"
}

// Start runs the websocket server
func (s *WsServer) Start() error {
	ws := mux.NewRouter()
	ws.Handle("/ws", s)

	go func() {
		s.logger.Info("starting web3 websocket server on " + s.Address[6:])
		err := http.ListenAndServe(s.Address[6:], ws)
		if err != nil {
			s.logger.Fatal("start web3 websocket server error ", err.Error())
		}
	}()

	return nil
}

func (s *WsServer) sendErrResponse(conn *websocket.Conn, msg string) {
	res := &ErrorResponseJSON{
		Jsonrpc: "2.0",
		Error: &ErrorMessageJSON{
			Code:    big.NewInt(-32600),
			Message: msg,
		},
		ID: nil,
	}
	err := conn.WriteJSON(res)
	if err != nil {
		s.logger.Error("websocket failed write message", "error", err)
	}
}

func (s *WsServer) readLoop(wsConn *websocket.Conn) {
	for {
		_, mb, err := wsConn.ReadMessage()
		if err != nil {
			_ = wsConn.Close()
			s.logger.Error("failed to read message; error", err)
			return
		}

		var msg map[string]interface{}
		err = json.Unmarshal(mb, &msg)
		if err != nil {
			s.sendErrResponse(wsConn, "invalid request")
			continue
		}

		// check if method == eth_subscribe or eth_unsubscribe
		method := msg["method"]
		if method.(string) == "eth_subscribe" {
			params := msg["params"].([]interface{})
			if len(params) == 0 {
				s.sendErrResponse(wsConn, "invalid parameters")
				continue
			}

			id, err := s.api.Subscribe(wsConn, params)
			if err != nil {
				s.sendErrResponse(wsConn, err.Error())
				continue
			}

			res := &SubscriptionResponseJSON{
				Jsonrpc: "2.0",
				ID:      1,
				Result:  id,
			}

			err = wsConn.WriteJSON(res)
			if err != nil {
				s.logger.Error("failed to write json response", err)
				continue
			}

			continue
		} else if method.(string) == "eth_unsubscribe" {
			ids, ok := msg["params"].([]interface{})
			if _, idok := ids[0].(string); !ok || !idok {
				s.sendErrResponse(wsConn, "invalid parameters")
				continue
			}

			ok = s.api.Unsubscribe(rpc.ID(ids[0].(string)))
			res := &SubscriptionResponseJSON{
				Jsonrpc: "2.0",
				ID:      1,
				Result:  ok,
			}

			err = wsConn.WriteJSON(res)
			if err != nil {
				s.logger.Error("failed to write json response", err)
				continue
			}

			continue
		}

		// otherwise, call the usual rpc server to respond
		err = s.tcpGetAndSendResponse(wsConn, mb)
		if err != nil {
			s.sendErrResponse(wsConn, err.Error())
		}
	}
}

// tcpGetAndSendResponse connects to the rest-server over tcp, posts a JSON-RPC request, and sends the response
// to the client over websockets
func (s *WsServer) tcpGetAndSendResponse(conn *websocket.Conn, mb []byte) error {
	addr := strings.Split(s.Address, "tcp://")
	if len(addr) != 2 {
		return fmt.Errorf("invalid laddr %s", s.Address)
	}

	tcpConn, err := net.Dial("tcp", addr[1])
	if err != nil {
		return fmt.Errorf("cannot connect to %s; %s", s.Address, err)
	}

	buf := &bytes.Buffer{}
	_, err = buf.Write(mb)
	if err != nil {
		return fmt.Errorf("failed to write message; %s", err)
	}

	req, err := http.NewRequest("POST", s.Address, buf)
	if err != nil {
		return fmt.Errorf("failed to request; %s", err)
	}

	req.Header.Set("Content-Type", "application/json;")
	err = req.Write(tcpConn)
	if err != nil {
		return fmt.Errorf("failed to write to rest-server; %s", err)
	}

	respBytes, err := ioutil.ReadAll(tcpConn)
	if err != nil {
		return fmt.Errorf("error reading response from rest-server; %s", err)
	}

	respbuf := &bytes.Buffer{}
	respbuf.Write(respBytes)
	resp, err := http.ReadResponse(bufio.NewReader(respbuf), req)
	if err != nil {
		return fmt.Errorf("could not read response; %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read body from response; %s", err)
	}

	var wsSend interface{}
	err = json.Unmarshal(body, &wsSend)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rest-server response; %s", err)
	}

	return conn.WriteJSON(wsSend)
}
