/*
	Copyright 2017 - 2018 OneLedger

	Query the chain for answers
*/
package shared

import (
	"bytes"
	"encoding/hex"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

func GetAccountKey(identity string) []byte {
	request := action.Message("Identity=" + identity)
	response := comm.Query("/accountKey", request)
	if response == nil || response.Response.Value == nil {
		log.Fatal("Failed to resolve Identity")
	}

	return response.Response.Value
}

func GetBalance(accountKey id.AccountKey) data.Coin {
	request := action.Message("accountKey=" + hex.EncodeToString(accountKey))
	response := comm.Query("/balance", request)

	if response == nil {
		log.Fatal("Failed to get Balance", "response", response, "key", "accountKey")
	}

	if response.Response.Value == nil || len(response.Response.Value) == 0 {
		log.Fatal("Failed to return Balance", "response", response, "key", "accountKey")
	}

	if bytes.Compare(response.Response.Value, []byte("null")) == 0 {
		log.Fatal("Failed to return Balance", "response", response, "key", "accountKey")
	}

	var balance data.Balance
	buffer, _ := comm.Deserialize(response.Response.Value, &balance)
	if buffer != nil {
		log.Debug("Deserialize", "buffer", buffer, "response", response, "value", response.Response.Value)
		return buffer.(*data.Balance).Amount
	}

	log.Fatal("Can't deserialize", "response", response)
	return data.NewCoin(0, "OLT")
}
