/*

 */

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type jsonRPCData struct {
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id"`
	JsonRpc string                 `json:"jsonrpc"`
}
type RPCResponse struct {
	Result map[string]interface{} `json:"result"`
	Error  map[string]interface{} `json:"error"`
}

func makeRPCcall(method string, params map[string]interface{}) (*RPCResponse, error) {

	url := "http://127.0.0.1:26602/jsonrpc"

	payload, _ := json.Marshal(&jsonRPCData{
		Method:  method,
		Params:  params,
		ID:      51,
		JsonRpc: "2.0",
	})

	req, _ := http.NewRequest("POST", url,
		bytes.NewReader(payload))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("json rpc error", err)
		return nil, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		panic(string(body))
	}

	fmt.Println(string(body))
	resp := RPCResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		fmt.Println(resp.Error)
		return nil, errors.New("rpc error")
	}

	return &resp, nil

}
