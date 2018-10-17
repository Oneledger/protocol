/*
	Copyright 2017-2018 OneLedger

	Return a Response to an Info messages
*/
package abci

import (
	"github.com/Oneledger/protocol/node/convert"
	"github.com/Oneledger/protocol/node/serial"
)

// Response arguments
type ResponseInfo struct {
	//Hashes int `json:"hashes"`
	//Tx    int `json:"txs"`
	Size int `json:"size"`
}

func NewResponseInfo(hashes int, txs int, size int) *ResponseInfo {
	return &ResponseInfo{
		//Hashes: hashes,
		//Tx:    txs,
		Size: size,
	}
}

// Convert to JSON
func (info *ResponseInfo) JSON() string {
	bytes, err := convert.ToJSON(info)
	if err != nil {
		// TODO: Replace this with real error handling
		panic("JSON conversion failed")
	}
	return string(bytes)
}

// Return as a Contract
func (info *ResponseInfo) JSONMessage() serial.Message {
	return serial.Message(info.JSON())
}
