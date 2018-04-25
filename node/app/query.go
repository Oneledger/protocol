/*
	Copyright 2017-2018 OneLedger

	Wrap a query.
*/
package app

type Query struct {
	account string
}

func (query Query) JSON() Message {
	return nil
}
