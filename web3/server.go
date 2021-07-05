package web3

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/log"

	"github.com/ethereum/go-ethereum/rpc"
)

// Server holds an RPC server that is served over HTTP
type Server struct {
	rpc    *rpc.Server
	logger *log.Logger
	cfg    *config.Server
	url    *url.URL
}

func NewServer(w io.Writer, config *config.Server) (*Server, error) {
	logger := log.NewLoggerWithPrefix(w, "rpc")

	url, err := url.Parse(config.Network.Web3Address)
	if err != nil {
		return nil, err
	}
	return &Server{
		rpc:    rpc.NewServer(),
		logger: logger,
		cfg:    config,
		url:    url,
	}, nil
}

func (s *Server) RegisterName(name string, receiver interface{}) error {
	return s.rpc.RegisterName(name, receiver)
}

func (s *Server) Start() error {
	var err error
	channel := make(chan error)

	// TODO: Add cors handling
	http.HandleFunc("/", s.rpc.ServeHTTP)

	go func(ch chan error) {
		defer s.rpc.Stop()

		uri := fmt.Sprintf("%s:%s", s.url.Hostname(), s.url.Port())
		s.logger.Info("starting Web3 RPC server on " + uri)
		err := http.ListenAndServe(uri, nil)
		if err != nil {
			s.logger.Fatalf("server: %s", err)
		}
		ch <- err
	}(channel)

	select {
	case err = <-channel:
	}

	return err
}

// func main() {
// 	server := rpc.NewServer()
// 	defer server.Stop()

// 	http.HandleFunc("/", server.ServeHTTP)
// 	log.Fatal(http.ListenAndServe(":12345", nil))
// }
