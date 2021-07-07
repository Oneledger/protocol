package eth

import (
	"math/big"

	rpctypes "github.com/Oneledger/protocol/web3/types"
	rpcutils "github.com/Oneledger/protocol/web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
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

	ethAcc, _ := svc.ctx.GetAccountKeeper().GetVersionedAccount(blockNum, address.Bytes())
	txLen := ethAcc.Sequence

	// for pending
	if height == -2 {
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
func (svc *Service) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
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
		return nil, nil
	}

	resBlock, err := svc.getTMClient().Block(&resTx.Height)
	if err != nil {
		svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "block not found")
		return nil, err
	}

	txIndex := hexutil.Uint64(resTx.Index)
	tx, err := rpctypes.LegacyRawBlockAndTxToEthTx(resBlock.Block, &resTx.Tx, (*big.Int)(&chainID), &txIndex)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := uint64(resTx.TxResult.GasUsed)
	if tx.TransactionIndex != nil && int(*tx.TransactionIndex) != 0 {
		cumulativeGasUsed += rpctypes.GetBlockCumulativeGas(resBlock.Block, int(*tx.TransactionIndex))
	}

	// Set status codes based on tx result
	status := hexutil.Uint(resTx.TxResult.Code)

	// TODO: Implement this
	logs := []*ethtypes.Log{}
	bloom := ethtypes.BytesToBloom(make([]byte, 6))

	// TODO: Add handle if tx type is smart contract
	contractAddress := common.Address{}
	if tx.To == nil {
		contractAddress = common.Address{}
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"logsBloom":         bloom,
		"logs":              logs,

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": tx.Hash,
		"contractAddress": contractAddress,
		"gasUsed":         hexutil.Uint64(resTx.TxResult.GasUsed),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        tx.BlockHash,
		"blockNumber":      tx.BlockNumber,
		"transactionIndex": tx.TransactionIndex,

		// sender and receiver (contract or EOA) addresses
		"from": tx.From,
		"to":   tx.To,
	}

	return receipt, nil
}
