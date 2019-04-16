package rpc

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"math/big"
	"net/http"
)

type EthRPCTestSuite struct {
	suite.Suite
	rpc *EthRPCClient
}

func (s *EthRPCTestSuite) SetupSuite() {
	s.rpc = NewEthRPCClient("http://127.0.0.1:8545")

	httpmock.Activate()
}

func (s *EthRPCTestSuite) TearDownSuite() {
	httpmock.Deactivate()
}

func (s *EthRPCTestSuite) TearDownTest() {
	httpmock.Reset()
}

func (s *EthRPCTestSuite) getBody(request *http.Request) []byte {
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	s.Require().Nil(err)

	return body
}

func (s *EthRPCTestSuite) methodEqual(body []byte, expected string) {
	value := gjson.GetBytes(body, "method").String()

	s.Require().Equal(expected, value)
}

func (s *EthRPCTestSuite) paramsEqual(body []byte, expected string) {
	value := gjson.GetBytes(body, "params").Raw
	if expected == "null" {
		s.Require().Equal(expected, value)
	} else {
		s.JSONEq(expected, value)
	}
}

func (s *EthRPCTestSuite) TestWeb3ClientVersion() {
	response := `{"jsonrpc":"2.0", "id":1, "result": "test comm"}`

	httpmock.RegisterResponder("POST", s.rpc.url, func(request *http.Request) (*http.Response, error) {
		body := s.getBody(request)
		s.methodEqual(body, "web3_clientVersion")
		s.paramsEqual(body, `null`)

		return httpmock.NewStringResponse(200, response), nil
	})

	v, err := s.rpc.Web3ClientVersion()
	s.Require().Nil(err)
	s.Require().Equal("test comm", v)
}

func (s *EthRPCTestSuite) registerResponse(result string, callback func([]byte)) {
	httpmock.Reset()
	response := fmt.Sprintf(`{"jsonrpc":"2.0", "id":1, "result": %s}`, result)
	httpmock.RegisterResponder("POST", s.rpc.url, func(request *http.Request) (*http.Response, error) {
		callback(s.getBody(request))
		return httpmock.NewStringResponse(200, response), nil
	})
}

func (s *EthRPCTestSuite) registerResponseError(err error) {
	httpmock.Reset()
	httpmock.RegisterResponder("POST", s.rpc.url, func(request *http.Request) (*http.Response, error) {
		return nil, err
	})
}

func (s *EthRPCTestSuite) TestEthSyncing() {
	s.registerResponseError(errors.New("Error"))
	syncing, err := s.rpc.EthSyncing()
	s.Require().NotNil(err)

	expected := &Syncing{
		IsSyncing:     false,
		CurrentBlock:  0,
		HighestBlock:  0,
		StartingBlock: 0,
	}
	s.registerResponse(`false`, func(body []byte) {
		s.methodEqual(body, "eth_syncing")
	})
	syncing, err = s.rpc.EthSyncing()

	s.Require().Nil(err)
	s.Require().Equal(expected, syncing)

	httpmock.Reset()
	s.registerResponse(`{
		"currentBlock": "0x8c3be",
		"highestBlock": "0x9bb3b",
		"startingBlock": "0x0"
	}`, func(body []byte) {})

	expected = &Syncing{
		IsSyncing:     true,
		CurrentBlock:  574398,
		HighestBlock:  637755,
		StartingBlock: 0,
	}
	syncing, err = s.rpc.EthSyncing()
	s.Require().Nil(err)
	s.Require().Equal(expected, syncing)
}

func (s *EthRPCTestSuite) TestEthGetBlockByHash() {
	// Test with transactions
	hash := "0x111"
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getBlockByHash")
		s.paramsEqual(body, `["0x111", true]`)
	})

	_, err := s.rpc.EthGetBlockByHash(hash, true)
	s.Require().Nil(err)

	httpmock.Reset()

	// Test without transactions
	hash = "0x222"
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getBlockByHash")
		s.paramsEqual(body, `["0x222", false]`)
	})

	_, err = s.rpc.EthGetBlockByHash(hash, false)
	s.Require().Nil(err)
}

