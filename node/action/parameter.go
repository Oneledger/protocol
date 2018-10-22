/*
	Copyright 2017 - 2018 OneLedger

	Handle arbitrary, but lossely typed parameters to the function calls
*/
package action

import (
	"github.com/Oneledger/protocol/node/chains/bitcoin"
	"github.com/Oneledger/protocol/node/chains/ethereum"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/btcsuite/btcutil"
)

type Parameter byte

func init() {
	serial.Register(Parameter(0))
}

// TODO: Move to parameter.go
const (
	ROLE Parameter = iota
	MY_ACCOUNT
	THEM_ACCOUNT

	AMOUNT
	EXCHANGE
	NONCE
	PREIMAGE

	PASSWORD
	ETHCONTRACT
	BTCCONTRACT

	LOCKTIME
	COUNT
	CHAINID
	NEXTCHAINNAME

	EVENTTYPE
	CHECK_EXTERNAL_CHAIN
	STOREMESSAGE
	STOREKEY
	STAGE
	OWNER
	TARGET
)

func GetInt(value FunctionValue) int {
	if value == nil {
		return 0
	}
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
	if value == nil {
		return ""
	}
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

func GetByte32(value FunctionValue) [32]byte {
	switch value.(type) {
	case [32]byte:
		return value.([32]byte)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return [32]byte{}
}

func GetETHContract(value FunctionValue) *ethereum.HTLContract {
	switch value.(type) {
	case *ethereum.HTLContract:
		return value.(*ethereum.HTLContract)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return nil
}

func GetBTCContract(value FunctionValue) *bitcoin.HTLContract {
	switch value.(type) {
	case *bitcoin.HTLContract:
		return value.(*bitcoin.HTLContract)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return nil
}

func GetType(value FunctionValue) Type {
	switch value.(type) {
	case Type:
		return value.(Type)
	default:
		log.Fatal("Bad Type Cast in FUnction Parameter", "value", value)
	}
	return INVALID
}

func GetChain(value FunctionValue) data.ChainType {
	switch value.(type) {
	case data.ChainType:
		return value.(data.ChainType)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return data.UNKNOWN
}
