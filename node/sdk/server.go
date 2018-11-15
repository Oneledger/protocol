package sdk

import (
	"errors"
	"fmt"
	"net"
	"regexp"

	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/tendermint/tendermint/libs/common"
	tlog "github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
)

// implements common.Service
type Server struct {
	listener net.Listener
	server   *grpc.Server
	quit     chan struct{}
	logger   tlog.Logger
}

// Ensure sdk.Server implements common.Service
var _ common.Service = &Server{}

// Error messages
const (
	HAS_STARTED     = "SDK gRPC Server has already started."
	ALREADY_STOPPED = "SDK gRPC Server already stopped."
)

// Pick out the port from a full address
func GetPort(addr string) string {
	automata := regexp.MustCompile(`.*?:.*?:(.*)`)
	groups := automata.FindStringSubmatch(addr)

	if groups == nil || len(groups) != 2 {
		log.Fatal("Failed to parse SDK address", "addr", addr)
	}
	return groups[1]
}

// Start a new SDK Server
func NewServer(addr string, sdkServer pb.SDKServer) (*Server, error) {

	listener, err := net.Listen("tcp", ":"+GetPort(addr))
	if err != nil {
		return nil, fmt.Errorf("Failed to start tcp listener on port %s err=%v", addr, err)
	}

	server := grpc.NewServer()
	// TODO: Configure for TLS here
	// TODO: Service registration should be configurable
	pb.RegisterSDKServer(server, sdkServer)
	return &Server{
		listener: listener,
		server:   server,
		logger:   log.GetLogger(),
	}, nil
}

// Is the server still running
func (server *Server) IsRunning() bool {
	return server.quit != nil
}

func (server *Server) Start() error {
	if server.IsRunning() {
		return errors.New(HAS_STARTED)
	}

	server.quit = make(chan struct{})

	err := server.OnStart()
	if err != nil {
		log.Debug("Server Failed")
		return errors.New("OnStart method returned an error value")
	}

	go func() {
		addr := server.listener.Addr()
		server.server.Serve(server.listener)

		log.Info("SDK Service listening", "Network", addr.Network(), "Name", addr.String())

		select {
		case _, ok := <-server.quit:
			if !ok {
				log.Info("Stopping", "Name", server.String())
				server.server.Stop()
				return
			}
		}
	}()

	return nil
}

func (server *Server) OnStart() error {
	return nil
}

func (server *Server) Stop() error {
	if server.IsRunning() {
		return errors.New(ALREADY_STOPPED)
	}
	server.Quit()
	server.OnStop()
	server.quit = nil
	return nil
}

func (server *Server) OnStop() {
	return
}

func (server *Server) Reset() error {
	if server.IsRunning() {
		server.Stop()
	}

	err := server.Start()
	if err != nil {
		return err
	}
	server.OnReset()
	return nil
}

func (server *Server) OnReset() error {
	return nil
}

func (server *Server) Quit() <-chan struct{} {
	close(server.quit)
	return server.quit
}

func (server *Server) String() string {
	addr := server.listener.Addr()
	return fmt.Sprintf("SDK.gRPC:%s:%s", addr.Network(), addr.String())
}

func (server *Server) SetLogger(logger tlog.Logger) {
	server.logger = logger
}
