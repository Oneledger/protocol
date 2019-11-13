/*

 */

package action

import (
	"os"

	hdwallet "github.com/Oneledger/hdkeychain"

	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"
	"github.com/btcsuite/btcd/chaincfg"
)

type JobsContext struct {
	Service *Service
	Logger  *log.Logger

	Trackers *bitcoin.TrackerStore

	HDWallet         hdwallet.HDWallet
	Params           *chaincfg.Params
	ValidatorAddress Address

	BlockCypherToken string

	LockScripts *bitcoin.LockScriptStore

	BTCNodeAddress string
	BTCRPCPort     string
	BTCRPCUsername string
	BTCRPCPassword string

	BTCChainnet string
}

func NewJobsContext(chainType string, svc *Service, trackers *bitcoin.TrackerStore, hdw hdwallet.HDWallet,
	valAddress keys.Address, bcyToken string, lStore *bitcoin.LockScriptStore,
	btcAddress, btcRPCPort, BTCRPCUsername, BTCRPCPassword, btcChain string,
) *JobsContext {

	var params *chaincfg.Params
	switch chainType {
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
		Service:          svc,
		Logger:           log.NewLoggerWithPrefix(w, "internal_jobs"),
		Trackers:         trackers,
		HDWallet:         hdw,
		Params:           params,
		ValidatorAddress: valAddress,
		BlockCypherToken: bcyToken,
		LockScripts:      lStore,

		BTCNodeAddress: btcAddress,
		BTCRPCPort:     btcRPCPort,
		BTCRPCUsername: BTCRPCUsername,
		BTCRPCPassword: BTCRPCPassword,

		BTCChainnet: btcChain,
	}

}
