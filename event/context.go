/*

 */

package event

import (
	"crypto/ecdsa"
	"os"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
)

type JobsContext struct {
	cfg     config.Server
	Service *Service
	Logger  *log.Logger

	Trackers   *bitcoin.TrackerStore
	Validators *identity.ValidatorStore

	BTCPrivKey       keys.PrivateKey
	ETHPrivKey       ecdsa.PrivateKey
	Params           *chaincfg.Params
	ValidatorAddress action.Address

	BlockCypherToken string

	LockScripts *bitcoin.LockScriptStore

	BTCNodeAddress string
	BTCRPCPort     string
	BTCRPCUsername string
	BTCRPCPassword string

	EthereumTrackers *ethereum.TrackerStore

	BTCChainnet string
}

func NewJobsContext(cfg config.Server, btcchainType string, svc *Service,
	trackers *bitcoin.TrackerStore, validators *identity.ValidatorStore,
	privKey *keys.PrivateKey, ethprivKey *ecdsa.PrivateKey,
	valAddress keys.Address, bcyToken string, lStore *bitcoin.LockScriptStore,
	btcAddress, btcRPCPort, BTCRPCUsername, BTCRPCPassword, btcChain string,
	ethTracker *ethereum.TrackerStore,
) *JobsContext {

	var params *chaincfg.Params
	switch btcchainType {
	case "mainnet":
		params = &chaincfg.MainNetParams
	case "testnet3":
		params = &chaincfg.TestNet3Params
	case "regtest":
		params = &chaincfg.RegressionNetParams
	case "simnet":
		params = &chaincfg.SimNetParams
	default:
		params = &chaincfg.TestNet3Params
	}

	w := os.Stdout

	return &JobsContext{
		cfg:              cfg,
		Service:          svc,
		Logger:           log.NewLoggerWithPrefix(w, "internal_jobs"),
		Trackers:         trackers,
		Validators:       validators,
		BTCPrivKey:       *privKey,
		ETHPrivKey:       *ethprivKey,
		Params:           params,
		ValidatorAddress: valAddress,
		BlockCypherToken: bcyToken,
		LockScripts:      lStore,
		BTCNodeAddress:   btcAddress,
		BTCRPCPort:       btcRPCPort,
		BTCRPCUsername:   BTCRPCUsername,
		BTCRPCPassword:   BTCRPCPassword,
		EthereumTrackers: ethTracker,
		BTCChainnet:      btcChain,
	}

}
