package vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// RLPLog used for properly RLP encoding of all field for the log
type RLPLog struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber" gencodec:"required"`
	// hash of the transaction
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
	// hash of the block in which the transaction was included
	BlockHash common.Hash `json:"blockHash" gencodec:"required"`
	// index of the log in the block
	Index uint `json:"logIndex" gencodec:"required"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed" gencodec:"required"`
}

// RLPLogConvert go-ethereum log in order to serialize properly with RLP
func RLPLogConvert(log types.Log) *RLPLog {
	rlpLog := RLPLog(log)
	return &rlpLog
}

// Encode encodes log rlp bytes
func (l *RLPLog) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(l)
}

// Decode decodes log rlp bytes
func (l *RLPLog) Decode(data []byte) (*types.Log, error) {
	err := rlp.DecodeBytes(data, l)
	if err != nil {
		return nil, err
	}
	log := types.Log(*l)
	return &log, nil
}
