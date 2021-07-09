package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	EarliestBlockNumber = 0
	LatestBlockNumber   = -1
	PendingBlockNumber  = -2
)

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"data"`
}

// Block represents a block returned to RPC clients.
type Block struct {
	Number           hexutil.Uint64      `json:"number"`
	Hash             common.Hash         `json:"hash"`
	ParentHash       common.Hash         `json:"parentHash"`
	Nonce            ethtypes.BlockNonce `json:"nonce"`
	Sha3Uncles       common.Hash         `json:"sha3Uncles"`
	LogsBloom        ethtypes.Bloom      `json:"logsBloom"`
	TransactionsRoot common.Hash         `json:"transactionsRoot"`
	StateRoot        common.Hash         `json:"stateRoot"`
	Miner            common.Address      `json:"miner"`
	MixHash          common.Hash         `json:"mixHash"`
	Difficulty       hexutil.Uint64      `json:"difficulty"`
	TotalDifficulty  hexutil.Uint64      `json:"totalDifficulty"`
	ExtraData        hexutil.Bytes       `json:"extraData"`
	Size             hexutil.Uint64      `json:"size"`
	GasLimit         hexutil.Uint64      `json:"gasLimit"`
	GasUsed          *hexutil.Big        `json:"gasUsed"`
	Timestamp        hexutil.Uint64      `json:"timestamp"`
	Uncles           []common.Hash       `json:"uncles"`
	ReceiptsRoot     common.Hash         `json:"receiptsRoot"`
	Transactions     []interface{}       `json:"transactions"`
}

// Transaction represents a transaction returned to RPC clients.
type Transaction struct {
	BlockHash        *common.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              hexutil.Uint64  `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex *hexutil.Uint64 `json:"transactionIndex"`
	Value            hexutil.Big     `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *common.Hash    `json:"r"`
	S                *common.Hash    `json:"s"`
}

// TransactionReceipt represents a mined transaction returned to RPC clients.
type TransactionReceipt struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	Status            hexutil.Uint64  `json:"status"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	LogsBloom         ethtypes.Bloom  `json:"logsBloom"`
	Logs              []*ethtypes.Log `json:"logs"`

	// Implementation fields: These fields are added by geth when processing a transaction.
	// They are stored in the chain database.
	TransactionHash common.Hash     `json:"transactionHash"`
	ContractAddress *common.Address `json:"contractAddress"`
	GasUsed         hexutil.Uint64  `json:"gasUsed"`

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        common.Hash    `json:"blockHash"`
	BlockNumber      hexutil.Big    `json:"blockNumber"`
	TransactionIndex hexutil.Uint64 `json:"transactionIndex"`

	// sender and receiver (contract or EOA) addresses
	From common.Address  `json:"from"`
	To   *common.Address `json:"to"`
}
