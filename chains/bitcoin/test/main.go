/*

 */

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Oneledger/protocol/chains/bitcoin"

	"github.com/Oneledger/protocol/storage"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

func main() {

	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18831",
		User:         "oltest01",
		Pass:         "olpass01",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	clt, err := rpcclient.New(connCfg, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	hashes, err := clt.Generate(101)
	fmt.Println("##################################################################################################")
	fmt.Println("1. Generating blocks", len(hashes), "err", err)
	time.Sleep(time.Second * 5)

	unspent, err := clt.ListUnspent()
	if err != nil {
		fmt.Println(err)
	}
	source := unspent[0]
	sourceHash, _ := chainhash.NewHashFromStr(unspent[0].TxID)

	keychain := NewKeyChain()

	validator1, err := keychain.GetBitcoinWIF(1, &chaincfg.RegressionNetParams)
	if err != nil {
		log.Fatal("validator 1 failed")
	}
	validator2, err := keychain.GetBitcoinWIF(2, &chaincfg.RegressionNetParams)
	if err != nil {
		log.Fatal("validator 2 failed")
	}
	validator3, err := keychain.GetBitcoinWIF(3, &chaincfg.RegressionNetParams)
	if err != nil {
		log.Fatal("validator 3 failed")
	}

	validators := []*btcutil.WIF{validator1, validator2, validator3}

	val1Key, _ := btcutil.NewAddressPubKey(validator1.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.RegressionNetParams)
	val2Key, _ := btcutil.NewAddressPubKey(validator2.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.RegressionNetParams)
	val3Key, _ := btcutil.NewAddressPubKey(validator3.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.RegressionNetParams)
	validatorPubKeys := []*btcutil.AddressPubKey{
		val1Key,
		val2Key,
		val3Key,
	}

	fmt.Println(validatorPubKeys)

	lockScript, lockScriptAddress, err := bitcoin.CreateMultiSigAddress(1, validatorPubKeys)

	lockScriptAddressTemp, _ := btcutil.NewAddressScriptHash(lockScript, &chaincfg.RegressionNetParams)
	lockScriptAddress, err = txscript.PayToAddrScript(lockScriptAddressTemp)
	if err != nil {
		log.Fatal(err)
	}

	validatorKeyDB := bitcoin2.NewKeyDB()
	validatorKeyDB.Add(btcutil.Address(val1Key), validator1.PrivKey)
	validatorKeyDB.Add(btcutil.Address(val2Key), validator2.PrivKey)
	validatorKeyDB.Add(btcutil.Address(val3Key), validator3.PrivKey)

	lockStore := bitcoin2.NewLockStore("primary", "/tmp/oneledger/btc/chaindriver",
		"goleveldb", storage.PERSISTENT)

	utxoInit := bitcoin2.NewUTXO(chainhash.Hash{}, 0, 0)

	err = lockStore.InitializeTracker("tracker1", *utxoInit)
	if err != nil {
		log.Fatal("inti tracker error", err)
	}

	cd := bitcoin2.NewChainDriver()

	amt := big.NewFloat(source.Amount)
	satoshiPerBitcoin := new(big.Int).Exp(big.NewInt(10), big.NewInt(8), nil)
	amtSatoshi := big.NewFloat(0).Mul(amt, big.NewFloat(0.0).SetInt(satoshiPerBitcoin))
	amtSatoshiInt, _ := amtSatoshi.Int64()

	sourceFunds := bitcoin2.NewUTXO(*sourceHash, source.Vout, amtSatoshiInt)
	lockTxBytes := cd.PrepareLock(utxoInit, sourceFunds, nil, lockScriptAddress)
	err = lockStore.UpdateUTXO("tracker1", *sourceFunds)
	if err != nil {
		//		log.Fatal("lockstore lock failed", err)
	}

	lockTx := wire.NewMsgTx(wire.TxVersion)
	lockTx.Deserialize(bytes.NewReader(lockTxBytes))

	lockTx, ok, err := clt.SignRawTransaction(lockTx)
	if !ok {
		fmt.Println(err)
		log.Fatal("node account signing failed")
	}
	if err != nil {
		log.Fatal("sign raw transaction failed", err)
	}

	userSign := lockTx.TxIn[0].SignatureScript
	lockTx = cd.AddUserLockSignature(lockTxBytes, userSign)

	fmt.Println(clt.GetBalance(""))
	lockHash, err := cd.BroadcastTx(lockTx, clt)

	clt.Generate(1)

	fmt.Println(clt.GetBalance(""))

	clt.Generate(10)

	abcd, err := clt.GetTransaction(lockHash)
	if abcd.Confirmations > 10 {
		lockStore.FinalizeTracker("tracker1")
	} else {
		log.Fatal("insufficient confirmations for lock")
	}

	randomGuy, _ := keychain.GetBitcoinWIF(10, &chaincfg.RegressionNetParams)
	randomGuyAddress, _ := btcutil.NewAddressPubKey(randomGuy.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.RegressionNetParams)
	utxo, err := lockStore.GetLatestUTXO("tracker1")
	if err != nil {
		// log.Fatal("{{{{", err)
	}
	utxo.Balance = lockTx.TxOut[0].Value
	utxo.TxID = *lockHash
	utxo.Index = 0

	redeemAmount := utxo.Balance / 2

	sourceAddress, err := btcutil.DecodeAddress(randomGuyAddress.EncodeAddress(), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	outAddress, _ := txscript.PayToAddrScript(sourceAddress)

	redeemTxBytes := cd.PrepareRedeem(utxo, outAddress, redeemAmount, lockScriptAddress)

	redeemTx := wire.NewMsgTx(wire.TxVersion)
	redeemTx.Deserialize(bytes.NewReader(redeemTxBytes))

	builder := txscript.NewScriptBuilder().AddOp(txscript.OP_FALSE)
	signed := 0
	for i := range validators {
		key := validators[i].PrivKey

		sig, err := txscript.RawTxInSignature(redeemTx, 0, lockScript, txscript.SigHashAll, key)
		if err != nil {
			fmt.Println(err, "RawTxInSignature")
		}

		builder.AddData(sig)
		signed++
		break
	}

	builder.AddFullData(lockScript)
	sigScript, err := builder.Script()
	//sigScript = append(sigScript, lockScript...)

	redeemTx = cd.AddLockSignature(redeemTxBytes, sigScript)

	buf := bytes.NewBuffer([]byte{})
	redeemTx.Serialize(buf)

	fmt.Println("before broadcast")
	fmt.Println(hex.EncodeToString(sigScript))
	fmt.Println(hex.EncodeToString(buf.Bytes()))

	fmt.Println(clt.GetBalance(""))
	redeemHash, err := cd.BroadcastTx(redeemTx, clt)
	fmt.Println(redeemHash, err)

	fmt.Println(clt.GetBalance(""))
	clt.Generate(1)

	fmt.Println(clt.GetBalance(""))

	// redeemTx.TxIn[0].SignatureScript = []byte{}
	vm, err := txscript.NewEngine(lockTx.TxOut[0].PkScript, redeemTx, 0, txscript.StandardVerifyFlags, nil, nil, lockTx.TxOut[0].Value)
	if err != nil {
		fmt.Println("new engine", err)
	}
	if err := vm.Execute(); err != nil {
		fmt.Println("vm Execute", err)
	}
}
