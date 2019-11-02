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

	wif, _ := btcutil.DecodeWIF("cRvMpZ6X4DVTb7G94msgfFiefhW1bRfW9193FWzNtnoTc8SDqjLy")

	address, _ := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)

	fmt.Println(address.EncodeAddress())
	txBytes, _ := hex.DecodeString("01000000018675bccc5cf57e10e2c75cb25187d25a3d8d1f2d1107d286df54ed84ef320a860100000000ffffffff016c752e000000000014aae651e577abfe1d951872de8b48232bb8787f7300000000")

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
