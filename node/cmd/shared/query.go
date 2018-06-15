/*
	Copyright 2017 - 2018 OneLedger

	Query the chain for answers
*/
package shared

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
)

func GetAccountKey(identity string) []byte {
	request := action.Message("Identity=" + identity)
	response := comm.Query("/accountKey", request)
	key := response.Response.Value
	return key
}
