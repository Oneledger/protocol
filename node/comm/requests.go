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

type SendArguments struct {
	Party        string
	CounterParty string
	Currency     string
	Amount       string
	Gas          string
	Fee          string
}

type SwapArguments struct {
	Party        string
	CounterParty string
	Amount       string
	Currency     string
	Fee          string
	Gas          string // TODO: Not sure this is necessary, unless the chain is like Ethereum
	Exchange     string
	Excurrency   string
	Nonce        int64
}

type ExSendArguments struct {
	SenderId        string
	ReceiverId      string
	SenderAddress   string
	ReceiverAddress string
	Currency        string
	Amount          string
	Gas             string
	Fee             string
	Chain           string
	ExGas           string
	ExFee           string
}

func init() {
	serial.Register(ApplyValidatorArguments{})
	serial.Register(ExSendArguments{})
	serial.Register(SendArguments{})
	serial.Register(SwapArguments{})
}
