package eth

import (
	"math/big"

	rpctypes "github.com/Oneledger/protocol/web3/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
		txLen = ethAcc.Sequence + rpctypes.GetPendingTxCountByAddress(svc.getTMClient(), address)
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
		pendingTx, err := rpctypes.GetPendingTx(svc.getTMClient(), hash, (*big.Int)(&chainID))
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
