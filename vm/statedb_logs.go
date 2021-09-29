package vm

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/kv"
)

// ----------------------------------------------------------------------------
// Transaction logs
// ----------------------------------------------------------------------------

// AddLog adds a new log to the state and sets the log metadata from the state.
func (s *CommitStateDB) AddLog(log *ethtypes.Log) {
	s.journal.append(addLogChange{txhash: s.thash})

	log.BlockHash = s.bhash
	log.TxHash = s.thash
	log.Index = s.logSize

	s.logs[s.thash] = append(s.logs[s.thash], log)
	s.logSize++
}

// GetTxLogs return current tx logs
func (s *CommitStateDB) GetTxLogs() []*ethtypes.Log {
	return s.logs[s.thash]
}

// updateBloom block bloom with tx log
func (s *CommitStateDB) updateBloom(log *ethtypes.Log) {
	s.bloom.Add(log.Address.Bytes(), s.bloomBuffer)
	for _, b := range log.Topics {
		s.bloom.Add(b[:], s.bloomBuffer)
	}
}

// GetBloomEvent for block ender
func (s *CommitStateDB) GetBloomEvent() *types.Event {
	if len(s.logs) == 0 {
		s.logger.Detailf("bloom for block %s - not generated\n", s.bhash)
		return nil
	}
	s.logger.Detailf("bloom for block %s generated with buffer %s and bloom %s\n", s.bhash, s.bloomBuffer, s.bloom.Bytes())
	return &types.Event{
		Type: "olvm",
		Attributes: []kv.Pair{
			{
				Key:   []byte("block.bloom"),
				Value: s.bloom.Bytes(),
			},
		},
	}
}
