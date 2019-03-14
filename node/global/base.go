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

	"github.com/Oneledger/protocol/node/persist"
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

	RpcAddress string // rpc address
	Transport  string // socket vs grpc

	AppAddress string // app address

	BTCAddress string // Bitcoin node Address port
	ETHAddress string // Ethereum node Address port

	SDKAddress string // SDK RPC address

	OLVMProtocol string // Config for the OLVM
	OLVMAddress  string

	Sequence int64 // replay protection

	TendermintRoot    string
	TendermintAddress string
	TendermintPubKey  string

	PersistentPeers    string
	Seeds              string
	SeedMode           bool
	P2PAddress         string
	ExternalP2PAddress string

	ConsensusNode *tmnode.Node

	//Minimum Fees
	MinSendFee     float64
	MinSwapFee     float64
	MinContractFee float64
	MinRegisterFee float64
}

func init() {
	Current = NewContext("OneLedger")
}

// Set the default values for any context variables here (and no where else)
func NewContext(name string) *Context {
	var debug = false
	if os.Getenv("OLDEBUG") == "true" {
		debug = true
	}

	return &Context{
		Debug:            debug,
		DisablePasswords: true,

		ConfigName:      "config", // TODO: needs to deal with client/server
		NodeName:        name,
		NodeAccountName: "",
		PaymentAccount:  "Payment",
		RootDir:         os.Getenv("OLDATA") + "/" + name,

		SDKAddress: "http://127.0.01:6900",

		// TODO: Should be params in the chain
		MinSendFee:     0.1,
		MinSwapFee:     0.1,
		MinContractFee: 0.1,
		MinRegisterFee: 0.1,
	}
}

func (context *Context) SetApplication(app persist.Access) persist.Access {
	context.Application = app
	return app
}

func (context *Context) SetConsensusNode(node *tmnode.Node) {
	context.ConsensusNode = node
}

func (context *Context) GetApplication() persist.Access {
	return context.Application
}

func ConsensusDir() string {
	result, _ := filepath.Abs(filepath.Join(Current.RootDir, "consensus"))
	return result
}

func DatabaseDir() string {
	result, _ := filepath.Abs(filepath.Join(Current.RootDir, "nodedata"))
	return result
}
