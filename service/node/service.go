package node

import (
	"fmt"
	"net/url"

	codes "github.com/Oneledger/protocol/status_codes"

	"github.com/Oneledger/protocol/log"

	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/tendermint/tendermint/p2p"
)

type Service struct {
	ctx    node.Context
	cfg    *config.Server
	logger *log.Logger
}

func NewService(ctx node.Context, cfg *config.Server, logger *log.Logger) *Service {
	return &Service{
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
	}
}

func Name() string {
	return "node"
}

func (svc *Service) NodeName(_ client.NodeNameRequest, reply *client.NodeNameReply) error {
	*reply = svc.ctx.NodeName
	return nil
}

func (svc *Service) Address(_ client.NodeAddressRequest, reply *client.NodeAddressReply) error {
	*reply = svc.ctx.Address()
	return nil
}

func (svc *Service) ID(req client.NodeIDRequest, reply *client.NodeIDReply) error {
	nodeKey, err := p2p.LoadNodeKey(svc.cfg.TMConfig().NodeKeyFile())
	if err != nil {
		return codes.ErrLoadingNodeKey
	}

	// silenced error because not present means false
	ip := p2pAddressFromCFG(svc.cfg)
	if req.ShouldShowIP {
		u, err := url.Parse(ip)
		if err != nil {
			return codes.ErrParsingAddress
		}
		out := fmt.Sprintf("%s@%s", nodeKey.ID(), u.Host)
		*reply = out
	} else {
		*reply = string(nodeKey.ID())
	}
	return nil
}

// This function returns the external p2p address if it exists, but falls back to the regular p2p address if it is
// not present from the config
func p2pAddressFromCFG(cfg *config.Server) string {
	extP2P := cfg.Network.ExternalP2PAddress
	if extP2P != "" {
		return cfg.Network.P2PAddress
	}

	return cfg.Network.ExternalP2PAddress
}
