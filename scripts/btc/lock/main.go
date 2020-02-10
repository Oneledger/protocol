/*

 */

package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Oneledger/protocol/client"
)

func main() {

	sourceBTCHash := "b7c8d66068fccecccd2a9b29bf59fc5eb2f4117378fe53791dc3ee618077d5df"
	var sourceBTCIndex uint32 = 0

	sourceBTCHash2 := "232a164b54952dbcadb287608500463cd75e6e44d9a936e5ff6cb88edb74f887"
	var sourceBTCIndex2 uint32 = 1

	wif := "cPwDvkefgLhtWiMJYapEm4MuKhxAiRDjNFZSUHkenChAZVooaVr5"

	txn, tname := prepareLock(sourceBTCHash, sourceBTCIndex, sourceBTCHash2, sourceBTCIndex2,
		1200000, 30, "mkW45toPFaa1uyNGV4TXEWWCxyuDC7BbKG")
	fmt.Println("Received response of PrepareLock")
	fmt.Println("Tracker for lock: ", tname)
	fmt.Println("BTC Unsigned Txn: ", hex.EncodeToString(txn))

	time.Sleep(3 * time.Second)

	fmt.Println(hex.EncodeToString(txn))
	// os.Exit(1)

	addrs := addressess()
	fmt.Println("Will lock to OLT Address: ", addrs[0])

	signedTxn := btcSign(txn, wif, 0)
	fmt.Println(hex.EncodeToString(signedTxn), "======================")
	signedTxn = btcSign(signedTxn, wif, 1)
	fmt.Println(hex.EncodeToString(signedTxn), "======================")

	rawTx := addSignature(base64.StdEncoding.EncodeToString(signedTxn),
		addrs[0], tname)

	signed, signer := sign(base64.StdEncoding.EncodeToString(rawTx), addrs[0])

	time.Sleep(30 * time.Second)
	result := broadcastCommit(base64.StdEncoding.EncodeToString(rawTx),
		base64.StdEncoding.EncodeToString(signed),
		signer)

	fmt.Println(result)
}

func prepareLock(txHash string, index uint32, txHash2 string, index2 uint32, amount int64, feeRate int64, return_address string) ([]byte, string) {

	inputs := []client.InputTransaction{
		{txHash, index},
		{txHash2, index2},
	}

	params := map[string]interface{}{
		"inputs":         inputs,
		"amount":         amount,
		"fee_rate":       feeRate,
		"return_address": return_address,
	}
	resp, err := makeRPCcall("btc.PrepareLock", params)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

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

func addSignature(txn, addr, trackerName string) []byte {
	params := map[string]interface{}{
		"txn":          txn,
		"address":      addr,
		"tracker_name": trackerName,
		"gasprice": map[string]interface{}{
			"currency": "OLT",
			"value":    "1000000000",
		},
		"gas": 800000,
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
