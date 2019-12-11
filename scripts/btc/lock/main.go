/*

 */

package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	sourceBTCHash := "08eb169a33cd53f150f412a5a11fa7f6412f8922fea0c0ba3a430d8b6d3d1d97"
	sourceBTCIndex := 0
	wif := "cMe6v24SwBTyhsBaJfieGvJGWcw6wgpE4anE8nX5hrQdP8vaBsym"

	txn, tname := prepareLock(sourceBTCHash, sourceBTCIndex)
	fmt.Println("Received response of PrepareLock")
	fmt.Println("Tracker for lock: ", tname)
	fmt.Println("BTC Unsigned Txn: ", hex.EncodeToString(txn))

	time.Sleep(30 * time.Second)

	fmt.Println(hex.EncodeToString(txn))
	// os.Exit(1)

	addrs := addressess()
	fmt.Println("Will lock to OLT Address: ", addrs[0])

	btcSignature := btcSign(txn, wif)
	rawTx := addSignature(base64.StdEncoding.EncodeToString(txn),
		base64.StdEncoding.EncodeToString(btcSignature),
		addrs[0], tname)

	signed, signer := sign(base64.StdEncoding.EncodeToString(rawTx), addrs[0])

	time.Sleep(30 * time.Second)
	result := broadcastCommit(base64.StdEncoding.EncodeToString(rawTx),
		base64.StdEncoding.EncodeToString(signed),
		signer)

	fmt.Println(result)
}

func prepareLock(txHash string, index int) ([]byte, string) {
	params := map[string]interface{}{
		"hash":     txHash,
		"index":    index,
		"fees_btc": 50000,
	}
	resp, err := makeRPCcall("btc.PrepareLock", params)
	if err != nil {
		panic(err)
	}

	txnHex, _ := resp.Result["txn"].(string)
	txn, err := hex.DecodeString(txnHex)
	if err != nil {
		panic(err)
	}

	tracker_name, _ := resp.Result["tracker_name"].(string)
	if len(txn) == 0 ||
		tracker_name == "" {

		panic("prepareLock failed")
	}

	return txn, tracker_name
}

func addressess() []string {
	resp, err := makeRPCcall("owner.ListAccountAddresses", map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	add, ok := resp.Result["addresses"].([]interface{})
	if !ok {
		panic("failed to get address")
	}
	strs := []string{}
	for i := range add {
		strs = append(strs, add[i].(string))
	}
	return strs
}

func addSignature(txn, sign, addr, trackerName string) []byte {
	params := map[string]interface{}{
		"txn":          txn,
		"signature":    sign,
		"address":      addr,
		"tracker_name": trackerName,
		"gasprice": map[string]interface{}{
			"currency": "OLT",
			"value":    "1000000000",
		},
		"gas": 400000,
	}
	resp, err := makeRPCcall("btc.AddUserSignatureAndProcessLock", params)
	if err != nil {
		panic(err)
	}

	oltTxnB64, _ := resp.Result["rawTx"].(string)
	oltTxn, err := base64.StdEncoding.DecodeString(oltTxnB64)
	if err != nil {
		panic(err)
	}

	return oltTxn
}

func sign(rawTx, address string) ([]byte, interface{}) {
	resp, err := makeRPCcall("owner.SignWithAddress",
		map[string]interface{}{
			"rawTx":   rawTx,
			"address": address,
		})
	if err != nil {
		panic(err)
	}

	signature, ok := resp.Result["signature"].(map[string]interface{})
	if !ok {
		fmt.Println(resp.Result)
		panic("failed to get signature")
	}
	signedStr := signature["Signed"].(string)
	signed, err := base64.StdEncoding.DecodeString(signedStr)
	if err != nil {
		panic(err)
	}

	signerStr := signature["Signer"]

	return signed, signerStr
}

func broadcastCommit(rawTx, signature string, pubKey interface{}) map[string]interface{} {
	resp, err := makeRPCcall("broadcast.TxCommit",
		map[string]interface{}{
			"rawTx":     rawTx,
			"signature": signature,
			"publicKey": pubKey,
		})
	if err != nil {
		panic(err)
	}

	return resp.Result
}
