// alias.go holds some type aliases from go-ethereum
package ethereum

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

type (
	Address      = common.Address
	TrackerName  = common.Hash
	TransactionHash = common.Hash
	Client       = ethclient.Client
	Transaction  = types.Transaction
	TransactOpts = bind.TransactOpts
	CallOpts     = bind.CallOpts
	Genesis      = core.Genesis
)

// Alias some helper funcs
var (
	IsHexAddress = common.IsHexAddress
	FromHex      = common.FromHex
	Wei          = params.Wei
	GWei         = params.GWei
	Ether        = params.Ether
)
