package rpc

import (
	"fmt"

	"github.com/Oneledger/protocol/node/chains/common"
)

// Web3ClientVersion returns the current client version.
func (rpc *EthRPCClient) Web3ClientVersion() (string, error) {
	var clientVersion string

	err := rpc.call("web3_clientVersion", &clientVersion)
	return clientVersion, err
}

// Web3Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (rpc *EthRPCClient) Web3Sha3(data []byte) (string, error) {
	var hash string

	err := rpc.call("web3_sha3", &hash, fmt.Sprintf("0x%x", data))
	return hash, err
}

// EthSyncing returns an object with data about the sync status or false.
func (rpc *EthRPCClient) EthSyncing() (*Syncing, error) {
	syncing := new(Syncing)
	err := rpc.call("eth_syncing", &syncing)
	if err != nil {
		return nil, err
	}

	return syncing, err
}

func (rpc *EthRPCClient) getBlock(method string, withTransactions bool, params ...interface{}) (*Block, error) {
	var response ProxyBlock
	if withTransactions {
		response = new(JsonBlockWithTransactions)
	} else {
		response = new(JsonBlockWithoutTransactions)
	}

	err := rpc.call(method, response, params...)
	if err != nil {
		return nil, err
	}
	block := response.toBlock()

	return &block, nil
}

// EthGetBlockByHash returns information about a block by hash.
func (rpc *EthRPCClient) EthGetBlockByHash(hash string, withTransactions bool) (*Block, error) {
	return rpc.getBlock("eth_getBlockByHash", withTransactions, hash, withTransactions)
}

// EthGetBlockByNumber returns information about a block by block number.
func (rpc *EthRPCClient) EthGetBlockByNumber(number int, withTransactions bool) (*Block, error) {
	return rpc.getBlock("eth_getBlockByNumber", withTransactions, common.IntToHex(number), withTransactions)
}

func (rpc *EthRPCClient) getTransaction(method string, params ...interface{}) (*Transaction, error) {
	transaction := new(Transaction)

	err := rpc.call(method, transaction, params...)
	return transaction, err
}

// EthGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (rpc *EthRPCClient) EthGetTransactionByHash(hash string) (*Transaction, error) {
	return rpc.getTransaction("eth_getTransactionByHash", hash)
}

// EthGetTransactionByBlockHashAndIndex returns information about a transaction by block hash and transaction index position.
func (rpc *EthRPCClient) EthGetTransactionByBlockHashAndIndex(blockHash string, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("eth_getTransactionByBlockHashAndIndex", blockHash, common.IntToHex(transactionIndex))
}

// EthGetTransactionByBlockNumberAndIndex returns information about a transaction by block number and transaction index position.
func (rpc *EthRPCClient) EthGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("eth_getTransactionByBlockNumberAndIndex", common.IntToHex(blockNumber), common.IntToHex(transactionIndex))
}

func (rpc *EthRPCClient) EthGetTransactionReceipt(hash string) (*Transaction, error) {
	return rpc.getTransaction("eth_getTransactionReceipt", hash)
}