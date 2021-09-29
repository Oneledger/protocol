package types

import (
	"bytes"
	"errors"
	"math/big"
	"strconv"

	"github.com/Oneledger/protocol/action"
	rpcclient "github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/utils"
	"github.com/Oneledger/protocol/vm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	abci "github.com/tendermint/tendermint/abci/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmrpccore "github.com/tendermint/tendermint/rpc/core"
	"github.com/tendermint/tendermint/store"
	tmtypes "github.com/tendermint/tendermint/types"
)

var jsonSerializer serialize.Serializer = serialize.GetSerializer(serialize.NETWORK)

// GetBlockBloom from tm events
func GetBlockBloom(events []abci.Event) vm.Bloom {
	bloom := vm.BytesToBloom(make([]byte, 6))
	for _, evt := range events {
		if evt.Type != "olvm" {
			continue
		}
		for _, attr := range evt.Attributes {
			if bytes.Equal(attr.Key, []byte("block.bloom")) {
				bloom = vm.BytesToBloom(attr.Value)
				return bloom
			}
		}
	}
	return bloom
}

type LogReceipt struct {
	Bloom           vm.Bloom
	Logs            []*ethtypes.Log
	Status          uint64
	ContractAddress *common.Address
}

// GetTxEthLogs substract logs from deliver response
func GetTxEthLogs(res *abci.ResponseDeliverTx, txIndex uint32) *LogReceipt {
	lr := &LogReceipt{
		Bloom:  vm.BytesToBloom(make([]byte, 6)),
		Logs:   make([]*ethtypes.Log, 0),
		Status: ethtypes.ReceiptStatusSuccessful,
	}
	for _, evt := range res.Events {
		for _, attr := range evt.Attributes {
			if bytes.Equal(attr.Key, []byte("tx.contract")) {
				contractAddress := common.BytesToAddress(attr.Value)
				lr.ContractAddress = &contractAddress
			} else if bytes.Contains(attr.Key, []byte("tx.logs")) {
				log, err := new(vm.RLPLog).Decode(attr.Value)
				if err == nil {
					log.TxIndex = uint(txIndex)
					lr.Logs = append(lr.Logs, log)
				}
			} else if bytes.Equal(attr.Key, []byte("tx.bloom")) {
				lr.Bloom = vm.BytesToBloom(attr.Value)
			} else if bytes.Equal(attr.Key, []byte("tx.status")) {
				if status, err := strconv.Atoi(string(attr.Value)); err == nil {
					lr.Status = uint64(status)
				}
			}
		}
	}
	return lr
}

// GetBlockCumulativeGas returns the cumulative gas used on a block up to a given
// transaction index. The returned gas used includes the gas from both the SDK and
// EVM module transactions.
func GetBlockCumulativeGas(block *tmtypes.Block, idx int) uint64 {
	var gasUsed uint64

	for i := 0; i < idx && i < len(block.Txs); i++ {
		tx, err := ParseLegacyTx(block.Txs[i])
		if err != nil {
			continue
		}

		gasUsed += uint64(tx.Fee.Gas)
	}
	return gasUsed
}

