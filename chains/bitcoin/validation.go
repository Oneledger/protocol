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

func ValidateLock(tx *wire.MsgTx, token, chainType string, lockScriptAddress []byte, currentBalance, lockAmount int64, isFirstlock bool) bool {

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
	if len(tx.TxOut) == 0 {
		fmt.Println("btc lock validate error, output is 0")
		return false
	}

	output := tx.TxOut[0].Value

	if len(tx.TxOut) > 2 {
		fmt.Println("btc lock validate error, tx output is more than 2")
		return false
	}

	if len(tx.TxOut) == 2 {
		output += tx.TxOut[1].Value
	}

	fees := input - output
	txSize := estimateTxSize(tx, isFirstlock)
	fees_per_byte := fees / int64(txSize)

	if fees_per_byte < 20 || fees_per_byte > 70 {

		fmt.Println("btc lock validate err, fees should be more than 20 per byte or less than 70 per byte")
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

	txSize := estimateTxSize(tx, false)
	fees_per_byte := fees / int64(txSize)

	if fees_per_byte < 20 || fees_per_byte > 70 {
		fmt.Println("redeem validate error, fees per byte should be more than 20 and less than 70")
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

func estimateTxSize(tx *wire.MsgTx, isFirstLock bool) int {

	if isFirstLock {
		return tx.SerializeSize()
	}

	sigScriptSize := 46 + 74*6 + 34*8
	return tx.SerializeSize() + sigScriptSize
}

func EstimateTxSizeBeforeUserSign(tx *wire.MsgTx, isFirstLock bool) int {

	p2pkhSigSize := 146

	inputSigsSize := p2pkhSigSize * len(tx.TxIn)

	if isFirstLock {

		return tx.SerializeSize() + inputSigsSize + 20
	}

	sigScriptSize := 46 + 74*6 + 34*8
	return tx.SerializeSize() + sigScriptSize + inputSigsSize + 20
}
