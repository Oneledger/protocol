package eth

import (
	"math/big"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/utils"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	rpcutils "github.com/Oneledger/protocol/web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (svc *Service) GetTransactionCount(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		return nil, err
	}

	svc.logger.Debug("eth_getTransactionCount", "address", address, "height", height)

	// getting actual block
	blockNum := svc.getStateHeight(height)

	var txLen uint64
	ethAcc, err := svc.ctx.GetAccountKeeper().GetVersionedAccount(address.Bytes(), blockNum)
	if err == nil {
		txLen = ethAcc.Sequence
	}

	// for pending
	if height == rpctypes.PendingBlockNumber {
		txLen = ethAcc.Sequence + rpcutils.GetPendingTxCountByAddress(svc.getTMClient(), address)
	}
	n := hexutil.Uint64(txLen)
	return &n, nil
}

// GetTransactionByHash returns the transaction identified by hash.
func (svc *Service) GetTransactionByHash(hash common.Hash) (*rpctypes.Transaction, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	chainID, err := svc.ChainId()
	if err != nil {
		svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "failed to get chainId")
		return nil, err
	}

	svc.logger.Debug("eth_getTransactionByHash", "hash", hash)
	resTx, err := svc.getTMClient().Tx(hash.Bytes(), false)
	if err != nil {
		// Try to get pending
		pendingTx, err := rpcutils.GetPendingTx(svc.getTMClient(), hash, (*big.Int)(&chainID))
		if err != nil {
			svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "tx not found")
			return nil, nil
		}
		return pendingTx, nil
	}

	resBlock, err := svc.getTMClient().Block(&resTx.Height)
	if err != nil {
		svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "block not found")
		return nil, err
	}

	txIndex := hexutil.Uint64(resTx.Index)
	return rpctypes.LegacyRawBlockAndTxToEthTx(resBlock.Block, &resTx.Tx, (*big.Int)(&chainID), &txIndex)
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (svc *Service) GetTransactionReceipt(hash common.Hash) (*rpctypes.TransactionReceipt, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.logger.Debug("eth_getTransactionReceipt", "hash", hash)
	resTx, err := svc.getTMClient().Tx(hash.Bytes(), false)
	if err != nil {
		return nil, nil
	}

	resBlock, err := svc.getTMClient().Block(&resTx.Height)
	if err != nil {
		svc.logger.Debug("eth_getTransactionReceipt", "hash", hash, "err", err)
		return nil, err
	}
	if resBlock.Block == nil {
		svc.logger.Debug("eth_getTransactionReceipt", "hash", hash, "block not found")
		return nil, nil
	}

	chainID := utils.HashToBigInt(resBlock.Block.ChainID)
	txIndex := hexutil.Uint64(resTx.Index)
	tx, err := rpctypes.LegacyRawBlockAndTxToEthTx(resBlock.Block, &resTx.Tx, chainID, &txIndex)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := uint64(resTx.TxResult.GasUsed)
	if tx.TransactionIndex != nil && int(*tx.TransactionIndex) != 0 {
		cumulativeGasUsed += rpctypes.GetBlockCumulativeGas(resBlock.Block, int(*tx.TransactionIndex))
	}

	// Set status codes based on tx result
	status := hexutil.Uint64(resTx.TxResult.Code)

	stateDB := action.NewCommitStateDB(svc.ctx.GetContractStore(), svc.ctx.GetAccountKeeper(), svc.logger)

	logs, err := stateDB.GetLogs(hash)
	if err != nil {
		return nil, err
	}
	// TODO: Implement bloom
	bloom := ethtypes.BytesToBloom(make([]byte, 6))

	var contractAddress *common.Address
	if tx.To == nil {
		contractAddress = rpctypes.GetContractAddress(&resTx.TxResult)
	}

	receipt := &rpctypes.TransactionReceipt{
		Status:            status,
		CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
		LogsBloom:         bloom,
		Logs:              logs,
		TransactionHash:   tx.Hash,
		ContractAddress:   contractAddress,
		GasUsed:           hexutil.Uint64(resTx.TxResult.GasUsed),
		BlockHash:         *tx.BlockHash,
		BlockNumber:       *tx.BlockNumber,
		TransactionIndex:  *tx.TransactionIndex,
		From:              tx.From,
		To:                tx.To,
	}

	return receipt, nil
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by block hash and index.
func (svc *Service) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint64) (*rpctypes.Transaction, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	svc.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash, "idx", idx)
	resBlock, err := svc.getTMClient().BlockByHash(hash.Bytes())
	if err != nil {
		svc.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash, "idx", idx, "block not found")
		return nil, err
	}
	return svc.getTransactionByBlockAndIndex(resBlock.Block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (svc *Service) GetTransactionByBlockNumberAndIndex(blockNrOrHash rpc.BlockNumberOrHash, idx hexutil.Uint64) (*rpctypes.Transaction, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.getTMClient(), blockNrOrHash)
	if err != nil {
		return nil, err
	}
	svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", height, "idx", idx)

	var (
		blockNum int64
	)

	switch height {
	case rpctypes.PendingBlockNumber:
		blockNum = svc.getStateHeight(height)
		svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", blockNum, "idx", idx, "for pending txs")

		result, err := svc.getTMClient().Block(&blockNum)
		if err != nil {
			return nil, err
		}
		if result.Block == nil {
			svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", blockNum, "idx", idx, "block not found with height")
			return nil, nil
		}

		unconfirmed, err := svc.getTMClient().UnconfirmedTxs(1000)
		if err != nil {
			svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", blockNum, "idx", idx, "failed to get unconfirmed txs", err)
			return nil, err
		}
		// return if index out of bounds
		if uint64(idx) >= uint64(len(unconfirmed.Txs)) {
			return nil, nil
		}
		chainID := utils.HashToBigInt(result.Block.ChainID)

		return rpctypes.LegacyRawBlockAndTxToEthTx(result.Block, &unconfirmed.Txs[idx], chainID, &idx)
	case rpctypes.LatestBlockNumber:
		blockNum = svc.getStateHeight(height)
	default:
		blockNum = height
	}

	result, err := svc.getTMClient().Block(&blockNum)
	if err != nil {
		return nil, err
	}
	return svc.getTransactionByBlockAndIndex(result.Block, idx)
}

func (svc *Service) getTransactionByBlockAndIndex(block *tmtypes.Block, idx hexutil.Uint64) (*rpctypes.Transaction, error) {
	if block == nil {
		svc.logger.Debug("getTransactionByBlockAndIndex", "block not found")
		return nil, nil
	}

	// return if index out of bounds
	if uint64(idx) >= uint64(len(block.Txs)) {
		return nil, nil
	}

	chainID := utils.HashToBigInt(block.ChainID)
	return rpctypes.LegacyRawBlockAndTxToEthTx(block, &block.Txs[idx], chainID, &idx)
}
