/*

 */

package event

import (
	"crypto/ecdsa"
	"os"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/Oneledger/protocol/action"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"
)

type JobsContext struct {
	Service *Service
	Logger  *log.Logger

	Trackers *bitcoin.TrackerStore

	BTCPrivKey       keys.PrivateKey
	ETHPrivKey       ecdsa.PrivateKey
	BTCParams        *chaincfg.Params
	ValidatorAddress action.Address

	BlockCypherToken string

	LockScripts *bitcoin.LockScriptStore

	BTCNodeAddress string
	BTCRPCPort     string
	BTCRPCUsername string
	BTCRPCPassword string

	EthereumTrackers   *ethereum.TrackerStore
	ETHContractABI     string // Replace 39,40,41 with ethchaindriverconfig
	ETHConnection      string
	ETHContractAddress string

	BTCChainnet string
}

func NewJobsContext(chainType string, svc *Service,
	trackers *bitcoin.TrackerStore, privKey *keys.PrivateKey, ethprivKey *ecdsa.PrivateKey,
	valAddress keys.Address, bcyToken string, lStore *bitcoin.LockScriptStore,
	btcAddress, btcRPCPort, BTCRPCUsername, BTCRPCPassword, btcChain string,
	ETHAbi string, ETHconn string, ETHContractaddress string, ethTracker *ethereum.TrackerStore,
) *JobsContext {

	params := bitcoin2.GetChainParams(chainType)

	w := os.Stdout

	return &JobsContext{
		Service:            svc,
		Logger:             log.NewLoggerWithPrefix(w, "internal_jobs"),
		Trackers:           trackers,
		BTCPrivKey:         *privKey,
		ETHPrivKey:         *ethprivKey,
		BTCParams:          params,
		ValidatorAddress:   valAddress,
		BlockCypherToken:   bcyToken,
		LockScripts:        lStore,
		BTCNodeAddress:     btcAddress,
		BTCRPCPort:         btcRPCPort,
		BTCRPCUsername:     BTCRPCUsername,
		BTCRPCPassword:     BTCRPCPassword,
		ETHConnection:      ETHconn,
		ETHContractAddress: ETHContractaddress,
		ETHContractABI:     ETHAbi,
		EthereumTrackers:   ethTracker,
		BTCChainnet:        btcChain,
	}

}
