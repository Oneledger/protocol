/*
	Copyright 2017-2018 OneLedger

	Return a Response to an Info messages
*/
package abci

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/convert"
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
	/*
		bytes, err := json.Marshal(info)
		if err != nil {
			panic("Marshal Failed")
		}
		return string(bytes)
	*/
}

// Return as a Contract
func (info *ResponseInfo) JSONMessage() comm.Message {
	return comm.Message(info.JSON())
}
