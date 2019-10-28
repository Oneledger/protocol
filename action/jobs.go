/*

 */

package action

import (
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
)

type JobsContext struct {
	Service *Service

	Trackers *bitcoin.TrackerStore

	BTCPrivKey       *btcec.PrivateKey
	Params           *chaincfg.Params
	ValidatorAddress Address

	BlockCypherToken string

	LockScripts *bitcoin.LockScriptStore
}

func NewJobsContext(chainType string, svc *Service, trackers *bitcoin.TrackerStore, privKey *btcec.PrivateKey,
	valAddress keys.Address, bcyToken string, lStore *bitcoin.LockScriptStore) *JobsContext {

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

	return &JobsContext{
		Service:          svc,
		Trackers:         trackers,
		BTCPrivKey:       privKey,
		Params:           params,
		ValidatorAddress: valAddress,
		BlockCypherToken: bcyToken,
		LockScripts:      lStore,
	}

}
