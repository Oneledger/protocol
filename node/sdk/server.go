package sdk

import (
	"errors"
	"fmt"
	"net"

	olog "github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
)

// implements common.Service
type Server struct {
	listener net.Listener
	server   *grpc.Server
	quit     chan struct{}
	logger   log.Logger
}

// Ensure sdk.Server implements common.Service
var _ common.Service = &Server{}

// Error messages
const (
	HAS_STARTED     = "SDK gRPC Server has already started."
	ALREADY_STOPPED = "SDK gRPC Server already stopped."
)

func NewServer(port int, sdkServer pb.SDKServer) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("Failed to start tcp listener on port :%d, %v", port, err)
	}

	server := grpc.NewServer()
	// TODO: Configure for TLS here
	// TODO: Service registration should be configurable
	pb.RegisterSDKServer(server, sdkServer)
	return &Server{
		listener: listener,
		server:   server,
		logger:   olog.GetLogger(),
	}, nil
}

func (s *Server) IsRunning() bool {
	return s.quit != nil
}

func (s *Server) Start() error {
	if s.IsRunning() {
		return errors.New(HAS_STARTED)
	}
	s.quit = make(chan struct{})
	err := s.OnStart()
	if err != nil {
		return errors.New("OnStart method returned an error value")
	}

	go func() {
		addr := s.listener.Addr()
		s.logger.Info(fmt.Sprintf("SDK Service listening on %s %s", addr.Network(), addr.String()))
		s.server.Serve(s.listener)
		select {
		case _, ok := <-s.quit:
			if !ok {
				s.logger.Info("Stopping %s", s.String())
				s.server.Stop()
				return
			}
		}
	}()

	return nil
}

func (s *Server) OnStart() error {
	return nil
}

func (s *Server) Stop() error {
	if s.IsRunning() {
		return errors.New(ALREADY_STOPPED)
	}
	s.Quit()
	s.OnStop()
	s.quit = nil
	return nil
}

func (s *Server) OnStop() {
	return
}

func (s *Server) Reset() error {
	if s.IsRunning() {
		s.Stop()
	}

	err := s.Start()
	if err != nil {
		return err
	}
	s.OnReset()
	return nil
}

func (s *Server) OnReset() error {
	return nil
}

func (s *Server) Quit() <-chan struct{} {
	close(s.quit)
	return s.quit
}

func (s *Server) String() string {
	a := s.listener.Addr()
	return fmt.Sprintf("SDK.gRPC:%s:%s", a.Network(), a.String())
}

func (s *Server) SetLogger(l log.Logger) {
	s.logger = l
}
