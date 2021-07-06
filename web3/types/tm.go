package types

import (
	"errors"
	"math/big"

	rpcclient "github.com/Oneledger/protocol/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
)

// EthHeaderFromTendermint is an util function that returns an Ethereum Header
// from a tendermint Header.
func EthHeaderFromTendermint(header tmtypes.Header) *ethtypes.Header {
	return &ethtypes.Header{
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash.Bytes()),
		UncleHash:   common.Hash{},
		Coinbase:    common.BytesToAddress(header.ProposerAddress.Bytes()),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      common.BytesToHash(header.DataHash),
		ReceiptHash: common.Hash{},
		Difficulty:  nil,
		Number:      big.NewInt(header.Height),
		Time:        uint64(header.Time.Unix()),
		Extra:       nil,
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
	}
}

// // RawTxToEthTx returns a evm MsgEthereum transaction from raw tx bytes.
// func RawTxToEthTx(bz []byte) (*evmtypes.MsgEthereumTx, error) {
// 	retu
// }

// BlockMaxGasFromConsensusParams returns the gas limit for the latest block from the chain consensus params.
func BlockMaxGasFromConsensusParams(tmClient rpcclient.Client) (int64, error) {
	resConsParams, err := tmClient.ConsensusParams(nil)
	if err != nil {
		return 0, err
	}

	gasLimit := resConsParams.ConsensusParams.Block.MaxGas
	if gasLimit == -1 {
		// Sets gas limit to max uint32 to not error with javascript dev tooling
		// This -1 value indicating no block gas limit is set to max uint64 with geth hexutils
		// which errors certain javascript dev tooling which only supports up to 53 bits
		gasLimit = int64(^uint32(0))
	}

	return gasLimit, nil
}

// EthBlockFromTendermint returns a JSON-RPC compatible Ethereum blockfrom a given Tendermint block.
func EthBlockFromTendermint(tmClient rpcclient.Client, block *tmtypes.Block, fullTx bool) (map[string]interface{}, error) {
	// TODO: Implement fullTx
	gasLimit, err := BlockMaxGasFromConsensusParams(tmClient)
	if err != nil {
		return nil, err
	}

	transactions, gasUsed, err := EthTransactionsFromTendermint(tmClient, block.Txs)
	if err != nil {
		return nil, err
	}

	// res, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", evmtypes.ModuleName, evmtypes.QueryBloom, block.Height))
	// if err != nil {
	// 	return nil, err
	// }

	// var bloomRes evmtypes.QueryBloomFilter
	// clientCtx.Codec.MustUnmarshalJSON(res, &bloomRes)

	// bloom := bloomRes.Bloom
	// TODO: Add bloom generation
	bloom := ethtypes.BytesToBloom(make([]byte, 6))

	return FormatBlock(block.Header, block.Size(), block.Hash(), gasLimit, gasUsed, transactions, bloom), nil
}

// EthTransactionsFromTendermint returns a slice of ethereum transaction hashes and the total gas usage from a set of
// tendermint block transactions.
func EthTransactionsFromTendermint(tmClient rpcclient.Client, txs []tmtypes.Tx) ([]common.Hash, *big.Int, error) {
	transactionHashes := []common.Hash{}
	gasUsed := big.NewInt(0)

	// for _, tx := range txs {
	// 	// ethTx, err := RawTxToEthTx(tx)
	// 	// if err != nil {
	// 	// 	// continue to next transaction in case it's not a MsgEthereumTx
	// 	// 	continue
	// 	// }
	// 	// TODO: Remove gas usage calculation if saving gasUsed per block
	// 	// gasUsed.Add(gasUsed, big.NewInt(int64(ethTx.GetGas())))
	// 	// transactionHashes = append(transactionHashes, common.BytesToHash(tx.Hash()))
	// }

	return transactionHashes, gasUsed, nil
}

// FormatBlock creates an ethereum block from a tendermint header and ethereum-formatted
// transactions.
func FormatBlock(
	header tmtypes.Header, size int, curBlockHash tmbytes.HexBytes, gasLimit int64,
	gasUsed *big.Int, transactions interface{}, bloom ethtypes.Bloom,
) map[string]interface{} {
	if len(header.DataHash) == 0 {
		header.DataHash = tmbytes.HexBytes(common.Hash{}.Bytes())
	}

	return map[string]interface{}{
		"number":           hexutil.Uint64(header.Height),
		"hash":             hexutil.Bytes(curBlockHash),
		"parentHash":       hexutil.Bytes(header.LastBlockID.Hash),
		"nonce":            hexutil.Uint64(0), // PoW specific
		"sha3Uncles":       common.Hash{},     // No uncles in Tendermint
		"logsBloom":        bloom,
		"transactionsRoot": hexutil.Bytes(header.DataHash),
		"stateRoot":        hexutil.Bytes(header.AppHash),
		"miner":            common.BytesToAddress(header.ProposerAddress.Bytes()),
		"mixHash":          common.Hash{},
		"difficulty":       0,
		"totalDifficulty":  0,
		"extraData":        hexutil.Uint64(0),
		"size":             hexutil.Uint64(size),
		"gasLimit":         hexutil.Uint64(gasLimit), // Static gas limit
		"gasUsed":          (*hexutil.Big)(gasUsed),
		"timestamp":        hexutil.Uint64(header.Time.Unix()),
		"transactions":     transactions.([]common.Hash),
		"uncles":           []string{},
		"receiptsRoot":     common.Hash{},
	}
}

func StateAndHeaderByNumberOrHash(tmClient rpcclient.Client, blockNrOrHash rpc.BlockNumberOrHash) (int64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return blockNr.Int64(), nil
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := tmClient.BlockByHash(hash.Bytes())
		if err != nil {
			return 0, err
		}
		if header == nil || header.Block == nil {
			return 0, errors.New("header for hash not found")
		}
		return header.Block.Header.Height, nil
	}
	return 0, errors.New("invalid arguments; neither block nor hash specified")
}
