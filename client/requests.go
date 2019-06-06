/*
	Copyright 2017-2018 OneLedger

	Cover over the arguments of client requests
*/
package client

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/data/keys"
)

type ApplyValidatorArguments struct {
	ValidatorAddress     keys.Address
	Name 		string
	Amount 		string
	Purge  		bool
	TmAddress   []byte
	TmPubKey    keys.Address
}

type SendArguments struct {
	Party        []byte
	CounterParty []byte
	Amount       action.Amount
	Fee          action.Amount
	Gas          int64
	CurrencyStr  string
	AmountStr    string
	FeeStr       string
}

type SwapArguments struct {
	Party        string
	CounterParty string
	Amount       float64
	Currency     string
	Exchange     float64
	Excurrency   string
	Nonce        int64

	Fee float64
	Gas int64
}

type ExSendArguments struct {
	SenderId        string
	ReceiverId      string
	SenderAddress   string
	ReceiverAddress string
	Currency        string
	Amount          float64

	Gas int64
	Fee float64

	Chain string
	ExGas string
	ExFee string
}
