package eth

import (
	"github.com/Oneledger/protocol/web3rpc"
	"math/big"

	"github.com/Oneledger/protocol/utils"
	"github.com/Oneledger/protocol/version"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProtocolVersion returns the supported Ethereum protocol version.
func (svc *web3rpc.Service) ProtocolVersion() string {
	svc.logger.Debug("eth_protocolVersion")
	return version.Protocol.String()
}

// ChainId returns the chain's identifier in hex format
func (svc *web3rpc.Service) ChainId() (hexutil.Big, error) {
	svc.logger.Debug("eth_chainId")
	blockResult, err := svc.getTMClient().Block(nil)
	if err != nil {
		return hexutil.Big(*big.NewInt(0)), err
	}
	chainID := utils.HashToBigInt(blockResult.Block.Header.ChainID)
	return hexutil.Big(*chainID), nil
}

// Syncing returns whether or not the current node is syncing with other peers. Returns false if not, or a struct
// outlining the state of the sync if it is.
func (svc *web3rpc.Service) Syncing() (interface{}, error) {
	svc.logger.Debug("eth_syncing")

	status, err := svc.getTMClient().Status()
	if err != nil {
		return false, err
	}

	if !status.SyncInfo.CatchingUp {
		return false, nil
	}

	return map[string]interface{}{
		// "startingBlock": nil, // NA
		"currentBlock": hexutil.Uint64(status.SyncInfo.LatestBlockHeight),
		// "highestBlock":  nil, // NA
		// "pulledStates":  nil, // NA
		// "knownStates":   nil, // NA
	}, nil
}

// Coinbase is the address that staking rewards will be send to (alias for Etherbase).
func (svc *web3rpc.Service) Coinbase() (common.Address, error) {
	svc.logger.Debug("eth_coinbase")

	pubKeyHander, err := svc.ctx.GetNodeContext().PubKey().GetHandler()
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(pubKeyHander.Bytes()), nil
}

// Mining returns whether or not this node is currently mining. Always false.
func (svc *web3rpc.Service) Mining() bool {
	svc.logger.Debug("eth_mining")
	// TODO: Add info in according if this node is validator and it is active
	return false
}

// Hashrate returns the current node's hashrate. Always 0.
func (svc *web3rpc.Service) Hashrate() hexutil.Uint64 {
	svc.logger.Debug("eth_hashrate")
	return 0
}

// GasPrice returns the current gas price
func (svc *web3rpc.Service) GasPrice() *hexutil.Big {
	svc.logger.Debug("eth_gasPrice")
	// NOTE: Maybe we could add the mean price in according to previous blocks
	// but not really need it right now
	out := big.NewInt(0)
	return (*hexutil.Big)(out)
}
