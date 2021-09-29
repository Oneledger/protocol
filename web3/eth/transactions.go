package eth

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/utils"
	"github.com/Oneledger/protocol/vm"
	rpctypes "github.com/Oneledger/protocol/web3/types"
	rpcutils "github.com/Oneledger/protocol/web3/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	abci "github.com/tendermint/tendermint/abci/types"
	mempl "github.com/tendermint/tendermint/mempool"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
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
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.GetBlockStore(), blockNrOrHash)
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
		txLen += rpcutils.GetPendingTxCountByAddress(svc.GetMempool(), address)
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
	tx, err := tmrpccore.Tx(nil, hash.Bytes(), false)
	if err != nil {
		// Try to get pending
		pendingTx, err := rpcutils.GetPendingTx(svc.GetMempool(), hash, (*big.Int)(&chainID))
		if err != nil {
			svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "tx not found")
			return nil, nil
		}
		return pendingTx, nil
	}

	block := svc.GetBlockStore().LoadBlock(tx.Height)
	if block == nil {
		svc.logger.Debug("eth_getTransactionByHash", "hash", hash, "block not found")
		return nil, err
	}

	txIndex := hexutil.Uint64(tx.Index)
	return rpctypes.LegacyRawBlockAndTxToEthTx(block, &tx.Tx, (*big.Int)(&chainID), &txIndex)
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (svc *Service) GetTransactionReceipt(hash common.Hash) (*rpctypes.TransactionReceipt, error) {
	svc.logger.Debug("eth_getTransactionReceipt", "hash", hash)
	tx, err := tmrpccore.Tx(nil, hash.Bytes(), false)
	if err != nil {
		return nil, nil
	}

	block := svc.GetBlockStore().LoadBlock(tx.Height)
	if block == nil {
		svc.logger.Debug("eth_getTransactionReceipt", "hash", hash, "block not found")
		return nil, nil
	}

	chainID := utils.HashToBigInt(block.ChainID)
	txIndex := hexutil.Uint64(tx.Index)
	oneTx, err := rpctypes.LegacyRawBlockAndTxToEthTx(block, &tx.Tx, chainID, &txIndex)
	if err != nil {
		return nil, err
	}

	cumulativeGasUsed := uint64(tx.TxResult.GasUsed)
	if tx.Index != 0 {
		cumulativeGasUsed += rpctypes.GetBlockCumulativeGas(block, int(tx.Index))
	}

	var (
		contractAddress *common.Address
		bloom           = vm.BytesToBloom(make([]byte, 6))
		logs            = make([]*ethtypes.Log, 0)
	)
	// Set status codes based on tx result
	status := ethtypes.ReceiptStatusSuccessful
	if tx.TxResult.GetCode() == 1 {
		status = ethtypes.ReceiptStatusFailed
	} else {
		logReceipt := rpctypes.GetTxEthLogs(&tx.TxResult, tx.Index)
		status = logReceipt.Status

		if status == ethtypes.ReceiptStatusSuccessful {
			if oneTx.To == nil {
				contractAddress = logReceipt.ContractAddress
			}
			logs = logReceipt.Logs
			bloom = logReceipt.Bloom
		}
	}

	receipt := &rpctypes.TransactionReceipt{
		Status:            hexutil.Uint64(status),
		CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
		LogsBloom:         bloom,
		Logs:              logs,
		TransactionHash:   oneTx.Hash,
		ContractAddress:   contractAddress,
		GasUsed:           hexutil.Uint64(tx.TxResult.GasUsed),
		BlockHash:         *oneTx.BlockHash,
		BlockNumber:       *oneTx.BlockNumber,
		TransactionIndex:  *oneTx.TransactionIndex,
		From:              oneTx.From,
		To:                oneTx.To,
	}

	return receipt, nil
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by block hash and index.
func (svc *Service) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint64) (*rpctypes.Transaction, error) {
	svc.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash, "idx", idx)
	block := svc.GetBlockStore().LoadBlockByHash(hash.Bytes())
	if block == nil {
		svc.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash, "idx", idx, "block not found")
		return nil, nil
	}
	return svc.getTransactionByBlockAndIndex(block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (svc *Service) GetTransactionByBlockNumberAndIndex(blockNrOrHash rpc.BlockNumberOrHash, idx hexutil.Uint64) (*rpctypes.Transaction, error) {
	height, err := rpctypes.StateAndHeaderByNumberOrHash(svc.GetBlockStore(), blockNrOrHash)
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

		block := svc.GetBlockStore().LoadBlock(blockNum)
		if block == nil {
			svc.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "height", blockNum, "idx", idx, "block not found with height")
			return nil, nil
		}

		txs := svc.GetMempool().ReapMaxTxs(50)

		// return if index out of bounds
		if uint64(idx) >= uint64(len(txs)) {
			return nil, nil
		}
		chainID := utils.HashToBigInt(block.ChainID)

		return rpctypes.LegacyRawBlockAndTxToEthTx(block, &txs[idx], chainID, &idx)
	case rpctypes.EarliestBlockNumber:
		blockNum = rpctypes.InitialBlockNumber
	case rpctypes.LatestBlockNumber:
		blockNum = svc.getState().Version()
	default:
		blockNum = height
	}

	block := svc.GetBlockStore().LoadBlock(blockNum)
	if block == nil {
		return nil, nil
	}
	return svc.getTransactionByBlockAndIndex(block, idx)
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
	config := svc.ctx.GetConfig()

	signedTx, err := rpcutils.EthToOLSignedTx(ethTx)
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
	mempool := svc.GetMempool()
	txHash := tmtypes.Tx(packet).Hash()
	if config.Node.UseAsync {
		svc.logger.Info("Use async mode to propagate tx", txHash)
		err = mempool.CheckTx(packet, nil, mempl.TxInfo{})
		if err != nil {
			return common.Hash{}, err
		}
	} else {
		svc.logger.Info("Use sync mode to propagate tx", txHash)
		resCh := make(chan *abci.Response, 1)
		err := mempool.CheckTx(packet, func(res *abci.Response) {
			resCh <- res
		}, mempl.TxInfo{})
		if err != nil {
			return common.Hash{}, deserializeError(err.Error())
		}
		res := <-resCh
		resBrodTx := res.GetCheckTx()
		if resBrodTx.Code != 0 {
			return common.Hash{}, deserializeError(resBrodTx.Log)
		}
	}
	return common.BytesToHash(txHash), nil
}

func deserializeError(errMsg string) error {
	unpackedData := struct {
		Msg string `json:"msg"`
	}{}
	err := jsonSerializer.Deserialize([]byte(errMsg), &unpackedData)
	if err != nil {
		return fmt.Errorf(errMsg)
	}
	return fmt.Errorf(unpackedData.Msg)
}
