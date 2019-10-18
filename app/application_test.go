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

func TestMain(m *testing.M) {

	os.Exit(teardown([]string{"test_dbpath"}))
}

// global setup
func setup(fn func() (*config.Server, *node.Context)) *App {
	cfg, nodeContext := fn()
	app, _ := NewApp(cfg, nodeContext)
	app.Start()
	return app
}

// remove test db dir after
func teardown(dbPaths []string) int {
	for _, v := range dbPaths {
		err := os.RemoveAll(v)
		if err != nil {
			errors.New("Remove test db file error")
		}
	}
	return 0
}

func setupForGeneralInfo() (*config.Server, *node.Context) {
	cfg := config.DefaultServerConfig()
	NodeConfig := &config.NodeConfig{
		NodeName: "test_node",
		DBDir:    "test_dbpath",
		DB:       "goleveldb",
	}
	cfg.Node = NodeConfig
	nodeContext := &node.Context{}
	return cfg, nodeContext
}

func setupForStart() (*config.Server, *node.Context) {
	cfg := config.DefaultServerConfig()
	networkConfig := &config.NetworkConfig{
		P2PAddress: "tcp://127.0.0.1:26601",
		RPCAddress: "tcp://127.0.0.1:26600",
		SDKAddress: "http://127.0.0.1:26603",
	}
	nodeConfig := &config.NodeConfig{
		NodeName: "test_node",
		FastSync: true,
		DBDir:    "test_dbpath",
		DB:       "goleveldb",
	}
	cfg.Node = nodeConfig
	cfg.Network = networkConfig
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
		app := setup(setupForGeneralInfo)
		assert.NotEmpty(t, app)
	})
}

// ABCI, Header, Node functions
func TestApp_ABCI(t *testing.T) {
	app := setup(setupForGeneralInfo)
	t.Run("ABCI function should return non empty ABCI", func(t *testing.T) {

		abci := app.ABCI()
		assert.NotEmpty(t, abci)
	})

	t.Run("Header should return empty app header", func(t *testing.T) {
		header := app.Header()
		assert.Empty(t, header)
	})

	t.Run("Node should return empty app node", func(t *testing.T) {
		node := app.Node()
		assert.Empty(t, node)
	})
}

func TestApp_setupstate(t *testing.T) {
	app := setup(setupForGeneralInfo)
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
	app := setup(setupForGeneralInfo)

	req := RequestInitChain{}
	// prepare for currencies
	currencies := balance.NewCurrencySet()
	currency := balance.Currency{
		Name: "VT",
	}
	err := currencies.Register(currency)
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
	app := setup(setupForGeneralInfo)

	context := app.Context.Action(&app.header, app.Context.check)
	assert.NotEmpty(t, context)
}

func TestContext_Accounts(t *testing.T) {
	app := setup(setupForGeneralInfo)

	accounts := app.Context.Accounts()
	assert.Empty(t, accounts.Accounts())
}

func TestContext_ValidatorCtx(t *testing.T) {
	app := setup(setupForGeneralInfo)

	validatorContext := app.Context.ValidatorCtx()
	assert.NotNil(t, validatorContext)
}

func TestContext_Balances(t *testing.T) {
	app := setup(setupForGeneralInfo)

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
	app := setup(setupForGeneralInfo)

	node := app.Context.Node()
	assert.Equal(t, "", node.NodeName)
}

func TestContext_Validators(t *testing.T) {
	app := setup(setupForGeneralInfo)

	vs := app.Context.Validators()
	validators, _ := vs.GetValidatorSet()
	assert.Empty(t, validators)
}
