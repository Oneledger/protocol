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
		log.Fatal("No Response from Node", "identity", identity)
	}

	value := response.Response.Value
	if value == nil || len(value) == 0 {
		log.Fatal("Key is Missing", "identity", identity)
	}

	key, status := hex.DecodeString(string(value))
	if status != nil {
		log.Fatal("Decode Failed", "identity", identity, "value", value)
	}

	return key
}

func GetSwapAddress(currencyName string) string {
	request := action.Message("currency=" + currencyName)
	response := comm.Query("/swapAddress", request)

	if response == nil || response.Response.Value == nil {
		log.Fatal("Failed to get address",  "chain", chain)
	}

	value := response.Response.Value
	if value == nil || len(value) == 0 {
		log.Fatal("Returned address is empty", "chain", chain )
	}

	return string(value)
}

func GetBalance(accountKey id.AccountKey) data.Coin {
	request := action.Message("accountKey=" + hex.EncodeToString(accountKey))

	// Send out a query
	response := comm.Query("/balance", request)
	if response == nil {
		log.Fatal("Failed to get Balance", "response", response, "key", accountKey)
	}

	// Check the response
	value := response.Response.Value
	if value == nil || len(value) == 0 {
		log.Fatal("Failed to return Balance", "response", response, "key", accountKey)
	}
	if bytes.Compare(value, []byte("null")) == 0 {
		log.Fatal("Null Balance", "response", response, "key", "accountKey")
	}

	// Convert to a balance
	var balance data.Balance
	buffer, status := comm.Deserialize(value, &balance)
	if status != nil {
		log.Fatal("Deserialize", "status", status, "value", value)
	}
	if buffer == nil {
		log.Fatal("Can't deserialize", "response", response)
	}

	log.Debug("Deserialize", "buffer", buffer, "response", response, "value", value)
	return buffer.(*data.Balance).Amount
}
