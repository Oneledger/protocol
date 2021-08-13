package websockets

import (
	"fmt"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"
	web3types "github.com/Oneledger/protocol/web3/types"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
)

// WsServer defines a server that handles Ethereum websockets.
type WsServer struct {
	Address string
	api     *PubSubAPI
	logger  *log.Logger
}

// NewServer creates a new websocket server instance.
func NewServer(ctx web3types.Web3Context, config *config.Server) *WsServer {
	return &WsServer{
		Address: config.Web3.WSAddress,
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
		s.logger.Info("starting web3 websocket server on " + s.Address[5:])
		err := http.ListenAndServe(s.Address[5:], ws)
		if err != nil {
			s.logger.Fatal("start web3 websocket server error ", err.Error())
		}
	}()

	return nil
}

func (s *WsServer) readLoop(wsConn *websocket.Conn) {
}