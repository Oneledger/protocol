package action

import (
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/serialize"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// ----------------------------------------------------------------------------
// Transaction logs
// ----------------------------------------------------------------------------

func MarshalLogs(logs []*ethtypes.Log) ([]byte, error) {
	return serialize.GetSerializer(serialize.PERSISTENT).Serialize(logs)
}

func UnmarshalLogs(in []byte) ([]*ethtypes.Log, error) {
	logs := []*ethtypes.Log{}
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(in, logs)
	return logs, err
}

// SetLogs sets the logs for a transaction in the KVStore.
func (s *CommitStateDB) SetLogs(hash ethcmn.Hash, logs []*ethtypes.Log) error {
	bz, err := MarshalLogs(logs)
	if err != nil {
		return err
	}

	s.contractStore.Set(evm.KeyPrefixLogs, hash.Bytes(), bz)
	s.logSize = uint(len(logs))
	return nil
}

// DeleteLogs removes the logs from the KVStore. It is used during journal.Revert.
func (s *CommitStateDB) DeleteLogs(hash ethcmn.Hash) {
	s.contractStore.Delete(evm.KeyPrefixLogs, hash.Bytes())
}

// AddLog adds a new log to the state and sets the log metadata from the state.
func (s *CommitStateDB) AddLog(log *ethtypes.Log) {
	s.journal.append(addLogChange{txhash: s.thash})

	log.TxHash = s.thash
	log.BlockHash = s.bhash
	log.TxIndex = uint(s.txIndex)
	log.Index = s.logSize

	logs, err := s.GetLogs(s.thash)
	if err != nil {
		// panic on unmarshal error
		panic(err)
	}

	if err = s.SetLogs(s.thash, append(logs, log)); err != nil {
		// panic on marshal error
		panic(err)
	}
}

// GetLogs returns the current logs for a given transaction hash from the KVStore.
func (s *CommitStateDB) GetLogs(hash ethcmn.Hash) ([]*ethtypes.Log, error) {
	bz, _ := s.contractStore.Get(evm.KeyPrefixLogs, hash.Bytes())
	if len(bz) == 0 {
		// return nil error if logs are not found
		return []*ethtypes.Log{}, nil
	}

	return UnmarshalLogs(bz)
}
