package vm

import (
	"math/big"

	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/serialize"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// ----------------------------------------------------------------------------
// Transaction logs
// ----------------------------------------------------------------------------

type BlockLogs struct {
	BlockHash   ethcmn.Hash
	BlockNumber uint64
	Logs        map[ethcmn.Hash][]*ethtypes.Log
}

func (bl *BlockLogs) MarshalLogs() ([]byte, error) {
	return serialize.GetSerializer(serialize.PERSISTENT).Serialize(bl)
}

func (bl *BlockLogs) UnmarshalLogs(in []byte) error {
	return serialize.GetSerializer(serialize.PERSISTENT).Deserialize(in, bl)
}

// UpdateLogs sets the logs for transactions in the store.
func (s *CommitStateDB) UpdateLogs(height uint64) error {
	bl := &BlockLogs{
		BlockHash:   s.bhash,
		BlockNumber: height,
		Logs:        s.logs,
	}
	bz, err := bl.MarshalLogs()
	if err != nil {
		return err
	}
	err = s.contractStore.Set(evm.KeyPrefixLogs, s.bhash.Bytes(), bz)
	if err != nil {
		return err
	}
	err = s.contractStore.Set(evm.KeyPrefixBloom, evm.BloomKey(height), s.bloom.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// UpdateBloom filter for block
func (s *CommitStateDB) UpdateBloom() {
	logs := s.logs[s.thash]
	// calculating bloom for the block
	bloomInt := big.NewInt(0).SetBytes(ethtypes.LogsBloom(logs))
	s.bloom.Or(s.bloom, bloomInt)
}

// AddLog adds a new log to the state and sets the log metadata from the state.
func (s *CommitStateDB) AddLog(log *ethtypes.Log) {
	s.journal.append(addLogChange{txhash: s.thash})

	log.BlockHash = s.bhash
	log.TxHash = s.thash
	log.TxIndex = uint(s.txIndex)
	log.Index = s.logSize

	s.logs[s.thash] = append(s.logs[s.thash], log)
	s.logSize++
}

// GetTxLogs return current tx logs
func (s *CommitStateDB) GetTxLogs() []*ethtypes.Log {
	return s.logs[s.thash]
}

// GetLogs returns the current logs for a given transaction hash from the store.
func (s *CommitStateDB) GetLogs(blockHash ethcmn.Hash) (*BlockLogs, error) {
	bz, _ := s.contractStore.Get(evm.KeyPrefixLogs, blockHash.Bytes())

	bl := &BlockLogs{}
	err := bl.UnmarshalLogs(bz)
	if err != nil {
		return nil, err
	}
	return bl, nil
}

// GetBlockBloom get bloom by block height
func (s *CommitStateDB) GetBlockBloom(height uint64) ethtypes.Bloom {
	bz, _ := s.contractStore.Get(evm.KeyPrefixBloom, evm.BloomKey(height))
	if len(bz) == 0 {
		return ethtypes.Bloom{}
	}
	return ethtypes.BytesToBloom(bz)
}
