/*
	Copyright 2017 - 2018 OneLedger

	Query the chain for answers
*/
package shared

import (
	"encoding/hex"
	"github.com/Oneledger/protocol/node/convert"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

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
	switch param := response.(type) {
	case id.AccountKey:
		return param

	case string:
		log.Warn("Query Error:", "err", param, "request", request, "response", response)
		return nil
	}
	log.Warn("Query Unknown Type:", "response", response)
	return nil
}

func GetCurrencyAddress(currencyName string, id string) []byte {
	request := action.Message("currency=" + currencyName + "|" + "id=" + id)
	response := comm.Query("/currencyAddress", request)
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
func GetBalance(accountKey id.AccountKey) *data.Balance {
	// TODO: This is wrong, should pass by type, not encode/decode
	request := action.Message("accountKey=" + hex.EncodeToString(accountKey))
	response := comm.Query("/balance", request)
	if response == nil {
		// New Accounts don't have a balance yet.
		result := data.NewBalance()
		return &result
	}
	if serial.GetBaseType(response).Kind() == reflect.String {
		log.Error("Error:", "response", response)
		return nil
	}
	balance := response.(*data.Balance)

	return balance

}

func GetTxByHash(hash []byte) *ctypes.ResultTx {
	response := comm.Tx(hash, true)
	if response == nil {
		log.Error("Search tx by hash failed", "hash", hash)
		return nil
	}
	return response
}

func GetTxByHeight(height int) *ctypes.ResultTxSearch {
	request := "tx.height=" + convert.GetString(height)

	response := comm.Search(request, true, 1, 100)
	if response == nil {
		log.Error("Search tx by height failed", "request", request)
	}
	return response
}

func GetTxByType(t string) *ctypes.ResultTxSearch {
	request := "tx.type" + t

	response := comm.Search(request, true, 1, 100)
	if response == nil {
		log.Error("Search tx by hash failed", "request", request)
		return nil
	}
	return response
}

func GetSequenceNumber(accountKey id.AccountKey) int64 {
	request := action.Message("AccountKey=" + hex.EncodeToString(accountKey))
	response := comm.Query("/sequenceNumber", request)
	if response == nil {
		log.Warn("Query returned nothing", "request", request)
		return -1
	}

	result := response.(id.SequenceRecord)
	return result.Sequence
}
