/*
	Copywrite 2017-2018 OneLedger
*/
package app

import "encoding/json"

// Response arguments
type ResponseInfo struct {
	Hashes int `json:"hashes"`
	Txs    int `json:"txs"`
}

func NewResponseInfo(hashes int, txs int) *ResponseInfo {
	return &ResponseInfo{
		Hashes: hashes,
		Txs:    txs,
	}
}

// Convert to JSON
func (info *ResponseInfo) Json() string {
	bytes, err := json.Marshal(info)
	if err != nil {
		panic("Marshal Failed")
	}
	return string(bytes)
}
