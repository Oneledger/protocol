package eth

import (
	"math/big"

	"github.com/Oneledger/protocol/utils"
	"github.com/Oneledger/protocol/version"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProtocolVersion returns the supported Ethereum protocol version.
func (svc *Service) ProtocolVersion() string {
	svc.logger.Debug("eth_protocolVersion")
	return version.Protocol.String()
}

// ChainId returns the chain's identifier in hex format
func (svc *Service) ChainId() (hexutil.Big, error) {
	svc.logger.Debug("eth_chainId")
	genesis := svc.ctx.GetGenesisDoc()
	if genesis == nil {
		return hexutil.Big(*big.NewInt(0)), nil
	}
	chainID := utils.HashToBigInt(genesis.ChainID)
	return hexutil.Big(*chainID), nil
}

// Syncing returns whether or not the current node is syncing with other peers. Returns false if not, or a struct
// outlining the state of the sync if it is.
func (svc *Service) Syncing() (interface{}, error) {
	svc.logger.Debug("eth_syncing")

	if !svc.ctx.GetConsensusReactor().FastSync() {
		return false, nil
	}

	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(1),
		"currentBlock":  hexutil.Uint64(svc.getState().Version()),
		"highestBlock":  nil, // NA
		"pulledStates":  nil, // NA
		"knownStates":   nil, // NA
	}, nil
}

// Coinbase is the address that staking rewards will be send to (alias for Etherbase).
func (svc *Service) Coinbase() (common.Address, error) {
	svc.logger.Debug("eth_coinbase")

	pubKeyHander, err := svc.ctx.GetNodeContext().PubKey().GetHandler()
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(pubKeyHander.Address().Bytes()), nil
}

// Mining returns whether or not this node is currently mining. Always false.
func (svc *Service) Mining() bool {
	svc.logger.Debug("eth_mining")
	// TODO: Add info in according if this node is validator and it is active
	return false
}

// Hashrate returns the current node's hashrate. Always 0.
func (svc *Service) Hashrate() hexutil.Uint64 {
	svc.logger.Debug("eth_hashrate")
	return 0
}

// GasPrice returns the current gas price
func (svc *Service) GasPrice() *hexutil.Big {
	svc.logger.Debug("eth_gasPrice")
	out := svc.ctx.GetFeePool().GetOpt().MinFee().Amount.BigInt()
	return (*hexutil.Big)(out)
}
