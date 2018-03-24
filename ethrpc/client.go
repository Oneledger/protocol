package ethrpc

import (
	"io/ioutil"
	"bytes"
	"fmt"
	"os"
	"log"
	"net/http"
	"encoding/json"

	"github.com/Oneledger/prototype-api/common"
)


// EthError - ethereum error
type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err EthError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type ethResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

type ethRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// EthRPC - Ethereum rpc client
type EthRPCClient struct {
	url    string
	client common.HttpClient
	log    common.Logger
	Debug  bool
}

// NewEthRPC create new rpc client with given url
func NewEthRPCClient(url string, options ...func(rpc *EthRPCClient)) *EthRPCClient {
	rpc := &EthRPCClient{
		url:    url,
		client: http.DefaultClient,
		log:    log.New(os.Stderr, "", log.LstdFlags),
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// Call returns raw response of method call
func (rpc *EthRPCClient) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if rpc.Debug {
		rpc.log.Println(fmt.Sprintf("%s\nRequest: %s\nResponse: %s\n", method, body, data))
	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

func (rpc *EthRPCClient) call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}


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
