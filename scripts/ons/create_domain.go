/*

 */

package main

import (
	"fmt"

	"github.com/Oneledger/protocol/client"
	"github.com/powerman/rpc-codec/jsonrpc2"
)

func main() {

	clt, err := jsonrpc2.Dial("tcp", "127.0.0.1:26602")
	if err != nil {
		fmt.Println(err)
	}

	req := client.ONSCreateRequest{Name: "alice.btc"}
	reply := &client.SignRawTxResponse{}

	err = clt.Call("/tx.ONS_CreateRawCreate", req, reply)
	if err != nil {
		fmt.Println("rpc error")
		fmt.Println(err)
	}
}
