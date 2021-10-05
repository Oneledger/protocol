package eth

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/action/olvm"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/utils"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	rpcutils "github.com/Oneledger/protocol/web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	jsonSerializer = serialize.GetSerializer(serialize.NETWORK)
)

// TODO: Move to the config
const (
	RPCTxFeeCap = 1 // olt
)

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (svc *Service) GetTransactionCount(address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.GetTMClient(), blockNrOrHash)
	if err != nil {
		return nil, nil
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
		txLen += rpcutils.GetPendingTxCountByAddress(svc.GetTMClient(), address)
	}
	n := hexutil.Uint64(txLen)
	return &n, nil
}

// GetTransactionByHash returns the transaction identified by hash.
func (svc *Service) GetTransactionByHash(hash common.Hash) (*rpctypes.Transaction, error) {
	chainID, err := svc.ChainId()
	if err != nil {
		svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "failed to get chainId")
		return nil, err
	}

	svc.logger.Debug("eth_getTransactionByHash", "hash", hash)
	resTx, err := svc.GetTMClient().Tx(hash.Bytes(), false)
	if err != nil {
		// Try to get pending
		pendingTx, err := rpcutils.GetPendingTx(svc.GetTMClient(), hash, (*big.Int)(&chainID))
		if err != nil {
			svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "tx not found")
			return nil, nil
		}
		return pendingTx, nil
	}

	resBlock, err := svc.GetTMClient().Block(&resTx.Height)
	if err != nil {
		svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "block not found")
		return nil, err
	}

	txIndex := hexutil.Uint64(resTx.Index)
	return rpctypes.LegacyRawBlockAndTxToEthTx(resBlock.Block, &resTx.Tx, (*big.Int)(&chainID), &txIndex)
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (svc *Service) GetTransactionReceipt(hash common.Hash) (*rpctypes.TransactionReceipt, error) {
	svc.logger.Debug("eth_getTransactionReceipt", "hash", hash)
	resTx, err := svc.GetTMClient().Tx(hash.Bytes(), false)
	if err != nil {
		return nil, nil
	}

	resBlock, err := svc.GetTMClient().Block(&resTx.Height)
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

	logs := rpctypes.GetTxEthLogs(&resTx.TxResult)
	if logs == nil {
		logs = make([]*ethtypes.Log, 0)
	}

	// Set status codes based on tx result
	status := ethtypes.ReceiptStatusSuccessful
	if resTx.TxResult.GetCode() == 1 {
		status = ethtypes.ReceiptStatusFailed
	}

	stateDB := svc.GetStateDB()

	bloom := stateDB.GetBlockBloom(uint64(resBlock.Block.Height))

	var contractAddress *common.Address
	if tx.To == nil {
		contractAddress = new(common.Address)
		*contractAddress = crypto.CreateAddress(tx.From, uint64(tx.Nonce))
	}

	receipt := &rpctypes.TransactionReceipt{
		Status:            hexutil.Uint64(status),
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
	svc.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash, "idx", idx)
	resBlock, err := svc.GetTMClient().BlockByHash(hash.Bytes())
	if err != nil {
		svc.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash, "idx", idx, "block not found")
		return nil, err
	}
	return svc.getTransactionByBlockAndIndex(resBlock.Block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (svc *Service) GetTransactionByBlockNumberAndIndex(blockNrOrHash rpc.BlockNumberOrHash, idx hexutil.Uint64) (*rpctypes.Transaction, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.GetTMClient(), blockNrOrHash)
	if err != nil {
		return nil, nil
	}
	svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", height, "idx", idx)

	var (
		blockNum int64
	)

	switch height {
	case rpctypes.PendingBlockNumber:
		blockNum = svc.getStateHeight(height)
		svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", blockNum, "idx", idx, "for pending txs")

		result, err := svc.GetTMClient().Block(&blockNum)
		if err != nil {
			return nil, err
		}
		if result.Block == nil {
			svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", blockNum, "idx", idx, "block not found with height")
			return nil, nil
		}

		unconfirmed, err := svc.GetTMClient().UnconfirmedTxs(1000)
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
	case rpctypes.EarliestBlockNumber:
		blockNum = rpctypes.InitialBlockNumber
	case rpctypes.LatestBlockNumber:
		blockNum = svc.getState().Version()
	default:
		blockNum = height
	}

	result, err := svc.GetTMClient().Block(&blockNum)
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

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (svc *Service) SendRawTransaction(input hexutil.Bytes) (common.Hash, error) {
	tx := new(ethtypes.Transaction)
	if err := tx.UnmarshalBinary(input); err != nil {
		return common.Hash{}, err
	}
	return svc.submitTransaction(tx)
}

// submitTransaction is a helper function that submits tx to tendermint pool and logs a message.
func (svc *Service) submitTransaction(tx *ethtypes.Transaction) (common.Hash, error) {
	// If the transaction fee cap is already specified, ensure the
	// fee of the given transaction is _reasonable_.
	if err := rpcutils.CheckTxFee(tx.GasPrice(), tx.Gas(), RPCTxFeeCap); err != nil {
		return common.Hash{}, err
	}

	if tx.Type() != ethtypes.LegacyTxType {
		return common.Hash{}, errors.New("only legacy transactions allowed over RPC")
	}

	if !tx.Protected() {
		// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
		return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
	}

	// Print a log with full tx details for manual investigations and interventions
	resBlock, err := svc.GetTMClient().Block(nil)
	if err != nil {
		return common.Hash{}, err
	}
	if resBlock.Block == nil {
		return common.Hash{}, err
	}
	signer := ethtypes.NewEIP155Signer(tx.ChainId())
	from, err := ethtypes.Sender(signer, tx)
	if err != nil {
		return common.Hash{}, err
	}

	chainID, err := svc.ChainId()
	if err != nil {
		return common.Hash{}, err
	}

	if signer.ChainID().Cmp(chainID.ToInt()) != 0 {
		return common.Hash{}, ethtypes.ErrInvalidChainId
	}

	txHash, err := svc.sendTx(tx)
	if err != nil {
		return common.Hash{}, err
	}

	if tx.To() == nil {
		addr := crypto.CreateAddress(from, tx.Nonce())
		svc.logger.Info("Submitted contract creation", "hash", txHash.Hex(), "from", from, "nonce", tx.Nonce(), "contract", addr.Hex(), "value", tx.Value())
	} else {
		svc.logger.Info("Submitted transaction", "hash", txHash.Hex(), "from", from, "nonce", tx.Nonce(), "recipient", tx.To(), "value", tx.Value())
	}
	return txHash, nil
}

// sendTx directly to the pool, all validation steps was done before
func (svc *Service) sendTx(ethTx *ethtypes.Transaction) (common.Hash, error) {
	var (
		resBrodTx *coretypes.ResultBroadcastTx
		err       error
	)
	config := svc.ctx.GetConfig()

	signedTx, err := rpcutils.EthToOLSignedTx(ethTx)
	if err != nil {
		return common.Hash{}, err
	}

	olvmTx := &olvm.Transaction{}
	err = olvmTx.Unmarshal(signedTx.RawTx.Data)
	if err != nil {
		return common.Hash{}, err
	}

	err = olvmTx.ValidateEthTx(svc.ctx.GetAccountKeeper(), ethTx, svc.GasPrice().ToInt(), true)
	if err != nil {
		return common.Hash{}, err
	}

	_, err = olvmTx.CheckIntrinsicGas(ethTx)
	if err != nil {
		return common.Hash{}, err
	}

	packet, err := jsonSerializer.Serialize(signedTx)
	if err != nil {
		return common.Hash{}, err
	}

	// support two scenarios
	// - for fast fullnode send throughtput (async)
	// - for web3 wallets to show an errors from validation (sync)
	if config.Node.UseAsync {
		svc.logger.Info("Use async mode to propagate tx", tmtypes.Tx(packet).Hash())
		resBrodTx, err = svc.GetTMClient().BroadcastTxAsync(packet)
		if err != nil {
			return common.Hash{}, err
		}
	} else {
		svc.logger.Info("Use sync mode to propagate tx", tmtypes.Tx(packet).Hash())
		resBrodTx, err = svc.GetTMClient().BroadcastTxSync(packet)
		if err != nil {
			return common.Hash{}, err
		}
		if resBrodTx.Code != 0 {
			unpackedData := struct {
				Msg string `json:"msg"`
			}{}

			err = jsonSerializer.Deserialize([]byte(resBrodTx.Log), &unpackedData)
			if err == nil {
				return common.Hash{}, fmt.Errorf(unpackedData.Msg)
			}
			return common.Hash{}, fmt.Errorf(resBrodTx.Log)
		}
	}
	return common.BytesToHash(resBrodTx.Hash), nil
}
