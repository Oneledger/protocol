package rpc

import (
	"context"
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"net/rpc"
	"net/url"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"github.com/powerman/rpc-codec/jsonrpc2"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"
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
	rpc           *rpc.Server
	http          *http.Server
	authenticator authHandler
	listener      net.Listener
	logger        *log.Logger
	// Request multiplexer
	mux *http.ServeMux
	cfg *config.Server
}

func NewServer(w io.Writer, config *config.Server) *Server {
	logger := log.NewLoggerWithPrefix(w, "rpc")
	return &Server{
		rpc:           rpc.NewServer(),
		http:          &http.Server{},
		authenticator: &rpcAuthHandler{},
		mux:           http.NewServeMux(),
		logger:        logger,
		cfg:           config,
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

type authHandler interface {
	http.Handler
	Authorized(respW http.ResponseWriter, req *http.Request) bool
}

var _ authHandler = &rpcAuthHandler{}

type rpcAuthHandler struct {
	rpcHandler http.Handler
	cfg        *config.Server
}

func (r rpcAuthHandler) ServeHTTP(respW http.ResponseWriter, req *http.Request) {
	if r.Authorized(respW, req) {
		r.rpcHandler.ServeHTTP(respW, req)
	}
}

func (r *rpcAuthHandler) Authorized(respW http.ResponseWriter, req *http.Request) bool {
	if r.cfg != nil && r.cfg.Node.Auth.RPCPrivateKey != "" {
		respErr := ""
		defer func() {
			if respErr != "" {
				http.Error(respW, respErr, 401)
			}
		}()

		respW.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		token := req.Header.Get("Authorization")

		data := base58.Decode(token)
		if len(data) <= 20 {
			respErr = "invalid token"
			return false
		}
		signData := data[:20]
		signature := data[20:]

		var keyData []byte
		//Get Private Key for signature verification.
		keyData, err := base64.StdEncoding.DecodeString(r.cfg.Node.Auth.RPCPrivateKey)
		if err != nil {
			respErr = err.Error()
			return false
		}

		privateKey, err := keys.GetPrivateKeyFromBytes(keyData, keys.ED25519)
		if err != nil {
			respErr = err.Error()
			return false
		}

		//Private Key Handler
		privateKeyHandler, err := privateKey.GetHandler()
		if err != nil {
			respErr = err.Error()
			return false
		}

		//Get public key from private key handler
		pubKey := privateKeyHandler.PubKey()
		pubKeyHandler, err := pubKey.GetHandler()
		if err != nil {
			respErr = err.Error()
			return false
		}

		//Verify Message with public key
		verified := pubKeyHandler.VerifyBytes(signData, signature)

		if !verified {
			respErr = "not authorized"
			return false
		}
	}

	return true
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

	//Register jsonrpc handler to authenticator.
	srv.authenticator = &rpcAuthHandler{jsonrpc2.HTTPHandler(srv.rpc), srv.cfg}

	// Register the handlers with our mux
	srv.mux.Handle(Path, srv.authenticator)
	srv.http.Handler = srv.mux
	srv.listener = l
	return nil
}

// Start spawns a new goroutine and listens on the configured address. Prepare
// should be called before calling this method
func (srv *Server) Start() error {

	channel := make(chan error)
	timeout := make(chan error)
	var err error

	if srv.listener == nil {
		return errors.New("no listener specified on server, was Prepare called?")
	}

	//Timeout Go routine
	go func() {
		time.Sleep(time.Duration(srv.cfg.Network.RPCStartTimeout) * time.Second)
		timeout <- nil
	}()

	go func(l net.Listener, ch chan error) {
		srv.logger.Info("starting RPC server on " + l.Addr().String())
		err := srv.http.Serve(l)
		if err != nil {
			srv.logger.Fatalf("server: %s", err)
		}
		ch <- err
	}(srv.listener, channel)

	select {
	case err = <-channel:
	case err = <-timeout:
	}

	return err
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
