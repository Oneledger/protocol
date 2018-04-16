/*
	Copyright 2017-2018 OneLedger

	Return a Response to an Info messages
*/
package app

import "encoding/json"

// Response arguments
type ResponseInfo struct {
	//Hashes int `json:"hashes"`
	//Txs    int `json:"txs"`
	Size int `json:"size"`
}

func NewResponseInfo(hashes int, txs int, size int) *ResponseInfo {
	return &ResponseInfo{
		//Hashes: hashes,
		//Txs:    txs,
		Size: size,
	}
}

// Convert to JSON
func (info *ResponseInfo) JSON() string {
	bytes, err := json.Marshal(info)
	if err != nil {
		panic("Marshal Failed")
	}
	return string(bytes)
}

// Return as a Message
func (info *ResponseInfo) JSONMessage() Message {
	return Message(info.JSON())
}
