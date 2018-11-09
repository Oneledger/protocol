/*
	Copyright 2017 - 2018 OneLedger

	Query the chain for answers
*/
package shared

import (
	"encoding/hex"
	"reflect"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
)

func GetAccountKey(identity string) []byte {
	request := action.Message("Identity=" + identity)
	response := comm.Query("/accountKey", request)
	if response == nil {
		log.Warn("Query returned nothing", "request", request)
		return nil
	}
	//result := response.([]uint8)
	result := response.(id.AccountKey)
	return result
}

func GetSwapAddress(currencyName string) []byte {
	request := action.Message("currency=" + currencyName)
	response := comm.Query("/swapAddress", request)
	if response == nil {
		return nil
	}

	/*
		value := response.Response.Value
		if value == nil || len(value) == 0 {
			log.Fatal("Returned address is empty", "chain", currencyName)
		}
	*/

	return response.([]byte)
}

func GetNodeName() string {
	request := action.Message("")
	response := comm.Query("/nodeName", request)
	if response == nil {
		return "Unknown"
	}
	return response.(string)
}

// TODO: Return a balance, not a coin
func GetBalance(accountKey id.AccountKey) *data.Coin {
	request := action.Message("accountKey=" + hex.EncodeToString(accountKey))
	response := comm.Query("/balance", request)
	if response == nil {
		// New Accounts don't have a balance yet.
		result := data.NewCoin(0, "OLT")
		return &result
	}
	if serial.GetBaseType(response).Kind() == reflect.String {
		log.Error("Error:", "response", response)
		return nil
	}
	balance := response.(*data.Balance)
	return &balance.Amount

}
