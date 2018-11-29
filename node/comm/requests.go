/*
	Copyright 2017-2018 OneLedger

	Cover over the arguments of client requests
*/
package comm

import "github.com/Oneledger/protocol/node/serial"

type ApplyValidatorArguments struct {
	Id           string
	Amount       string
}

func init() {
	serial.Register(ApplyValidatorArguments{})
}
