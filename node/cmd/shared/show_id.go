package shared

import (
	"fmt"
	"net/url"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/consensus"
	"github.com/Oneledger/protocol/node/global"
	"github.com/tendermint/tendermint/p2p"
)

func ShowNodeID(cfg *config.Server, shouldShowIP bool) (result string, err error) {
	configuration, err := consensus.ParseConfig(global.Current.Config)
	if err != nil {
		return result, err
	}
	nodeKey, err := p2p.LoadNodeKey(configuration.CFG.NodeKeyFile())
	if err != nil {
		return result, err
	}

	ip := configuration.CFG.P2P.ExternalAddress
	if shouldShowIP {
		u, err := url.Parse(ip)
		if err != nil {
			return result, err
		}
		return fmt.Sprintf("%s@%s", nodeKey.ID(), u.Host), nil
	} else {
		return string(nodeKey.ID()), nil
	}
}