// LegacyRawBlockAndTxToEthTx returns a eth Transaction compatible from the legacy tx structure.
func LegacyRawBlockAndTxToEthTx(tmBlock *tmtypes.Block, tmTx *tmtypes.Tx, chainID *big.Int, txIndex *hexutil.Uint64) (*Transaction, error) {
	var (
		to           *common.Address
		value        hexutil.Big   = hexutil.Big(*big.NewInt(0))
		input        hexutil.Bytes = make(hexutil.Bytes, 0)
		unpackedData               = struct {
			ChainID *big.Int       `json:"chainID"`
			To      *keys.Address  `json:"to"`
			Amount  *action.Amount `json:"amount"`
			Data    []byte         `json:"data"`
			Nonce   uint64         `json:"nonce"`
		}{}
		nonce       hexutil.Uint64
		blockNumber *hexutil.Big
		blockHash   *common.Hash
		r           *common.Hash
		s           *common.Hash
		v           *big.Int
	)

	lTx, err := ParseLegacyTx(*tmTx)
	if err != nil {
		return nil, err
	}

	if tmBlock != nil {
		blockHash = new(common.Hash)
		*blockHash = common.BytesToHash(tmBlock.Hash())
		blockNumber = (*hexutil.Big)(big.NewInt(tmBlock.Height))
	}

	from := common.Address{}
	// If signatures found, means that we have a sender, so taking first one
	if len(lTx.Signatures) > 0 {
		actSig := lTx.Signatures[0]
		sig := actSig.Signed
		pubKeyHandler, err := actSig.Signer.GetHandler()
		if err != nil {
			return nil, err
		}
		from = common.BytesToAddress(pubKeyHandler.Address().Bytes())

		var tmpV byte
		if len(sig) == 65 {
			tmpV = sig[len(sig)-1:][0]
		} else {
			// for legacy support
			tmpV = byte(int(sig[0]) % 2)
		}

		r = new(common.Hash)
		*r = common.BytesToHash(sig[:32])

		s = new(common.Hash)
		*s = common.BytesToHash(sig[32:64])

		v = new(big.Int).SetBytes([]byte{tmpV + 27})
	}

	err = jsonSerializer.Deserialize(lTx.Data, &unpackedData)
	if err == nil {
		if unpackedData.To != nil {
			to = new(common.Address)
			*to = common.BytesToAddress(unpackedData.To.Bytes())
		} else if unpackedData.To == nil && unpackedData.ChainID == nil {
			// just set as zero
			to = new(common.Address)
		}
		if unpackedData.Amount != nil {
			value = hexutil.Big(*unpackedData.Amount.Value.BigInt())
		}
		if len(unpackedData.Data) > 0 {
			input = (hexutil.Bytes)(unpackedData.Data)
		}
		nonce = hexutil.Uint64(unpackedData.Nonce)
	}

	return &Transaction{
		BlockHash:        blockHash,
		BlockNumber:      blockNumber,
		From:             from,
		Gas:              hexutil.Uint64(lTx.Fee.Gas),
		GasPrice:         (*hexutil.Big)(&lTx.Fee.Price.Value),
		Hash:             common.BytesToHash(tmTx.Hash()),
		Input:            input,
		Nonce:            nonce,
		To:               to,
		TransactionIndex: txIndex,
		Value:            value,
		V:                (*hexutil.Big)(v),
		R:                r,
		S:                s,
	}, nil
}

