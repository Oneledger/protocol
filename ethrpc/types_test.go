package ethrpc

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHexIntUnmarshal(t *testing.T) {
	test := struct {
		ID hexInt `json:"id"`
	}{}

	data := []byte(`{"id": "0x1cc348"}`)
	err := json.Unmarshal(data, &test)

	require.Nil(t, err)
	require.Equal(t, hexInt(1885000), test.ID)
}

func TestHexBigUnmarshal(t *testing.T) {
	test := struct {
		ID hexBig `json:"id"`
	}{}

	data := []byte(`{"id": "0x51248487c7466b7062d"}`)
	err := json.Unmarshal(data, &test)

	require.Nil(t, err)
	b := big.Int{}
	b.SetString("23949082357483433297453", 10)

	require.Equal(t, hexBig(b), test.ID)
}

func TestSyncingUnmarshal(t *testing.T) {
	syncing := new(Syncing)
	err := json.Unmarshal([]byte("0"), syncing)
	require.NotNil(t, err)

	data := []byte(`{
		"startingBlock": "0x384",
		"currentBlock": "0x386",
		"highestBlock": "0x454"
	}`)

	err = json.Unmarshal(data, syncing)
	require.Nil(t, err)
	require.True(t, syncing.IsSyncing)
	require.Equal(t, 900, syncing.StartingBlock)
	require.Equal(t, 902, syncing.CurrentBlock)
	require.Equal(t, 1108, syncing.HighestBlock)
}

func TestTransactionUnmarshal(t *testing.T) {
	tx := new(Transaction)
	err := json.Unmarshal([]byte("111"), tx)
	require.NotNil(t, err)

	data := []byte(`{
        "blockHash": "0x3003694478c108eaec173afcb55eafbb754a0b204567329f623438727ffa90d8",
        "blockNumber": "0x83319",
        "from": "0x201354729f8d0f8b64e9a0c353c672c6a66b3857",
        "gas": "0x15f90",
        "gasPrice": "0x4a817c800",
        "hash": "0xfc7dcd42eb0b7898af2f52f7c5af3bd03cdf71ab8b3ed5b3d3a3ff0d91343cbe",
        "input": "0xe1fa8e8425f1af44eb895e4900b8be35d9fdc28744a6ef491c46ec8601990e12a58af0ed",
        "nonce": "0x6ba1",
        "to": "0xd10e3be2bc8f959bc8c41cf65f60de721cf89adf",
        "transactionIndex": "0x3",
        "value": "0x0"
    }`)

	err = json.Unmarshal(data, tx)

	require.Nil(t, err)
	require.Equal(t, "0x3003694478c108eaec173afcb55eafbb754a0b204567329f623438727ffa90d8", tx.BlockHash)
	require.Equal(t, 537369, *tx.BlockNumber)
	require.Equal(t, "0x201354729f8d0f8b64e9a0c353c672c6a66b3857", tx.From)
	require.Equal(t, 90000, tx.Gas)
	require.Equal(t, *big.NewInt(20000000000), tx.GasPrice)
	require.Equal(t, "0xfc7dcd42eb0b7898af2f52f7c5af3bd03cdf71ab8b3ed5b3d3a3ff0d91343cbe", tx.Hash)
	require.Equal(t, "0xe1fa8e8425f1af44eb895e4900b8be35d9fdc28744a6ef491c46ec8601990e12a58af0ed", tx.Input)
	require.Equal(t, 27553, tx.Nonce)
	require.Equal(t, "0xd10e3be2bc8f959bc8c41cf65f60de721cf89adf", tx.To)
	require.Equal(t, 3, *tx.TransactionIndex)
	require.Equal(t, *big.NewInt(0), tx.Value)
}