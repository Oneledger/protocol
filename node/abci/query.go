/*
	Copyright 2017-2018 OneLedger

	Wrap a query.
*/
package abci

import "github.com/Oneledger/protocol/node/serial"

type Query struct {
	account string
}

func (query Query) JSON() serial.Message {
	return nil
}
