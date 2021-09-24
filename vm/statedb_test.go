package vm

import (
	"os"
	"testing"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/evm"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
)

func TestCommitStateDB(t *testing.T) {

	db := db.NewDB("test", db.MemDBBackend, "")

	balances := balance.NewStore("tb", storage.NewState(storage.NewChainState("balance", db)))

	currencies := balance.NewCurrencySet()
	currency := balance.Currency{
		Name:    "OLT",
		Chain:   chain.Type(1),
		Decimal: int64(18),
	}
	currencies.Register(currency)

	logger := log.NewLoggerWithPrefix(os.Stdout, "Test-Logger")

	stateDB := NewCommitStateDB(
		evm.NewContractStore(storage.NewState(storage.NewChainState("contracts", db))),
		balance.NewNesterAccountKeeper(
			storage.NewState(storage.NewChainState("keeper", db)),
			balances,
			currencies,
		),
		logger,
	)
	stateDB.SetHeightHash(1, ethcmn.Hash{})

	t.Run("test snapshot and it is ok", func(t *testing.T) {
		stateDB.SetHeightHash(1, ethcmn.Hash{})
		defer stateDB.Reset()

		snap := stateDB.Snapshot()
		assert.Equal(t, 0, len(stateDB.journal.entries))

		stateDB.AddAddressToAccessList(ethcmn.Address{1})
		assert.Equal(t, 1, len(stateDB.journal.entries))

		stateDB.RevertToSnapshot(snap)
		assert.Equal(t, 0, len(stateDB.journal.entries))
	})

	t.Run("test copy state db and it is ok", func(t *testing.T) {
		stateDB.SetHeightHash(1, ethcmn.Hash{})
		defer stateDB.Reset()

		stateDB.Snapshot()
		assert.Equal(t, 0, len(stateDB.journal.entries))
		assert.Equal(t, 1, len(stateDB.validRevisions))

		stateDB.AddAddressToAccessList(ethcmn.Address{1})
		assert.Equal(t, 1, len(stateDB.journal.entries))
		assert.Equal(t, 1, len(stateDB.validRevisions))

		stateDBCopy := stateDB.Copy()
		assert.Equal(t, 1, len(stateDB.journal.entries))
		assert.Equal(t, 1, len(stateDB.validRevisions))

		assert.Equal(t, 0, len(stateDBCopy.journal.entries))
		assert.Equal(t, 0, len(stateDBCopy.validRevisions))
	})

	t.Run("test log addition and it is ok", func(t *testing.T) {
		stateDB.SetHeightHash(1, ethcmn.Hash{})
		defer stateDB.Reset()

		stateDB.thash = ethcmn.BytesToHash([]byte{1, 2})
		assert.Equal(t, 0, int(stateDB.logSize))
		assert.Equal(t, 0, len(stateDB.logs))

		logTestCase := [][]int{
			{1, 1, 1},
			{2, 1, 2},
			{3, 1, 3},
		}

		for _, testCase := range logTestCase {
			stateDB.AddLog(&ethtypes.Log{})
			assert.Equal(t, testCase[0], int(stateDB.logSize), "Should increase for new log")
			assert.Equal(t, testCase[1], len(stateDB.logs), "Always the same as logs storing to txhash key")
			assert.Equal(t, testCase[2], len(stateDB.logs[stateDB.thash]), "Must add a new log")
		}
	})
}
