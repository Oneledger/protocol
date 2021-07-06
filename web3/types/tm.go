package types

import (
	"errors"
	"math/big"

	"github.com/Oneledger/protocol/action"
	rpcclient "github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ParseLegacyTx is used to parse the signed tx for old OneLedger tx types
func ParseLegacyTx(tmTx tmtypes.Tx) (*action.SignedTx, error) {
	tx := &action.SignedTx{}

	err := serialize.GetSerializer(serialize.NETWORK).Deserialize(tmTx, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

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
func EthBlockFromTendermint(tmClient rpcclient.Client, block *tmtypes.Block, fullTx bool) (*Block, error) {
	// TODO: Implement fullTx
	gasLimit, err := BlockMaxGasFromConsensusParams(tmClient)
	if err != nil {
		return nil, err
	}

	transactions, gasUsed, err := EthTransactionsFromTendermint(tmClient, block.Txs, fullTx)
	if err != nil {
		return nil, err
	}

	// TODO: Add get bloom from storage when bloom filter will be implemented
	bloom := ethtypes.BytesToBloom(make([]byte, 6))

	return FormatBlock(block.Header, block.Size(), block.Hash(), gasLimit, gasUsed, transactions, bloom), nil
}

// EthTransactionsFromTendermint returns a slice of ethereum transaction hashes and the total gas usage from a set of
// tendermint block transactions.
func EthTransactionsFromTendermint(tmClient rpcclient.Client, txs []tmtypes.Tx, fullTx bool) ([]common.Hash, *big.Int, error) {
	transactionHashes := []common.Hash{}
	gasUsed := big.NewInt(0)

	for _, tx := range txs {
		// first parse legacy tx
		lTx, err := ParseLegacyTx(tx)
		if err != nil {
			// means tx is not legacy and we need to check is tx is ethereum
			// TODO: Add ethereum tx check when it will be released
			continue
		}
		// TODO: Remove gas usage calculation if saving gasUsed per block
		gasUsed.Add(gasUsed, big.NewInt(int64(lTx.Fee.Gas)))
		// TODO: Add full tx handle
		transactionHashes = append(transactionHashes, common.BytesToHash(tx.Hash()))
	}

	return transactionHashes, gasUsed, nil
}

// FormatBlock creates an ethereum block from a tendermint header and ethereum-formatted
// transactions.
func FormatBlock(
	header tmtypes.Header, size int, curBlockHash tmbytes.HexBytes, gasLimit int64,
	gasUsed *big.Int, transactions interface{}, bloom ethtypes.Bloom,
) *Block {
	if len(header.DataHash) == 0 {
		header.DataHash = tmbytes.HexBytes(common.Hash{}.Bytes())
	}

	return &Block{
		Number:           hexutil.Uint64(header.Height),
		Hash:             hexutil.Bytes(curBlockHash),
		ParentHash:       hexutil.Bytes(header.LastBlockID.Hash),
		Nonce:            hexutil.Uint64(0), // PoW specific
		Sha3Uncles:       common.Hash{},     // No uncles in Tendermint
		LogsBloom:        bloom,
		TransactionsRoot: hexutil.Bytes(header.DataHash),
		StateRoot:        hexutil.Bytes(header.AppHash),
		Miner:            common.BytesToAddress(header.ProposerAddress.Bytes()),
		MixHash:          common.Hash{},
		Difficulty:       0,
		TotalDifficulty:  0,
		ExtraData:        hexutil.Uint64(0),
		Size:             hexutil.Uint64(size),
		GasLimit:         hexutil.Uint64(gasLimit), // Static gas limit
		GasUsed:          (*hexutil.Big)(gasUsed),
		Timestamp:        hexutil.Uint64(header.Time.Unix()),
		Transactions:     transactions.([]common.Hash),
		Uncles:           make([]common.Hash, 0),
		ReceiptsRoot:     common.Hash{},
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
