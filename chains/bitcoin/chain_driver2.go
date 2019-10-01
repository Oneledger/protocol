/*

 */

package bitcoin

import (
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type ChainDriver2 interface {

}


type btcChainDriver struct {

}

func (cd *btcChainDriver) PrepareLock(store *bitcoin.TrackerStore, inputHash chainhash.Hash, fessSatoshi) ([]byte, error) {

	tracker, err := store.GetTrackerForLock()
	if err != nil {
		return nil, err
	}

	tracker.LatestUTXO
}