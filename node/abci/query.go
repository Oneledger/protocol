/*
	Copyright 2017-2018 OneLedger

	Wrap a query.
*/
package abci

import "github.com/Oneledger/protocol/node/comm"

type Query struct {
	account string
}

func (query Query) JSON() comm.Message {
	return nil
}
