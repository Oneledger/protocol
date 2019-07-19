package app

import (
	"errors"
	"os"
	"testing"

	"github.com/Oneledger/protocol/data/balance"

	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/config"
	"github.com/stretchr/testify/assert"
)

// global setup
func setup() []string {
	testDBs := []string{"test_dbpath"}
	return testDBs
}

// remove test db dir after
func teardown(dbPaths []string) {
	for _, v := range dbPaths {
		err := os.RemoveAll(v)
		if err != nil {
			errors.New("Remove test db file error")
		}
	}
}

func setupForGeneralInfo() (*config.Server, *node.Context) {
	NodeConfig := &config.NodeConfig{
		NodeName: "test_node",
		DBDir:    "test_dbpath",
		DB:       "goleveldb",
	}
	cfg := &config.Server{
		Node: NodeConfig,
	}
	nodeContext := &node.Context{}
	return cfg, nodeContext
}

func setupForStart() (*config.Server, *node.Context) {
	networkConfig := &config.NetworkConfig{
		P2PAddress: "tcp://127.0.0.1:26601",
		RPCAddress: "tcp://127.0.0.1:26600",
		SDKAddress: "http://127.0.0.1:26603",
	}
	consensus := &config.ConsensusConfig{}
	p2pconfig := &config.P2PConfig{}
	mempool := &config.MempoolConfig{}
	nodeConfig := &config.NodeConfig{
		NodeName: "test_node",
		FastSync: true,
		DBDir:    "test_dbpath",
		DB:       "goleveldb",
	}
	cfg := &config.Server{
		Node:      nodeConfig,
		Network:   networkConfig,
		Consensus: consensus,
		P2P:       p2pconfig,
		Mempool:   mempool,
	}
	nodeContext := &node.Context{}
	return cfg, nodeContext
}

func TestNewApp(t *testing.T) {
	t.Run("should return an error with config input as nil", func(t *testing.T) {
		nodeContext := &node.Context{}
		_, error := NewApp(nil, nodeContext)
		assert.Error(t, error)
	})

	t.Run("should return no error and non-empty new app", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)
		cfg, nodeContext := setupForGeneralInfo()
		app, error := NewApp(cfg, nodeContext)
		if assert.NoError(t, error) {
			assert.NotEmpty(t, app)
		}
	})
}

// ABCI, Header, Node functions
func TestApp_ABCI(t *testing.T) {
	t.Run("ABCI function should return non empty ABCI", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)
		cfg, nodeContext := setupForGeneralInfo()
		app, _ := NewApp(cfg, nodeContext)
		abci := app.ABCI()
		assert.NotEmpty(t, abci)
	})

	t.Run("Header should return empty app header", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)
		cfg, nodeContext := setupForGeneralInfo()
		app, _ := NewApp(cfg, nodeContext)
		header := app.Header()
		assert.Empty(t, header)
	})

	t.Run("Node should return empty app node", func(t *testing.T) {
		testDB := setup()
		defer teardown(testDB)
		cfg, nodeContext := setupForGeneralInfo()
		app, _ := NewApp(cfg, nodeContext)
		node := app.Node()
		assert.Empty(t, node)
	})
}

// TODO : depend on starting a new node and tendermint server
//func TestApp_Start(t *testing.T) {
//	testDB := setup()
//	testDB = append(testDB, "./consensus")
//
//	cfg, nodeContext := setupForStart()
//	app, _ := NewApp(cfg, nodeContext)
//
//	err := os.MkdirAll("./consensus/config/", 0777)
//	assert.Nil(t, err)
//	pvKeyFilePath := "./consensus/config/priv_validator_key.json"
//	pvKeyFile ,err := os.Create(pvKeyFilePath)
//	assert.Nil(t, err)
//
//	err = os.MkdirAll("./consensus/data/", 0777)
//	assert.Nil(t, err)
//	pvStateFilePath := "./consensus/data/priv_validator_state.json"
//	pvStateFile ,err := os.Create(pvStateFilePath)
//	assert.Nil(t, err)
//
//	defer teardown(testDB)
//	defer pvKeyFile.Close()
//	defer pvStateFile.Close()
//
//	err = app.Start()
//	assert.Error(t, err)
//}

func TestApp_setupstate(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, _ := NewApp(cfg, nodeContext)
	t.Run("should return an error with empty []byte as input", func(t *testing.T) {
		testdata := []byte("")
		err := app.setupState(testdata)
		assert.Error(t, err)
	})

	// TODO : need a proper stateBytes test data
	//t.Run("should return no error when given a proper stateBytes data", func(t *testing.T) {
	//})
}

func TestApp_setupValidators(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, err := NewApp(cfg, nodeContext)
	assert.Nil(t, err)

	req := RequestInitChain{}
	// prepare for currencies
	currencies := balance.NewCurrencyList()
	currency := balance.Currency{
		Name: "VT",
	}
	err = currencies.Register(currency)
	assert.Nil(t, err)

	validators, _ := app.setupValidators(req, currencies)
	assert.Empty(t, validators)
}

// TODO: start a tendermint node to test this function for now, but need to mock a tendermint server to do proper unit test
//func TestApp_rpcStarter(t *testing.T) {
//	testDB := setup()
//	defer teardown(testDB)
//	cfg, nodeContext := setupForStart()
//	app, err := NewApp(cfg, nodeContext)
//	assert.Nil(t, err)
//
//	_, err = app.rpcStarter()
//	assert.Error(t, err)
//}

func TestContext_Action(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, err := NewApp(cfg, nodeContext)
	assert.Nil(t, err)

	context := app.Context.Action(&app.header)
	assert.NotEmpty(t, context)
}

func TestContext_Accounts(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, err := NewApp(cfg, nodeContext)
	assert.Nil(t, err)

	accounts := app.Context.Accounts()
	assert.Empty(t, accounts.Accounts())
}

func TestContext_ValidatorCtx(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, err := NewApp(cfg, nodeContext)
	assert.Nil(t, err)

	validatorContext := app.Context.ValidatorCtx()
	assert.NotNil(t, validatorContext)
}

func TestContext_Balances(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, err := NewApp(cfg, nodeContext)
	assert.Nil(t, err)

	balance := app.Context.Balances()
	assert.Equal(t, 0, balance.Currencies().Len())
}

// TODO: start a tendermint node to test this function for now, but need to mock a tendermint server to do proper unit test
//func TestContext_Services(t *testing.T) {
//	testDB := setup()
//	defer teardown(testDB)
//	cfg, nodeContext := setupForStart()
//	app, err := NewApp(cfg, nodeContext)
//	assert.Nil(t, err)
//
//	_, err = app.Context.Services()
//	assert.Error(t, err)
//
//	app.Close()
//}

func TestContext_Node(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, err := NewApp(cfg, nodeContext)
	assert.Nil(t, err)

	node := app.Context.Node()
	assert.Equal(t, "", node.NodeName)
}

func TestContext_Validators(t *testing.T) {
	testDB := setup()
	defer teardown(testDB)
	cfg, nodeContext := setupForGeneralInfo()
	app, err := NewApp(cfg, nodeContext)
	assert.Nil(t, err)

	vs := app.Context.Validators()
	assert.Empty(t, vs.Hash)
	assert.Empty(t, vs.LastHash)
	assert.Equal(t, int8(0), vs.TreeHeight)
}
