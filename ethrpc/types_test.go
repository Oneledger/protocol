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

func TestBlockWithoutTransactionsUnmarshal(t *testing.T) {
	block := new(Block)
	err := json.Unmarshal([]byte("222"), block)
	require.NotNil(t, err)

	data := []byte(`{
		"number": "0x1b4",
		"hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		"parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
		"nonce": "0xe04d296d2460cfb8472af2c5fd05b5a214109c25688d3704aed5484f9a7792f2",
		"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"logsBloom": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
		"miner": "0x4e65fda2159562a496f9f3522f89122a3088497a",
		"difficulty": "0x27f07",
		"totalDifficulty":  "0x27f07",
		"extraData": "0x4554482e45544846414e532e4f52472d4243323041423836",
		"size": "0x27f07",
		"gasLimit": "0x9f759",
		"gasUsed": "0x9f759",
		"timestamp": "0x54e34e8e",
		"transactions": ["0xfc7dcd42eb0b7898af2f52f7c5af3bd03cdf71ab8b3ed5b3d3a3ff0d91343cbe","0xecd8a21609fa852c08249f6c767b7097481da34b9f8d2aae70067918955b4e69"],
		"uncles": ["0x1606e5", "0xd5145a9"]
	}`)
	
	pb := new(JsonBlockWithoutTransactions)
	err = json.Unmarshal(data, &pb)
	blockOne := pb.toBlock()

	require.Nil(t, err)
	require.Equal(t, 436, blockOne.Number)
	require.Equal(t, "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331", blockOne.Hash)
	require.Equal(t, "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5", blockOne.ParentHash)
	require.Equal(t, "0xe04d296d2460cfb8472af2c5fd05b5a214109c25688d3704aed5484f9a7792f2", blockOne.Nonce)
	require.Equal(t, "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347", blockOne.Sha3Uncles)
	require.Equal(t, "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331", blockOne.LogsBloom)
	require.Equal(t, "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421", blockOne.TransactionsRoot)
	require.Equal(t, "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff", blockOne.StateRoot)
	require.Equal(t, "0x4e65fda2159562a496f9f3522f89122a3088497a", blockOne.Miner)
	require.Equal(t, *big.NewInt(163591), blockOne.Difficulty)
	require.Equal(t, *big.NewInt(163591), blockOne.TotalDifficulty)
	require.Equal(t, "0x4554482e45544846414e532e4f52472d4243323041423836", blockOne.ExtraData)
	require.Equal(t, 163591, blockOne.Size)
	require.Equal(t, 653145, blockOne.GasLimit)
	require.Equal(t, 653145, blockOne.GasUsed)
	require.Equal(t, 1424182926, blockOne.Timestamp)
	require.Equal(t, 2, len(blockOne.Transactions))
	require.Equal(t, "0xfc7dcd42eb0b7898af2f52f7c5af3bd03cdf71ab8b3ed5b3d3a3ff0d91343cbe", blockOne.Transactions[0].Hash)
	require.Equal(t, "0xd5145a9", blockOne.Uncles[1])
}

func TestBlockWithTransactionsUnmarshal(t *testing.T) {
	data := []byte(`{
		"number": "0x1b4",
		"hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		"parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
		"nonce": "0xe04d296d2460cfb8472af2c5fd05b5a214109c25688d3704aed5484f9a7792f2",
		"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"logsBloom": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
		"miner": "0x4e65fda2159562a496f9f3522f89122a3088497a",
		"difficulty": "0x27f07",
		"totalDifficulty":  "0x27f07",
		"extraData": "0x4554482e45544846414e532e4f52472d4243323041423836",
		"size": "0x27f07",
		"gasLimit": "0x9f759",
		"gasUsed": "0x9f759",
		"timestamp": "0x54e34e8e",
		"uncles": ["0x1606e5", "0xd5145a9"],
		"transactions": [{
			"hash": "0xfc7dcd42eb0b7898af2f52f7c5af3bd03cdf71ab8b3ed5b3d3a3ff0d91343cbe",
			"nonce": "0x6ba1",
			"blockHash": "0x3003694478c108eaec173afcb55eafbb754a0b204567329f623438727ffa90d8",
			"blockNumber": "0x83319",
			"transactionIndex": "0x3",
			"from": "0x201354729f8d0f8b64e9a0c353c672c6a66b3857",
			"to": "0xd10e3be2bc8f959bc8c41cf65f60de721cf89adf",
			"value": "0x0",
			"gas": "0x15f90",
			"gasPrice": "0x4a817c800",
			"input": "0xe1fa8e8425f1af44eb895e4900b8be35d9fdc28744a6ef491c46ec8601990e12a58af0ed"
		}]
	}`)

	pb := new(JsonBlockWithTransactions)
	err := json.Unmarshal(data, &pb)
	blockTwo := pb.toBlock()

	require.Nil(t, err)
	require.Equal(t, 1, len(blockTwo.Transactions))
	require.Equal(t, "0xfc7dcd42eb0b7898af2f52f7c5af3bd03cdf71ab8b3ed5b3d3a3ff0d91343cbe", blockTwo.Transactions[0].Hash)
}