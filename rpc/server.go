package rpc

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
	"time"

	"github.com/Oneledger/protocol/log"
	"github.com/pkg/errors"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

// The http path used for our rpc handlers
const (
	// All incoming requests must show this MIME type in the header
	ContentType = "application/json"

	PathJSON = "/jsonrpc"
	PathGOB  = "/rpc/gob"
	Path     = PathJSON
)

// Server holds an RPC server that is served over HTTP
type Server struct {
	rpc      *rpc.Server
	http     *http.Server
	listener net.Listener
	logger   *log.Logger
	// Request multiplexer
	mux *http.ServeMux
}

func NewServer(w io.Writer) *Server {
	logger := log.NewLoggerWithPrefix(w, "rpc")
	return &Server{
		rpc:    rpc.NewServer(),
		http:   &http.Server{},
		mux:    http.NewServeMux(),
		logger: logger,
	}
}

// Register creates a service on the Server with the given name.
// The criteria of a service method is the same as defined in the net/rpc package:
// - the method's type is exported.
// - the method is exported.
// - the method has two arguments, both exported (or builtin) types.
// - the method's second argument is a pointer.
// - the method has return type error.
func (srv *Server) Register(name string, rcvr interface{}) error {
	return srv.rpc.RegisterName(name, rcvr)
}

// RegisterRestfulMap registers all restful API functions in a map on the Server
func (srv *Server) RegisterRestfulMap(routerMap map[string]http.HandlerFunc) {
	for path, handlerFun := range routerMap {
		srv.mux.HandleFunc(path, handlerFun)
	}
}

// Prepare injects all the data necessary for serving over the specified URL.
// It  prepares a net.Listener over the specified URL, and registers all methods
// inside the given receiver. After this method is called, the Start function
// is ready to be called.
func (srv *Server) Prepare(u *url.URL) error {
	if u == nil {
		return errors.New("no URL was provided")
	} else if u.Port() == "" {
		return errors.New("no port was provided")
	}

	l, err := net.Listen("tcp", u.Host)
	if err != nil {
		return errors.Wrap(err, "invalid URL provided, failed to create listener")
	}

	// Register the handlers with our mux
	srv.mux.Handle(Path, jsonrpc2.HTTPHandler(srv.rpc))
	srv.http.Handler = srv.mux
	srv.listener = l
	return nil
}

// Start spawns a new goroutine and listens on the configured address. Prepare
// should be called before calling this method
func (srv *Server) Start() error {

	if srv.listener == nil {
		return errors.New("no listener specified on server, was Prepare called?")
	}
	go func(l net.Listener) {
		srv.logger.Info("starting RPC server on " + l.Addr().String())
		err := srv.http.Serve(l)
		if err != nil {
			srv.logger.Fatalf("server: %s", err)
		}
	}(srv.listener)

	return nil
}

// Close terminates the underlying HTTP server and listener
func (srv *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.logger.Info("closing server")
	err := srv.http.Shutdown(ctx)
	if err != nil {
		srv.logger.Error("Error shutting down", err)
	}
}
