/*

 */

package bitcoin

import (
	"bytes"
	"fmt"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func ValidateLock(tx *wire.MsgTx, token, chainType string, trackerPrevTxID *chainhash.Hash,
	lockScriptAddress []byte, currentBalance, lockAmount int64) bool {

	if !(len(tx.TxIn) == 1 || len(tx.TxIn) == 2) {
		fmt.Println("btc lock validate err, TxIn should be 1 or 2")
		return false
	}

	if len(tx.TxOut) != 1 {

		fmt.Println("btc lock validate err, TxOut should be 1 ")
		return false
	}

	// 1.0
	if trackerPrevTxID != nil {
		if tx.TxIn[0].PreviousOutPoint.Hash != *trackerPrevTxID {

			fmt.Println("btc lock validate err, input 0 hash should match tracker hash")
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

			fmt.Println("btc lock validate err, error finding txIn", i, err)
			return false
		}

		if txIn.Outputs[index].SpentBy != "" {

			fmt.Println("btc lock validate err, not spendable txIn", i)
			return false
		}

		input += int64(txIn.Outputs[index].Value)
	}

	if lockAmount > (input - currentBalance) {
		fmt.Println("btc lock validate err, input should be sum of curr bal and lock amount")
		return false
	}

	// 4
	output := tx.TxOut[0].Value
	fees := input - output

	if fees < 40-000 {

		fmt.Println("btc lock validate err, fees should be more than 40000")
		return false
	}

	// 5
	if !bytes.Equal(tx.TxOut[0].PkScript, lockScriptAddress) {

		fmt.Println("btc lock validate err, txout pkscript must be lockscript address")
		return false
	}

	return true
}

func ValidateRedeem(tx *wire.MsgTx, token, chainType string, trackerPrevTxID *chainhash.Hash,
	lockScriptAddress []byte, currentBalance, redeemAmount int64) bool {

	if !(len(tx.TxIn) == 1) {
		fmt.Println("redeem validate err, TxIn should be 1")
		return false
	}

	if len(tx.TxOut) != 2 {
		fmt.Println("redeem validate err, TxOut should be 2")
		return false
	}

	// 1.0

	if tx.TxIn[0].PreviousOutPoint.Hash != *trackerPrevTxID {

		fmt.Println("redeem validate err, first txn in should be tracker current TxId")
		return false
	}

	// 2, 3
	var input int64
	for i := range tx.TxIn {
		h := tx.TxIn[i].PreviousOutPoint.Hash
		index := tx.TxIn[i].PreviousOutPoint.Index

		btc := gobcy.API{token, "btc", chainType}

		txIn, err := btc.GetTX(h.String(), nil)
		if err != nil {
			fmt.Println(i, "redeem validate err, TxIn must exist on chain", err)
			return false
		}

		if txIn.Outputs[index].SpentBy != "" {
			fmt.Println(i, "redeem validate err, TxIn must be spendable")
			return false
		}

		input += int64(txIn.Outputs[index].Value)
	}

	// 4
	output := tx.TxOut[0].Value + tx.TxOut[1].Value
	fees := input - output

	if fees < 40-000 {
		fmt.Println("redeem validate error, fess should be > 40,000")
		return false
	}

	diff := input - tx.TxOut[0].Value
	if diff != redeemAmount {
		fmt.Println("redeem validate err, next tracker balance minus current tracker balance should be equal to redeem amount")
		return false
	}
	if input != currentBalance {
		fmt.Println("redeem validate err, input should be equal to currentbalance")
		return false
	}

	// 5
	if !bytes.Equal(tx.TxOut[0].PkScript, lockScriptAddress) {
		fmt.Println("redeem validate err, first tx out should be tracker lockscript address")
		return false
	}

	return true
}
