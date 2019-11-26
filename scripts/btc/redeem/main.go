/*

 */

package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func main() {

	addrs := addressess()
	txn, tname := prepareRedeem(addrs[0], "mkW45toPFaa1uyNGV4TXEWWCxyuDC7BbKG", 300000)
	fmt.Println("Tracker for lock: ", tname)

	os.Exit(1)
	signed, signer := sign(base64.StdEncoding.EncodeToString(txn), addrs[0])

	result := broadcastCommit(base64.StdEncoding.EncodeToString(txn),
		base64.StdEncoding.EncodeToString(signed),
		signer)

	fmt.Println(result)
}

func prepareRedeem(address, btcAddress string, amount int64) ([]byte, string) {

	params := map[string]interface{}{
		"address":     address,
		"btc_address": btcAddress,
		"amount":      amount,
		"fees_btc":    50000,
		"gasprice": map[string]interface{}{
			"currency": "OLT",
			"value":    "1000000000",
		},
		"gas": 400000,
	}
	resp, err := makeRPCcall("btc.PrepareRedeem", params)
	if err != nil {
		panic(err)
	}

	txnHex, _ := resp.Result["rawTx"].(string)
	txn, err := base64.StdEncoding.DecodeString(txnHex)
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
