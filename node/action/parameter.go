/*
	Copyright 2017 - 2018 OneLedger

	Handle arbitrary, but lossely typed parameters to the function calls
*/
package action

import "github.com/Oneledger/protocol/node/log"

func GetInt(value FunctionValue) int {
	switch value.(type) {
	case int:
		return value.(int)
	default:
		log.Fatal("Bad Type Cast in Function Parameter", "value", value)
	}
	return 0
}
