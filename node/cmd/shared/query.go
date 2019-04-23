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

func GetAccountKey(ctx comm.ClientContext, identity string) id.AccountKey {
	request := action.Message("Identity=" + identity)
	response := ctx.Query("/accountKey", request)
	if response == nil {
		log.Warn("Query returned nothing", "request", request)
		return nil
	}
	switch param := response.(type) {
	case *id.AccountKey:
		return *param

	case string:
		log.Warn("Query Error:", "err", param, "request", request, "response", response)
		return nil
	}
	log.Warn("Query Unknown Type:", "response", response)
	return nil
}

func GetCurrencyAddress(ctx comm.ClientContext, currencyName string, id string) []byte {
	request := action.Message("currency=" + currencyName + "|" + "id=" + id)
	response := ctx.Query("/currencyAddress", request)
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

func GetNodeName(ctx comm.ClientContext) string {
	request := action.Message("")
	response := ctx.Query("/nodeName", request)
	if response == nil {
		return "Unknown"
	}
	return response.(string)
}

// TODO: Return a balance, not a coin
func GetBalance(ctx comm.ClientContext, accountKey id.AccountKey) *data.Balance {
	// TODO: This is wrong, should pass by type, not encode/decode
	request := action.Message("accountKey=" + hex.EncodeToString(accountKey))
	response := ctx.Query("/balance", request)
	if response == nil {
		// New Accounts don't have a balance yet.
		result := data.NewBalance()
		return result
	}
	if serial.GetBaseType(response).Kind() == reflect.String {
		log.Error("Error:", "response", response)
		return nil
	}
	balance := response.(*data.BalanceData).Primitive()
	return balance.(*data.Balance)
}

func GetTxByHash(ctx comm.ClientContext, hash []byte) *ctypes.ResultTx {
	response := ctx.Tx(hash, true)
	if response == nil {
		log.Error("Search tx by hash failed", "hash", hash)
		return nil
	}
	return response
}

func GetTxByHeight(ctx comm.ClientContext, height int) *ctypes.ResultTxSearch {
	request := "tx.height=" + convert.GetString(height)

	response := ctx.Search(request, true, 1, 100)
	if response == nil {
		log.Error("Search tx by height failed", "request", request)
	}
	return response
}

func GetTxByType(ctx comm.ClientContext, txType string) *ctypes.ResultTxSearch {
	request := "tx.type" + txType

	response := ctx.Search(request, true, 1, 100)
	if response == nil {
		log.Error("Search tx by hash failed", "request", request)
		return nil
	}
	return response
}

func GetSequenceNumber(ctx comm.ClientContext, accountKey id.AccountKey) int64 {
	request := action.Message("AccountKey=" + hex.EncodeToString(accountKey))
	response := ctx.Query("/sequenceNumber", request)
	if response == nil {
		log.Warn("Query returned nothing", "request", request)
		return -1
	}
	switch value := response.(type) {
	case id.SequenceRecord:
		return value.Sequence
	case string:
		Console.Error(value)
	}

	return 0
}
