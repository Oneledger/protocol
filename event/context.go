/*

 */

package event

import (
	"crypto/ecdsa"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Oneledger/protocol/action"
	bitcoin2 "github.com/Oneledger/protocol/chains/bitcoin"
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

	BTCPrivKey *keys.PrivateKey
	ETHPrivKey *keys.PrivateKey
	BTCParams  *chaincfg.Params

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

func NewJobsContext(cfg config.Server, btcChainType string, svc *Service,
	trackers *bitcoin.TrackerStore, validators *identity.ValidatorStore,
	privKey *keys.PrivateKey, ethprivKey *keys.PrivateKey,
	valAddress keys.Address, bcyToken string, lStore *bitcoin.LockScriptStore,
	btcAddress, btcRPCPort, BTCRPCUsername, BTCRPCPassword string,
	ethTracker *ethereum.TrackerStore,
) *JobsContext {

	params := bitcoin2.GetChainParams(btcChainType)

	w := os.Stdout

	return &JobsContext{
		cfg:              cfg,
		Service:          svc,
		Logger:           log.NewLoggerWithPrefix(w, "internal_jobs"),
		Trackers:         trackers,
		Validators:       validators,
		BTCPrivKey:       privKey,
		ETHPrivKey:       ethprivKey,
		BTCParams:        params,
		ValidatorAddress: valAddress,
		BlockCypherToken: bcyToken,
		LockScripts:      lStore,
		BTCNodeAddress:   btcAddress,
		BTCRPCPort:       btcRPCPort,
		BTCRPCUsername:   BTCRPCUsername,
		BTCRPCPassword:   BTCRPCPassword,
		EthereumTrackers: ethTracker,
		BTCChainnet:      btcChainType,
	}

}

func (jc *JobsContext) GetValidatorETHAddress() common.Address {
	privkey := keys.ETHSECP256K1TOECDSA(jc.ETHPrivKey.Data)

	pubkey := privkey.Public()
	ecdsapubkey, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		jc.Logger.Error("failed to cast pubkey", pubkey)
		return common.Address{}
	}
	addr := crypto.PubkeyToAddress(*ecdsapubkey)
	return addr
}

func (jc *JobsContext) GetValidatorETHPrivKey() *ecdsa.PrivateKey {
	privkey := keys.ETHSECP256K1TOECDSA(jc.ETHPrivKey.Data)

	return privkey
}
