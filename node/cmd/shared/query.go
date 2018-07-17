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
		log.Error("No Response from Node", "identity", identity)
		return nil
	}

	value := response.Response.Value
	if value == nil || len(value) == 0 {
		log.Error("Key is Missing", "identity", identity)
		return nil
	}

	key, status := hex.DecodeString(string(value))
	if status != nil {
		log.Error("Decode Failed", "identity", identity, "value", value)
		return nil
	}

	return key
}

func GetSwapAddress(currencyName string) []byte {
	request := action.Message("currency=" + currencyName)
	response := comm.Query("/swapAddress", request)

	if response == nil || response.Response.Value == nil {
		log.Fatal("Failed to get address",  "chain", currencyName)
	}

	value := response.Response.Value
	if value == nil || len(value) == 0 {
		log.Fatal("Returned address is empty", "chain", currencyName )
	}

	return value
}

func GetBalance(accountKey id.AccountKey) data.Coin {
	request := action.Message("accountKey=" + hex.EncodeToString(accountKey))

	// Send out a query
	response := comm.Query("/balance", request)
	if response == nil {
		log.Error("Failed to get Balance", "response", response, "key", accountKey)
		return nil
	}

	// Check the response
	value := response.Response.Value
	if value == nil || len(value) == 0 {
		log.Error("Failed to return Balance", "response", response, "key", accountKey)
		return nil
	}
	if bytes.Compare(value, []byte("null")) == 0 {
		log.Error("Null Balance", "response", response, "key", "accountKey")
		return nil
	}

	// Convert to a balance
	var balance data.Balance
	buffer, status := comm.Deserialize(value, &balance)
	if status != nil {
		log.Error("Deserialize", "status", status, "value", value)
		return nil
	}
	if buffer == nil {
		log.Error("Can't deserialize", "response", response)
		return nil
	}

	log.Debug("Deserialize", "buffer", buffer, "response", response, "value", value)
	return &(buffer.(*data.Balance).Amount)
}
