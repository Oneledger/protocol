/*
	Copyright 2017 - 2018 OneLedger

	Handle arbitrary, but lossely typed parameters to the function calls
*/
package action

import (
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcutil"
	"github.com/Oneledger/protocol/node/id"
)

func GetInt(value FunctionValue) int {
	switch value.(type) {
	case int:
		return value.(int)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return 0
}

func GetInt64(value FunctionValue) int64 {
	switch value.(type) {
	case int64:
		return value.(int64)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return 0
}

func GetAmount(value FunctionValue) btcutil.Amount {
	switch value.(type) {
	case btcutil.Amount:
		return value.(btcutil.Amount)
	default:
		log.Fatal("Bad Type Cast in Function Parameter")
	}
	return 0
}

func GetRole(value FunctionValue) Role {
	switch value.(type) {
	case Role:
		return value.(Role)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return 0
}

func GetString(value FunctionValue) string {
	switch value.(type) {
	case string:
		return value.(string)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return ""
}

func GetCoin(value FunctionValue) data.Coin {
	switch value.(type) {
	case data.Coin:
		return value.(data.Coin)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return data.Coin{}
}

func GetParty(value FunctionValue) Party {
	switch value.(type) {
	case Party:
		return value.(Party)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return Party{nil, nil}
}

func GetAccountKey(value FunctionValue) id.AccountKey {
	switch value.(type) {
	case id.AccountKey:
		return value.(id.AccountKey)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return nil
}

func GetBytes(value FunctionValue) []byte {
	switch value.(type) {
	case []byte:
		return value.([]byte)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return []byte(nil)
}

func GetETHContract(value FunctionValue) *ethereum.HtlContract {
	switch value.(type) {
	case ethereum.HtlContract:
		return value.(*ethereum.HtlContract)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return nil
}
