/*

 */

package internal

import (
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type CheckFinalityData struct {
	txID *chainhash.Hash
}

func CheckFinality(data interface{}) {

	checkFinData := data.(*CheckFinalityData)

	cd := bitcoin.NewChainDriver("abcd")
	isFinalized, err := cd.CheckFinality(checkFinData.txID)
	if err != nil {
		return
	}

	if !isFinalized {

		return
	}

}
