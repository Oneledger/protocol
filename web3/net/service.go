package net

import (
	"fmt"
	"github.com/Oneledger/protocol/service"
	"github.com/Oneledger/protocol/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// PublicNetAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicNetAPI struct {
	networkVersion string
	peerCount      hexutil.Big
}

// NewService creates an instance of the public Net Web3 API.
func NewService(clientCtx *service.Context) *PublicNetAPI {
	peerCount := big.NewInt(int64(len(clientCtx.Cfg.P2P.PersistentPeers)))
	return &PublicNetAPI{
		networkVersion: clientCtx.Cfg.ChainID(),
		peerCount:      hexutil.Big(*peerCount),
	}
}

func (svc *PublicNetAPI) Listening() bool {
	// TODO: how to fetch if the node is open to network connection??
	return true
}

func (svc *PublicNetAPI) PeerCount() hexutil.Big {
	// TODO: how to fetch the current connection peers?
	return svc.peerCount
}

func (svc *PublicNetAPI) Version() string {
	return fmt.Sprintf("%s", utils.HashToBigInt(svc.networkVersion))
}
