/*

 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type jsonRPCData struct {
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id"`
	JsonRpc float64                `json:"jsonrpc"`
}
type RPCResponse struct {
	Result map[string]interface{} `json:"result"`
}

func makeRPCcall(method string, params map[string]interface{}) (*RPCResponse, error) {

	url := "http://127.0.0.1:26602/jsonrpc"

	payload, _ := json.Marshal(&jsonRPCData{
		Method:  method,
		Params:  params,
		ID:      123,
		JsonRpc: 2.0,
	})

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("json rpc error", err)
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	resp := RPCResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil

}
