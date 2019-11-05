/*

 */

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func main() {

	wif, _ := btcutil.DecodeWIF("cSxM9B2KMPFa5k8cC8VnMN5jyWG2FH3e5RCKQ2bpWbjbQvX6tW1j")

	address, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)

	fmt.Println(address.EncodeAddress())
	txBytes, _ := hex.DecodeString("0100000001d6e7493a71bd2566d530b94ca6548199553190496ba0548a90905e879f921f4c0100000000ffffffff0140290f0000000000140678fa5af71cdcfbfa849cb439bae319ba27f0ff00000000")

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

	fmt.Println(hex.EncodeToString(txBytes))

	fmt.Println(err)

	fmt.Println(hex.EncodeToString(sigScript))
	fmt.Println("above sigscript")

	flags := txscript.StandardVerifyFlags
	vm, err := txscript.NewEngine(sc, tx, 0, flags, nil, nil, tx.TxOut[0].Value)

	err = vm.Execute()
	fmt.Println(err)

}
