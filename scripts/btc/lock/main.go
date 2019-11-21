/*

 */

package main

import "fmt"

func main() {
	sourceBTCHash := ""
	sourceBTCIndex := 0

	prepareLock(sourceBTCHash, sourceBTCIndex)
}

func prepareLock(txHash string, index int) {
	params := map[string]interface{}{
		"hash":     txHash,
		"index":    index,
		"btc_fees": 50000,
	}
	resp, err := makeRPCcall("btc.PrepareLock", params)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Result["txn"])
	fmt.Println(resp.Result["tracker_name"])
}