// ParseLegacyTx is used to parse the signed tx for old OneLedger tx types
func ParseLegacyTx(tmTx tmtypes.Tx) (*action.SignedTx, error) {
	tx := &action.SignedTx{}

	err := jsonSerializer.Deserialize(tmTx, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// EthHeaderFromTendermint is an util function that returns an Ethereum Header
// from a tendermint Block header.
func EthHeaderFromTendermint(block *tmtypes.Block) (*Header, error) {
	results, err := tmrpccore.BlockResults(nil, &block.Height)
	if err != nil {
		return nil, err
	}
	gasLimit, err := BlockMaxGasFromConsensusParams(nil)
	if err != nil {
		return nil, err
	}
	_, gasUsed, err := EthTransactionsFromTendermint(block, false)
	if err != nil {
		return nil, err
	}
	header := block.Header
	ethHeader := &Header{
		Hash:        common.BytesToHash(block.Hash()),
		ParentHash:  common.BytesToHash(header.LastBlockID.Hash),
		UncleHash:   common.Hash{},
		Coinbase:    common.BytesToAddress(header.ProposerAddress.Bytes()),
		Root:        common.BytesToHash(header.AppHash),
		TxHash:      common.BytesToHash(header.DataHash),
		ReceiptHash: common.Hash{},
		Difficulty:  vm.DefaultDifficulty,
		Number:      big.NewInt(header.Height),
		Time:        uint64(header.Time.Unix()),
		Extra:       common.Hex2Bytes(""),
		MixDigest:   common.Hash{},
		Nonce:       ethtypes.BlockNonce{},
		GasLimit:    new(big.Int).SetUint64(gasLimit),
		GasUsed:     new(big.Int).SetUint64(gasUsed),
		Bloom:       GetBlockBloom(results.EndBlockEvents),
		Size:        uint64(block.Size()),
	}
	return ethHeader, nil
}

// BlockMaxGasFromConsensusParams returns the gas limit for the latest block from the chain consensus params.
func BlockMaxGasFromConsensusParams(tmClient rpcclient.Client) (uint64, error) {
	// vm.DefaultBlockGasLimit will not be supported by javascript
	// which errors certain javascript dev tooling which only supports up to 53 bits
	return uint64(^uint32(0)), nil
}

// EthBlockFromTendermint returns a JSON-RPC compatible Ethereum blockfrom a given Tendermint block.
func EthBlockFromTendermint(block *tmtypes.Block, fullTx bool) (*Block, error) {
	gasLimit, err := BlockMaxGasFromConsensusParams(nil)
	if err != nil {
		return nil, err
	}

	transactions, gasUsed, err := EthTransactionsFromTendermint(block, fullTx)
	if err != nil {
		return nil, err
	}

	header := block.Header

	if len(header.DataHash) == 0 {
		header.DataHash = tmbytes.HexBytes(common.Hash{}.Bytes())
	}

	return &Block{
		Number:           hexutil.Uint64(header.Height),
		Hash:             common.BytesToHash(block.Hash()),
		ParentHash:       common.BytesToHash(header.LastBlockID.Hash),
		Nonce:            ethtypes.BlockNonce{}, // PoW specific
		Sha3Uncles:       common.Hash{},         // No uncles in Tendermint
		LogsBloom:        vm.BytesToBloom(make([]byte, 6)),
		TransactionsRoot: common.BytesToHash(header.DataHash),
		StateRoot:        common.BytesToHash(header.AppHash),
		Miner:            common.BytesToAddress(header.ProposerAddress.Bytes()),
		MixHash:          common.Hash{},
		Difficulty:       hexutil.Uint64(vm.DefaultDifficulty.Uint64()),
		TotalDifficulty:  hexutil.Uint64(vm.DefaultDifficulty.Uint64()),
		ExtraData:        common.Hex2Bytes(""),
		Size:             hexutil.Uint64(block.Size()),
		GasLimit:         (*hexutil.Big)(new(big.Int).SetUint64(gasLimit)), // Static gas limit
		GasUsed:          (*hexutil.Big)(new(big.Int).SetUint64(gasUsed)),
		Timestamp:        hexutil.Uint64(header.Time.Unix()),
		Transactions:     transactions,
		Uncles:           make([]common.Hash, 0),
		ReceiptsRoot:     common.Hash{},
	}, nil
}

// EthTransactionsFromTendermint returns a slice of ethereum transaction hashes and the total gas usage from a set of
// tendermint block transactions.
func EthTransactionsFromTendermint(block *tmtypes.Block, fullTx bool) ([]interface{}, uint64, error) {
	transactions := make([]interface{}, 0)
	gasUsed := uint64(0)

	chainID := utils.HashToBigInt(block.ChainID)

	for i, tx := range block.Txs {
		if !fullTx {
			// first parse legacy tx
			lTx, err := ParseLegacyTx(tx)
			if err != nil {
				// means tx is not legacy and we need to check is tx is ethereum
				continue
			}
			gasUsed += uint64(lTx.Fee.Gas)
			transactions = append(transactions, common.BytesToHash(tx.Hash()))
		} else {
			index := hexutil.Uint64(i)
			fTx, err := LegacyRawBlockAndTxToEthTx(block, &tx, chainID, &index)
			if err != nil {
				continue
			}
			gasUsed += uint64(fTx.Gas)
			transactions = append(transactions, fTx)
		}
	}

	return transactions, gasUsed, nil
}

func StateAndHeaderByNumberOrHash(blockStore *store.BlockStore, blockNrOrHash rpc.BlockNumberOrHash) (int64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return blockNr.Int64(), nil
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		block := blockStore.LoadBlockByHash(hash.Bytes())
		if block == nil {
			return 0, errors.New("header for hash not found")
		}
		return block.Height, nil
	}
	return 0, errors.New("invalid arguments; neither block nor hash specified")
}
