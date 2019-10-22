/*

 */

package action

import (
	"github.com/Oneledger/protocol/data/bitcoin"
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

	LockScripts bitcoin.LockScriptStore
}
