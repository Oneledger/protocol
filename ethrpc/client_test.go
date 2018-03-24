package ethrpc

import (
	"fmt"
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/suite"
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
	response := `{"jsonrpc":"2.0", "id":1, "result": "test client"}`

	httpmock.RegisterResponder("POST", s.rpc.url, func(request *http.Request) (*http.Response, error) {
		body := s.getBody(request)
		s.methodEqual(body, "web3_clientVersion")
		s.paramsEqual(body, `null`)

		return httpmock.NewStringResponse(200, response), nil
	})

	v, err := s.rpc.Web3ClientVersion()
	s.Require().Nil(err)
	s.Require().Equal("test client", v)
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