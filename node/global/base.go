/*
	Copyright 2017-2018 OneLedger

	The overall running context. Initialized right away, but is mutable.

	Contains the main variables.

	Precedence:
		- Default values
	 	- Environment variables (like $OLROOT)
		- Configuration files
		- Command line arguments
		- Overrides
*/
package global

import (
	"os"
	"path/filepath"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/persist"
	"github.com/mitchellh/go-homedir"
	tmnode "github.com/tendermint/tendermint/node"
)

var Current *Context

type Context struct {
	Application persist.Access // Global Access to the application when it is running

	Debug            bool // DEBUG flag
	DisablePasswords bool // DEBUG flag

	ConfigName      string // The Name of the config file (without extension)
	NodeName        string // Name of this instance
	NodeAccountName string // TODO: Should be a list of accounts
	PaymentAccount  string
	NodeIdentity    string
	RootDir         string // Working directory for this instance

	// Do we really need this?
	Sequence int64 // replay protection

	// This should be set to a function, shouldn't be specifiable
	TendermintAddress string
	TendermintPubKey  string

	PersistentPeers []string
	Seeds           string
	SeedMode        bool
	P2PAddress      string

	ConsensusNode *tmnode.Node

	//Minimum Fees
	MinSendFee     float64
	MinSwapFee     float64
	MinContractFee float64
	MinRegisterFee float64

	Config *config.Server
}

func init() {
	Current = NewContext("OneLedger")
}

// Set the default values for any context variables here
func NewContext(name string) *Context {
	var debug = false
	if os.Getenv("OLDEBUG") == "true" {
		debug = true
	}

	defaultRootDir := os.Getenv("OLDATA")
	if defaultRootDir == "" {
		home, _ := homedir.Dir()
		defaultRootDir = filepath.Join(home, ".olfullnode")
	}

	return &Context{
		Debug:            debug,
		DisablePasswords: true,

		ConfigName:      "config", // TODO: needs to deal with client/server
		NodeName:        name,
		NodeAccountName: "",
		PaymentAccount:  "Payment",
		RootDir:         defaultRootDir,

		// TODO: Should be params in the chain
		MinSendFee:     0.1,
		MinSwapFee:     0.1,
		MinContractFee: 0.1,
		MinRegisterFee: 0.1,
		Config:         config.DefaultServerConfig(),
	}
}

func (ctx *Context) SetApplication(app persist.Access) persist.Access {
	ctx.Application = app
	return app
}

func (ctx *Context) SetConsensusNode(node *tmnode.Node) {
	ctx.ConsensusNode = node
}

func (ctx *Context) GetApplication() persist.Access {
	return ctx.Application
}

func (ctx *Context) ConsensusDir() string {
	result, _ := filepath.Abs(filepath.Join(Current.RootDir, "consensus"))
	return result
}

func (ctx *Context) DatabaseDir() string {
	result, _ := filepath.Abs(filepath.Join(Current.RootDir, "nodedata"))
	return result
}

// ReadConfig looks for a configuration file in the rootdir
func (ctx *Context) ReadConfig() error {
	var cfg config.Server
	err := cfg.ReadFile(filepath.Join(ctx.RootDir, config.FileName))
	if err != nil {
		return err
	}
	ctx.Config = &cfg
	return nil
}

// SaveConfig writes the currently loaded configuration onto a config.toml file in the root directory
func (ctx *Context) SaveConfig() error {
	return ctx.Config.SaveFile(ctx.RootDir)
}