func (s *EthRPCTestSuite) TestEthGetBlockByNumber() {
	// Test with transactions
	number := 3274863
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getBlockByNumber")
		s.paramsEqual(body, `["0x31f86f", true]`)
	})

	_, err := s.rpc.EthGetBlockByNumber(number, true)
	s.Require().Nil(err)

	httpmock.Reset()

	// Test without transactions
	number = 14322
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getBlockByNumber")
		s.paramsEqual(body, `["0x37f2", false]`)
	})

	_, err = s.rpc.EthGetBlockByNumber(number, false)
	s.Require().Nil(err)
}

func newBigInt(s string) big.Int {
	i, _ := new(big.Int).SetString(s, 10)
	return *i
}

func (s *EthRPCTestSuite) TestGetTransaction() {
	result := `{
        "blockHash": "0x8b0404b2e5173e7abdbfc98f521d50808486ccaff3cd0a6344e0bb6c7aa8cef0",
        "blockNumber": "0x4109ed",
        "from": "0xe3a7ca9d2306b0dc900ea618648bed9ec6cb1106",
        "gas": "0x3d090",
        "gasPrice": "0xee6b2800",
        "hash": "0x3068bb24a6c65a80eb350b89b2ef2f4d0605f59e5d07fd3467eb76511c4408e7",
        "input": "0x522",
        "nonce": "0xa8",
        "to": "0x8d12a197cb00d4747a1fe03395095ce2a5cc6819",
        "transactionIndex": "0x98",
        "value": "0x9184e72a000"
    }`
	s.registerResponse(result, func(body []byte) {
		s.methodEqual(body, "trx")
	})

	transaction, err := s.rpc.getTransaction("trx")
	s.Require().Nil(err)
	s.Require().NotNil(transaction)
	s.Require().Equal("0x3068bb24a6c65a80eb350b89b2ef2f4d0605f59e5d07fd3467eb76511c4408e7", transaction.Hash)
	s.Require().Equal(168, transaction.Nonce)
	s.Require().Equal("0x8b0404b2e5173e7abdbfc98f521d50808486ccaff3cd0a6344e0bb6c7aa8cef0", transaction.BlockHash)
	s.Require().Equal(4262381, *transaction.BlockNumber)
	s.Require().Equal(152, *transaction.TransactionIndex)
	s.Require().Equal("0xe3a7ca9d2306b0dc900ea618648bed9ec6cb1106", transaction.From)
	s.Require().Equal("0x8d12a197cb00d4747a1fe03395095ce2a5cc6819", transaction.To)
	s.Require().Equal(newBigInt("10000000000000"), transaction.Value)
	s.Require().Equal(250000, transaction.Gas)
	s.Require().Equal(newBigInt("4000000000"), transaction.GasPrice)
	s.Require().Equal("0x522", transaction.Input)
}

func (s *EthRPCTestSuite) TestEthGetTransactionByHash() {
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getTransactionByHash")
		s.paramsEqual(body, `["0x123"]`)
	})

	t, err := s.rpc.EthGetTransactionByHash("0x123")
	s.Require().Nil(err)
	s.Require().NotNil(t)
}

func (s *EthRPCTestSuite) TestEthGetTransactionByBlockHashAndIndex() {
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getTransactionByBlockHashAndIndex")
		s.paramsEqual(body, `["0x623", "0x12"]`)
	})

	t, err := s.rpc.EthGetTransactionByBlockHashAndIndex("0x623", 18)
	s.Require().Nil(err)
	s.Require().NotNil(t)
}

func (s *EthRPCTestSuite) TestEthGetTransactionByBlockNumberAndIndex() {
	s.registerResponse(`{}`, func(body []byte) {
		s.methodEqual(body, "eth_getTransactionByBlockNumberAndIndex")
		s.paramsEqual(body, `["0x1f537da", "0xa"]`)
	})

	t, err := s.rpc.EthGetTransactionByBlockNumberAndIndex(32847834, 10)
	s.Require().Nil(err)
	s.Require().NotNil(t)
}
