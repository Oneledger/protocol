package eth

import (
	"math/big"

	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	tmtypes "github.com/tendermint/tendermint/types"
)

// BlockNumber returns the current block number.
func (svc *Service) BlockNumber() hexutil.Big {
	svc.logger.Debug("eth_blockNumber")
	height := svc.getState().Version()
	blockNumber := big.NewInt(height)
	return hexutil.Big(*blockNumber)
}

func (svc *Service) blockWithBloom(tmBlock *tmtypes.Block, fullTx bool) (*rpctypes.Block, error) {
	block, err := rpctypes.EthBlockFromTendermint(tmBlock, fullTx)
	if err != nil {
		return nil, err
	}
	block.LogsBloom = svc.GetStateDB().GetBlockBloom(uint64(block.Number))
	return block, nil
}

// GetBlockByHash returns the block identified by hash.
func (svc *Service) GetBlockByHash(hash common.Hash, fullTx bool) (*rpctypes.Block, error) {
	svc.logger.Debug("eth_getBlockByHash", "hash", hash, "fullTx", fullTx)

	result, err := svc.GetTMClient().BlockByHash(hash.Bytes())
	if err != nil {
		svc.logger.Debug("eth_getBlockByHash", "block err", err)
		return nil, nil
	}
	if result.Block == nil {
		svc.logger.Debug("eth_getBlockByHash", "block not found with hash", common.Bytes2Hex(hash.Bytes()))
		return nil, nil
	}
	return svc.blockWithBloom(result.Block, fullTx)
}

// GetBlockByNumber returns the block identified by number.
func (svc *Service) GetBlockByNumber(blockNrOrHash rpc.BlockNumberOrHash, fullTx bool) (*rpctypes.Block, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.GetTMClient(), blockNrOrHash)
	if err != nil {
		svc.logger.Debug("eth_getBlockByNumber", "block err", err)
		return nil, nil
	}
	svc.logger.Debug("eth_getBlockByNumber", "height", height, "fullTx", fullTx)

	var blockNum int64
	switch height {
	case rpctypes.LatestBlockNumber, rpctypes.PendingBlockNumber:
		blockNum = svc.getState().Version()
	case rpctypes.EarliestBlockNumber:
		blockNum = rpctypes.InitialBlockNumber
	default:
		blockNum = height
	}
	result, err := svc.GetTMClient().Block(&blockNum)
	if err != nil {
		svc.logger.Debug("eth_getBlockByNumber", "block err", err)
		return nil, nil
	}
	if result.Block == nil {
		svc.logger.Debug("eth_getBlockByNumber", "block not found with height", blockNum)
		return nil, nil
	}
	return svc.blockWithBloom(result.Block, fullTx)
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (svc *Service) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	svc.logger.Debug("eth_getBlockTransactionCountByHash", "hash", hash)

	result, err := svc.GetTMClient().BlockByHash(hash.Bytes())
	if err != nil {
		return nil
	}
	if result.Block == nil {
		svc.logger.Debug("eth_getBlockTransactionCountByHash", "block not found with hash", common.Bytes2Hex(hash.Bytes()))
		return nil
	}
	n := hexutil.Uint(len(result.Block.Txs))
	return &n
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by its height.
func (svc *Service) GetBlockTransactionCountByNumber(blockNrOrHash rpc.BlockNumberOrHash) *hexutil.Uint {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.GetTMClient(), blockNrOrHash)
	if err != nil {
		return nil
	}

	var (
		blockNum int64
		txsLen   int
	)

	switch height {
	case rpctypes.PendingBlockNumber:
		blockNum = svc.getState().Version()
		svc.logger.Debug("eth_getBlockTransactionCountByNumber", "height", blockNum, "for pending txs")

		result, err := svc.GetTMClient().Block(&blockNum)
		if err != nil {
			return nil
		}
		if result.Block == nil {
			svc.logger.Debug("eth_getBlockTransactionCountByNumber", "block not found with height", blockNum)
			return nil
		}

		unconfirmed, err := svc.GetTMClient().UnconfirmedTxs(1000)
		if err != nil {
			svc.logger.Debug("eth_getBlockTransactionCountByNumber", "failed to get unconfirmed txs", err)
			return nil
		}
		txsLen = len(result.Block.Txs) + len(unconfirmed.Txs)
	case rpctypes.LatestBlockNumber:
		blockNum = svc.getState().Version()
		svc.logger.Debug("eth_getBlockTransactionCountByNumber", "height", blockNum, "for last txs")

		result, err := svc.GetTMClient().Block(&blockNum)
		if err != nil {
			return nil
		}
		if result.Block == nil {
			svc.logger.Debug("eth_getBlockTransactionCountByNumber", "block not found with height", blockNum)
			return nil
		}
		txsLen = len(result.Block.Txs)
	case rpctypes.EarliestBlockNumber:
		blockNum = rpctypes.InitialBlockNumber
		svc.logger.Debug("eth_getBlockTransactionCountByNumber", "height", blockNum, "for last txs")

		result, err := svc.GetTMClient().Block(&blockNum)
		if err != nil {
			return nil
		}
		if result.Block == nil {
			svc.logger.Debug("eth_getBlockTransactionCountByNumber", "block not found with height", blockNum)
			return nil
		}
		txsLen = len(result.Block.Txs)
	default:
		blockNum = height
		svc.logger.Debug("eth_getBlockTransactionCountByNumber", "height", blockNum)

		result, err := svc.GetTMClient().Block(&blockNum)
		if err != nil {
			return nil
		}
		if result.Block == nil {
			svc.logger.Debug("eth_getBlockTransactionCountByNumber", "block not found with height", blockNum)
			return nil
		}
		txsLen = len(result.Block.Txs)
	}

	svc.logger.Debug("eth_getBlockTransactionCountByNumber", "height", blockNum, "txsLen", txsLen)
	txCount := hexutil.Uint(txsLen)
	return &txCount
}

// GetUncleCountByBlockHash returns the number of uncles in the block idenfied by hash. Always zero.
func (svc *Service) GetUncleCountByBlockHash(_ common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block idenfied by number. Always zero.
func (svc *Service) GetUncleCountByBlockNumber(_ rpc.BlockNumberOrHash) hexutil.Uint {
	return 0
}

// GetUncleByBlockHashAndIndex returns the uncle identified by hash and index. Always returns nil.
func (svc *Service) GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) *rpctypes.Block {
	return nil
}

// GetUncleByBlockNumberAndIndex returns the uncle identified by number and index. Always returns nil.
func (svc *Service) GetUncleByBlockNumberAndIndex(number hexutil.Uint, idx hexutil.Uint) *rpctypes.Block {
	return nil
}
