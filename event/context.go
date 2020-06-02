/*

 */

package event

import (
	"crypto/ecdsa"
	"github.com/Oneledger/protocol/data/governance"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

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

	BTCPrivKey *keys.PrivateKey
	ETHPrivKey *keys.PrivateKey

	ValidatorAddress action.Address
	LockScripts      *bitcoin.LockScriptStore
	EthereumTrackers *ethereum.TrackerStore

	ProposalMaster *governance.ProposalMasterStore
}

func NewJobsContext(cfg config.Server,
	svc *Service,
	trackers *bitcoin.TrackerStore,
	validators *identity.ValidatorStore,
	privKey *keys.PrivateKey,
	ethprivKey *keys.PrivateKey,
	valAddress keys.Address,
	lStore *bitcoin.LockScriptStore,
	ethTracker *ethereum.TrackerStore,
	proposalMaster *governance.ProposalMasterStore,
	logger *log.Logger,
) *JobsContext {

	return &JobsContext{
		cfg:              cfg,
		Service:          svc,
		Logger:           logger,
		Trackers:         trackers,
		Validators:       validators,
		BTCPrivKey:       privKey,
		ETHPrivKey:       ethprivKey,
		ValidatorAddress: valAddress,
		LockScripts:      lStore,
		EthereumTrackers: ethTracker,
		ProposalMaster:   proposalMaster,
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
