/*

 */

package bitcoin

import (
	"bytes"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func ValidateLock(tx *wire.MsgTx, token, chainType string, trackerPrevTxID *chainhash.Hash,
	lockScriptAddress []byte) bool {

	if !(len(tx.TxIn) == 1 || len(tx.TxIn) == 2) {
		return false
	}

	if len(tx.TxOut) != 1 {
		return false
	}

	// 1.0
	if trackerPrevTxID != nil {
		if tx.TxIn[0].PreviousOutPoint.Hash != *trackerPrevTxID {
			return false
		}
	}

	// 2, 3
	var input int64
	for i := range tx.TxIn {
		h := tx.TxIn[i].PreviousOutPoint.Hash
		index := tx.TxIn[i].PreviousOutPoint.Index

		btc := gobcy.API{token, "btc", chainType}

		txIn, err := btc.GetTX(h.String(), nil)
		if err != nil {
			return false
		}

		if txIn.Outputs[index].SpentBy != "" {
			return false
		}

		input += int64(txIn.Outputs[index].Value)
	}

	// 4
	output := tx.TxOut[0].Value
	fees := input - output

	if fees < 40-000 {
		return false
	}

	// 5
	if !bytes.Equal(tx.TxOut[0].PkScript, lockScriptAddress) {
		return false
	}

	return true
}
