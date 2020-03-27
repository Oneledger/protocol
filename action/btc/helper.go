/*

 */

package btc

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"

	"github.com/Oneledger/protocol/data/bitcoin"
)

func ValidateExtLockStructure(tracker *bitcoin.Tracker, tx *wire.MsgTx, params *chaincfg.Params) bool {

	// Validate outputs

	// Tx out should be single or two at a time
	if len(tx.TxOut) > 2 || len(tx.TxOut) < 1 {
		return false
	}

	// tracker ProcessLockScriptAddress should be equal to txout address
	if !bytes.Equal(tx.TxOut[0].PkScript, tracker.ProcessLockScriptAddress) {
		return false
	}

	// if there is return address is present
	if len(tx.TxOut) == 2 {

		// it must be a valid btc address
		// TODO

	}

	// Validate Inputs

	// if this is not a first transaction for tracker
	if tracker.CurrentTxId != nil {

		if len(tx.TxIn) < 2 {
			return false
		}

		// then the first input should be tracker balance
		if !tx.TxIn[0].PreviousOutPoint.Hash.IsEqual(tracker.CurrentTxId) {
			return false
		}

		// all user inputs must be signed
		for i := range tx.TxIn {
			if i == 0 {
				continue
			}

			if len(tx.TxIn[i].SignatureScript) == 0 {
				return false
			}
		}

	} else {
		// if this is the first transaction for the tracker

		if len(tx.TxIn) < 1 {
			return false
		}

		// all user inputs must be signed
		for i := range tx.TxIn {
			if len(tx.TxIn[i].SignatureScript) == 0 {
				return false
			}
		}
	}

	return true
}
