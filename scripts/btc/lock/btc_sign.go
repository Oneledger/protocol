/*

 */

package main

import (
	"bytes"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func btcSign(txBytes []byte, wifStr string) []byte {

	wif, _ := btcutil.DecodeWIF(wifStr)

	tx := wire.NewMsgTx(wire.TxVersion)
	buf := bytes.NewBuffer(txBytes)
	tx.Deserialize(buf)

	sc, _ := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(btcutil.Hash160(wif.PrivKey.PubKey().SerializeCompressed())).AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).Script()

	sig, err := txscript.RawTxInSignature(tx, 0, sc, txscript.SigHashAll, wif.PrivKey)
	if err != nil {
		panic(err)
	}

	sigScript, _ := txscript.NewScriptBuilder().
		AddData(sig).AddData(wif.PrivKey.PubKey().SerializeCompressed()).
		Script()

	tx.TxIn[0].SignatureScript = sigScript

	buf = bytes.NewBuffer(nil)
	tx.Serialize(buf)
	txBytes = buf.Bytes()

	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sc, tx, 0, flags, nil, nil, tx.TxOut[0].Value)
	if err != nil {
		panic(err)
	}

	err = vm.Execute()
	if err != nil {
		panic(err)
	}

	return sigScript
}
